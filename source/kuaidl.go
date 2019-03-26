package source

import (
	"fmt"
	"github.com/dashengbuqi/proxypool/models"
	"github.com/gocolly/colly"
	"strconv"
	"strings"
	"time"
)

func Kuaidl() (result []*models.ProxyItem) {
	c := colly.NewCollector()
	c.AllowURLRevisit = false

	// On every a element which has href attribute call callback
	c.OnHTML(".table-striped", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, el *colly.HTMLElement) {
			if i > 0 {
				rs := &models.ProxyItem{}
				el.ForEach("td", func(i int, ele *colly.HTMLElement) {
					switch i {
					case 0:
						rs.Addr = ele.Text
					case 1:
						rs.Port, _ = strconv.ParseInt(ele.Text, 10, 64)
					case 3:
						rs.Scheme = strings.ToLower(ele.Text)
					case 5:
						rs.Speed = extractSpeed(ele.Text)
					case 6:
						rs.UpdatedBy = time.Now().Unix()
					}
				})
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
	c.Visit("http://www.kuaidaili.com/free/inha/")
	return
}
