package users

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateUserReq(e echo.Context) (*CreateUserRequest, error) {
	user := CreateUserRequest{}
	err := json.NewDecoder(e.Request().Body).Decode(&user)
	defer e.Request().Body.Close()
	if err != nil {
		return nil, err
	}

	if !IsValidRequest(user) {
		e.Response().WriteHeader(http.StatusBadRequest)
		return nil, json.NewEncoder(e.Response()).Encode(ValidationErrs)
	}

	return &user, nil
}

var ValidationErrs []ValidationResponse

func IsValidRequest(user CreateUserRequest) bool {
	v := validator.New()
	err := v.Struct(user)

	if err == nil {
		return true
	}

	for _, er := range err.(validator.ValidationErrors) {
		errors := ValidationResponse{
			Message: er.Error(),
		}
		ValidationErrs = append(ValidationErrs, errors)
	}

	return false

}
