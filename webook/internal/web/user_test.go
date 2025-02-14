package web

import (
	"GoInAction/webook/internal/domain"
	"GoInAction/webook/internal/service"
	svcmocks "GoInAction/webook/internal/service/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUserHandler_SignUpV1(t *testing.T) {

	testCases := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		reqBuilder func(t *testing.T) *http.Request
		wantCode   int
		wantBody   string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil)
				return userSvc, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body, err := json.Marshal(SingUpReq{
					Email:           "test@qq.com",
					Password:        "hello#world123",
					ConfirmPassword: "hello#world123",
				})
				require.NoError(t, err)
				req := httptest.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantCode: http.StatusOK,
			wantBody: "注册成功",
		},
		{
			name: "绑定请求参数失败",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				return nil, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/users/signup", nil)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantCode: http.StatusOK,
			wantBody: "系统错误",
		},
		{
			name: "邮箱格式不对",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				return nil, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body, err := json.Marshal(SingUpReq{
					Email:           "invalid",
					Password:        "hello#world123",
					ConfirmPassword: "hello#world123",
				})
				require.NoError(t, err)
				req := httptest.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantCode: http.StatusOK,
			wantBody: "你的邮箱格式不对",
		},
		{
			name: "两次密码不一致",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				return nil, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body, err := json.Marshal(SingUpReq{
					Email:           "test@qq.com",
					Password:        "Hello#world123",
					ConfirmPassword: "Hello#world124",
				})
				require.NoError(t, err)
				req := httptest.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantCode: http.StatusOK,
			wantBody: "两次输入的密码不一致",
		},

		{
			name: "密码格式不对",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				return nil, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body, err := json.Marshal(SingUpReq{
					Email:           "test@qq.com",
					Password:        "hello",
					ConfirmPassword: "hello",
				})
				require.NoError(t, err)
				req := httptest.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantCode: http.StatusOK,
			wantBody: "密码必须大于8位，包含数字字符，且包含符号",
		},

		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "test@qq.com",
					Password: "Hello#world123",
				}).Return(service.ErrUserDuplicateEmail)
				return userSvc, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body, err := json.Marshal(SingUpReq{
					Email:           "test@qq.com",
					Password:        "Hello#world123",
					ConfirmPassword: "Hello#world123",
				})
				require.NoError(t, err)
				req := httptest.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantCode: http.StatusOK,
			wantBody: "邮箱冲突",
		},
		{
			name: "注册时发生系统错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(errors.New("system error"))
				return userSvc, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body, err := json.Marshal(SingUpReq{
					Email:           "test@qq.com",
					Password:        "Hello#world123",
					ConfirmPassword: "Hello#world123",
				})
				require.NoError(t, err)
				req := httptest.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantCode: http.StatusOK,
			wantBody: "系统错误",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)

			defer ctrl.Finish()

			// 构造 handler
			userSvc, codeSvc := tc.mock(ctrl)

			hdl := NewUserHandler(userSvc, codeSvc)

			// 准备服务器
			server := gin.Default()
			hdl.RegisterUsersRoutes(server)

			// 准备请求
			req := tc.reqBuilder(t)

			// 准备记录响应
			recorder := httptest.NewRecorder()

			// 执行
			server.ServeHTTP(recorder, req)

			// 验证
			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantBody, recorder.Body.String())
		})
	}

}

// 测试登录功能
func TestUserHandler_Login(t *testing.T) {
	testCases := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		reqBuilder func(t *testing.T) *http.Request
		wantCode   int
		wantBody   string
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().Login(gomock.Any(), "test@qq.com", "Hello#world123").
					Return(domain.User{
						Id:    123,
						Email: "test@qq.com",
					}, nil)
				return userSvc, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body, err := json.Marshal(LoginReq{
					Email:    "test@qq.com",
					Password: "Hello#world123",
				})
				require.NoError(t, err)
				req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantCode: http.StatusOK,
			wantBody: "登录成功",
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(domain.User{}, service.ErrUserNotFound)
				return userSvc, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body, err := json.Marshal(LoginReq{
					Email:    "notexist@qq.com",
					Password: "Hello#world123",
				})
				require.NoError(t, err)
				req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantCode: http.StatusOK,
			wantBody: "用户不存在",
		},
		{
			name: "密码错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(domain.User{}, service.ErrInvalidUserOrPassword)
				return userSvc, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body, err := json.Marshal(LoginReq{
					Email:    "test@qq.com",
					Password: "wrong",
				})
				require.NoError(t, err)
				req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantCode: http.StatusOK,
			wantBody: "用户名或密码不对",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userSvc, codeSvc := tc.mock(ctrl)
			hdl := NewUserHandler(userSvc, codeSvc)

			server := gin.Default()
			// 添加session中间件
			store := memstore.NewStore([]byte("secret"))
			server.Use(sessions.Sessions("mysession", store))
			hdl.RegisterUsersRoutes(server)

			req := tc.reqBuilder(t)
			recorder := httptest.NewRecorder()

			server.ServeHTTP(recorder, req)

			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantBody, recorder.Body.String())

			// 验证session
			if tc.name == "登录成功" {
				// 获取response中的cookie
				cookies := recorder.Result().Cookies()
				// 找到session cookie
				var cookie *http.Cookie
				for _, c := range cookies {
					if c.Name == "mysession" {
						cookie = c
						break
					}
				}
				require.NotNil(t, cookie)

				// 验证session中的userId
				s, err := store.Get(req, "mysession")
				require.NoError(t, err)
				val := s.Values["userId"]
				require.Equal(t, int64(123), val)
			}
		})
	}
}

// 测试发送验证码
func TestUserHandler_SendLoginSMSCode(t *testing.T) {
	testCases := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		reqBuilder func(t *testing.T) *http.Request
		wantCode   int
		wantResp   Result
	}{
		{
			name: "发送成功",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				codeSvc.EXPECT().Send(gomock.Any(), "login", "13800138000").Return(nil)
				return nil, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body, err := json.Marshal(map[string]string{
					"phone": "13800138000",
				})
				require.NoError(t, err)
				req := httptest.NewRequest(http.MethodPost, "/users/login_sms/code/send", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantCode: http.StatusOK,
			wantResp: Result{
				Code: 0,
				Msg:  "发送成功",
			},
		},
		{
			name: "发送太频繁",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				codeSvc.EXPECT().Send(gomock.Any(), "login", gomock.Any()).
					Return(service.ErrCodeSendTooMany)
				return nil, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body, err := json.Marshal(map[string]string{
					"phone": "13800138000",
				})
				require.NoError(t, err)
				req := httptest.NewRequest(http.MethodPost, "/users/login_sms/code/send", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantCode: http.StatusOK,
			wantResp: Result{
				Code: 100000,
				Msg:  "发送太频繁，请稍后重试",
			},
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				codeSvc.EXPECT().Send(gomock.Any(), "login", gomock.Any()).
					Return(errors.New("system error"))
				return nil, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body, err := json.Marshal(map[string]string{
					"phone": "13800138000",
				})
				require.NoError(t, err)
				req := httptest.NewRequest(http.MethodPost, "/users/login_sms/code/send", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantCode: http.StatusOK,
			wantResp: Result{
				Code: 100000,
				Msg:  "系统错误",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userSvc, codeSvc := tc.mock(ctrl)
			hdl := NewUserHandler(userSvc, codeSvc)

			server := gin.Default()
			hdl.RegisterUsersRoutes(server)

			req := tc.reqBuilder(t)
			recorder := httptest.NewRecorder()

			server.ServeHTTP(recorder, req)

			assert.Equal(t, tc.wantCode, recorder.Code)
			var result Result
			err := json.Unmarshal(recorder.Body.Bytes(), &result)
			require.NoError(t, err)
			assert.Equal(t, tc.wantResp, result)
		})
	}
}

// 测试验证码登录
func TestUserHandler_LoginSMSCode(t *testing.T) {
	testCases := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		reqBuilder func(t *testing.T) *http.Request
		wantCode   int
		wantResp   Result
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				userSvc := svcmocks.NewMockUserService(ctrl)

				codeSvc.EXPECT().Verify(gomock.Any(), "login", "123456", "13800138000").
					Return(true, nil)
				userSvc.EXPECT().FindOrCreate(gomock.Any(), "13800138000").
					Return(domain.User{Id: 123}, nil)

				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body, err := json.Marshal(map[string]string{
					"phone": "13800138000",
					"code":  "123456",
				})
				require.NoError(t, err)
				req := httptest.NewRequest(http.MethodPost, "/users/login_sms", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantCode: http.StatusOK,
			wantResp: Result{
				Code: 0,
				Msg:  "验证成功",
			},
		},
		{
			name: "验证码错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				codeSvc.EXPECT().Verify(gomock.Any(), "login", gomock.Any(), gomock.Any()).
					Return(false, nil)
				return nil, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body, err := json.Marshal(map[string]string{
					"phone": "13800138000",
					"code":  "invalid",
				})
				require.NoError(t, err)
				req := httptest.NewRequest(http.MethodPost, "/users/login_sms", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantCode: http.StatusOK,
			wantResp: Result{
				Code: 100000,
				Msg:  "验证码错误",
			},
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				codeSvc.EXPECT().Verify(gomock.Any(), "login", gomock.Any(), gomock.Any()).
					Return(false, errors.New("system error"))
				return nil, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body, err := json.Marshal(map[string]string{
					"phone": "13800138000",
					"code":  "error",
				})
				require.NoError(t, err)
				req := httptest.NewRequest(http.MethodPost, "/users/login_sms", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantCode: http.StatusOK,
			wantResp: Result{
				Code: 100000,
				Msg:  "系统错误",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userSvc, codeSvc := tc.mock(ctrl)
			hdl := NewUserHandler(userSvc, codeSvc)

			server := gin.Default()
			hdl.RegisterUsersRoutes(server)

			req := tc.reqBuilder(t)
			recorder := httptest.NewRecorder()

			server.ServeHTTP(recorder, req)

			assert.Equal(t, tc.wantCode, recorder.Code)
			var result Result
			err := json.Unmarshal(recorder.Body.Bytes(), &result)
			require.NoError(t, err)
			assert.Equal(t, tc.wantResp, result)
		})
	}
}

// 测试JWT登录
func TestUserHandler_LoginJWT(t *testing.T) {
	testCases := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		reqBuilder func(t *testing.T) *http.Request
		wantCode   int
		wantBody   string
		wantToken  bool
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().Login(gomock.Any(), "test@qq.com", "Hello#world123").
					Return(domain.User{
						Id:    123,
						Email: "test@qq.com",
					}, nil)
				return userSvc, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body, err := json.Marshal(LoginReq{
					Email:    "test@qq.com",
					Password: "Hello#world123",
				})
				require.NoError(t, err)
				req := httptest.NewRequest(http.MethodPost, "/users/loginJWT", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantCode:  http.StatusOK,
			wantBody:  "登录成功",
			wantToken: true,
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(domain.User{}, service.ErrUserNotFound)
				return userSvc, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body, err := json.Marshal(LoginReq{
					Email:    "notexist@qq.com",
					Password: "Hello#world123",
				})
				require.NoError(t, err)
				req := httptest.NewRequest(http.MethodPost, "/users/loginJWT", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantCode:  http.StatusOK,
			wantBody:  "用户不存在",
			wantToken: false,
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(domain.User{}, errors.New("system error"))
				return userSvc, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body, err := json.Marshal(LoginReq{
					Email:    "test@qq.com",
					Password: "Hello#world123",
				})
				require.NoError(t, err)
				req := httptest.NewRequest(http.MethodPost, "/users/loginJWT", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			wantCode:  http.StatusOK,
			wantBody:  "系统错误",
			wantToken: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userSvc, codeSvc := tc.mock(ctrl)
			hdl := NewUserHandler(userSvc, codeSvc)

			server := gin.Default()
			hdl.RegisterUsersRoutes(server)

			req := tc.reqBuilder(t)
			recorder := httptest.NewRecorder()

			server.ServeHTTP(recorder, req)

			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantBody, recorder.Body.String())
			if tc.wantToken {
				assert.NotEmpty(t, recorder.Header().Get("x-jwt-token"))
			} else {
				assert.Empty(t, recorder.Header().Get("x-jwt-token"))
			}
		})
	}
}

// 测试个人资料
func TestUserHandler_Profile(t *testing.T) {
	testCases := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		setupAuth  func(ctx *gin.Context) // 修改为 gin.Context
		wantCode   int
		wantResult domain.User
	}{
		{
			name: "获取成功",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().Profile(gomock.Any(), int64(123)).
					Return(domain.User{
						Id:    123,
						Email: "test@qq.com",
					}, nil)
				return userSvc, nil
			},
			setupAuth: func(ctx *gin.Context) {
				// 直接设置用户ID，模拟中间件认证后的状态
				ctx.Set("userId", int64(123))
			},
			wantCode: http.StatusOK,
			wantResult: domain.User{
				Id:    123,
				Email: "test@qq.com",
			},
		},
		{
			name: "未授权",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				return nil, nil
			},
			setupAuth: func(ctx *gin.Context) {
				// 不设置认证信息，模拟未登录状态
			},
			wantCode: http.StatusUnauthorized,
		},
		{
			name: "查询用户失败",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().Profile(gomock.Any(), int64(123)).
					Return(domain.User{}, errors.New("db error"))
				return userSvc, nil
			},
			setupAuth: func(ctx *gin.Context) {
				ctx.Set("userId", int64(123))
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userSvc, codeSvc := tc.mock(ctrl)
			hdl := NewUserHandler(userSvc, codeSvc)

			server := gin.Default()
			// 添加一个测试专用的中间件来模拟认证
			server.Use(func(ctx *gin.Context) {
				tc.setupAuth(ctx)
				ctx.Next()
			})
			hdl.RegisterUsersRoutes(server)

			req := httptest.NewRequest(http.MethodGet, "/users/profile", nil)
			recorder := httptest.NewRecorder()

			server.ServeHTTP(recorder, req)

			assert.Equal(t, tc.wantCode, recorder.Code)
			if tc.wantCode == http.StatusOK {
				var result domain.User
				err := json.Unmarshal(recorder.Body.Bytes(), &result)
				require.NoError(t, err)
				assert.Equal(t, tc.wantResult, result)
			}
		})
	}
}

type SingUpReq struct {
	Email           string `json:"email"`
	ConfirmPassword string `json:"confirmPassword"`
	Password        string `json:"password"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
