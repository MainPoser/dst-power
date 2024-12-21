package model

import (
	"database/sql"
)

type Model struct {
	ID          uint         `gorm:"column:id;primaryKey;type:int(11)"`
	CreatedTime sql.NullTime `gorm:"column:created_time;type:timestamp"`
	UpdatedTime sql.NullTime `gorm:"column:updated_time;type:timestamp"`
}
