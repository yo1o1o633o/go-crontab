package main

import (
	"flag"
	"github.com/yo1o1o633o/go-crontab/worker"
)
var configFile string

func initArgs() {
	flag.StringVar(&configFile, "config", "E:/project/src/github.com/yo1o1o633o/go-crontab/worker/main/config.json", "加载config.json配置文件")
}

func main() {
	var (
		err error
	)
	if err = worker.InitConfig(configFile); err != nil {
		goto ERR
	}

	if err = worker.InitExecutor(); err != nil {
		goto ERR
	}

	if err = worker.InitScheduler(); err != nil {
		goto ERR
	}

	if err = worker.InitMgr(); err != nil {
		goto ERR
	}
	ERR:
}
