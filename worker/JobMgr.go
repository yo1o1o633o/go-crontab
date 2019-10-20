package worker

import (
	"context"
	"github.com/yo1o1o633o/go-crontab/common"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"time"
)

var G_WorkMgr *WorkMgr

type WorkMgr struct {
	client *clientv3.Client
	kv clientv3.KV
	lease clientv3.Lease
	watcher clientv3.Watcher
}

func InitMgr() (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv clientv3.KV
		lease clientv3.Lease
		watcher clientv3.Watcher
	)
	config = clientv3.Config{
		Endpoints: WorkConfig.EtcdEndpoints,
		DialTimeout: time.Duration(WorkConfig.EtcdDialTimeout) * time.Millisecond,
	}

	if client, err = clientv3.New(config); err != nil {
		return
	}

	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)
	watcher = clientv3.NewWatcher(client)

	G_WorkMgr = &WorkMgr{
		client:  client,
		kv:      kv,
		lease:   lease,
		watcher: watcher,
	}
	err = G_WorkMgr.watchJobs()
	return
}

func (WorkMgr *WorkMgr) watchJobs() (err error) {
	var (
		getResp *clientv3.GetResponse
		kvPair *mvccpb.KeyValue
		job *common.Job
		watchStartRevision int64
		watchChan clientv3.WatchChan
		watchResp clientv3.WatchResponse
		watchEvent *clientv3.Event
		jobEvent *common.JobEvent
		jobName string
	)
	if getResp, err = G_WorkMgr.kv.Get(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix()); err != nil {
		return
	}

	for _, kvPair = range getResp.Kvs {
		if job, err = common.UnpackJob(kvPair.Value); err == nil {
			jobEvent = common.BuildJobEvent(common.JOB_ENEVT_SAVE, job)
			G_Scheduer.PushJobEvent(jobEvent)
		}
	}

	go func() {
		watchStartRevision = getResp.Header.Revision + 1
		// 监听版本变化
		watchChan = G_WorkMgr.watcher.Watch(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix())
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT:
					if job, err = common.UnpackJob(watchEvent.Kv.Value); err != nil {
						continue
					}
					jobEvent = common.BuildJobEvent(common.JOB_ENEVT_SAVE, job)
				case mvccpb.DELETE:
					jobName = common.ExtractJobName(string(watchEvent.Kv.Key))
					job = &common.Job{
						Name: jobName,
					}
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_DELETE, job)
				}
				// 变化事件推到调度协程
				G_Scheduer.PushJobEvent(jobEvent)
			}
		}
	}()
	// 监听协程
	return
}
