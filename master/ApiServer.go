package master

import (
	"net"
	"net/http"
	"strconv"
	"time"
)

type ApiServer struct {
	httpServer *http.Server
}

// 单例对象
var G_apiServer *ApiServer

func InitApiServer() (err error){
	var (
		mux *http.ServeMux
		listener net.Listener
		httpServer *http.Server
	)
	// 初始化路由
	mux = http.NewServeMux()
	mux.HandleFunc("/jobs/save", handleJobSave)

	// 启动TCP监听
	if listener, err = net.Listen("tcp", ":" + strconv.Itoa(G_config.ApiPort)); err != nil {
		return
	}

	// 创建HTTP服务
	httpServer = &http.Server{
		Handler:           mux,
		ReadTimeout:       time.Duration(G_config.ApiReadTimeout) * time.Millisecond,
		WriteTimeout:      time.Duration(G_config.ApiWriteTimeout) * time.Millisecond,
	}

	// 赋值单例
	G_apiServer = &ApiServer{httpServer:httpServer}

	// 启动服务端
	go httpServer.Serve(listener)
	return
}

func handleJobSave(w http.ResponseWriter, r *http.Request) {

}
