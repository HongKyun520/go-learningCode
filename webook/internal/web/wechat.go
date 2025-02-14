package web

import (
	"GoInAction/webook/internal/service"
	"GoInAction/webook/internal/service/oauth2/wechat"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
)

type OAuth2WechatHandler struct {
	svc     wechat.Service
	userSvc service.UserService
	jwtHandler
	key             []byte
	stateCookieName string
}

func NewOAuth2WechatHandler(svc wechat.Service, userService service.UserService) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:             svc,
		userSvc:         userService,
		key:             []byte("k6CswdUm75WKcbM68UQUuxVsHSpTCwgK"),
		stateCookieName: "jwt-state",
	}
}

func (o *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", o.AuthURL)
	g.Any("/callback", o.Callback)
}

func (o *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {

	state := uuid.New()

	// 获取授权URL
	val, err := o.svc.AuthURL(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "构造跳转URL失败",
			Code: 5,
		})
		return
	}
	err = o.setStateCookie(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "系统错误",
			Code: 5,
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: val,
	})
}

func (o *OAuth2WechatHandler) Callback(ctx *gin.Context) {

	err := o.verifyStateCookie(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "系统错误",
			Code: 5,
		})
		return
	}
	code := ctx.Query("code")
	// state := ctx.Query("state")
	wechatInfo, err := o.svc.VerifyCode(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "验证码有误",
			Code: 4,
		})
		return
	}

	// 校验用户是否登录 （注册/denglu）
	u, err := o.userSvc.FindOrCreateByWechat(ctx, wechatInfo)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "系统错误",
			Code: 5,
		})
		return
	}

	o.setJWTToken(ctx, u.Id)
	// 返回成功
	ctx.JSON(http.StatusOK, Result{
		Msg: "登录成功",
	})

}

func (o *OAuth2WechatHandler) setStateCookie(ctx *gin.Context, state string) error {

	claims := StateClaims{
		State: state,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(o.key)
	if err != nil {
		return err
	}
	ctx.SetCookie(o.stateCookieName, tokenString, 600, "/oauth2/wechat/callback", "", false, true)
	return nil
}

func (o *OAuth2WechatHandler) verifyStateCookie(ctx *gin.Context) error {
	state := ctx.Query("state")
	ck, err := ctx.Cookie(o.stateCookieName)
	if err != nil {
		return fmt.Errorf("获取cookie失败")
	}

	var sc StateClaims
	_, err = jwt.ParseWithClaims(ck, &sc, func(token *jwt.Token) (interface{}, error) {
		return o.key, nil
	})
	if err != nil {
		return fmt.Errorf("解析token失败")
	}

	if state != sc.State {
		return fmt.Errorf("state不匹配")
	}
	return nil
}

type StateClaims struct {
	State string
	jwt.RegisteredClaims
}
