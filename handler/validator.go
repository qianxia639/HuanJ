package handler

import (
	"Rejuv/internal/utils"

	"github.com/go-playground/validator/v10"
)

// 校验性别是否支持
var validGender = func(fieldLevel validator.FieldLevel) bool {
	if gender, ok := fieldLevel.Field().Interface().(int16); ok {
		return utils.IsSupportedGender(gender)
	}
	return false
}
