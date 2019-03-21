package models

import "sync"

type ProxyItem struct {
	Addr   string `json:"addr"`
	Scheme string `json:"scheme"`
	Speed  int64  `json:"speed"`
}

type ProxyItems struct {
	items []*ProxyItem
	mu    sync.RWMutex
}
