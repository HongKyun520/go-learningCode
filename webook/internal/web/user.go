package web

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 正则表达式
const (
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`

	userIdKey = "userId"
	bizLogin  = "login"
)

// UserHandler 在它上面定义跟user有关的路由

// handler定义在哪，就在哪进行路由的注册

type UserHandler struct {
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

// 预编译正则表达式
func NewUserHandler() *UserHandler {
	return &UserHandler{
		emailExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
	}
}

// 路由注册 gin
func (u *UserHandler) RegisterUsersRoutes(ctx *gin.Engine) {
	group := ctx.Group("/user")

	group.POST("/signUp", u.SignUp)
	group.POST("/login", u.Login)
	group.POST("/edit", u.Edit)
	group.GET("/profit", u.Profile)
}

// SignUp 注册
func (u *UserHandler) SignUp(ctx *gin.Context) {

	// 定义结构体接受字段
	type SingUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}

	var req SingUpReq
	// gin的参数绑定
	// 解析错了，就会直接写回一个4xx的说明
	if err := ctx.Bind(&req); err != nil {
		return
	}

	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	if !ok {
		ctx.String(http.StatusOK, "你的邮箱格式不对")
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入的密码不一致")
		return
	}

	ok, err = u.passwordExp.MatchString(req.Email)
	if err != nil {
		// 记录日志
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	if !ok {
		ctx.String(http.StatusOK, "密码必须大于8位，包含数字字符")
		return
	}

	ctx.String(http.StatusOK, "注册成功")
	// 数据库操作
	fmt.Printf("%v", req)

}

// Login 登录
func (u *UserHandler) Login(ctx *gin.Context) {

}

// Edit 编辑
func (u *UserHandler) Edit(ctx *gin.Context) {

}

func (u *UserHandler) Profile(ctx *gin.Context) {

}
