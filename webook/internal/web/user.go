package web

import (
	"GoInAction/webook/internal/domain"
	"GoInAction/webook/internal/service"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 正则表达式
const (
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 和上面比起来，用 ` 看起来就比较清爽   限制长度 8-72
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,72}$`

	userIdKey = "userId"
	bizLogin  = "login"
)

// UserHandler 在它上面定义跟user有关的路由

// handler定义在哪，就在哪进行路由的注册

type UserHandler struct {
	svc         *service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

// 预编译正则表达式
func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		svc:         svc,
		emailExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
	}
}

// gin的路由注册
func (u *UserHandler) RegisterUsersRoutes(ctx *gin.Engine) {
	group := ctx.Group("/users")
	group.POST("/signup", u.SignUp)
	group.POST("/login", u.Login)
	group.POST("/edit", u.Edit)
	group.GET("/profit", u.Profile)
}

// SignUp 注册方法
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
	// 正则表达式有误
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

	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		// 记录日志
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	if !ok {
		ctx.String(http.StatusOK, "密码必须大于8位，包含数字字符，且包含符号")
		return
	}

	// 对象转换，然后调用一下svc的方法
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if err == service.ErrUserDuplicateEmail {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.String(http.StatusOK, "注册成功")
	// 数据库操作
	fmt.Printf("%v", req)
}

// Login 登录方法
func (u *UserHandler) Login(ctx *gin.Context) {

	// 定义结构体接受
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	// 手动绑定
	if err := ctx.Bind(&req); err != nil {
		return
	}

	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}

	if err == service.ErrUserNotFound {
		ctx.String(http.StatusOK, "用户不存在")
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	// 登录成功后，在当前session中设置用户id
	sess := sessions.Default(ctx)
	sess.Set("userId", user.Id)
	sess.Save()

	// 设置serssion
	ctx.String(http.StatusOK, "登录成功")
	return

}

// Edit 编辑方法
func (u *UserHandler) Edit(ctx *gin.Context) {
	ctx.String(http.StatusOK, "这是你的Edit")
}

// Profile
func (u *UserHandler) Profile(ctx *gin.Context) {
	ctx.String(http.StatusOK, "这是你的Profile")
}
