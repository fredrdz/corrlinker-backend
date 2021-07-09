// Package corrlinker provides ...
package feeds

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"runtime"

	"github.com/chromedp/chromedp"
)

type config struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func readConfig() (*config, error) {
	_, filePath, _, _ := runtime.Caller(0)
	pwd := filePath[:len(filePath)-10]
	txt, err := ioutil.ReadFile(pwd + "/config.json")
	if err != nil {
		return nil, err
	}
	var cfg = new(config)
	if err := json.Unmarshal(txt, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func Mailbox() {
	cfg, err := readConfig()
	if err != nil {
		log.Fatalf("Could not read config file: %v", err)
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", false),
		// chromedp.UserDataDir(`/Users/lobotech/Library/Application Support/Google/Chrome/Profile 1`),
		chromedp.Flag("use-mock-keychain", true),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("enable-automation", true),
		chromedp.Flag("restore-on-startup", false),
		chromedp.Flag("disable-extensions", true),
		// chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("hide-scrollbars", true),
		chromedp.Flag("mute-audio", true),
		chromedp.UserAgent(randomUserAgent()),
		chromedp.WindowSize(1024, 768),
		// chromedp.ExecPath(`/Applications/Google Chrome.app/Contents/MacOS/Google Chrome`),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	if err := chromedp.Run(ctx, loginMailbox(cfg)); err != nil && err != context.Canceled && err != context.DeadlineExceeded {
		log.Fatal(err)
	}

	readInboxPage(ctx, cfg)
}
