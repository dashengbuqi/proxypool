package crawl

import (
	"github.com/dashengbuqi/proxypool/models"
	"github.com/dashengbuqi/proxypool/source"
	"log"
	"sync"
)

//启动爬虫
func Run() {
	var wg sync.WaitGroup
	log.Println("启动爬虫..")
	funs := []func() []*models.ProxyItem{
		source.Feiyi,
	}

	for _, f := range funs {
		wg.Add(1)
		go func(f func() []*models.ProxyItem) {
			temp := f()
			for _, v := range temp {
				log.Println(v)
			}
			wg.Done()
		}(f)
	}
	log.Println("结束爬虫..")
	wg.Wait()
}
