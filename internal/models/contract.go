package models

import "time"

type ContractRecord struct {
	ID       uint      `gorm:"primaryKey"`
	DemandId int64     `gorm:"demand_id;index;not null"`
	TokenId  int64     `gorm:"token_id;index"`              // 链上的token索引号
	TokenURI string    `gorm:"type:varchar(256);token_uri"` // 其实就是nft存放在ipfs上的地址
	DidURI   string    `gorm:"type:varchar(256);did_uri"`   // nft中存放的did doc文档在ipfs上的地址
	Did      string    `gorm:"type:varchar(512);did"`       // did doc 中的 did 标识
	SignTime time.Time `gorm:"sign_time"`                   // 用户授权的时间
}
