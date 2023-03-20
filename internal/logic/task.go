package logic

import (
	"patronus/internal/dao/mysql"
	"patronus/internal/models"
	"patronus/internal/schema"
	"patronus/pkg/snowflake"
)

const (
	StatusTaskInit = 1 // 初始状态
)

func CreateTask(demandId int64, user *models.User) error {
	task := models.Task{
		TaskId:   snowflake.GenID(),
		DemandId: demandId,
		UserId:   user.ID, // 这里是关联的User的ID（外键），不是雪雪花ID
		Status:   StatusTaskInit,
	}
	return mysql.InsertTask(&task)
	// todo：任务对象创建后开始执行后台算法
}

func ListTasks(demandId int64) *schema.TaskListResp {
	tasks := mysql.GetAllTasksByDemandId(demandId)

	var taskList []schema.TaskItemResp
	for idx := range tasks {
		item := tasks[idx]
		taskList = append(taskList, schema.TaskItemResp{
			Index:     idx + 1,
			// Username:  item.User.Username,
			Status:    item.Status,
			TaskId:    item.TaskId,
			CreatedAt: item.CreatedAt,
		})
	}

	total := mysql.GetAllTasksCountByDemandId(demandId)

	return &schema.TaskListResp{
		Tasks: taskList,
		Total: total,
	}
}
