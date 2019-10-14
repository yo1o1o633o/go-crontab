package main

import (
	"flag"
	"fmt"
	"github.com/yo1o1o633o/go-crontab/common"
	"github.com/yo1o1o633o/go-crontab/master"
	"runtime"
)

var configFile string

func initArgs() {
	flag.StringVar(&configFile, "config", "E:/project/src/github.com/yo1o1o633o/go-crontab/master/main/config.json", "加载config.json配置文件")
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

	if err = common.InitJobMgr(); err != nil {
		goto ERR
	}

	if err = master.InitApiServer(); err != nil {
		goto ERR
	}

	ERR:
		fmt.Println(err)
}
