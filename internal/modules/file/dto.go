package file

import "time"

type UploadRequest struct {
	Category string
}

type UploadAvatarRequest struct{}

type ListFilesQuery struct {
	Page       int    `form:"page" binding:"omitempty,min=1" default:"1" example:"1"`
	PageSize   int    `form:"page_size" binding:"omitempty,min=1,max=100" default:"10" example:"10"`
	Keyword    string `form:"keyword" example:"report"`
	Category   string `form:"category" example:"document"`
	Storage    string `form:"storage" example:"local"`
	UploaderID *uint  `form:"uploader_id" swaggertype:"integer" example:"1"`
	StartTime  string `form:"start_time" example:"2026-04-01T00:00:00Z"`
	EndTime    string `form:"end_time" example:"2026-04-30T23:59:59Z"`
}

type FileResponse struct {
	ID           uint      `json:"id" example:"1"`
	UploaderID   *uint     `json:"uploader_id,omitempty" swaggertype:"integer" example:"1"`
	Storage      string    `json:"storage" example:"local"`
	Category     string    `json:"category" example:"document"`
	OriginalName string    `json:"original_name" example:"report.pdf"`
	ObjectName   string    `json:"object_name" example:"1b9bcb9f0f7c4e0e9b48bb86d3ce28cb.pdf"`
	Ext          string    `json:"ext" example:".pdf"`
	MIMEType     string    `json:"mime_type" example:"application/pdf"`
	Size         int64     `json:"size" example:"102400"`
	URL          string    `json:"url" example:"/uploads/2026/04/14/1b9bcb9f0f7c4e0e9b48bb86d3ce28cb.pdf"`
	CreatedAt    time.Time `json:"created_at" swaggertype:"string" example:"2026-04-14T12:00:00Z"`
}

type FilePageResult struct {
	Items    []FileResponse `json:"items"`
	Total    int64          `json:"total" example:"1"`
	Page     int            `json:"page" example:"1"`
	PageSize int            `json:"page_size" example:"10"`
}

func toResponse(item File) FileResponse {
	return FileResponse{
		ID:           item.ID,
		UploaderID:   item.UploaderID,
		Storage:      item.Storage,
		Category:     item.Category,
		OriginalName: item.OriginalName,
		ObjectName:   item.ObjectName,
		Ext:          item.Ext,
		MIMEType:     item.MIMEType,
		Size:         item.Size,
		URL:          item.URL,
		CreatedAt:    item.CreatedAt,
	}
}
