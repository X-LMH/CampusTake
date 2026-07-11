package db

import "gorm.io/gorm"

const (
	// 默认页号与页大小
	DefaultPage     = 1
	DefaultPageSize = 10

	// 页大小边界
	MinPageSize = 1
	MaxPageSize = 100
)

// Paginate GORM 自动分页 Scope
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 1. 基础防护：页号至少为 DefaultPage
		if page < DefaultPage {
			page = DefaultPage
		}
		// 2. 最大/最小页大小限制（滤镜逻辑）
		switch {
		case pageSize > MaxPageSize:
			pageSize = MaxPageSize
		case pageSize < MinPageSize:
			pageSize = DefaultPageSize
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
