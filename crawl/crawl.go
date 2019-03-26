package crawl

import (
	"github.com/dashengbuqi/proxypool/models"
	"github.com/dashengbuqi/proxypool/source"
	"log"
	"sync"
)

//启动爬虫
func Run(inIpChan chan<- *models.IProxyItem) {
	var wg sync.WaitGroup
	log.Println("启动爬虫..")
	funs := []func() []*models.ProxyItem{
		source.Kuaidl,
		source.Feiyi,
		source.IP66,
		source.PLP,
		source.IP89,
	}

	for _, f := range funs {
		wg.Add(1)
		go func(f func() []*models.ProxyItem) {
			temp := f()
			proxyItemMap := make(map[int]*models.ProxyItem)
			for k, v := range temp {
				proxyItemMap[k] = v
			}
			if len(proxyItemMap) > 0 {
				pi := models.NewProxyItem(proxyItemMap)
				inIpChan <- &pi
			}
			wg.Done()
		}(f)
	}
	log.Println("结束爬虫..")
	wg.Wait()
}
