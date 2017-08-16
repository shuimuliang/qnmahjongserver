package msg

import (
	"qnmahjong/cache"
	"qnmahjong/def"
	"qnmahjong/msg/mj/hb"
	"qnmahjong/msg/mj/kf"
	"qnmahjong/msg/mj/tdh"
	"qnmahjong/msg/mj/zz"
	"qnmahjong/pf"
	"qnmahjong/util"

	"github.com/olahol/melody"
)

// LogicHandle handle logic request
func LogicHandle(msg []byte, s *melody.Session) (err error) {
	defer util.Stack()

	abs := &pf.AbsMessage{}
	err = abs.Unmarshal(msg)
	if err != nil {
		return
	}

	token := abs.Token
	msgID := abs.MsgID
	msgBody := abs.MsgBody
	c, success := util.ValidateToken(token)

	// 检查token
	if !success {
		var bytes []byte
		errorRecv := &pf.ErrorRecv{
			Status: def.StatusErrorToken,
		}

		bytes, err = errorRecv.Marshal()
		if err != nil {
			return
		}

		abs.MsgID = int32(pf.Error)
		abs.MsgBody = bytes
		bytes, err = abs.Marshal()
		if err != nil {
			return
		}

		err = s.WriteBinary(bytes)
		if err == nil {
			util.LogRecv(msgID, 0, 0, errorRecv, "Token")
		}
		return
	}

	// 获取玩家id
	playerID := c.PlayerID

	// auth
	if msgID == int32(pf.Auth) {
		send := &pf.AuthSend{}
		err = send.Unmarshal(msgBody)
		if err != nil {
			return
		}

		// 设置id
		s.Set("PlayerID", playerID)
		util.LogSend(msgID, playerID, 0, send, "Auth")
		err = handleAuth(s, c, abs)
		return
	}

	// dirty
	if msgID == int32(pf.Dirty) {
		send := &pf.DirtySend{}
		err = send.Unmarshal(msgBody)
		if err != nil {
			return
		}

		util.LogSend(msgID, playerID, 0, send, "Dirty")
		err = handleDirty(s, c, abs)
		if err != nil {
			return
		}
		err = handleBroadcast(s, c, abs)
		return
	}

	// 检查player
	mjPlayer := cache.GetMjPlayer(playerID)
	if mjPlayer == nil {
		var bytes []byte
		recv := &pf.ErrorRecv{
			Status: def.StatusErrorPlayer,
		}

		bytes, err = recv.Marshal()
		if err != nil {
			return
		}

		abs.MsgID = int32(pf.Error)
		abs.MsgBody = bytes
		bytes, err = abs.Marshal()
		if err != nil {
			return
		}

		err = s.WriteBinary(bytes)
		if err == nil {
			util.LogRecv(msgID, playerID, 0, recv, "Player")
		}
		return
	}

	// 获取房间id
	roomID := mjPlayer.RoomID
	mjType := mjPlayer.MJType

	// 检查room
	if roomID != 0 {
		exist := true
		switch mjType {
		case def.MjTypeHb:
			exist = hb.RoomCache.ExistPlayer(playerID, roomID)
		case def.MjTypeZz:
			exist = zz.RoomCache.ExistPlayer(playerID, roomID)
		case def.MjTypeTdh:
			exist = tdh.RoomCache.ExistPlayer(playerID, roomID)
		case def.MjTypeKf:
			exist = kf.RoomCache.ExistPlayer(playerID, roomID)
		}

		if !exist {
			var bytes []byte
			recv := &pf.ErrorRecv{
				Status: def.StatusErrorRoom,
			}

			bytes, err = recv.Marshal()
			if err != nil {
				return
			}

			abs.MsgID = int32(pf.Error)
			abs.MsgBody = bytes
			bytes, err = abs.Marshal()
			if err != nil {
				return
			}

			err = s.WriteBinary(bytes)
			if err == nil {
				util.LogRecv(msgID, playerID, roomID, recv, "Room")
			}
			return
		}
	}

	// 更新token
	token, success = util.CreateToken(c)
	if success {
		abs.Token = token
	}

	switch msgID {
	// 大厅相关
	case int32(pf.Logout):
		send := &pf.LogoutSend{}
		err = send.Unmarshal(msgBody)
		if err != nil {
			return
		}

		util.LogSend(msgID, playerID, roomID, send, "Logout")
		err = handleLogout(s, mjPlayer, abs)
	case int32(pf.Feedback):
		send := &pf.FeedbackSend{}
		err = send.Unmarshal(msgBody)
		if err != nil {
			return
		}

		util.LogSend(msgID, playerID, roomID, send, "Feedback")
		err = handleFeedback(s, mjPlayer, send, abs)
	// 邀请码支付相关
	case int32(pf.InviteCode):
		send := &pf.InviteCodeSend{}
		err = send.Unmarshal(msgBody)
		if err != nil {
			return
		}

		util.LogSend(msgID, playerID, roomID, send, "InviteCode")
		err = handleInviteCode(s, mjPlayer, send, abs)
	case int32(pf.InviteList):
		send := &pf.InviteListSend{}
		err = send.Unmarshal(msgBody)
		if err != nil {
			return
		}

		util.LogSend(msgID, playerID, roomID, send, "InviteList")
		err = handleInviteList(s, mjPlayer, abs)
	case int32(pf.InviteAward):
		send := &pf.InviteAwardSend{}
		err = send.Unmarshal(msgBody)
		if err != nil {
			return
		}

		util.LogSend(msgID, playerID, roomID, send, "InviteAward")
		err = handleInviteAward(s, mjPlayer, send, abs)
	case int32(pf.GoodsList):
		send := &pf.GoodsListSend{}
		err = send.Unmarshal(msgBody)
		if err != nil {
			return
		}

		util.LogSend(msgID, playerID, roomID, send, "GoodsList")
		err = handleGoodsList(s, mjPlayer, abs)
	case int32(pf.OrderApplyNew):
		send := &pf.OrderApplyNewSend{}
		err = send.Unmarshal(msgBody)
		if err != nil {
			return
		}

		util.LogSend(msgID, playerID, roomID, send, "OrderApplyNew")
		err = handleOrderApplyNew(s, mjPlayer, send, abs)
	// 麻将房间相关
	case int32(pf.EnterRoom):
		send := &pf.EnterRoomSend{}
		err = send.Unmarshal(msgBody)
		if err != nil {
			return
		}

		util.LogSend(msgID, playerID, roomID, send, "EnterRoom")
		if send.RoomID == 0 {
			// 创建房间
			mjType = send.MjType
		} else if send.Relink == 0 {
			// 进入房间
			mjType = send.RoomID / 10000
		} else {
			// 断线重连
			mjType = send.RoomID / 10000
		}
		switch mjType {
		case def.MjTypeHb:
			err = hb.HandleEnterRoom(s, mjPlayer, send, abs)
		case def.MjTypeZz:
			err = zz.HandleEnterRoom(s, mjPlayer, send, abs)
		case def.MjTypeTdh:
			err = tdh.HandleEnterRoom(s, mjPlayer, send, abs)
		case def.MjTypeKf:
			err = kf.HandleEnterRoom(s, mjPlayer, send, abs)
		default:
			err = hb.HandleEnterRoom(s, mjPlayer, send, abs)
		}
	case int32(pf.ExitRoom):
		send := &pf.ExitRoomSend{}
		err = send.Unmarshal(msgBody)
		if err != nil {
			return
		}

		util.LogSend(msgID, playerID, roomID, send, "ExitRoom")
		switch mjType {
		case def.MjTypeHb:
			err = hb.HandleExitRoom(s, mjPlayer, abs)
		case def.MjTypeZz:
			err = zz.HandleExitRoom(s, mjPlayer, abs)
		case def.MjTypeTdh:
			err = tdh.HandleExitRoom(s, mjPlayer, abs)
		case def.MjTypeKf:
			err = kf.HandleExitRoom(s, mjPlayer, abs)
		}
	case int32(pf.CloseRoom):
		send := &pf.CloseRoomSend{}
		err = send.Unmarshal(msgBody)
		if err != nil {
			return
		}

		util.LogSend(msgID, playerID, roomID, send, "CloseRoom")
		switch mjType {
		case def.MjTypeHb:
			err = hb.HandleCloseRoom(s, mjPlayer, abs)
		case def.MjTypeZz:
			err = zz.HandleCloseRoom(s, mjPlayer, abs)
		case def.MjTypeTdh:
			err = tdh.HandleCloseRoom(s, mjPlayer, abs)
		case def.MjTypeKf:
			err = kf.HandleCloseRoom(s, mjPlayer, abs)
		}
	case int32(pf.VoteClose):
		send := &pf.VoteCloseSend{}
		err = send.Unmarshal(msgBody)
		if err != nil {
			return
		}

		util.LogSend(msgID, playerID, roomID, send, "VoteClose")
		switch mjType {
		case def.MjTypeHb:
			err = hb.HandleVoteClose(s, mjPlayer, send, abs)
		case def.MjTypeZz:
			err = zz.HandleVoteClose(s, mjPlayer, send, abs)
		case def.MjTypeTdh:
			err = tdh.HandleVoteClose(s, mjPlayer, send, abs)
		case def.MjTypeKf:
			err = kf.HandleVoteClose(s, mjPlayer, send, abs)
		}
	case int32(pf.GameChat):
		send := &pf.GameChatSend{}
		err = send.Unmarshal(msgBody)
		if err != nil {
			return
		}

		util.LogSend(msgID, playerID, roomID, send, "GameChat")
		switch mjType {
		case def.MjTypeHb:
			err = hb.HandleGameChat(s, mjPlayer, send, abs)
		case def.MjTypeZz:
			err = zz.HandleGameChat(s, mjPlayer, send, abs)
		case def.MjTypeTdh:
			err = tdh.HandleGameChat(s, mjPlayer, send, abs)
		case def.MjTypeKf:
			err = kf.HandleGameChat(s, mjPlayer, send, abs)
		}
	// 麻将牌局相关
	case int32(pf.PrepareGame):
		send := &pf.PrepareGameSend{}
		err = send.Unmarshal(msgBody)
		if err != nil {
			return
		}

		util.LogSend(msgID, playerID, roomID, send, "PrepareGame")
		switch mjType {
		case def.MjTypeHb:
			err = hb.HandlePrepareGame(s, mjPlayer, abs)
		case def.MjTypeZz:
			err = zz.HandlePrepareGame(s, mjPlayer, abs)
		case def.MjTypeTdh:
			err = tdh.HandlePrepareGame(s, mjPlayer, abs)
		case def.MjTypeKf:
			err = kf.HandlePrepareGame(s, mjPlayer, abs)
		}
	case int32(pf.CancelPrepare):
		send := &pf.CancelPrepareSend{}
		err = send.Unmarshal(msgBody)
		if err != nil {
			return
		}

		util.LogSend(msgID, playerID, roomID, send, "CancelPrepare")
		switch mjType {
		case def.MjTypeHb:
			err = hb.HandleCancelPrepare(s, mjPlayer, abs)
		case def.MjTypeZz:
			err = zz.HandleCancelPrepare(s, mjPlayer, abs)
		case def.MjTypeTdh:
			err = tdh.HandleCancelPrepare(s, mjPlayer, abs)
		case def.MjTypeKf:
			err = kf.HandleCancelPrepare(s, mjPlayer, abs)
		}
	case int32(pf.GameRecord):
		send := &pf.GameRecordSend{}
		err = send.Unmarshal(msgBody)
		if err != nil {
			return
		}

		util.LogSend(msgID, playerID, roomID, send, "GameRecord")
		err = handleGameRecord(s, mjPlayer, abs)
	case int32(pf.BuyPao):
		send := &pf.BuyPaoSend{}
		err = send.Unmarshal(msgBody)
		if err != nil {
			return
		}

		util.LogSend(msgID, playerID, roomID, send, "BuyPao")
		switch mjType {
		case def.MjTypeHb:
			// 河北麻将不带买跑
		case def.MjTypeZz:
			err = zz.HandleBuyPao(s, mjPlayer, send, abs)
		case def.MjTypeTdh:
			err = tdh.HandleBuyPao(s, mjPlayer, send, abs)
		case def.MjTypeKf:
			err = kf.HandleBuyPao(s, mjPlayer, send, abs)
		}
	// 麻将逻辑相关
	case int32(pf.Discard):
		send := &pf.DiscardSend{}
		err = send.Unmarshal(msgBody)
		if err != nil {
			return
		}

		util.LogSend(msgID, playerID, roomID, send, "Discard")
		switch mjType {
		case def.MjTypeHb:
			err = hb.HandleDiscard(s, mjPlayer, send, abs)
		case def.MjTypeZz:
			err = zz.HandleDiscard(s, mjPlayer, send, abs)
		case def.MjTypeTdh:
			err = tdh.HandleDiscard(s, mjPlayer, send, abs)
		case def.MjTypeKf:
			err = kf.HandleDiscard(s, mjPlayer, send, abs)
		}
	case int32(pf.Operation):
		send := &pf.OperationSend{}
		err = send.Unmarshal(msgBody)
		if err != nil {
			return
		}

		util.LogSend(msgID, playerID, roomID, send, "Operation")
		switch mjType {
		case def.MjTypeHb:
			err = hb.HandleOperation(s, mjPlayer, send, abs)
		case def.MjTypeZz:
			err = zz.HandleOperation(s, mjPlayer, send, abs)
		case def.MjTypeTdh:
			err = tdh.HandleOperation(s, mjPlayer, send, abs)
		case def.MjTypeKf:
			err = kf.HandleOperation(s, mjPlayer, send, abs)
		}
	default:
		err = def.ErrHandleLogic
	}
	return
}
