// Package corrlinker provides ...
package corrlinker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

type inboxMessage struct {
	From            string `json:"from"`
	Subject         string `json:"subject"`
	SentDate        string `json:"sent_date"`
	Message         string `json:"message"`
	MessageID       uint64 `json:"message_id"`
	IsUnreadMessage bool   `json:"is_unread_message"`
}

var inboxMessages []inboxMessage

func readInboxPage(ctx context.Context, cfg *config) {
	// var pageCountstr string
	var pages []*cdp.Node
	if err := chromedp.Run(ctx, chromedp.Tasks{
		// chromedp.Emulate(device.IPadProlandscape),
		chromedp.Navigate(cfg.URL + "/Inbox.aspx"),
		chromedp.WaitVisible(`#ctl00_mainContentPlaceHolder_MessagesViewPanel`, chromedp.ByID),
		chromedp.SetValue(`#ctl00_mainContentPlaceHolder_startDateTextBox`, "1/1/2021", chromedp.ByID),
		chromedp.SetValue(`#ctl00_mainContentPlaceHolder_endDateTextBox`, "12/31/2021", chromedp.ByID),
		chromedp.Click(`#ctl00_mainContentPlaceHolder_updateButton`, chromedp.ByID),
		chromedp.Sleep(1 * time.Second),
		chromedp.Nodes(`#ctl00_mainContentPlaceHolder_inboxGridView tbody tr.Pager td table tbody tr td`, &pages, chromedp.ByQueryAll),
		// chromedp.ActionFunc(func(c context.Context) error {
		// log.Println("Found rows", len(pages))
		// loop through all the td nodes
		// const sel = `/html/body/table/tbody/tr[%d]/td[%d]`
		// var firstCol, secondCol string
		// for i := 1; i <= len(rows); i++ {
		//     if err := chromedp.Run(ctx,
		//         chromedp.Text(fmt.Sprintf(sel, i, 1), &firstCol),
		//         chromedp.Text(fmt.Sprintf(sel, i, 2), &secondCol),
		//     ); err != nil {
		//         log.Fatal(err)
		//     }
		//     fmt.Printf("got first column value = %s, from row %d\n", firstCol, i)
		//     fmt.Printf("got second column value = %s, from row %d\n", secondCol, i)
		// }
		// return nil
		// }),
		// chromedp.Text(`#ctl00_mainContentPlaceHolder_inboxGridView tbody tr.Pager td table tbody tr td:nth-child(%v) a`, &pageCountstr, chromedp.ByID),
	}); err != nil && err != context.Canceled && err != context.DeadlineExceeded {
		log.Fatal(err)
	}
	// pageCount, _ := strconv.Atoi(pageCountstr)
	pageCount := len(pages)
	for pageNum := 1; pageNum <= pageCount; pageNum++ {
		err := readInbox(ctx, pageCount, pageNum)
		if err != nil && err != context.Canceled && err != context.DeadlineExceeded {
			log.Fatalf("COULD NOT READ INBOX:\n %v", err)
		}
	}

	outputJSON, err := json.Marshal(&inboxMessages)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(outputJSON))
}

func readInbox(ctx context.Context, pageCount, pageNum int) error {
	var inboxContent string

	// read inbox page
	if err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.WaitVisible(`#ctl00_mainContentPlaceHolder_inboxGridView`, chromedp.ByID),
		chromedp.OuterHTML(`#ctl00_mainContentPlaceHolder_inboxGridView`, &inboxContent, chromedp.ByID),
	}); err != nil && err != context.Canceled && err != context.DeadlineExceeded {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(inboxContent))
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("table tr").Has(`#ctl00_mainContentPlaceHolder_inboxGridView .MessageDataGrid`).Each(func(_ int, tr *goquery.Selection) {
		m := inboxMessage{}

		from := tr.Find("th span").Text()
		subject := tr.Find("td span").Eq(1).Text()
		sentDate := strings.TrimSpace(tr.Find("td").Eq(2).Text())
		messageID, _ := tr.Find("td span").Attr("messageid")
		messageIDint, _ := strconv.ParseUint(messageID, 10, 64)
		isUnreadMessage, _ := tr.Find("td span").Attr("isunreadmessage")
		isUnreadMessageBool, _ := strconv.ParseBool(isUnreadMessage)

		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		readLink, _ := tr.Find("td a[id*=readLink]").Attr("id")
		var msgFound, msgMatch string
		if err := chromedp.Run(ctx, chromedp.Tasks{
			chromedp.Click(readLink, chromedp.ByID),
			chromedp.WaitVisible(`#ctl00_mainContentPlaceHolder_messagePanel`, chromedp.ByID),
			chromedp.TextContent(`#ctl00_mainContentPlaceHolder_messageTextBox`, &msgFound, chromedp.ByID),
			chromedp.ActionFunc(func(c context.Context) error {
				// parse message
				// regex match on ex: "lastname, firstname on 5/18/2021 5:21 PM wrote"
				r := regexp.MustCompile(`((?s)\A.*?)(?-s:.+?on\s\d{1,2}/\d{1,2}/\d{4}\s\d{1,2}:\d{1,2}\s[Pp|Aa][Mm]\swrote|\z)`)
				msgMatch = r.FindStringSubmatch(msgFound)[1]
				return nil
			}),
			chromedp.Click(`#ctl00_mainContentPlaceHolder_cancelButton`, chromedp.ByID),
			chromedp.WaitVisible(`#ctl00_mainContentPlaceHolder_MessagesViewPanel`),
		}); err != nil {
			log.Fatalf("COULD NOT READ MESSAGE:\n %v", err)
		}
		// log.Printf("\nFROM: %s\nSUBJECT: %s\nSENT DATE: %s\nMESSAGE ID: %v\nUNREAD: %v\nMESSAGE: %v\n\n",
		//     from, subject, sentDate, messageIDint, isUnreadMessageBool, msgMatch)

		m.From = from
		m.Subject = subject
		m.SentDate = sentDate
		m.Message = msgMatch
		m.MessageID = messageIDint
		m.IsUnreadMessage = isUnreadMessageBool
		inboxMessages = append(inboxMessages, m)
		// fmt.Println(message)
	})
	// fmt.Println(inboxMessages)

	if pageNum == pageCount {
		return err
	} else {
		if err := chromedp.Run(ctx, chromedp.Tasks{
			chromedp.Click(fmt.Sprintf(`#ctl00_mainContentPlaceHolder_inboxGridView tbody tr.Pager td table tbody tr td:nth-child(%v) a`,
				pageNum+1), chromedp.ByID),
			chromedp.Sleep(1 * time.Second),
		}); err != nil && err != context.Canceled && err != context.DeadlineExceeded {
			log.Fatal(err)
		}
	}
	return err
}
