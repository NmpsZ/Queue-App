package entity

import (
	"time"

	"gorm.io/gorm"
)

type Queue struct {
	ID        uint           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	QueueNo   string         `gorm:"type:varchar(50);not null;unique" json:"queue_no"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Phone     string         `gorm:"type:varchar(20);not null" json:"phone"`
	Status    string         `gorm:"type:varchar(20);default:'waiting'" json:"status"`
	QRCode    string         `gorm:"type:text" json:"qr_code"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
