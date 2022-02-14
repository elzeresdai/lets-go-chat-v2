package users

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/maxchagin/go-memorycache-example"
	"lets-go-chat-v2/internal/auth"
	"lets-go-chat-v2/internal/customerrors"
	"lets-go-chat-v2/internal/handlers"
	"lets-go-chat-v2/internal/middleware"
	"lets-go-chat-v2/pkg/hasher"
	"lets-go-chat-v2/pkg/logging"
	"net/http"
	"strings"
)

type Handler struct {
	Logger     *logging.Logger
	Repository RepositoryInterface
	Cache      memorycache.Cache
}

func NewHandler(repository RepositoryInterface, logger *logging.Logger, cache memorycache.Cache) handlers.HandlerInterface {
	return &Handler{
		Repository: repository,
		Logger:     logger,
		Cache:      cache,
	}
}

func (h *Handler) Register(e *echo.Echo) {
	e.POST("/user", middleware.ErrorMiddleware(h.CreateUser))
	e.POST("/user/login", middleware.ErrorMiddleware(h.LoginUser))
	e.GET("user/active", middleware.ErrorMiddleware(h.ActiveUsers))
}

func (h *Handler) CreateUser(e echo.Context) error {
	user, err, er := CreateUserReq(e)
	if err != nil {
		h.Logger.Error(err)
		return err
	}
	if er != nil {
		for _, ers := range er {
			e.Response().WriteHeader(http.StatusNotFound)
			e.Response().Write((*customerrors.AppError).Marshal(ers))
		}
		return nil
	}
	_, exist, ers := h.Repository.GetUser(e.Request().Context(), user.UserName)
	if ers != nil {
		return err
	}
	if exist {
		e.Response().WriteHeader(http.StatusBadRequest)
		resp := "User is already exist"
		json.NewEncoder(e.Response()).Encode(resp)
		return nil
	}

	newUser, err := h.Repository.CreateUser(e.Request().Context(), user)
	if err != nil {
		h.Logger.Error(err)
		return err
	}
	e.Response().WriteHeader(http.StatusOK)
	resp := CreateUserResponse{
		newUser.ID,
		newUser.UserName,
	}

	json.NewEncoder(e.Response()).Encode(resp)
	return nil
}

func (h *Handler) LoginUser(e echo.Context) error {
	user, err, er := CreateUserReq(e)
	if err != nil {
		h.Logger.Error(err)
		return err
	}
	if er != nil {
		for _, ers := range er {
			e.Response().WriteHeader(http.StatusNotFound)
			e.Response().Write((*customerrors.AppError).Marshal(ers))
		}
		return nil
	}
	login, exist, err := h.Repository.GetUser(e.Request().Context(), user.UserName)
	if err != nil {
		h.Logger.Error(err)
	}
	if !exist {
		e.Response().WriteHeader(http.StatusBadRequest)
		err := customerrors.NewAppError(
			nil,
			"User Not Found",
			"",
			"400",
		)
		h.Logger.Error(err)
		json.NewEncoder(e.Response()).Encode(err)

	}
	if !hasher.CheckPasswordHash(user.Password, login[0].PasswordHash) {
		e.Response().WriteHeader(http.StatusBadRequest)
		err := customerrors.NewAppError(
			nil,
			"Invalid Password",
			"",
			"400",
		)
		h.Logger.Error(err)
		json.NewEncoder(e.Response()).Encode(err)
		return nil
	}

	token, _ := auth.CreateJWTToken(login[0].UserName, login[0].ID)
	GetWSLink(*login[0], e, token)

	return nil
}
func (h *Handler) ActiveUsers(e echo.Context) error {
	memCache, err := h.Cache.Get("activeUsers")
	counter := 0
	if err && memCache != nil {
		arr := strings.Split(memCache.(string), ":")
		for _, value := range arr {
			if value != "" {
				counter++
			}
		}
	}
	resp := ActiveUsersResponse{counter}
	e.Response().WriteHeader(http.StatusOK)
	json.NewEncoder(e.Response()).Encode(resp)
	return nil
}
