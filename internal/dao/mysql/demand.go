package mysql

import (
	"patronus/internal/models"

	"gorm.io/gorm"
)

func InsertDemand(demand *models.Demand) error {
	return db.Create(demand).Error
}

func GetAllDemands() []models.Demand {
	var demands []models.Demand
	db.Order("created_at DESC").Find(&demands)
	return demands
}

func GetAllDemandsCount() int64 {
	var total int64
	db.Model(&models.Demand{}).Count(&total)
	return total
}

func GetDemandDetail(demandId int64) (*models.Demand, error) {
	var demand models.Demand
	if err := db.Where("demand_id=?", demandId).First(&demand).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = ErrDemandNotExist
		}
		return nil, err
	}
	return &demand, nil
}

func UpdateDemand(demand *models.Demand) error {
	return db.Model(&models.Demand{}).Where("demand_id=?", demand.DemandId).Updates(*demand).Error
}

func UpdateDemandStatus(demandId int64, status int) error {
	return db.Model(&models.Demand{}).Where("demand_id=?", demandId).Updates(models.Demand{Status: status}).Error
}

func UpdateDemandContract(demandId int64, address, tx string) error {
	return db.Model(&models.Demand{}).Where("demand_id=?", demandId).Updates(models.Demand{
		ContractAddr: address,
		ContractTx:   tx,
		Status:       models.DemandStatusPublishing,
	}).Error
}

func CheckDemandInitStatus(demandId int64) error {
	var ds models.Demand
	if err := db.Select("status").Where("demand_id=?", demandId).Find(&ds).Error; err != nil {
		return err
	}
	// 如果不是草稿状态，则不能发布
	if ds.Status != models.DemandStatusInit {
		return ErrDemandStatusNotInit
	}
	return nil
}

func GetAllDemandsByStatusAndCategory(status int, category string) []models.Demand {
	var demands []models.Demand
	db.Order("created_at DESC").Where("status=? AND category=? AND contract_addr != ''", status, category).Find(&demands)
	return demands
}

func GetAllDemandsByStatus(status int) []models.Demand {
	var demands []models.Demand
	db.Order("created_at DESC").Where("status=?", status).Find(&demands)
	return demands
}

func GetAllDemandsByStatusCount(status int) int64 {
	var total int64
	db.Model(&models.Demand{}).Where("status=?", status).Count(&total)
	return total
}
