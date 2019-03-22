package models

type ProxyItem struct {
	Addr      string `json:"addr"`
	Scheme    string `json:"scheme"`
	Port      int64  `json:"port"`
	Speed     int64  `json:"speed"`
	UpdatedBy int64  `json:"updated_by"`
}

func NewProxyItem() *ProxyItem {
	return &ProxyItem{}
}
