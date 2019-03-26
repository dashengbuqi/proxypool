package consumer

import (
	"github.com/dashengbuqi/proxypool/models"
)

//从通道中读取代理ip
func ProxyIpSpeed(ipChan <-chan *models.IProxyItem) {
	for {
		select {
		case item := <-ipChan:
			go (*item).CheckSpeedAndUpdate()
		}
	}
}
