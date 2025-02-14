package wechat

import (
	"GoInAction/webook/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Service interface {
	AuthURL(ctx context.Context, state string) (string, error)

	VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error)
}

const authURL = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_userinfo&state=%s#wechat_redirect"

type wechatService struct {
	appID     string
	appSecret string
	client    *http.Client
}

var redirectURL = url.PathEscape("https://meoying.com/oauth2/wechat/callback")

func NewWechatService(appID string, appSecret string) Service {
	return &wechatService{
		appID:     appID,
		appSecret: appSecret,
		client:    http.DefaultClient,
	}
}

func (svc *wechatService) AuthURL(ctx context.Context, state string) (string, error) {
	//const authURL = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_userinfo&state=%s#wechat_redirect"
	const authURL = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect"
	return fmt.Sprintf(authURL, svc.appID, redirectURL, state), nil
}

func (svc *wechatService) VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error) {

	accessTokenUrl := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", svc.appID, svc.appSecret, code)

	req, err := http.NewRequestWithContext(ctx, "GET", accessTokenUrl, nil)
	if err != nil {
		return domain.WechatInfo{}, err
	}

	httpResp, err := svc.client.Do(req)
	if err != nil {
		return domain.WechatInfo{}, err
	}

	// 解析
	var res Result
	err = json.NewDecoder(httpResp.Body).Decode(&res)
	if err != nil {
		return domain.WechatInfo{}, err
	}

	if res.ErrCode != 0 {
		return domain.WechatInfo{},
			fmt.Errorf("调用微信接口失败")
	}

	return domain.WechatInfo{
		OpenId:  res.Openid,
		UnionId: res.UnionId,
	}, nil
}

type Result struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionId      string `json:"unionid"`

	// 错误返回
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
