package master

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	ApiPort int `json:"apiPort"`
	ApiReadTimeout int `json:"apiReadTimeout"`
	ApiWriteTimeout int `json:"apiWriteTimeout"`
	EtcdEndPoints []string `json:"etcdEndPoints"`
	EtcdDialTimeout int `json:"etcdDialTimeout"`
}

// 单例保存配置信息
var G_config *Config

func InitConfig(filename string) (err error) {
	var (
		content []byte
		conf Config
	)
	log.Printf("读取传入配置文件")
	// 读取传入配置文件
	if content, err = ioutil.ReadFile(filename); err != nil {
		log.Printf("配置文件读取失败")
		return
	}

	log.Printf("序列化配置文件")
	// 根据结构体反序列化json
	if err = json.Unmarshal(content, &conf); err != nil {
		log.Printf("序列化配置文件失败")
		return
	}
	// 传入单例
	G_config = &conf
	return
}
