package worker

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var WorkConfig *Config

type Config struct {
	EtcdEndpoints []string `json:"etcdEndpoints"`
	EtcdDialTimeout int `json:"etcdDialTimeout"`
}

func InitConfig(filename string) (err error) {
	var (
		content []byte
		conf Config
	)
	log.Printf("读取传入配置文件")
	if content, err = ioutil.ReadFile(filename); err != nil {
		log.Printf("配置文件读取失败. ERR: " + err.Error())
		return
	}

	log.Printf("序列化配置文件")
	if err = json.Unmarshal(content, &conf); err != nil {
		log.Printf("序列化配置文件失败. ERR: " + err.Error())
		return
	}

	WorkConfig = &conf
	return
}