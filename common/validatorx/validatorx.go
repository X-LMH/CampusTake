package validatorx

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	once     sync.Once
)

func V() *validator.Validate {
	once.Do(func() {
		validate = validator.New()

		// 注册字段名提取：优先 label，其次 json，最后结构体字段名
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			if label := fld.Tag.Get("label"); label != "" && label != "-" {
				return label
			}

			if jsonTag := fld.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
				name := strings.Split(jsonTag, ",")[0]
				if name != "" {
					return name
				}
			}

			return fld.Name
		})
	})
	return validate
}

// ValidateStruct 对外统一入口
func ValidateStruct(v interface{}) error {
	return V().Struct(v)
}

// TranslateValidationError 返回所有中文错误，以；分割
func TranslateValidationError(err error) string {
	if err == nil {
		return ""
	}

	verrs, ok := err.(validator.ValidationErrors)
	if !ok || len(verrs) == 0 {
		return "参数校验失败"
	}

	msgs := make([]string, 0, len(verrs))
	for _, fe := range verrs {
		msgs = append(msgs, formatError(fe))
	}

	return strings.Join(msgs, "；")
}

// 你给的格式化逻辑增强版：支持更多 tag + label 字段名
func formatError(err validator.FieldError) string {
	field := err.Field() // 这里会拿到 RegisterTagNameFunc 处理后的名字（label/json）

	switch err.Tag() {
	case "required":
		return field + "不能为空"
	case "eqfield":
		return field + "输入不一致"
	case "min":
		// 字符串长度/数值最小值
		return fmt.Sprintf("%s太短，最少为%s", field, err.Param())
	case "max":
		return fmt.Sprintf("%s太长，最多为%s", field, err.Param())
	case "len":
		return fmt.Sprintf("%s长度必须为%s", field, err.Param())
	case "email":
		return field + "格式不正确"
	case "oneof":
		return fmt.Sprintf("%s必须是[%s]中的一个", field, err.Param())
	case "gte":
		return fmt.Sprintf("%s不能小于%s", field, err.Param())
	case "lte":
		return fmt.Sprintf("%s不能大于%s", field, err.Param())
	default:
		return field + "校验失败"
	}
}
