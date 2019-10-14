package master

import (
	"encoding/json"
	"io/ioutil"
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

	// 读取传入配置文件
	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	// 根据结构体反序列化json
	if err = json.Unmarshal(content, &conf); err != nil {
		return
	}
	// 传入单例
	G_config = &conf
	return
}
