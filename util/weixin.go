package util

import (
	"encoding/json"
	"io/ioutil"
	"qnmahjong/def"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

// WXExchange ...
type WXExchange struct {
	AppID     string `json:"appid"`
	Secret    string `json:"secret"`
	Code      string `json:"code"`
	GrantType string `json:"grant_type"`
}

// WXRefresh ...
type WXRefresh struct {
	AppID        string `json:"appid"`
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

// WXAuth ...
type WXAuth struct {
	AccessToken string `json:"access_token"`
	OpenID      string `json:"openid"`
}

// WXToken ...
type WXToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int32  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
}

// WXUserInfo ...
type WXUserInfo struct {
	OpenID     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int32    `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	Headimgurl string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	Unionid    string   `json:"unionid"`
}

// WXStatus ...
type WXStatus struct {
	Errcode int32  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// WXTokenExchange 通过code获取access_token
func (e WXExchange) WXTokenExchange() (token WXToken, err error) {
	url := "https://api.weixin.qq.com/sns/oauth2/access_token?" +
		"appid=" + e.AppID +
		"&secret=" + e.Secret +
		"&code=" + e.Code +
		"&grant_type=" + e.GrantType
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(result, &token)
	if err != nil {
		return
	}

	if token.AccessToken == "" {
		status := WXStatus{}
		err = json.Unmarshal(result, &status)
		if err != nil {
			return
		}

		log.WithFields(log.Fields{
			"status": status,
		}).Error(def.ErrWeixinLogin)
		err = def.ErrWeixinLogin
	}
	return
}

// WXTokenRefresh 刷新或续期access_token使用
func (r WXRefresh) WXTokenRefresh() (token WXToken, err error) {
	url := "https://api.weixin.qq.com/sns/oauth2/refresh_token?" +
		"appid=" + r.AppID +
		"&grant_type=" + r.GrantType +
		"&refresh_token=" + r.RefreshToken
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(result, &token)
	if err != nil {
		return
	}

	if token.AccessToken == "" {
		status := WXStatus{}
		err = json.Unmarshal(result, &status)
		if err != nil {
			return
		}

		log.WithFields(log.Fields{
			"status": status,
		}).Error(def.ErrWeixinLogin)
		err = def.ErrWeixinLogin
	}
	return
}

// WXTokenCheck 检验授权凭证（access_token）是否有效
func (a WXAuth) WXTokenCheck() (status WXStatus, err error) {
	url := "https://api.weixin.qq.com/sns/auth?" +
		"access_token=" + a.AccessToken +
		"&openid=" + a.OpenID
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(result, &status)
	if err != nil {
		return
	}

	if status.Errcode != 0 {
		log.WithFields(log.Fields{
			"status": status,
		}).Error(def.ErrWeixinLogin)
		err = def.ErrWeixinLogin
	}
	return
}

// WXTokenUserInfo 获取用户个人信息（UnionID机制）
func (a WXAuth) WXTokenUserInfo() (userInfo WXUserInfo, err error) {
	url := "https://api.weixin.qq.com/sns/userinfo?" +
		"access_token=" + a.AccessToken +
		"&openid=" + a.OpenID
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(result, &userInfo)
	if err != nil {
		return
	}

	if userInfo.OpenID == "" {
		status := WXStatus{}
		err = json.Unmarshal(result, &status)
		if err != nil {
			return
		}

		log.WithFields(log.Fields{
			"status": status,
		}).Error(def.ErrWeixinLogin)
		err = def.ErrWeixinLogin
	}
	return
}
