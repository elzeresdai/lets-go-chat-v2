package users

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"lets-go-chat-v2/internal/handlers"
	"lets-go-chat-v2/internal/middleware"
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

func (h *handler) LoginUser(ctx echo.Context) error {
	return nil
}
