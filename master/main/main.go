package main

import (
	"flag"
	"fmt"
	"github.com/yo1o1o633o/go-crontab/master"
	"runtime"
	"time"
)

var configFile string

func initArgs() {
	flag.StringVar(&configFile, "config", "F:/goProject/src/github.com/yo1o1o633o/go-crontab/master/main/config.json", "加载config.json配置文件")
}

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		err error
	)
	initArgs()
	initEnv()

	if err = master.InitConfig(configFile); err != nil {
		goto ERR
	}

	if err = master.InitJobMgr(); err != nil {
		goto ERR
	}

	if err = master.InitApiServer(); err != nil {
		goto ERR
	}

	for {
		time.Sleep(1)
	}
	ERR:
		fmt.Println(err)
}
