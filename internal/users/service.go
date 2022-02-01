package users

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
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

func GenerateToken(userId uuid.UUID) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userId
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}

func GetWSLink(userId uuid.UUID, userName string, e echo.Context) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(err)
	}
	token, _ := GenerateToken(userId)
	resp := LoginUserResponse{
		"ws://" + os.Getenv("APP_URL") + ":" + os.Getenv("port") + "/chat/ws.rtm.start/" + token,
	}
	json.NewEncoder(e.Response()).Encode(resp)
}
