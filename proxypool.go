package main

import (
	"github.com/dashengbuqi/proxypool/crawl"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		crawl.Run()
		wg.Done()
	}()
	wg.Wait()
}
