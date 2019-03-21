package http

import (
	"github.com/dashengbuqi/proxypool/config"
	"github.com/kataras/golog"
	"net/http"
)

var (
	logger  = golog.Default
	version string
)

func init() {
	version = config.Setting("APP_VER").(string)
}

func Run() {

	mux := http.NewServeMux()

	mux.HandleFunc(version+"/ip", ProxyHander)
	mux.HandleFunc(version+"/https", FindHttpsHandler)
	mux.HandleFunc(version+"/http", FindHttpHandler)

	AppAddr := config.Setting("AppAddr").(string)
	AppPort := config.Setting("AppPort").(string)

	logger.Infof("Starting Http Server http://%s:%s", AppAddr, AppPort)
	http.ListenAndServe(AppAddr+":"+AppPort, mux)
}

func ProxyHander(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("content-type", "application/json")

	}
}

func FindHttpHandler(w http.ResponseWriter, r *http.Request) {

}

func FindHttpsHandler(w http.ResponseWriter, r *http.Request) {

}
