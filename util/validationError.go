package util

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ApiError struct {
	Field string `json:"field"`
	Message   string `json:"message"`
	Detail string `json:"detail"`

}
func ErrorValidator(err error) map[string]any{
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]ApiError, len(ve))
		for i, fe := range ve {
			out[i] = ApiError{fe.Field(), msgForTag(fe.Tag(),fe.Error()),fe.Error()}
		}
		return gin.H{"errors": out}
	}
	return  gin.H{"errors": []ApiError{}}
}

func msgForTag(tag string,defaultError string) string {
	switch tag {
	case "required":
			return "Trường này không được để trống"
	case "min":
			return "Trường này chưa đủ độ dài tối thiểu"
	case "max":
			return "Trường này vượt quá độ dài tối đa"
	case "alphanum":
			return "Trường này chỉ được chứa chữ và số"
	case "email":
			return "Email không hợp lệ"
	case "currency":
			return "Loại tiền phải thuộc các giá trị sau : USD, EUR,"
		default:
			return defaultError
	}
}