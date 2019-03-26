package http

import (
	"encoding/json"
	"github.com/dashengbuqi/proxypool/config"
	"github.com/dashengbuqi/proxypool/models"
	"github.com/kataras/golog"
	"net/http"
)

var (
	logger = golog.Default
)

const VERSION = "/v1"

func Run() {

	mux := http.NewServeMux()

	mux.HandleFunc(VERSION+"/ip", ProxyHander)
	mux.HandleFunc(VERSION+"/https", FindHttpsHandler)
	mux.HandleFunc(VERSION+"/http", FindHttpHandler)

	AppAddr := config.Setting("AppAddr").(string)
	AppPort := config.Setting("AppPort").(string)

	logger.Info(AppAddr)

	logger.Infof("Starting Http Server http://%s:%s", AppAddr, AppPort)
	http.ListenAndServe(":"+AppPort, mux)
}

func ProxyHander(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("content-type", "application/json")
		b, err := json.Marshal(models.RandomProxy())
		if err != nil {
			return
		}
		w.Write(b)
	}
}

func FindHttpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("content-type", "application/json")
		b, err := json.Marshal(models.RandomHttpProxy())
		if err != nil {
			return
		}
		w.Write(b)
	} else {
		w.WriteHeader(201)
	}
}

func FindHttpsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("content-type", "application/json")
		b, err := json.Marshal(models.RandomHttpsProxy())
		if err != nil {
			return
		}
		w.Write(b)
	}
}
