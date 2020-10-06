package db

type BaseModel struct {
	ID         int64  `gorm:"column:id;primary_key"`
	CreateTime int64  `gorm:"column:create_time;"`
	CreateBy   string `gorm:"column:create_by;"`
	ModifyTime int64  `gorm:"column:modify_time;"`
	ModifyBy   string `gorm:"column:modify_by;"`
	DeleteTime int64  `gorm:"column:delete_time;"`
	DeleteBy   string `gorm:"column:delete_by;"`
	IsDelete   uint8  `gorm:"column:is_delete;"`
}

type SimpleBaseModel struct {
	ID         int64  `gorm:"column:id;primary_key"`
	CreateTime int64  `gorm:"column:create_time;"`
	CreateBy   string `gorm:"column:create_by;"`
}
