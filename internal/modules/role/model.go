package role

import "time"

type Role struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"size:50;not null"`
	Code      string `gorm:"uniqueIndex;size:50;not null"`
	Remark    string `gorm:"size:200"`
	Status    int8   `gorm:"default:1;comment:1=启用 0=禁用"`
	Menus     []Menu `gorm:"many2many:role_menus;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

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
