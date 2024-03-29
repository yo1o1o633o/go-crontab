package master

import (
	"context"
	"encoding/json"
	"github.com/yo1o1o633o/go-crontab/common"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"log"
	"time"
)

type JobMgr struct {
	client *clientv3.Client
	kv clientv3.KV
	lease clientv3.Lease
}

var G_jobMgr *JobMgr

func InitJobMgr() (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv clientv3.KV
		lease clientv3.Lease
	)
	log.Printf("初始化ETCD配置信息")
	// 初始化配置
	config = clientv3.Config{
		Endpoints:   G_config.EtcdEndPoints,
		DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Millisecond,
	}
	log.Printf("建立ETCD连接")
	// 建立连接
	if client, err = clientv3.New(config); err != nil {
		log.Printf("建立ETCD连接失败")
		return
	}
	log.Printf("建立ETCD连接KV对象")
	// 获取KV和lease
	kv = clientv3.NewKV(client)
	log.Printf("建立ETCD连接Lease对象")
	lease = clientv3.NewLease(client)

	log.Printf("全局保存ETCD对象信息")
	G_jobMgr = &JobMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}
	return
}

// 保存任务
func (JobMgr *JobMgr) SaveJob(job *common.Job) (oldJob *common.Job, err error) {
	var (
		jobKey string
		jobValue []byte
		putResp *clientv3.PutResponse
		oldJobObj common.Job
	)
	// 保存到etcd中的任务的key
	jobKey = common.JOB_SAVE_DIR + job.Name
	log.Printf("将任务数据序列化未JSON结构")
	// 将任务信息job结构体数据序列化成json
	if jobValue, err = json.Marshal(job); err != nil {
		log.Printf("将任务数据序列化未JSON结构, ERR: " + err.Error())
		return
	}
	log.Printf("保存数据至ETCD")
	// 保存到etcd, 同时获取旧值
	if putResp, err = JobMgr.kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV()); err != nil {
		log.Printf("保存数据至ETCD异常. ERR: " + err.Error())
		return
	}
	// 保存成功返回旧值
	if putResp.PrevKv != nil {
		log.Printf("反序列化ETCD旧值")
		if err = json.Unmarshal(putResp.PrevKv.Value, &oldJobObj); err != nil {
			log.Printf("反序列化ETCD旧值. ERR: " + err.Error())
			return
		}
		oldJob = &oldJobObj
	}
	return
}

// 删除任务
func (JobMgr *JobMgr) DeleteJob(name string) (oldJob *common.Job, err error) {
	var (
		jobKey string
		delResp *clientv3.DeleteResponse
		oldJobObj common.Job
	)
	// 拼接key
	jobKey = common.JOB_SAVE_DIR + name

	if delResp, err = JobMgr.kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV()); err != nil {
		return
	}

	if len(delResp.PrevKvs) != 0 {
		if err = json.Unmarshal(delResp.PrevKvs[0].Value, &oldJobObj); err != nil {
			return
		}
		oldJob = &oldJobObj
	}
	return
}

// 任务列表
func (JobMgr *JobMgr) ListJob() (jobList []*common.Job, err error) {
	var (
		jobKey string
		getResp *clientv3.GetResponse
		kvPair *mvccpb.KeyValue
		job *common.Job
	)
	jobKey = common.JOB_SAVE_DIR
	if getResp, err = JobMgr.kv.Get(context.TODO(), jobKey, clientv3.WithPrefix()); err != nil {
		return
	}

	// 初始化数组空间
	jobList = make([]*common.Job, 0)

	for _, kvPair = range getResp.Kvs {
		job = &common.Job{}
		if err = json.Unmarshal(kvPair.Value, job); err != nil {
			continue
		}
		// 要重新给jobList赋值, 因为当数组空间发生变化时, 内存会重新分配
		jobList = append(jobList, job)
	}
	return
}

// 杀死任务
func (JobMgr *JobMgr) killJob(name string) (err error) {
	var (
		killKey string
		lease *clientv3.LeaseGrantResponse
	)
	killKey = "/cron/jobs/" + name

	if lease, err = JobMgr.lease.Grant(context.TODO(), 1); err != nil {
		return
	}
	if _, err = JobMgr.kv.Put(context.TODO(), killKey,"", clientv3.WithLease(lease.ID)); err != nil {
		return
	}
	return
}

