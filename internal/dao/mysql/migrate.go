package mysql

import (
	models2 "patronus/internal/models"
)

func AutoMigrate() {
	db.AutoMigrate(&models2.User{}, &models2.Demand{}, &models2.Task{})
}
