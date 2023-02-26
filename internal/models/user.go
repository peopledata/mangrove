package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	UserId   int64  `gorm:"uniqueIndex:idx_user_id;not null"`
	Username string `gorm:"type:varchar(255);unique;not null"`
	Password string `gorm:"type:varchar(255);not null"`
	Created  int64  `gorm:"autoCreateTime"` // 使用时间戳秒数填充创建时间
	Updated  int64  `gorm:"autoUpdateTime"` // 使用时间戳秒数填充更新时间
}

// TableName ... 指定 User 结构体对应的数据表为 user
func (u User) TableName() string {
	return "user"
}
