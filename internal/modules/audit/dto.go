package audit

import "time"

type OperationLogEntry struct {
	RequestID    string
	UserID       *uint
	Username     string
	Method       string
	Path         string
	IP           string
	StatusCode   int
	RequestBody  string
	ErrorMessage string
	LatencyMS    int64
}

type LoginLogEntry struct {
	RequestID  string
	UserID     *uint
	Username   string
	IP         string
	Region     string
	UserAgent  string
	Success    bool
	FailReason string
}

type ListOperationLogsQuery struct {
	Page      int    `form:"page" binding:"omitempty,min=1" default:"1" example:"1"`
	PageSize  int    `form:"page_size" binding:"omitempty,min=1,max=100" default:"10" example:"10"`
	Username  string `form:"username" example:"admin"`
	Method    string `form:"method" example:"POST"`
	Path      string `form:"path" example:"/api/v1/auth/login"`
	Status    *int   `form:"status" binding:"omitempty,min=100,max=599" swaggertype:"integer" example:"200"`
	StartTime string `form:"start_time" example:"2026-04-01T00:00:00Z"`
	EndTime   string `form:"end_time" example:"2026-04-11T23:59:59Z"`
}

type ListLoginLogsQuery struct {
	Page      int    `form:"page" binding:"omitempty,min=1" default:"1" example:"1"`
	PageSize  int    `form:"page_size" binding:"omitempty,min=1,max=100" default:"10" example:"10"`
	Username  string `form:"username" example:"admin"`
	IP        string `form:"ip" example:"127.0.0.1"`
	Success   *bool  `form:"success" swaggertype:"boolean" example:"true"`
	StartTime string `form:"start_time" example:"2026-04-01T00:00:00Z"`
	EndTime   string `form:"end_time" example:"2026-04-11T23:59:59Z"`
}

type OperationLogResponse struct {
	ID           uint      `json:"id" example:"1"`
	RequestID    string    `json:"request_id" example:"20260411122824-2d95416e434d"`
	UserID       *uint     `json:"user_id,omitempty" swaggertype:"integer" example:"1"`
	Username     string    `json:"username" example:"admin"`
	Method       string    `json:"method" example:"POST"`
	Path         string    `json:"path" example:"/api/v1/auth/login"`
	Module       string    `json:"module" example:"认证授权"`
	Action       string    `json:"action" example:"用户登录"`
	IP           string    `json:"ip" example:"127.0.0.1"`
	StatusCode   int       `json:"status_code" example:"200"`
	Result       string    `json:"result" example:"success"`
	RequestBody  string    `json:"request_body" example:"{\"username\":\"admin\"}"`
	ErrorMessage string    `json:"error_message" example:""`
	LatencyMS    int64     `json:"latency_ms" example:"12"`
	CreatedAt    time.Time `json:"created_at" swaggertype:"string" example:"2026-04-11T12:00:00Z"`
}

type LoginLogResponse struct {
	ID         uint      `json:"id" example:"1"`
	RequestID  string    `json:"request_id" example:"20260411122824-2d95416e434d"`
	UserID     *uint     `json:"user_id,omitempty" swaggertype:"integer" example:"1"`
	Username   string    `json:"username" example:"admin"`
	IP         string    `json:"ip" example:"127.0.0.1"`
	Region     string    `json:"region" example:"上海"`
	UserAgent  string    `json:"user_agent" example:"Mozilla/5.0"`
	Success    bool      `json:"success" example:"true"`
	FailReason string    `json:"fail_reason" example:""`
	CreatedAt  time.Time `json:"created_at" swaggertype:"string" example:"2026-04-11T12:00:00Z"`
}

type OperationLogPageResult struct {
	Items    []OperationLogResponse `json:"items"`
	Total    int64                  `json:"total" example:"1"`
	Page     int                    `json:"page" example:"1"`
	PageSize int                    `json:"page_size" example:"10"`
}

type LoginLogPageResult struct {
	Items    []LoginLogResponse `json:"items"`
	Total    int64              `json:"total" example:"1"`
	Page     int                `json:"page" example:"1"`
	PageSize int                `json:"page_size" example:"10"`
}
