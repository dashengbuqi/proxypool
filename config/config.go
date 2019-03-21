package config

import (
	"github.com/go-ini/ini"
	"github.com/kataras/golog"
	"os"
	"path"
	"reflect"
	"runtime"
	"strings"
	"time"
)

type Context struct {
	//App settings
	AppVer  string
	AppName string
	AppURL  string
	AppPath string
	AppAddr string
	AppPort string

	//Global setting objects
	DebugMode bool
	IsWindows bool

	//Security settings
	InstallLock bool // true mean installed

	// OAuth2
	SessionExpires time.Duration
}

var (
	globalConfigure *Context
	logger          *golog.Logger
)

func NewContext() *Context {
	if globalConfigure != nil {
		return globalConfigure
	}
	isWindows := runtime.GOOS == "windows"
	appPath := appPath()
	workDir, err := workDir(appPath)
	if err != nil {
		logger.Fatalf("Fail to get work directory: %v", err.Error())
	}
	iniFile := path.Join(workDir, "ini/app.ini")

	conf, err := ini.Load(iniFile)
	if err != nil {
		logger.Fatalf("Fail to parse %s: %v", iniFile, err.Error())
	}

	conf.NameMapper = ini.AllCapsUnderscore

	InstallLock := conf.Section("security").Key("INSTALL_LOCK").MustBool(false)

	server := conf.Section("server")

	AppName := conf.Section("").Key("APP_NAME").MustString("ProxyPool")
	AppVar := conf.Section("").Key("APP_VAR").MustString("v1")
	AppUrl := server.Key("ROOT_URL").MustString("http://ip:3000/")
	if AppUrl[len(AppUrl)-1] != '/' {
		AppUrl += "/"
	}
	AppAddr := server.Key("HTTP_ADDR").MustString("0.0.0.0")
	AppPort := server.Key("HTTP_PORT").MustString("3000")
	SessionExpires := server.Key("SESSION_EXPIRES").MustDuration(time.Hour * 24 * 7)
	return &Context{
		AppVer:         AppVar,
		AppName:        AppName,
		AppURL:         AppUrl,
		AppPath:        appPath,
		AppAddr:        AppAddr,
		AppPort:        AppPort,
		IsWindows:      isWindows,
		InstallLock:    InstallLock,
		SessionExpires: SessionExpires,
	}
}

func Setting(cField string) interface{} {
	context := NewContext()

	t := reflect.TypeOf(context)
	v := reflect.ValueOf(context)

	var fields = make(map[string]interface{})

	for i := 0; i < t.NumField(); i++ {
		fields[t.Field(i).Name] = v.Field(i).Interface()
	}

	if _, ok := fields[cField]; ok {
		return fields[cField]
	}
	return nil
}

func workDir(ap string) (string, error) {
	wd := os.Getenv("ALIGN_WORK_DIR")
	if len(wd) > 0 {
		return wd, nil
	}

	lastPathChr := strings.LastIndex(ap, "/")
	if lastPathChr == -1 {
		return ap, nil
	}
	return ap[:lastPathChr], nil
}

func appPath() string {

	appPath, err := findPath()
	if err != nil {
		logger.Fatalf("Fail find Application Path:%v\n", err)
	}
	return strings.Replace(appPath, "\\", "/", -1)
}
