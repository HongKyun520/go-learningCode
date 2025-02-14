package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type jwtHandler struct {
	signMethod jwt.SigningMethod
	signKey    []byte
}

func NewJwtHandler() jwtHandler {
	return jwtHandler{signMethod: jwt.SigningMethodHS512, signKey: []byte("k6CswdUm75WKcbM68UQUuxVsHSpTCwgK")}
}

// setJWTToken 设置JWT token
// 该方法会同时设置access token和refresh token
// access token 过期时间为30分钟
// ctx 是gin的上下文
// uId 是用户ID
// 返回error表示设置过程中的错误
func (u *jwtHandler) setJWTToken(ctx *gin.Context, uId int64) error {
	// 先设置refresh token
	err := u.setRefreshToken(ctx, uId)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return err
	}

	// 构造JWT claims
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			// 设置登录态时间为30分钟
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		// 设置用户信息
		Uid:       uId,
		UserAgent: ctx.Request.UserAgent(),
	}

	// 生成JWT token
	token := jwt.NewWithClaims(u.signMethod, claims)
	tokenStr, err := token.SignedString(u.signKey)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return err
	}
	fmt.Println(tokenStr)
	// 设置到header中
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

func (u *jwtHandler) setRefreshToken(ctx *gin.Context, uId int64) error {
	claims := RefreshTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			// 设置7天过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		Uid: uId,
	}

	refreshToken := jwt.NewWithClaims(u.signMethod, claims)
	refreshTokenStr, err := refreshToken.SignedString(u.signKey)
	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", refreshTokenStr)
	return nil
}

type RefreshTokenClaims struct {
	jwt.RegisteredClaims
	Uid int64
}

type UserClaims struct {
	jwt.RegisteredClaims

	// 声明自己放进token的数据
	Uid       int64
	UserAgent string
}
