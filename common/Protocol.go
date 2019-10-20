package common

import (
	"encoding/json"
	"github.com/gorhill/cronexpr"
	"strings"
	"time"
)

// 定时任务
type Job struct {
	Name string `json:"name"`
	Command string `json:"command"`
	CronExpr string `json:"cronExpr"`
}

// 任务调度计划
type JobSchedulerPlan struct {
	Job *Job					// 任务信息
	Expr *cronexpr.Expression	// 解析好的crontab表达式
	NextTime time.Time			// 下次执行时间
}

// 任务执行状态
type JobExecuteInfo struct {
	Job *Job
	PlanTime time.Time // 理论执行时间
	RealTime time.Time // 实际执行时间
}

// 任务执行结果
type JobExecuteResult struct {
	ExecutorInfo *JobExecuteInfo
	Output []byte
	Err error
	StartTime time.Time
	EndTime time.Time
}

// 定义通用返回格式
type Response struct {
	Errno int `json:"errno"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}

type JobEvent struct {
	EventType int
	Job *Job
}

func BuildResponse (errno int, msg string, data interface{}) (res []byte, err error) {
	// 定义结构体对象
	var (
		response Response
	)
	response.Errno = errno
	response.Msg = msg
	response.Data = data

	res, err = json.Marshal(response)
	return
}

// 反序列化Job 传入字节数组,返回对应Job结构体
func UnpackJob(value []byte) (res *Job, err error) {
	var (
		job *Job
	)
	job = &Job{}
	if err = json.Unmarshal(value, &job); err != nil {
		return
	}
	res = job
	return
}

func BuildJobEvent(eventType int, job *Job) (jobEvent *JobEvent) {
	return &JobEvent{
		EventType: eventType,
		Job: job,
	}
}

func ExtractJobName(jobString string) (jobName string) {
	return strings.TrimPrefix(JOB_SAVE_DIR, jobString)
}

func BuildJobSchedulerPlan(job *Job) (jobSchedulerPlan *JobSchedulerPlan, err error) {
	var (
		expr *cronexpr.Expression
	)
	if expr, err = cronexpr.Parse(job.CronExpr); err != nil {
		return
	}

	jobSchedulerPlan = &JobSchedulerPlan{
		Job:      job,
		Expr:     expr,
		NextTime: expr.Next(time.Now()),
	}
	return
}

func BuildJobExecuteInfo(jobSchedulerPlan *JobSchedulerPlan) (jobExecuteInfo *JobExecuteInfo) {
	jobExecuteInfo = &JobExecuteInfo{
		Job:      jobSchedulerPlan.Job,
		PlanTime: jobSchedulerPlan.NextTime,
		RealTime: time.Now(),
	}
	return
}