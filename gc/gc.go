package gc

import (
	"fmt"
	"github.com/dashengbuqi/proxypool/models"
	"sync"
	"time"
)

//删除超过1天未更新的IP地址
func Run() {
	var wg sync.WaitGroup

	ticker := time.NewTicker(time.Hour * 24)
	wg.Add(1)
	go func(t *time.Ticker) {
		for {
			select {
			case <-t.C:
				fmt.Println("开始执行ip回收...")
				models.RemoveProxyItem(1)
			}
		}
		wg.Done()
	}(ticker)

	wg.Wait()
}
