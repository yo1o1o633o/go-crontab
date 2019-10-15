package common

import "encoding/json"

// 定时任务
type Job struct {
	Name string `json:"name"`
	Command string `json:"command"`
	CronExpr string `json:"cronExpr"`
}

// 定义通用返回格式
type Response struct {
	Errno int `json:"errno"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
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
