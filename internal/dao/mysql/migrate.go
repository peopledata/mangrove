package mysql

import (
	"mangrove/internal/models"
)

func AutoMigrate() {
	db.AutoMigrate(
		&models.User{},
		&models.Demand{},
		&models.Task{},
		&models.ContractRecord{})
}
