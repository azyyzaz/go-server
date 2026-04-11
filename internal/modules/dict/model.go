package dict

import "time"

type DictType struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"size:100;not null"`
	Code      string `gorm:"uniqueIndex;size:100;not null"`
	Remark    string `gorm:"size:255"`
	Status    int8   `gorm:"default:1"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (DictType) TableName() string {
	return "dict_types"
}

type DictData struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	TypeID    uint   `gorm:"index;not null"`
	Label     string `gorm:"size:100;not null"`
	Value     string `gorm:"size:100;not null"`
	Sort      int
	Status    int8   `gorm:"default:1"`
	Remark    string `gorm:"size:255"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (DictData) TableName() string {
	return "dict_data"
}
