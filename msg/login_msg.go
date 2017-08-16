package msg

import (
	"qnmahjong/db"
	"qnmahjong/db/dao"
	"qnmahjong/def"
	"qnmahjong/pf"
	"qnmahjong/util"
	"time"

	"github.com/labstack/echo"
	"github.com/spf13/viper"
)

// LoginHandle handle login request
func LoginHandle(msg []byte, c echo.Context) (recv []byte, err error) {
	defer util.Stack()

	absMessage := &pf.AbsMessage{}
	err = absMessage.Unmarshal(msg)
	if err != nil {
		return
	}

	msgID := absMessage.GetMsgID()
	msgBody := absMessage.GetMsgBody()

	switch msgID {
	case int32(pf.Login):
		loginSend := &pf.LoginSend{}
		err = loginSend.Unmarshal(msgBody)
		if err != nil {
			return
		}

		token, id, loginRecv := handleLogin(loginSend, c)
		recv, err = loginRecv.Marshal()
		if err != nil {
			return
		}

		absMessage.Token = token
		util.LogSend(msgID, id, 0, loginSend, "Login")
		util.LogRecv(msgID, id, 0, loginRecv, "Login")
	default:
		err = def.ErrHandleLogin
		return
	}

	absMessage.MsgBody = recv
	recv, err = absMessage.Marshal()
	return
}

func handleLogin(send *pf.LoginSend, c echo.Context) (t string, id int32, recv *pf.LoginRecv) {
	recv = &pf.LoginRecv{
		Status:      def.StatusOK,
		LogicServer: viper.GetString("logic_address"),
	}

	switch send.GetLoginType() {
	case def.LoginTypeWX:

		WeixinAppID := ""
		WeixinSecret := ""
		switch send.Channel {
		case def.ChannelIOSHB, def.ChannelAndroidHB:
			WeixinAppID = def.WeixinAppIDHB
			WeixinSecret = def.WeixinSecretHB
		case def.ChannelIOSHN, def.ChannelAndroidHN:
			WeixinAppID = def.WeixinAppIDHN
			WeixinSecret = def.WeixinSecretHN
		}

		var token util.WXToken
		var userInfo util.WXUserInfo
		var err error

		if send.GetMachineID() == "test1" {
			token.OpenID = def.QuickLoginOpenID
			userInfo.Headimgurl = def.QuickLoginHeadimgurl
		} else if send.GetSession() == def.QuickLoginSession {
			token.OpenID = def.QuickLoginOpenID
			userInfo.Headimgurl = def.QuickLoginHeadimgurl
		} else if send.GetPlayerID() == def.QuickLoginPlayerID {
			token.OpenID = def.QuickLoginOpenID
			userInfo.Headimgurl = def.QuickLoginHeadimgurl
		} else {
			loginID := send.GetPlayerID()
			if loginID != 0 {
				dbPlayer, err := dao.PlayerByPlayerID(db.Pool, loginID)
				if err == nil {
					if send.GetRefreshToken() != dbPlayer.RefreshToken {
						recv.Status = def.StatusErrorWeixinExpires
						return
					}
					refresh := util.WXRefresh{
						AppID:        WeixinAppID,
						GrantType:    "refresh_token",
						RefreshToken: send.GetRefreshToken(),
					}

					token, err = refresh.WXTokenRefresh()
				} else {
					recv.Status = def.StatusErrorWeixinExpires
					return
				}
			} else {
				exchange := util.WXExchange{
					AppID:     WeixinAppID,
					Secret:    WeixinSecret,
					Code:      send.GetSession(),
					GrantType: "authorization_code",
				}
				token, err = exchange.WXTokenExchange()
			}

			if err != nil {
				recv.Status = def.StatusErrorWeixin
				return
			}

			auth := util.WXAuth{
				AccessToken: token.AccessToken,
				OpenID:      token.OpenID,
			}

			userInfo, err = auth.WXTokenUserInfo()
			if err != nil {
				recv.Status = def.StatusErrorWeixin
				return
			}
		}

		player, err := dao.PlayerByOpenid(db.Pool, token.OpenID)
		// 数据库不存在该用户
		if err != nil {
			expiresIn := time.Now().Add(time.Second * time.Duration(token.ExpiresIn))
			player = &dao.Player{
				HighID:       0,
				InviteAward:  0,
				Openid:       token.OpenID,
				AccessToken:  token.AccessToken,
				ExpiresIn:    expiresIn,
				RefreshToken: token.RefreshToken,
				Nickname:     userInfo.Nickname,
				Sex:          userInfo.Sex,
				Province:     userInfo.Province,
				City:         userInfo.City,
				Country:      userInfo.Country,
				Headimgurl:   userInfo.Headimgurl,
				Unionid:      userInfo.Unionid,
				Coins:        def.RegisterAward,
				Cards:        0,
			}

			// 注册用户
			err = player.Insert(db.Pool)
			if err != nil {
				util.LogError(err, "player", player, id, def.ErrInsertPlayer)
				recv.Status = def.StatusErrorLogin
				return
			}

			t := &dao.Treasure{
				PlayerID:   player.PlayerID,
				Reason:     def.Register,
				Coins:      def.RegisterAward,
				Cards:      0,
				ChangeTime: time.Now(),
			}

			// 注册奖励
			err = t.Insert(db.Pool)
			if err != nil {
				util.LogError(err, "treasure", t, id, def.ErrInsertTreasure)
			}

			register := &dao.Register{
				PlayerID:        player.PlayerID,
				RegisterChannel: send.GetChannel(),
				RegisterVersion: send.GetVersion(),
				RegisterType:    send.GetLoginType(),
				RegisterIP:      c.RealIP(),
				RegisterTime:    time.Now(),
				RegisterMachine: send.GetMachineID(),
			}

			// 注册记录
			err = register.Insert(db.Pool)
			if err != nil {
				util.LogError(err, "register", register, id, def.ErrInsertRegister)
			}
		}

		login := &dao.Login{
			PlayerID:     player.PlayerID,
			LoginChannel: send.GetChannel(),
			LoginVersion: send.GetVersion(),
			LoginType:    send.GetLoginType(),
			LoginIP:      c.RealIP(),
			LoginTime:    time.Now(),
			LoginMachine: send.GetMachineID(),
		}

		// 登录记录
		err = login.Insert(db.Pool)
		if err != nil {
			util.LogError(err, "login", login, player.PlayerID, def.ErrInsertLogin)
		}

		expiresIn := time.Now().Add(time.Second * time.Duration(token.ExpiresIn))
		player.AccessToken = token.AccessToken
		player.ExpiresIn = expiresIn
		player.RefreshToken = token.RefreshToken
		player.Nickname = userInfo.Nickname
		player.Sex = userInfo.Sex
		player.Province = userInfo.Province
		player.City = userInfo.City
		player.Country = userInfo.Country
		player.Headimgurl = userInfo.Headimgurl
		player.Unionid = userInfo.Unionid

		// 刷新玩家信息
		err = player.Update(db.Pool)
		if err != nil {
			util.LogError(err, "player", player, id, def.ErrUpdatePlayer)
		}

		var success bool
		var claims = &util.Claims{
			PlayerID:  player.PlayerID,
			Channel:   send.GetChannel(),
			Version:   send.GetVersion(),
			LoginType: send.GetLoginType(),
		}
		t, success = util.CreateToken(claims)
		if !success {
			recv.Status = def.StatusErrorLogin
			return
		}
		id = player.PlayerID
	default:
		recv.Status = def.StatusErrorLoginType
	}
	return
}
