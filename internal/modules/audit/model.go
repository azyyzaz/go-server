package audit

import "time"

type OperationLog struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	RequestID    string `gorm:"size:64;index"`
	UserID       *uint  `gorm:"index"`
	Username     string `gorm:"size:100;index"`
	Method       string `gorm:"size:10;index"`
	Path         string `gorm:"size:255;index"`
	Module       string `gorm:"size:100"`
	Action       string `gorm:"size:100"`
	IP           string `gorm:"size:64"`
	StatusCode   int    `gorm:"index"`
	Result       string `gorm:"size:20"`
	RequestBody  string `gorm:"type:text"`
	ErrorMessage string `gorm:"size:255"`
	LatencyMS    int64
	CreatedAt    time.Time `gorm:"index"`
}

func (OperationLog) TableName() string {
	return "operation_logs"
}

type LoginLog struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	RequestID  string    `gorm:"size:64;index"`
	UserID     *uint     `gorm:"index"`
	Username   string    `gorm:"size:100;index"`
	IP         string    `gorm:"size:64;index"`
	Region     string    `gorm:"size:100"`
	UserAgent  string    `gorm:"size:512"`
	Success    bool      `gorm:"index"`
	FailReason string    `gorm:"size:255"`
	CreatedAt  time.Time `gorm:"index"`
}

func (LoginLog) TableName() string {
	return "login_logs"
}
