package libs

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"log"
	"net/url"
	"path"
)

func Screenshot(urlStr string, content string) (string, error) {
	myUrl, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	outPath := fmt.Sprintf("./images/%s.jpeg", path.Base(myUrl.Path))

	chromedp.Run(ctx,
		chromedp.Navigate(urlStr),
		chromedp.ActionFunc(func(ctxt context.Context) error {
			_, _, contentRect, err := page.GetLayoutMetrics().Do(ctxt)
			if err != nil {
				return err
			}

			v := page.Viewport{
				X:      40,
				Y:      250,
				Width:  contentRect.Width - 80,
				Height: 380,
				Scale:  1,
			}
			log.Printf("Capture %#v", v)
			buf, err := page.CaptureScreenshot().WithClip(&v).Do(ctxt)
			if err != nil {
				return err
			}
			log.Printf("Write %v", outPath)
			ioutil.WriteFile(outPath, buf, 0644)
			return nil
		}))
	return outPath, nil
}
