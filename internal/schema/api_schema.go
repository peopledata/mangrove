package schema

import "time"

type RespData struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type RespTaskList struct {
	Code int           `json:"code"`
	Msg  string        `json:"msg"`
	Data *TaskListResp `json:"data"`
}

type AlgoRecordItem struct {
	Index   int     `json:"index"`
	Name    string  `json:"name"`
	Balance float32 `json:"balance"`
}

type RespAlgoRecordList struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data []AlgoRecordItem `json:"data"`
}

type ContractItem struct {
	ID       uint      `json:"id"`
	Did      string    `json:"did"`
	SignTime time.Time `json:"sign_time"`
}

type RespContract struct {
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
	Records  []ContractItem `json:"records"`
	Total    int            `json:"total"`
}

type RespContractList struct {
	Code int           `json:"code"`
	Msg  string        `json:"msg"`
	Data *RespContract `json:"data"`
}
