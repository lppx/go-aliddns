package config

import (
	"github.com/spf13/viper"
	"log"
)

var G Config

type  Config struct {
	Ali AliConfig
}


type AliConfig struct {
	AccessId	string
	AccessKey	string
	MainDomain	string
	SubDomain	string
	TimeStep	int64
}

func  InitConfig()  {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("读取配置文件失败：%v", err)
	}
	viper.Unmarshal(&G)
}