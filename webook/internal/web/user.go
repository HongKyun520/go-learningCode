package web

import (
	"GoInAction/webook/internal/domain"
	"GoInAction/webook/internal/service"
	"fmt"
	"net/http"
	"strings"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
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
	jwtHandler
	svc         service.UserService
	codeSvc     service.CodeService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

// 预编译正则表达式
func NewUserHandler(svc service.UserService, codeSvc service.CodeService) *UserHandler {
	return &UserHandler{
		jwtHandler:  NewJwtHandler(),
		svc:         svc,
		emailExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		codeSvc:     codeSvc,
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
	group.POST("/login_sms/code/send", u.SendLoginSMSCode)
	//group.POST("/login_sms/code/verify", u.VerifyLoginSMSCode)
	group.POST("/login_sms", u.LoginSMSCode)
	group.POST("/refresh_token", u.RefreshToken)
}

/*
*
验证码发送
*/
func (u *UserHandler) SendLoginSMSCode(ctx *gin.Context) {

	type Req struct {
		Phone string `json:"phone"`
	}

	const biz = "login"
	var req Req
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := u.codeSvc.Send(ctx, biz, req.Phone)

	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Code: 000000,
			Msg:  "发送成功",
		})

	case service.ErrCodeSendTooMany:
		zap.L().Warn("验证码发送太频繁")
		ctx.JSON(http.StatusOK, Result{
			Code: 100000,
			Msg:  "发送太频繁，请稍后重试",
		})

	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 100000,
			Msg:  "系统错误",
		})

	}
}

func (u *UserHandler) RefreshToken(ctx *gin.Context) {

	// 约定好，refresh token 放在Authorization中
	tokenStr := ExtractToken(ctx)
	var rc RefreshTokenClaims
	token, err := jwt.ParseWithClaims(tokenStr, &rc, func(token *jwt.Token) (interface{}, error) {
		return u.signKey, nil
	})

	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if token == nil || !token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	u.setJWTToken(ctx, rc.Uid)
	ctx.JSON(http.StatusOK, Result{
		Msg: "刷新成功",
	})

}

func (u *UserHandler) LoginSMSCode(ctx *gin.Context) {

	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}

	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	ok, err := u.codeSvc.Verify(ctx, bizLogin, req.Code, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 100000,
			Msg:  "系统错误",
		})
		zap.L().Error("验证码验证失败", zap.Error(err))
		return
	}

	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 100000,
			Msg:  "验证码错误",
		})
		return
	}

	// 根据手机号查uid
	user, err := u.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 100000,
			Msg:  "系统错误",
		})
		return
	}

	// 设置jwt
	if err = u.setJWTToken(ctx, user.Id); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 100000,
			Msg:  "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Code: 000000,
		Msg:  "验证成功",
	})

}

// SignUp 用户注册接口
// 1. 校验邮箱格式是否合法
// 2. 校验密码是否符合要求:
//   - 长度大于8位
//   - 包含数字和字母
//   - 包含特殊符号
//
// 3. 校验两次输入密码是否一致
// 4. 调用 service 层完成注册
func (u *UserHandler) SignUp(ctx *gin.Context) {
	// 请求参数结构体
	type SingUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}

	var req SingUpReq
	// 绑定请求参数
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	// 校验邮箱格式
	if ok, _ := u.emailExp.MatchString(req.Email); !ok {
		ctx.String(http.StatusOK, "你的邮箱格式不对")
		return
	}

	// 校验两次密码是否一致
	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入的密码不一致")
		return
	}

	// 校验密码强度
	if ok, _ := u.passwordExp.MatchString(req.Password); !ok {
		ctx.String(http.StatusOK, "密码必须大于8位，包含数字字符，且包含符号")
		return
	}

	// 调用 service 层完成注册
	err := u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	// 处理注册结果
	if err == service.ErrUserDuplicateEmail {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.String(http.StatusOK, "注册成功")
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

	if err = u.setJWTToken(ctx, user.Id); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	fmt.Println(user)
	// 设置serssion
	ctx.String(http.StatusOK, "登录成功")

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

	ctx.String(http.StatusOK, "这是你的EditJWT")
}

// Profile
// 查看用户个人信息
// Profile 查看用户个人信息
// @Summary 获取用户个人资料
// @Description 根据用户ID获取用户的个人资料信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {object} domain.User "用户信息"
// @Failure 401 "未授权"
// @Failure 500 "服务器内部错误"
// @Router /users/profile [get]
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

func ExtractToken(ctx *gin.Context) string {
	authCode := ctx.GetHeader("Authorization")
	if authCode == "" {
		return ""
	}

	segs := strings.Split(authCode, " ")
	if len(segs) != 2 || segs[0] != "Bearer" {
		return ""
	}

	tokenStr := segs[1]
	return tokenStr
}
