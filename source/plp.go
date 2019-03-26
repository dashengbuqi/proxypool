package source

import (
	"fmt"
	"github.com/dashengbuqi/proxypool/models"
	"github.com/gocolly/colly"
	"strconv"
	"time"
)

func PLP() (result []*models.ProxyItem) {
	c := colly.NewCollector()
	c.AllowURLRevisit = false

	// On every a element which has href attribute call callback
	c.OnHTML(".bg .cells", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, el *colly.HTMLElement) {
			if i > 0 {
				rs := &models.ProxyItem{}
				el.ForEach("td", func(i int, ele *colly.HTMLElement) {
					switch i {
					case 1:
						rs.Addr = ele.Text
					case 2:
						rs.Port, _ = strconv.ParseInt(ele.Text, 10, 64)
					case 5:
						rs.Scheme = extractScheme(ele.Text)
					}
				})
				rs.Speed = -1
				rs.UpdatedBy = time.Now().Unix()
				fmt.Println(rs.Speed, rs.Scheme, rs.Port, rs.Addr, rs.UpdatedBy)
				result = append(result, rs)
			}
		})
	})

	c.OnHTML("#listnav", func(e *colly.HTMLElement) {
		pageUrl := e.ChildAttr("li a", "href")

		c.Visit(e.Request.AbsoluteURL(pageUrl))

	})
	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	c.Visit("https://list.proxylistplus.com/Fresh-HTTP-Proxy-List-1")
	return
}

func extractScheme(str string) string {
	if len(str) == 2 {
		return "HTTP"
	}
	return "HTTPS"
}
