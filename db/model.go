package db

import "time"

type BaseModel struct {
	ID         int64     `gorm:"column:id;primary_key"`
	CreateTime time.Time `gorm:"column:create_time;"`
	CreateBy   string    `gorm:"column:create_by;"`
	ModifyTime time.Time `gorm:"column:modify_time;"`
	ModifyBy   string    `gorm:"column:modify_by;"`
	DeleteTime time.Time `gorm:"column:delete_time;"`
	DeleteBy   string    `gorm:"column:delete_by;"`
	IsDelete   uint8     `gorm:"column:is_delete;"`
}

type SimpleBaseModel struct {
	ID         int64     `gorm:"column:id;primary_key"`
	CreateTime time.Time `gorm:"column:create_time;"`
	CreateBy   string    `gorm:"column:create_by;"`
}
