package mysql

import (
	"patronus/internal/models"

	"gorm.io/gorm"
)

func InsertContractRecord(record *models.ContractRecord) error {
	return db.Create(record).Error
}

func GetContractRecordByTokenId(demandId, tokenId int64) (*models.ContractRecord, error) {
	var record models.ContractRecord
	if err := db.Where("demand_id=? AND token_id=?", demandId, tokenId).First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = ErrContractRecordNotExist
		}
		return nil, err
	}
	return &record, nil
}

func GetContractRecordsByDemandId(demandId int64, page, pageSize int) []models.ContractRecord {
	var records []models.ContractRecord
	db.Order("sign_time DESC").Offset((page-1)*pageSize).Limit(pageSize).Where("demand_id=?", demandId).Find(&records)
	return records
}

func GetContractRecordsByDemandIdCount(demandId int64) int64 {
	var total int64
	db.Model(&models.ContractRecord{}).Where("demand_id=?", demandId).Count(&total)
	return total
}
