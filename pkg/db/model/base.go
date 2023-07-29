package model

import "time"

type BaseModel struct {
	ID            uint
	CreatedAt     *time.Time
	LastUpdatedAt *time.Time
	CreatedBy     uint
	LastUpdatedBy uint
	notAudit      bool // 不自动更新审计字段
}
