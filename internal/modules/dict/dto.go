package dict

import "time"

type ListDictTypesQuery struct {
	Name string `form:"name"`
	Code string `form:"code"`
}

type CreateDictTypeRequest struct {
	Name   string `json:"name" binding:"required,min=1,max=100"`
	Code   string `json:"code" binding:"required,min=1,max=100"`
	Remark string `json:"remark"`
	Status int8   `json:"status" binding:"oneof=0 1"`
}

type UpdateDictTypeRequest struct {
	Name   string `json:"name" binding:"required,min=1,max=100"`
	Code   string `json:"code" binding:"required,min=1,max=100"`
	Remark string `json:"remark"`
	Status int8   `json:"status" binding:"oneof=0 1"`
}

type DictTypeResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Remark    string    `json:"remark"`
	Status    int8      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type ListDictDataQuery struct {
	TypeID   *uint  `form:"type_id"`
	TypeCode string `form:"type_code"`
	Label    string `form:"label"`
	Status   *int8  `form:"status" binding:"omitempty,oneof=0 1"`
}

type CreateDictDataRequest struct {
	TypeID uint   `json:"type_id" binding:"required"`
	Label  string `json:"label" binding:"required,min=1,max=100"`
	Value  string `json:"value" binding:"required,min=1,max=100"`
	Sort   int    `json:"sort"`
	Status int8   `json:"status" binding:"oneof=0 1"`
	Remark string `json:"remark"`
}

type UpdateDictDataRequest struct {
	TypeID uint   `json:"type_id" binding:"required"`
	Label  string `json:"label" binding:"required,min=1,max=100"`
	Value  string `json:"value" binding:"required,min=1,max=100"`
	Sort   int    `json:"sort"`
	Status int8   `json:"status" binding:"oneof=0 1"`
	Remark string `json:"remark"`
}

type DictDataResponse struct {
	ID        uint      `json:"id"`
	TypeID    uint      `json:"type_id"`
	TypeCode  string    `json:"type_code"`
	Label     string    `json:"label"`
	Value     string    `json:"value"`
	Sort      int       `json:"sort"`
	Status    int8      `json:"status"`
	Remark    string    `json:"remark"`
	CreatedAt time.Time `json:"created_at"`
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
