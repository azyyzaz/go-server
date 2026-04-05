package role

import "time"

// Role GORM model — table: roles
type Role struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"size:50;not null"`
	Code      string    `gorm:"uniqueIndex;size:50;not null"`
	Remark    string    `gorm:"size:200"`
	Status    int8      `gorm:"default:1;comment:1=启用 0=禁用"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
