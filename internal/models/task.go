package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type JSON json.RawMessage

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = JSON(result)
	return err
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}

// 算法任务执行结果记录
type Task struct {
	ID       uint  `gorm:"primaryKey"`
	TaskId   int64 `gorm:"task_id;uniqueIndex:idx_task_id;not null"` // 任务ID
	DemandId int64 `gorm:"demand_id;index;not null"`                 // 关联的需求ID
	UserId   uint  `gorm:"user_id;index;not null"`                   // 执行任务的用户ID
	// User        User         `gorm:"foreignKey:UserId"`
	Status      int          `gorm:"index"` // 状态（运行中：1、运行成功：2、运行失败：3）
	Result      JSON         // 执行结果
	CompletedAt sql.NullTime `gorm:"completed_time"` // 任务完成时间
	CreatedAt   time.Time    `gorm:"created_time"`
	UpdatedAt   time.Time    `gorm:"updated_time"`
}
