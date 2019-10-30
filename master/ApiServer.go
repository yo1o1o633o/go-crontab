package master

import (
	"encoding/json"
	"fmt"
	"github.com/yo1o1o633o/go-crontab/common"
	"log"
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
		staticDir http.Dir
		staticHandler http.Handler
	)
	log.Printf("初始化路由组")
	// 初始化路由
	mux = http.NewServeMux()
	mux.HandleFunc("/jobs/save", handleJobSave)
	mux.HandleFunc("/jobs/delete", handleJobDelete)
	mux.HandleFunc("/jobs/list", handleJobList)
	mux.HandleFunc("/jobs/kill", handleJobKill)

	staticDir = http.Dir("E:/project/src/github.com/yo1o1o633o/go-crontab/master/main/webroot")
	staticHandler = http.FileServer(staticDir)
	mux.Handle("/", http.StripPrefix("/", staticHandler))

	log.Printf("初始化TCP服务监听")
	// 启动TCP监听
	if listener, err = net.Listen("tcp", ":" + strconv.Itoa(G_config.ApiPort)); err != nil {
		log.Printf("初始化TCP服务监听.异常信息: " + err.Error())
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

	log.Printf("启动TCP服务监听")
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
	log.Printf("解析POST请求入参数据")
	if err = r.ParseForm(); err != nil {
		log.Printf("解析POST请求入参数据失败, ERR: " + err.Error())
		return
	}
	fmt.Println(r.PostForm.Get("job"))
	// 获取入参job字段
	postJob = r.PostForm.Get("job")

	fmt.Println(postJob)
	log.Printf("反序列化入参信息")
	// 反序列化入参,入参是json格式, 反序列到job结构体中保存
	if err = json.Unmarshal([]byte(postJob), &job); err != nil {
		log.Printf("反序列化入参信息, ERR: " + err.Error())
		return
	}

	// 将序列化后的结构体数据保存到etcd
	if oldJob, err = G_jobMgr.SaveJob(&job); err != nil {
		goto ERR
	}
	// 返回成功
	if res, err = common.BuildResponse(200, "success", oldJob); err == nil {
		w.Write(res)
	}
	return
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
	return
	ERR:
		if res, err = common.BuildResponse(10000, err.Error(), nil); err == nil {
			w.Write(res)
		}
}

func handleJobList(w http.ResponseWriter, r *http.Request) {
	var (
		jobList []*common.Job
		err error
		res []byte
	)
	if jobList, err = G_jobMgr.ListJob(); err != nil {
		fmt.Println(err)
		goto ERR
	}
	if res, err = common.BuildResponse(200, "success", jobList); err == nil {
		w.Write(res)
	}
	return
	ERR:
		if res, err = common.BuildResponse(10000, err.Error(), ""); err == nil {
			w.Write(res)
		}
}

// 杀死任务
func handleJobKill(w http.ResponseWriter, r *http.Request) {
	var (
		name string
		err error
		res []byte
	)
	if err = r.ParseForm(); err != nil {
		goto ERR
	}
	name = r.PostForm.Get("name")
	if err = G_jobMgr.killJob(name); err != nil {
		goto ERR
	}
	if res, err = common.BuildResponse(200, "success", nil); err == nil {
		w.Write(res)
	}
	return
	ERR:
		if res, err = common.BuildResponse(10000, err.Error(), nil); err != nil {
			w.Write(res)
		}
}
