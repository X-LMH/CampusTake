package db

import (
	"errors"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// IsDuplicateErr 判断是否是唯一键冲突
func IsDuplicateErr(err error) bool {
	if err == nil {
		return false
	}

	// 1. 优先判断 GORM 通用错误
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	// 2. 判断 MySQL 驱动特定错误
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return true
	}

	return false
}
