// Package corrlinker provides ...
package corrlinker

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/chromedp/chromedp"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:56.0) Gecko/20100101 Firefox/56.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Safari/537.36",
}

func randomUserAgent() string {
	rand.Seed(time.Now().Unix())
	randNum := rand.Int() % len(userAgents)
	return userAgents[randNum]
}

func getScrapeClient(proxyString interface{}) *http.Client {
	switch v := proxyString.(type) {
	case string:
		proxyUrl, _ := url.Parse(v)
		return &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	default:
		return &http.Client{}
	}
}

func scrapeClientRequest(searchURL string, proxyString interface{}) (*http.Response, error) {
	baseClient := getScrapeClient(proxyString)
	req, _ := http.NewRequest("GET", searchURL, nil)
	req.Header.Set("User-Agent", randomUserAgent())

	res, err := baseClient.Do(req)
	if res.StatusCode != 200 {
		err := fmt.Errorf("scraper received a non-200 status code suggesting a ban")
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return res, nil
}

func loginMailbox(cfg *config) chromedp.Tasks {
	return chromedp.Tasks{
		// chromedp.Emulate(device.IPadProlandscape),
		chromedp.Navigate(cfg.URL + "/Login.aspx"),
		chromedp.WaitVisible(`#ctl00_mainContentPlaceHolder_loginUserNameTextBox`, chromedp.ByID),
		chromedp.WaitVisible(`#ctl00_mainContentPlaceHolder_loginPasswordTextBox`, chromedp.ByID),
		chromedp.SendKeys(`#ctl00_mainContentPlaceHolder_loginUserNameTextBox`, cfg.Username, chromedp.ByID),
		chromedp.SendKeys(`#ctl00_mainContentPlaceHolder_loginPasswordTextBox`, cfg.Password, chromedp.ByID),
		chromedp.Click(`#ctl00_mainContentPlaceHolder_loginButton`, chromedp.ByID),
		chromedp.WaitVisible(`#ctl00_mainContentPlaceHolder_captchaPanel`, chromedp.ByID),
		chromedp.WaitReady(`iframe`, chromedp.ByQuery),
		chromedp.Sleep(300 * time.Millisecond),
		chromedp.Evaluate(`document.querySelector('iframe').contentWindow.document.getElementById("checkbox").click();`, nil),
		chromedp.WaitEnabled(`#ctl00_mainContentPlaceHolder_captchaFNameLNameSubmitButton`, chromedp.ByID),
		chromedp.Click(`#ctl00_mainContentPlaceHolder_captchaFNameLNameSubmitButton`, chromedp.ByID),
		chromedp.WaitNotPresent(`#ctl00_mainContentPlaceHolder_captchaFNameLNameSubmitButton`, chromedp.ByID),
	}
}
