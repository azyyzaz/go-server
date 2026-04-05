package user

import (
	"time"

	"go-server/internal/modules/role"
)

// User GORM model — table: users
type User struct {
	ID        uint        `gorm:"primaryKey;autoIncrement"`
	Username  string      `gorm:"uniqueIndex;size:50;not null"`
	Password  string      `gorm:"size:255;not null"`
	Name      string      `gorm:"size:100;not null"`
	Email     string      `gorm:"uniqueIndex;size:191;not null"`
	Phone     string      `gorm:"size:20"`
	Status    int8        `gorm:"default:1;comment:1=active 0=disabled"`
	Roles     []role.Role `gorm:"many2many:user_roles;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
