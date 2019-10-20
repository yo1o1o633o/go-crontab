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
		)
		result = &common.JobExecuteResult{
			ExecutorInfo: info,
			Output:       make([]byte, 0),
		}

		// 开始执行时间
		result.StartTime = time.Now()

		cmd = exec.CommandContext(context.TODO(), "/bin/bash", "-c", info.Job.Command)
		outPut, err = cmd.CombinedOutput()

		// 结束执行时间
		result.EndTime = time.Now()
		result.Output = outPut
		result.Err = err

		G_Scheduer.PushJobResult(result)
	}()
}

// 初始化执行器
func InitExecutor() (err error) {
	G_executor = &Executor{}
	return
}