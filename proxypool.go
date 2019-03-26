package main

import (
	"github.com/dashengbuqi/proxypool/consumer"
	"github.com/dashengbuqi/proxypool/crawl"
	"github.com/dashengbuqi/proxypool/http"
	"github.com/dashengbuqi/proxypool/models"
	"runtime"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	runtime.GOMAXPROCS(runtime.NumCPU())

	//开始http服务
	wg.Add(1)
	go func() {
		http.Run()
		wg.Done()
	}()

	//IP回收
	/*wg.Add(1)
	go func() {
		gc.Run()
		wg.Done()
	}()*/

	proxyIPChan := make(chan *models.IProxyItem, 1000)

	//开始http服务
	wg.Add(1)
	go func() {
		http.Run()
		wg.Done()
	}()
	//消费消息
	wg.Add(1)
	go func() {
		consumer.ProxyIpSpeed(proxyIPChan)
		wg.Done()
	}()
	//爬取IP地址
	wg.Add(1)
	go func() {
		for {
			if len(proxyIPChan) < 100 {
				crawl.Run(proxyIPChan)
			}
			time.Sleep(10 * time.Minute)
		}
		wg.Done()
	}()
	wg.Wait()
}
