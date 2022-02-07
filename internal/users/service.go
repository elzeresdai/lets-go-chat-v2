package users

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"lets-go-chat-v2/pkg/utils/cache"
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

func GetWSLink(user User, e echo.Context, token string) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(err)
	}
	resp := LoginUserResponse{
		"ws://" + os.Getenv("APP_URL") + ":" + os.Getenv("port") + "/chat/ws.rtm.start?token=" + token,
	}
	saveToCache(token, user.UserName)
	json.NewEncoder(e.Response()).Encode(resp)
}

func saveToCache(token string, userName string) {
	cacheConnections, err := cache.Cache.Get("webSocketUsers")
	if !err || cacheConnections == nil {
		cacheConnections = ""
	}
	cache.Cache.Set("activeUsers", cacheConnections.(string)+userName+"!"+token+":", 10*time.Minute)
}
