package users

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"lets-go-chat-v2/internal/customerrors"
	"lets-go-chat-v2/internal/handlers"
	"lets-go-chat-v2/internal/middleware"
	"lets-go-chat-v2/pkg/hasher"
	"lets-go-chat-v2/pkg/logging"
	"net/http"
)

type handler struct {
	logger     *logging.Logger
	repository RepositoryInterface
}

func NewHandler(repository RepositoryInterface, logger *logging.Logger) handlers.HandlerInterface {
	return &handler{
		repository: repository,
		logger:     logger,
	}
}

func (h *handler) Register(e *echo.Echo) {
	e.POST("/users", middleware.ErrorMiddleware(h.CreateUser))
	e.POST("/users/login", middleware.ErrorMiddleware(h.LoginUser))
	e.GET("user/active", middleware.ErrorMiddleware(h.ActiveUsers))
}

func (h *handler) CreateUser(e echo.Context) error {
	user, err := CreateUserReq(e)
	if err != nil {
		h.logger.Error(err)
		return err
	}
	_, exist, er := h.repository.GetUser(e.Request().Context(), user.UserName)
	if er != nil {
		return err
	}
	if exist {
		e.Response().WriteHeader(http.StatusBadRequest)
		resp := "User is already exist"
		json.NewEncoder(e.Response()).Encode(resp)
		return nil
	}

	newUser, err := h.repository.CreateUser(e.Request().Context(), user)
	if err != nil {
		h.logger.Error(err)
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

func (h *handler) LoginUser(e echo.Context) error {
	user, err := CreateUserReq(e)
	if err != nil {
		h.logger.Error(err)
		return err
	}
	login, exist, _ := h.repository.GetUser(e.Request().Context(), user.UserName)
	if !exist {
		e.Response().WriteHeader(http.StatusBadRequest)
		err := customerrors.NewAppError(
			nil,
			"User Not Found",
			"",
			"400",
		)
		h.logger.Error(err)
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
		h.logger.Error(err)
		json.NewEncoder(e.Response()).Encode(err)
	}

	GetWSLink(login[0].ID, user.UserName, e)

	return nil
}
func (h *handler) ActiveUsers(e echo.Context) error {
	return nil
}