package worker

import (
	"encoding/json"
	"io/ioutil"
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
	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	if err = json.Unmarshal(content, &conf); err != nil {
		return
	}

	WorkConfig = &conf
	return
}