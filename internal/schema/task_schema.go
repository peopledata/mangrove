package schema

import "time"

type TaskItemResp struct {
	Index     int       `json:"index"`
	TaskId    string    `json:"task_id"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
	Status    int       `json:"status"`
}

type TaskListResp struct {
	Tasks []TaskItemResp `json:"tasks"`
	Total int64          `json:"total"`
}

type CreateTaskReq struct {
	UserId   uint   `json:"user_id"`
	Username string `json:"username"`
}
