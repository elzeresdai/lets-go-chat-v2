package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/maxchagin/go-memorycache-example"
	"github.com/stretchr/testify/assert"
	"lets-go-chat-v2/internal/auth"
	mock_users "lets-go-chat-v2/internal/tests/mocks"
	"lets-go-chat-v2/internal/users"
	"lets-go-chat-v2/pkg/hasher"
	"lets-go-chat-v2/pkg/logging"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestHandler_LoginUser(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Log(err)
	}
	type mockBehavior func(repo *mock_users.MockRepositoryInterface, user *users.CreateUserRequest)
	type expectedResponse func(request *users.CreateUserRequest) string

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            users.CreateUserRequest
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody expectedResponse
	}{
		{
			name:      "Ok",
			inputBody: `{"username": "username",  "password": "qwerty123"}`,
			inputUser: users.CreateUserRequest{
				UserName: "username",
				Password: "qwerty123",
			},
			mockBehavior: func(repo *mock_users.MockRepositoryInterface, request *users.CreateUserRequest) {
				ctx := context.TODO()

				hashed, _ := hasher.HashPassword(request.Password)
				us := users.User{
					ID:           uuid.MustParse("4c0cce96-8d12-11ec-b909-0242ac120002"),
					UserName:     request.UserName,
					PasswordHash: hashed,
				}

				var s []*users.User
				s = append(s, &us)
				repo.EXPECT().GetUser(ctx, request.UserName).Return(s, true, nil)

			},
			expectedStatusCode: 200,
			expectedResponseBody: func(req *users.CreateUserRequest) string {
				hashed, _ := hasher.HashPassword(req.Password)
				us := users.User{
					ID:           uuid.MustParse("4c0cce96-8d12-11ec-b909-0242ac120002"),
					UserName:     req.UserName,
					PasswordHash: hashed,
				}
				token, _ := auth.CreateJWTToken(us.UserName, us.ID)
				resp := &users.LoginUserResponse{
					Url: "ws://" + os.Getenv("APP_URL") + ":" + os.Getenv("port") + "/chat/ws.rtm.start?token=" + token,
				}
				res, _ := json.Marshal(resp)
				return string(res) + "\n"
			},
		},
		{
			name:      "Validation errors",
			inputBody: `{"username": "usr",  "password": "qwer"}`,
			inputUser: users.CreateUserRequest{
				UserName: "usr",
				Password: "qwer",
			},
			mockBehavior: func(repo *mock_users.MockRepositoryInterface, user *users.CreateUserRequest) {
			},
			expectedStatusCode: 404,
			expectedResponseBody: func(request *users.CreateUserRequest) string {
				resp := users.IsValidRequest(request)
				var out string
				for _, res := range resp {
					re, _ := json.Marshal(res)
					out += string(re)
				}
				return out
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_users.NewMockRepositoryInterface(c)
			test.mockBehavior(repo, &test.inputUser)
			logger := logging.GetLogger()
			handler := users.Handler{Logger: logger, Repository: repo}

			// Init Endpoint
			e.POST("/user/login", handler.LoginUser)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/user/login",
				bytes.NewBufferString(test.inputBody))

			// Make Request
			e.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody(&test.inputUser))
		})
	}
}

func TestHandler_CreateUser(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Log(err)
	}
	type mockBehavior func(repo *mock_users.MockRepositoryInterface, user *users.CreateUserRequest)
	type expectedResponse func(user *users.CreateUserRequest) string

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            users.CreateUserRequest
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody expectedResponse
	}{
		{
			name:      "Ok",
			inputBody: `{"username": "username",  "password": "qwerty123"}`,
			inputUser: users.CreateUserRequest{
				UserName: "username",
				Password: "qwerty123",
			},
			mockBehavior: func(repo *mock_users.MockRepositoryInterface, request *users.CreateUserRequest) {
				ctx := context.TODO()

				hashed, _ := hasher.HashPassword(request.Password)
				us := users.User{
					ID:           uuid.MustParse("4c0cce96-8d12-11ec-b909-0242ac120002"),
					UserName:     request.UserName,
					PasswordHash: hashed,
				}
				var s []*users.User
				repo.EXPECT().GetUser(ctx, request.UserName).Return(s, false, nil)
				repo.EXPECT().CreateUser(ctx, request).Return(&us, nil)

			},
			expectedStatusCode: 200,
			expectedResponseBody: func(user *users.CreateUserRequest) string {
				us := users.CreateUserResponse{
					Id:       uuid.MustParse("4c0cce96-8d12-11ec-b909-0242ac120002"),
					UserName: user.UserName,
				}
				resp, _ := json.Marshal(us)
				return string(resp) + "\n"
			},
		},
		{
			name:      "User Exist",
			inputBody: `{"username": "username",  "password": "qwerty123"}`,
			inputUser: users.CreateUserRequest{
				UserName: "username",
				Password: "qwerty123",
			},
			mockBehavior: func(repo *mock_users.MockRepositoryInterface, request *users.CreateUserRequest) {
				ctx := context.TODO()

				hashed, _ := hasher.HashPassword(request.Password)
				us := users.User{
					ID:           uuid.MustParse("4c0cce96-8d12-11ec-b909-0242ac120002"),
					UserName:     request.UserName,
					PasswordHash: hashed,
				}
				var s []*users.User
				s = append(s, &us)
				repo.EXPECT().GetUser(ctx, request.UserName).Return(s, true, nil)

			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: func(user *users.CreateUserRequest) string {
				resp := "User is already exist"
				jsn, _ := json.Marshal(resp)
				return string(jsn) + "\n"
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_users.NewMockRepositoryInterface(c)
			test.mockBehavior(repo, &test.inputUser)
			logger := logging.GetLogger()
			handler := users.Handler{Logger: logger, Repository: repo}

			// Init Endpoint
			e.POST("/user", handler.CreateUser)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/user",
				bytes.NewBufferString(test.inputBody))

			// Make Request
			e.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody(&test.inputUser))
		})
	}
}

func TestHandler_ActiveUsers(t *testing.T) {
	type mockBehavior func()

	tests := []struct {
		name                 string
		expectedStatusCode   int
		expectedResponseBody string
		mockBehavior         mockBehavior
		wantErr              bool
	}{
		{
			name:               "ActiveUsers",
			wantErr:            false,
			expectedStatusCode: 200,
			mockBehavior: func() {
				cache := memorycache.New(5*time.Minute, 5*time.Minute)
				cache.Set("activeUsers", "testLogin!testToken:", 10*time.Second)
			},
			expectedResponseBody: "{}\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()
			repo := mock_users.NewMockRepositoryInterface(c)
			logger := logging.GetLogger()
			handler := users.Handler{Logger: logger, Repository: repo, Cache: memorycache.Cache{}}
			test.mockBehavior()
			// Init Endpoint
			e.POST("/user/active", handler.ActiveUsers)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/user/active", nil)

			// Make Request
			e.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})

	}
}
