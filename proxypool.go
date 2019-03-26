package main

import (
	"github.com/dashengbuqi/proxypool/consumer"
	"github.com/dashengbuqi/proxypool/crawl"
	"github.com/dashengbuqi/proxypool/models"
	"runtime"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	runtime.GOMAXPROCS(runtime.NumCPU())

	proxyIPChan := make(chan *models.IProxyItem, 1000)

	wg.Add(1)
	go func() {
		consumer.ProxyIpSpeed(proxyIPChan)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		crawl.Run(proxyIPChan)
		wg.Done()
	}()
	wg.Wait()
}
