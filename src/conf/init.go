package conf

import (
	"github.com/robfig/config"
	"log"
)

var Settings *config.Config

func init() {
	Settings = GetSettings()
}

func GetSettings() *config.Config {
	if Settings == nil {
		settings, err := config.ReadDefault("config.cfg")
		if err != nil {
			log.Fatalf("package conf %s \n %#v", "配置文件初始化出错", err)
		}
		Settings = settings
		return settings
	} else {
		return Settings
	}
}
