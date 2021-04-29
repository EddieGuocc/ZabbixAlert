package setting

import (
	"github.com/go-ini/ini"
	log "github.com/sirupsen/logrus"
)

var cfg *ini.File

type App struct {
	Debug            string
	LogSavePath      string
	CompanyId        string
	AgentId          string
	CompanySecret    string
	WeChatTokenUrl   string
	WeChatMessageUrl string
	Safe             uint
	ToAllUser        string
	ToUser           string
}

var AppSetting = &App{}

func Setup() {
	var err error
	cfg, err = ini.Load("/etc/ytalert.ini")
	if err != nil {
		log.Errorln("Fail to parse 'ytalert.ini': %v", err)
		return
	}
	mapTo("conf", AppSetting)
}

func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Errorln("Convert conf file encounter an err: %v", err)
		return
	}
}
