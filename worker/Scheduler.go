package worker

import (
	"github.com/yo1o1o633o/go-crontab/common"
	"time"
)

type Scheduler struct {
	jobEventChan chan *common.JobEvent
	jobPlanTable map[string]*common.JobSchedulerPlan
	jobExecutingTable map[string]*common.JobExecuteInfo
	jobResultChan chan *common.JobExecuteResult
}

var G_Scheduer *Scheduler

func InitScheduler() (err error) {
	G_Scheduer = &Scheduler{
		jobEventChan: make(chan *common.JobEvent, 1000),
		jobPlanTable: make(map[string]*common.JobSchedulerPlan),
		jobExecutingTable: make(map[string]*common.JobExecuteInfo),
		jobResultChan: make(chan *common.JobExecuteResult, 1000),
	}
	go G_Scheduer.schedulerLoop()
	return
}

func (Scheduler *Scheduler) schedulerLoop() {
	var (
		jobEvent *common.JobEvent
		schedulerAfter time.Duration
		schedulerTimer *time.Timer
		jobResult *common.JobExecuteResult
	)

	// 初始化一次
	schedulerAfter = Scheduler.TryScheduler()

	// 调度延迟器
	schedulerTimer = time.NewTimer(schedulerAfter)

	for {
		select {
		case jobEvent = <- Scheduler.jobEventChan:
			// 根据事件类型处理
			Scheduler.handleJobEvent(jobEvent)
		case <- schedulerTimer.C:
		case jobResult = <- Scheduler.jobResultChan:
			Scheduler.handleJobResult(jobResult)
		}
		// 调度一次任务
		schedulerAfter = Scheduler.TryScheduler()
		schedulerTimer.Reset(schedulerAfter)
	}
}

func (Scheduler *Scheduler) handleJobEvent(jobEvent *common.JobEvent) {
	var (
		jobSchedulerPlan *common.JobSchedulerPlan
		err error
		jobExisted bool
	)
	switch jobEvent.EventType {
	case common.JOB_ENEVT_SAVE:
		if jobSchedulerPlan, err = common.BuildJobSchedulerPlan(jobEvent.Job); err != nil {
			return
		}
		Scheduler.jobPlanTable[jobEvent.Job.Name] = jobSchedulerPlan
	case common.JOB_EVENT_DELETE:
		if jobSchedulerPlan, jobExisted = Scheduler.jobPlanTable[jobEvent.Job.Name]; jobExisted {
			delete(Scheduler.jobPlanTable, jobEvent.Job.Name)
		}
	}
}

// 处理任务结果
func (Scheduler *Scheduler) handleJobResult(result *common.JobExecuteResult) {
	// 删除任务
	delete(Scheduler.jobExecutingTable, result.ExecutorInfo.Job.Name)
}

// 尝试执行任务
func (Scheduler *Scheduler) TryStartJob(jobPlan *common.JobSchedulerPlan) {
	var (
		jobExecuteInfo *common.JobExecuteInfo
		jobExecuting bool
	)
	// 执行的任务是否已经在任务执行表中
	if jobExecuteInfo, jobExecuting = Scheduler.jobExecutingTable[jobPlan.Job.Name]; jobExecuting {
		return
	}
	jobExecuteInfo = common.BuildJobExecuteInfo(jobPlan)

	// 保存执行状态
	Scheduler.jobExecutingTable[jobPlan.Job.Name] = jobExecuteInfo

	// 执行任务
	G_executor.ExecutorJob(jobExecuteInfo)
}

func (Scheduler *Scheduler) TryScheduler() (schedulerAfter time.Duration) {
	var (
		jobPlan *common.JobSchedulerPlan
		now time.Time
		nearTime *time.Time
	)

	if len(Scheduler.jobPlanTable) == 0 {
		schedulerAfter = 1 * time.Second
	}

	now = time.Now()
	for _, jobPlan = range Scheduler.jobPlanTable {
		if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now) {
			// 尝试执行任务
			Scheduler.TryStartJob(jobPlan)
			jobPlan.NextTime = jobPlan.Expr.Next(now)  // 更新下次执行时间
		}
		// 最近一个要过期任务为空或下次循环发现有更近要过期的任务
		if nearTime == nil || jobPlan.NextTime.Before(*nearTime){
			nearTime = &jobPlan.NextTime
		}
	}
	// 下次调度时间
	//schedulerAfter = (nearTime).Sub(now)
	return
}

func (Scheduler *Scheduler) PushJobEvent(jobEvent *common.JobEvent) {
	Scheduler.jobEventChan <- jobEvent
}

func (Scheduler *Scheduler) PushJobResult(jobResult *common.JobExecuteResult) {
	Scheduler.jobResultChan <- jobResult
}