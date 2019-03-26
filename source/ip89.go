package source

import (
	"github.com/dashengbuqi/proxypool/models"
	"github.com/go-clog/clog"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func IP89() (result []*models.ProxyItem) {
	var ExprIP = regexp.MustCompile(`((25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\:([0-9]+)`)
	pollURL := "http://www.89ip.cn/tqdl.html?api=1&num=100&port=&address=%E7%BE%8E%E5%9B%BD&isp="

	resp, err := http.Get(pollURL)
	if err != nil {
		clog.Warn(err.Error())
		return
	}

	if resp.StatusCode != 200 {
		clog.Warn(err.Error())
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	bodyIPs := string(body)
	ips := ExprIP.FindAllString(bodyIPs, 100)

	for index := 0; index < len(ips); index++ {
		rs := &models.ProxyItem{}
		proxy := strings.TrimSpace(ips[index])
		addr := proxy[:strings.LastIndex(proxy, ":")]
		port := proxy[strings.LastIndex(proxy, ":")+1:]

		rs.Addr = addr
		rs.Port, _ = strconv.ParseInt(port, 10, 64)
		rs.Scheme = "http"
		rs.Speed = -1

		result = append(result, rs)
	}
	return
}
