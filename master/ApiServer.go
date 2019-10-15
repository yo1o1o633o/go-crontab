package master

import (
	"encoding/json"
	"github.com/yo1o1o633o/go-crontab/common"
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
	mux.HandleFunc("/jobs/delete", handleJobDelete)

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

// 保存接口
func handleJobSave(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		postJob string
		job common.Job
		oldJob *common.Job
		res []byte
	)
	// 解析入参
	if err = r.ParseForm(); err != nil {
		goto ERR
	}

	// 获取入参job字段
	postJob = r.PostForm.Get("job")

	// 反序列化入参,入参是json格式, 反序列到job结构体中保存
	if err = json.Unmarshal([]byte(postJob), &job); err != nil {
		goto ERR
	}

	// 将序列化后的结构体数据保存到etcd
	if oldJob, err = G_jobMgr.SaveJob(&job); err != nil {
		goto ERR
	}

	// 返回成功
	if res, err = common.BuildResponse(200, "success", oldJob); err == nil {
		w.Write(res)
	}

	ERR:
		// 返回异常
		if res, err = common.BuildResponse(10000, err.Error(), nil); err == nil {
			w.Write(res)
		}
}

// 删除接口
func handleJobDelete(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		name string
		oldJob *common.Job
		res []byte
	)
	if err = r.ParseForm(); err != nil {
		goto ERR
	}
	// 获取要删除的任务名
	name = r.PostForm.Get("name")
	// 在etcd中删除该key
	if oldJob, err = G_jobMgr.DeleteJob(name); err != nil {
		goto ERR
	}
	// 返回成功
	if res, err = common.BuildResponse(200, "success", oldJob); err == nil {
		w.Write(res)
	}
	ERR:
		if res, err = common.BuildResponse(10000, err.Error(), nil); err != nil {
			w.Write(res)
		}
}
