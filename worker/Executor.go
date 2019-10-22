package worker

import (
	"context"
	"github.com/yo1o1o633o/go-crontab/common"
	"os/exec"
	"time"
)

// 任务执行器
type Executor struct {

}

var G_executor *Executor

// 执行一个任务
func (Executor *Executor) ExecutorJob(info *common.JobExecuteInfo) {
	go func() {
		var (
			cmd *exec.Cmd
			outPut []byte
			err error
			result *common.JobExecuteResult
			jobLock *JobLock
		)
		result = &common.JobExecuteResult{
			ExecutorInfo: info,
			Output:       make([]byte, 0),
		}

		// 获取分布式锁
		jobLock = G_WorkMgr.CreateJobLock(info.Job.Name)
		// 开始执行时间
		result.StartTime = time.Now()

		err = jobLock.TryLock()
		defer jobLock.UnLock()
		if err != nil {
			// 抢锁失败
			result.Err = err
			result.EndTime = time.Now()
		} else {
			// 重置开始执行时间
			result.StartTime = time.Now()

			cmd = exec.CommandContext(context.TODO(), "/bin/bash", "-c", info.Job.Command)
			outPut, err = cmd.CombinedOutput()

			// 结束执行时间
			result.EndTime = time.Now()
			result.Output = outPut
			result.Err = err
		}
		G_Scheduer.PushJobResult(result)
	}()
}

// 初始化执行器
func InitExecutor() (err error) {
	G_executor = &Executor{}
	return
}