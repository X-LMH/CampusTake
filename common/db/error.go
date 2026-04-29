package db

import (
	"errors"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// IsDuplicateErr 判断是否是唯一键冲突
func IsDuplicateErr(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}

	// 兼容 gorm（不同版本/驱动可能用到）
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	return false
}
