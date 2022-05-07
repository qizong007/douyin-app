package conf

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

const ConfigPath = "./conf/conf.yaml"

//定义全局配置
var Config = GetConf()

//定义全局配置变量
type Conf struct {
	Redis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	}
	MYSQL struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Addr     string `yaml:"addr"`
		Database string `yaml:"database"`
	}
	Jwt struct {
		TokenExpireDuration int    `yaml:"token_expire_duration"` //小时为单位
		Secret              string `yaml:"secret"`
	}
	Server struct {
		Port string `yaml:"port"`
	}
}

//获取配置
func GetConf() *Conf {
	var c = Conf{}
	yamlFile, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		log.Println(err.Error())
	}

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Println(err.Error())
	}
	return &c
}
