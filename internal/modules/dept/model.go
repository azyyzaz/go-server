package dept

import "time"

type Dept struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	ParentID  *uint  `gorm:"index"`
	Name      string `gorm:"size:100;not null"`
	Leader    string `gorm:"size:100"`
	Phone     string `gorm:"size:20"`
	Email     string `gorm:"size:191"`
	Sort      int
	Status    int8 `gorm:"default:1"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Dept) TableName() string {
	return "depts"
}
