package main

import (
	"flag"
	"fmt"
	"github.com/yo1o1o633o/go-crontab/worker"
	"time"
)
var configFile string

func initArgs() {
	flag.StringVar(&configFile, "config", "./src/github.com/yo1o1o633o/go-crontab/worker/main/worker.json", "加载config.json配置文件")
}

func main() {
	var (
		err error
	)
	initArgs()

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
	for {
		time.Sleep(1)
	}
	ERR:
		fmt.Println(err)
}
