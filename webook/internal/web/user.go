package web

import (
	"GoInAction/webook/internal/domain"
	"GoInAction/webook/internal/service"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
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
	group.POST("/loginJWT", u.LoginJWT)
	group.POST("/edit", u.Edit)
	group.POST("/editJWT", u.EditJWT)
	group.GET("/profile", u.Profile)
}

// SignUp 注册方法
/**
1、校验密码格式、邮箱格式
2、入库
*/
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
	// 校验情况
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "系统错误")
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
/**
  检验数据库是否存在，账号密码是否正确
  在session中设置userId
*/
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
	sess.Options(sessions.Options{
		MaxAge: 30,
	})
	sess.Save()

	// 设置serssion
	ctx.String(http.StatusOK, "登录成功")
	return

}

func (u *UserHandler) LoginJWT(ctx *gin.Context) {

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

	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},

		Uid:       user.Id,
		UserAgent: ctx.Request.UserAgent(),
	}

	// 登录成功后，生成JWT 并设置user信息
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("k6CswdUm75WKcbM68UQUuxVsHSpTCwgK"))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
	}
	fmt.Println(tokenStr)
	ctx.Header("x-jwt-token", tokenStr)
	fmt.Println(user)

	// 设置serssion
	ctx.String(http.StatusOK, "登录成功")
	return

}

// Edit 编辑方法
// 允许用户补充个人信息

/**
允许用户补充个人信息
1、昵称：字符串，你需要考虑允许的长度
2、生日：前端输入为 1992-01-01 这种字符串。
3、个人简介：一段文本，你需要考虑允许的长度。
*/

func (u *UserHandler) Edit(ctx *gin.Context) {

	type EditReq struct {
		Id       string `json:"id"`
		Password string `json:"password" validate:"required"`
		Nickname string `json:"nickname" validate:"required"`
		Birthday string `json:"birthday" validate:"required"`
		Profile  string `json:"profile" validate:"required"`
	}

	var req EditReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.String(http.StatusOK, "这是你的Edit")
}

func (u *UserHandler) EditJWT(ctx *gin.Context) {

	// 获取uid
	c, _ := ctx.Get("userId")

	// 断言
	uid, ok := c.(int64)
	if !ok {
		ctx.String(http.StatusOK, "用户不存在")
		return
	}

	fmt.Println(uid)

	type EditReq struct {
		UserId   string `json:"id"`
		Password string `json:"password" validate:"required"`
		Nickname string `json:"nickname" validate:"required"`
		Birthday string `json:"birthday" validate:"required"`
		Profile  string `json:"profile" validate:"required"`
	}

	var req EditReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.String(http.StatusOK, "这是你的Edit")
}

// Profile
// 查看用户个人信息
func (u *UserHandler) Profile(ctx *gin.Context) {

	id := ctx.GetInt64("userId")
	if id == 0 {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// 获取用户信息
	profile, err := u.svc.Profile(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, profile)
}

type UserClaims struct {
	jwt.RegisteredClaims

	// 声明自己放进token的数据
	Uid       int64
	UserAgent string
}
