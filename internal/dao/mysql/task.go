package mysql

import (
	"patronus/internal/models"
)

func InsertTask(task *models.Task) error {
	return db.Create(task).Error
}

func GetAllTasksByDemandId(id int64) []models.Task {
	var tasks []models.Task
	db.Where("demand_id=?", id).Preload("User").Find(&tasks)
	return tasks
}

func GetAllTasksCountByDemandId(id int64) int64 {
	var total int64
	db.Model(&models.Task{}).Where("demand_id=?", id).Count(&total)
	return total
}
