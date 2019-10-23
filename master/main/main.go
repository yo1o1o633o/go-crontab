package main

import (
	"flag"
	"fmt"
	"github.com/yo1o1o633o/go-crontab/master"
	"log"
	"runtime"
	"time"
)

var configFile string

func initLog() {
	log.SetFlags(log.Ldate|log.Llongfile)
}

func initArgs() {
	flag.StringVar(&configFile, "config", "F:/goProject/src/github.com/yo1o1o633o/go-crontab/master/main/config.json", "加载config.json配置文件")
}

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	initLog()
	var (
		err error
	)
	initArgs()
	log.Printf("加载config.json配置文件")
	initEnv()
	log.Printf("初始化进程")
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
