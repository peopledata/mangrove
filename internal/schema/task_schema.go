package schema

import "time"

type TaskItemResp struct {
	Index     int       `json:"index"`
	TaskId    int64     `json:"task_id"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
	Status    int       `json:"status"`
}

type TaskListResp struct {
	Tasks []TaskItemResp `json:"tasks"`
	Total int64          `json:"total"`
}
