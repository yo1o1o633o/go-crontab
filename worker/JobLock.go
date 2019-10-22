package worker

import (
	"context"
	"github.com/yo1o1o633o/go-crontab/common"
	"go.etcd.io/etcd/clientv3"
)

type JobLock struct {
	kv clientv3.KV
	lease clientv3.Lease
	jobName string
	cancelFunc context.CancelFunc
	leaseId clientv3.LeaseID
	isLock bool
}

// 初始化锁
func InitJobLock(jobName string, kv clientv3.KV, lease clientv3.Lease) (jobLock *JobLock) {
	jobLock = &JobLock{
		kv:      kv,
		lease:   lease,
		jobName: jobName,
	}
	return
}

// 尝试上锁
func (jobLock *JobLock) TryLock() (err error) {
	var (
		leaseGrandResp *clientv3.LeaseGrantResponse
		cancelCtx context.Context
		cancelFunc context.CancelFunc
		leaseId clientv3.LeaseID
		keepRespChan <- chan *clientv3.LeaseKeepAliveResponse
		txn clientv3.Txn
		lockKey string
		txnResp *clientv3.TxnResponse
	)
	// 创建5秒租约
	if leaseGrandResp, err = jobLock.lease.Grant(context.TODO(), 5); err != nil {
		return
	}

	// 用于取消自动续租
	cancelCtx, cancelFunc = context.WithCancel(context.TODO())

	// 租约ID
	leaseId = leaseGrandResp.ID

	// 自动续租
	if keepRespChan, err = jobLock.lease.KeepAlive(cancelCtx, leaseId); err != nil {
		goto FAIL
	}

	// 处理续租应答协程
	go func() {
		var (
			keepResp *clientv3.LeaseKeepAliveResponse
		)
		for {
			select {
			case keepResp = <- keepRespChan:
				if keepResp == nil {
					goto END
				}
			}
		}
		END:
	}()

	// 创建事务
	txn = jobLock.kv.Txn(context.TODO())
	// 定义锁路径
	lockKey = common.JOB_LOCK_DIR + jobLock.jobName

	// 定义事务
	txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey, "", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet(lockKey))

	// 提交事务
	if txnResp, err = txn.Commit(); err != nil {
		goto FAIL
	}

	if !txnResp.Succeeded {
		// 锁被占用,取消租约
		err = common.ERR_LOCK_ALREADY_REQUIRED
		goto FAIL
	}

	// 成功上锁
	jobLock.leaseId = leaseId
	jobLock.cancelFunc = cancelFunc
	jobLock.isLock = true

	return

	FAIL:
		// 取消续租
		cancelFunc()
		_, _ = jobLock.lease.Revoke(context.TODO(), leaseId)
		return
}

// 释放锁
func (jobLock *JobLock) UnLock() {
	if jobLock.isLock {
		// 干掉租约协程
		jobLock.cancelFunc()
		// 干掉租约
		_, _ = jobLock.lease.Revoke(context.TODO(), jobLock.leaseId)
	}
}
