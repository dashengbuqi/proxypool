package application

import (
	"github.com/dashengbuqi/proxypool/crawl"
	"github.com/henson/proxypool/api"
	"github.com/kataras/golog"
	"sync"
)

type Application struct {
	Logger *golog.Logger
}

func NewApplication() *Application {
	return &Application{
		Logger: golog.Default,
	}
}

func (a *Application) Run() {
	var wg sync.WaitGroup
	//启动http
	wg.Add(1)
	go func() {
		api.Run(a.Logger)
		wg.Done()
	}()
	//启动爬虫
	go func() {
		crawl.NewProxySource().Run()
		wg.Done()
	}()
	wg.Wait()
}
