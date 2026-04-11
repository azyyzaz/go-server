package menu

import "time"

type Menu struct {
	ID         uint   `gorm:"primaryKey;autoIncrement"`
	ParentID   *uint  `gorm:"index"`
	Name       string `gorm:"size:100;not null"`
	Type       string `gorm:"size:20;not null"`
	Path       string `gorm:"size:200"`
	Component  string `gorm:"size:255"`
	Permission string `gorm:"size:100"`
	Sort       int
	Visible    int8 `gorm:"default:1"`
	Status     int8 `gorm:"default:1"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (Menu) TableName() string {
	return "menus"
}
