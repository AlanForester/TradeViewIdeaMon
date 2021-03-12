package libs

import (
	"fmt"
	sql2 "github.com/AlexCollin/TradeViewIdeaMon/sql"
	"github.com/gocolly/colly/v2"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"sync"
)

type Scraper struct {
	Count int
	Wg    sync.WaitGroup
}

func (m *Scraper) GetLastIdeas(newPost chan sql2.Post, link string) []sql2.Post {
	var res []sql2.Post
	c := colly.NewCollector()

	defer m.Wg.Wait()

	c.OnHTML("a.tv-widget-idea__title[href]", func(e *colly.HTMLElement) {
		if m.Count >= 3 {
			return
		}

		m.Count++
		e.Request.Visit(e.Attr("href"))
	})

	c.OnHTML(".tv-chart-view__section:has(.tv-chart-view__header)", func(e *colly.HTMLElement) {
		m.Wg.Add(1)
		go func() {
			curUrl := e.Request.URL.String()
			data := sql2.Post{}

			db := sql2.DB.First(&data, "url = ?", curUrl)
			if db.RowsAffected == 0 {
				data.Url = curUrl
				data.Title = e.ChildText("h1.tv-chart-view__title-name")
				data.Author = e.ChildText(".tv-chart-view__title-user-name")
				data.Tp = e.ChildText(".tv-chart-view__title-icons")
				data.Pair = e.ChildText("a.tv-chart-view__symbol-link.tv-chart-view__symbol--desc")
				data.Date, _ = e.DOM.Find(".tv-chart-view__title-time").Attr("data-timestamp")
				data.Descr = e.ChildText(".tv-chart-view__description")
				vid := e.ChildAttr("video.tv-chart-view__video", "src")
				log.Printf("Video %v", vid)
				if vid != "" {
					fn, err := m.DownloadFile(vid)
					if err != nil {
						log.Printf("Error on download video: %v", err)
					} else {
						data.Video = fn
					}
				} else {
					path, err := Screenshot(curUrl, "")
					if err != nil {
						log.Printf("Error on download image: %v", err)
					} else {
						data.Video = path
					}
				}
				res = append(res, data)
				sql2.DB.Save(&data)

				var author sql2.Author
				db := sql2.DB.First(&author, "name = ?", data.Author)
				if db.RowsAffected == 0 {
					author.Name = data.Author
					sql2.DB.Save(&author)
				}

				newPost <- data
			}
			m.Wg.Done()
		}()
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(link)
	m.Wg.Wait()

	return res
}

func (m *Scraper) DownloadFile(fileUrl string) (string, error) {

	parsedUrl, _ := url.Parse(fileUrl)

	// Get the data
	resp, err := http.Get(fileUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Create the file
	outPath := fmt.Sprintf("./video/%s", path.Base(parsedUrl.Path))
	out, err := os.Create(outPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return outPath, err
}
