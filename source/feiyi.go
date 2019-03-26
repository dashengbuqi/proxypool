package source

import (
	"fmt"
	"github.com/dashengbuqi/proxypool/models"
	"github.com/gocolly/colly"
	"regexp"
	"strconv"
	"strings"
)

//feiyi get ip from feiyiproxy.com
func Feiyi() (result []*models.ProxyItem) {
	c := colly.NewCollector()

	// On every a element which has href attribute call callback
	c.OnHTML(".et_pb_code_1 table", func(e *colly.HTMLElement) {
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
					case 6:
						rs.Speed = extractSpeed(ele.Text)
					case 7:
						rs.UpdatedBy = 900
					}
				})
				result = append(result, rs)
			}
		})
	})
	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	c.Visit("http://www.feiyiproxy.com/?page_id=1457")
	return
}

func extractSpeed(oritext string) int64 {
	reg := regexp.MustCompile(`[0-9]\d*\.[0-9]\d*`)
	temp := reg.FindString(oritext)
	if len(temp) > 0 {
		speed, _ := strconv.ParseInt(temp, 10, 64)
		return speed
	}
	return -1
}
