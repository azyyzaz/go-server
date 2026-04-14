package file

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	ID           uint           `gorm:"primaryKey"`
	UploaderID   *uint          `gorm:"column:uploader_id"`
	Storage      string         `gorm:"size:20;not null"`
	Category     string         `gorm:"size:50;not null"`
	OriginalName string         `gorm:"size:255;not null"`
	ObjectName   string         `gorm:"size:255;not null"`
	Ext          string         `gorm:"size:20;not null"`
	MIMEType     string         `gorm:"column:mime_type;size:100;not null"`
	Size         int64          `gorm:"not null"`
	Bucket       string         `gorm:"size:100;not null"`
	Path         string         `gorm:"size:500;not null"`
	URL          string         `gorm:"size:500;not null"`
	CreatedAt    time.Time      `gorm:"not null"`
	UpdatedAt    time.Time      `gorm:"not null"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (File) TableName() string {
	return "files"
}
