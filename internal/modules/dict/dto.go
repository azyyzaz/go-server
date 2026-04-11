package dict

import "time"

type ListDictTypesQuery struct {
	Name string `form:"name" example:"用户状态"`
	Code string `form:"code" example:"user_status"`
}

type CreateDictTypeRequest struct {
	Name   string `json:"name" binding:"required,min=1,max=100" example:"用户状态"`
	Code   string `json:"code" binding:"required,min=1,max=100" example:"user_status"`
	Remark string `json:"remark" example:"用户启用禁用状态"`
	Status int8   `json:"status" binding:"oneof=0 1" example:"1"`
}

type UpdateDictTypeRequest struct {
	Name   string `json:"name" binding:"required,min=1,max=100" example:"用户状态"`
	Code   string `json:"code" binding:"required,min=1,max=100" example:"user_status"`
	Remark string `json:"remark" example:"用户启用禁用状态"`
	Status int8   `json:"status" binding:"oneof=0 1" example:"1"`
}

type DictTypeResponse struct {
	ID        uint      `json:"id" example:"1"`
	Name      string    `json:"name" example:"用户状态"`
	Code      string    `json:"code" example:"user_status"`
	Remark    string    `json:"remark" example:"用户启用禁用状态"`
	Status    int8      `json:"status" example:"1"`
	CreatedAt time.Time `json:"created_at" swaggertype:"string" example:"2026-04-11T12:00:00Z"`
}

type ListDictDataQuery struct {
	TypeID   *uint  `form:"type_id" swaggertype:"integer" example:"1"`
	TypeCode string `form:"type_code" example:"user_status"`
	Label    string `form:"label" example:"启用"`
	Status   *int8  `form:"status" binding:"omitempty,oneof=0 1" swaggertype:"integer" example:"1"`
}

type CreateDictDataRequest struct {
	TypeID uint   `json:"type_id" binding:"required" example:"1"`
	Label  string `json:"label" binding:"required,min=1,max=100" example:"启用"`
	Value  string `json:"value" binding:"required,min=1,max=100" example:"1"`
	Sort   int    `json:"sort" example:"1"`
	Status int8   `json:"status" binding:"oneof=0 1" example:"1"`
	Remark string `json:"remark" example:"启用状态"`
}

type UpdateDictDataRequest struct {
	TypeID uint   `json:"type_id" binding:"required" example:"1"`
	Label  string `json:"label" binding:"required,min=1,max=100" example:"启用"`
	Value  string `json:"value" binding:"required,min=1,max=100" example:"1"`
	Sort   int    `json:"sort" example:"1"`
	Status int8   `json:"status" binding:"oneof=0 1" example:"1"`
	Remark string `json:"remark" example:"启用状态"`
}

type DictDataResponse struct {
	ID        uint      `json:"id" example:"1"`
	TypeID    uint      `json:"type_id" example:"1"`
	TypeCode  string    `json:"type_code" example:"user_status"`
	Label     string    `json:"label" example:"启用"`
	Value     string    `json:"value" example:"1"`
	Sort      int       `json:"sort" example:"1"`
	Status    int8      `json:"status" example:"1"`
	Remark    string    `json:"remark" example:"启用状态"`
	CreatedAt time.Time `json:"created_at" swaggertype:"string" example:"2026-04-11T12:00:00Z"`
}

func toTypeResponse(item DictType) DictTypeResponse {
	return DictTypeResponse{
		ID:        item.ID,
		Name:      item.Name,
		Code:      item.Code,
		Remark:    item.Remark,
		Status:    item.Status,
		CreatedAt: item.CreatedAt,
	}
}
