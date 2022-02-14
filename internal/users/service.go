package users

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"lets-go-chat-v2/internal/customerrors"
	"lets-go-chat-v2/pkg/utils/cache"
	"os"
	"time"
)

func CreateUserReq(e echo.Context) (*CreateUserRequest, error, []*customerrors.AppError) {
	user := CreateUserRequest{}
	err := json.NewDecoder(e.Request().Body).Decode(&user)
	defer e.Request().Body.Close()
	if err != nil {
		return nil, err, nil
	}
	er := IsValidRequest(&user)
	if er != nil {
		return nil, nil, er
	}

	return &user, nil, nil
}

func IsValidRequest(user *CreateUserRequest) []*customerrors.AppError {
	v := validator.New()
	err := v.Struct(user)
	if err == nil {
		return nil
	}
	if err != nil {

		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]*customerrors.AppError, len(ve))
			for i, fe := range ve {
				out[i] = &customerrors.AppError{
					Code:             "404",
					Message:          msgForTag(fe.Tag()),
					DeveloperMessage: fe.Error(),
				}
			}
			return out

		}

	}
	return nil
}

func GetWSLink(user User, e echo.Context, token string) {
	err := godotenv.Load("../../.env")
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

func msgForTag(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "userName":
		return "Invalid userName"
	case "password":
		return "Invalid password"
	}
	return ""

}
