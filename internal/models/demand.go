package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	DemandStatusInit          = 1 + iota // 初识状态：草稿
	DemandStatusPublished                // 已发布
	DemandStatusCompleted                // 已完成
	DemandStatusClosed                   // 已下架
	DemandStatusPublishing               // 发布中
	DemandStatusPublishFailed            // 发布失败
	DemandStatusRunning                  // 运行中
)

type Demand struct {
	ID             uint           `gorm:"primaryKey"`
	DemandId       int64          `gorm:"demand_id;uniqueIndex:idx_demand_id;not null"`
	UserId         uint           `gorm:"user_id;index;not null"` // 创建需求的用户ID
	User           User           `gorm:"foreignKey:UserId"`
	Name           string         `gorm:"type:varchar(255);not null"` // 需求名称
	Brief          string         `gorm:"brief;varchar(512)"`         // 简介
	ValidAt        time.Time      `gorm:"valid_time"`                 // 需求有效期
	Status         int            `gorm:"index"`                      // 状态（草稿：1、已发布：2、已完成：3、已下架：4）
	Category       string         `gorm:"index"`                      // 数据分类（bank）
	Content        string         `gorm:"varchar(512)"`               // 数据内容
	NeedUsers      int            `gorm:"need_users"`                 // 所需用户数
	UseTimes       int            `gorm:"use_times;default:1"`        // 数据使用次数
	ExistingUsers  int            `gorm:"existing_users;default:0"`   // 已有用户数
	AvailableTimes int            `gorm:"available_times;default:0"`  // 数据可用次数
	Purpose        string         `gorm:"text;not null"`              // 数据用途
	Algorithm      string         `gorm:"varchar(512);not null"`      // 算法文件镜像地址
	Agreement      string         `gorm:"text;not null"`              // 协议内容
	ContractToken  string         `gorm:"contract_token"`             // 合约名称：比如 PeopleDataBank
	ContractSymbol string         `gorm:"contract_symbol"`            // 合约标识：比如 PDB
	ContractAddr   string         `gorm:"contract_addr;varchar(255)"` // 部署后的合约地址
	ContractTx     string         `gorm:"contract_tx;varchar(255)"`   // 部署合约的交易ID
	CreatedAt      time.Time      `gorm:"created_time"`
	UpdatedAt      time.Time      `gorm:"updated_time"`
	DeletedAt      gorm.DeletedAt `gorm:"deleted_time;index"`
}
