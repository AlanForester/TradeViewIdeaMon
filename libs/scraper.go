package libs

import (
	"fmt"
	sql2 "github.com/AlexCollin/TradeViewIdeaMon/sql"
	"github.com/gocolly/colly/v2"
	"sync"
)

type Scraper struct {
	Count int
	Wg    sync.WaitGroup
}

func (m *Scraper) GetLastIdeas(newPost chan sql2.Post) []sql2.Post {
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
				path, err := Screenshot(curUrl, "")
				if err == nil {
					data.Image = path
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

	c.Visit("https://www.tradingview.com/ideas")
	m.Wg.Wait()

	return res
}
