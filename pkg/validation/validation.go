package validation

//import (
//	"fmt"
//	"github.com/go-playground/validator/v10"
//)
//
//func IsValidRequest() bool {
//	v := validator.New()
//	_ = v.RegisterValidation("passwd", func(fl validator.FieldLevel) bool {
//		return len(fl.Field().String()) > 6
//	})
//	err := v.Struct(user)
//
//	for _, e := range err.(validator.ValidationErrors) {
//		fmt.Println(e)
//	}
//
//	return false
//
//}
