package models

import (
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/dashengbuqi/proxypool/persistence"
	"github.com/parnurzeal/gorequest"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

const (
	DATABASE   = "sensor_spider"
	COLLECTION = "addr_pool"
)

type ProxyItem struct {
	Addr      string `json:"addr"`
	Scheme    string `json:"scheme"`
	Port      int64  `json:"port"`
	Speed     int64  `json:"speed"`
	UpdatedBy int64  `json:"updated_by" bson:"updated_by"`
}

type IProxyItem interface {
	CheckSpeedAndUpdate()
}

func NewProxyItem(items map[int]*ProxyItem) IProxyItem {
	return &MongodbProxyItem{
		Items: items,
	}
}

type MongodbProxyItem struct {
	Items map[int]*ProxyItem
	mu    sync.RWMutex
}

func RandomProxy() *ProxyItem {
	query := bson.M{"speed": bson.M{"$lt": 1000}}

	var result []*ProxyItem
	err := persistence.FindAll(DATABASE, COLLECTION, query, nil, &result)
	if err != nil {
		return nil
	}
	ips := len(result)
	if ips == 0 {
		return nil
	}
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(ips)

	return result[randomNum]
}

func RandomHttpProxy() *ProxyItem {
	query := bson.M{"scheme": "http", "speed": bson.M{"$lt": 1000}}

	var result []*ProxyItem
	err := persistence.FindAll(DATABASE, COLLECTION, query, nil, &result)
	if err != nil {
		return nil
	}
	ips := len(result)
	if ips == 0 {
		return nil
	}
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(ips)

	return result[randomNum]
}

func RandomHttpsProxy() *ProxyItem {
	query := bson.M{"scheme": "https", "speed": bson.M{"$lt": 1000}}

	var result []*ProxyItem
	err := persistence.FindAll(DATABASE, COLLECTION, query, nil, &result)
	if err != nil {
		return nil
	}
	ips := len(result)

	if ips == 0 {
		return nil
	}
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(ips)

	return result[randomNum]
}

//检查代理速度
func (m *MongodbProxyItem) CheckSpeedAndUpdate() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, pItem := range m.Items {
		if speed, ok := testIPSpeed(pItem); ok == false {
			continue
		} else {
			pItem.Speed = speed
		}
		proxyItem := &ProxyItem{
			Addr:      pItem.Addr,
			Scheme:    pItem.Scheme,
			Port:      pItem.Port,
			Speed:     pItem.Speed,
			UpdatedBy: time.Now().Unix(),
		}
		err := proxyItem.insertOrUpdate()
		if err != nil {
			fmt.Println(err.Error())
		}
	}

}

func (item *ProxyItem) insertOrUpdate() error {
	//检查ip，端口是否存在
	query := bson.M{"addr": item.Addr, "port": item.Port, "scheme": item.Scheme}

	var result ProxyItem
	err := persistence.FindOne(DATABASE, COLLECTION, query, nil, &result)

	if err != nil {
		if err1 := item.insertItem(); err1 != nil {
			return errors.New("增加ip失败" + err1.Error())
		}
	} else {
		if err1 := item.updateItem(&result); err1 != nil {
			return errors.New("更新ip失败" + err1.Error())
		}
	}
	return nil
}

func (item *ProxyItem) updateItem(old *ProxyItem) error {
	if item.Scheme == old.Scheme && item.Addr == old.Addr && item.Port == old.Port && item.Speed == old.Speed {
		return nil
	}
	err := persistence.Update(DATABASE, COLLECTION, bson.M{"addr": item.Addr, "port": item.Port, "scheme": item.Scheme}, item)
	if err != nil {
		return err
	}
	return nil
}

func (item *ProxyItem) insertItem() error {
	err := persistence.Insert(DATABASE, COLLECTION, item)
	if err != nil {
		return err
	}
	return nil
}

func testIPSpeed(item *ProxyItem) (int64, bool) {
	var speedUrl string
	var speedIP string
	if item.Scheme == "https" {
		speedIP = "https://" + item.Addr + ":" + strconv.FormatInt(item.Port, 10)
		speedUrl = "https://httpbin.org/get"
	} else {
		speedIP = "https://" + item.Addr + ":" + strconv.FormatInt(item.Port, 10)
		speedUrl = "http://httpbin.org/get"
	}
	fmt.Println(speedIP)

	begin := time.Now()
	resp, _, errs := gorequest.New().Proxy(speedIP).Get(speedUrl).End()
	if errs != nil {
		fmt.Printf("[CheckIP] testIP = %s, pollURL = %s: Error = %v", speedIP, speedUrl, errs)
		return 0, false
	}

	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		_, err := simplejson.NewFromReader(resp.Body)
		if err != nil {
			fmt.Printf("[CheckIP] testIP = %s, pollURL = %s: Error = %v", speedIP, speedUrl, err)
			return 0, false
		}

		//harrybi 计算该代理的速度，单位毫秒
		Speed := time.Now().Sub(begin).Nanoseconds() / 1000 / 1000
		return Speed, true
	}
	return 0, false
}

func (m *MongodbProxyItem) UpdateOrInsert() error {
	return nil
}

//删除指定时间内的数据
func RemoveProxyItem() error {
	tenAgoTime := time.Now().AddDate(0, 0, -1).Unix()

	//检查ip，端口是否存在

	//query := bson.M{"$or":[{"speed":bson.M{"$gt": 1000}},{"updatedby":}]}

	query := bson.M{"$or": []bson.M{bson.M{"speed": bson.M{"$gt": 1000}}, bson.M{"updated_by": bson.M{"$lt": tenAgoTime}}}}

	type tempProxyItem struct {
		Id        bson.ObjectId `json:"id" bson:"_id"`
		Addr      string        `json:"addr"`
		Scheme    string        `json:"scheme"`
		Port      int64         `json:"port"`
		Speed     int64         `json:"speed"`
		UpdatedBy int64         `json:"updated_by"`
	}

	var result []*tempProxyItem

	err := persistence.FindAll(DATABASE, COLLECTION, query, nil, &result)

	if err != nil {
		fmt.Printf(err.Error())
		return err
	} else {
		for _, v := range result {
			fmt.Println(v.Id)
			delErr := persistence.Remove(DATABASE, COLLECTION, bson.M{"_id": v.Id})
			if delErr != nil {
				fmt.Printf("IP回收失败, " + delErr.Error())
			}
		}
		fmt.Printf("成功回收了%d条代理IP", len(result))
	}
	return nil
}
