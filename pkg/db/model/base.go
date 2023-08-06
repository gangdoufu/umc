package model

import "time"

type BaseModel struct {
	ID            uint
	CreatedAt     *time.Time `json:"created_at"`
	LastUpdatedAt *time.Time `json:"last_updated_at"`
	CreatedBy     uint       `json:"created_by"`
	LastUpdatedBy uint       `json:"last_updated_by"`
	notAudit      bool       // 不自动更新审计字段
}
