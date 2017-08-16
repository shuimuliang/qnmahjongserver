package msg

import (
	"qnmahjong/cache"
	"qnmahjong/db"
	"qnmahjong/db/dao"
	"qnmahjong/def"
	"qnmahjong/pay"
	"qnmahjong/pf"
	"qnmahjong/redis"
	"qnmahjong/util"
	"time"

	"github.com/olahol/melody"
)

func handleLogout(s *melody.Session, mjPlayer *cache.MjPlayer, abs *pf.AbsMessage) (err error) {
	playerID := mjPlayer.PlayerID
	roomID := mjPlayer.RoomID
	msgID := abs.MsgID

	var bytes []byte
	recv := &pf.LogoutRecv{
		Status: def.StatusOK,
	}

	bytes, err = recv.Marshal()
	if err != nil {
		return
	}

	abs.MsgBody = bytes
	bytes, err = abs.Marshal()
	if err != nil {
		return
	}

	err = s.WriteBinary(bytes)
	if err == nil {
		util.LogRecv(msgID, playerID, roomID, recv, "Logout")
	}
	return
}

func handleAuth(s *melody.Session, c *util.Claims, abs *pf.AbsMessage) (err error) {
	playerID := c.PlayerID
	roomID := int32(0)
	msgID := abs.MsgID

	var bytes []byte
	recv := &pf.AuthRecv{
		Status: def.StatusOK,
	}

	bytes, err = recv.Marshal()
	if err != nil {
		return
	}

	abs.MsgBody = bytes
	bytes, err = abs.Marshal()
	if err != nil {
		return
	}

	err = s.WriteBinary(bytes)
	if err == nil {
		util.LogRecv(msgID, playerID, roomID, recv, "Auth")
	}
	return
}

func handleDirty(s *melody.Session, c *util.Claims, abs *pf.AbsMessage) (err error) {
	playerID := c.PlayerID
	roomID := int32(0)
	msgID := abs.MsgID

	var bytes []byte
	recv := &pf.DirtyRecv{
		Status: def.StatusOK,
	}

	cache.HandleDirty(playerID, s, c, recv)

	bytes, err = recv.Marshal()
	if err != nil {
		return
	}

	abs.MsgBody = bytes
	bytes, err = abs.Marshal()
	if err != nil {
		return
	}

	err = s.WriteBinary(bytes)
	if err == nil {
		util.LogRecv(msgID, playerID, roomID, recv, "Dirty")
	}
	return
}

func handleBroadcast(s *melody.Session, c *util.Claims, abs *pf.AbsMessage) (err error) {
	playerID := c.PlayerID
	roomID := int32(0)
	msgID := abs.MsgID

	var bytes []byte
	recv := &pf.BroadcastRecv{
		Status: def.StatusOK,
		Text:   []string{"千胜游戏仅供娱乐，严禁任何形式的赌博行为！祝大家游戏愉快。", "欢迎来到千胜河南麻将，代理咨询微信号：poker888mm。"},
	}
	bytes, err = recv.Marshal()
	if err != nil {
		return
	}

	abs.MsgID = int32(pf.Broadcast)
	abs.MsgBody = bytes
	bytes, err = abs.Marshal()
	if err != nil {
		return
	}

	err = s.WriteBinary(bytes)
	if err == nil {
		util.LogRecv(msgID, playerID, roomID, recv, "Broadcast")
	}
	return
}

func handleFeedback(s *melody.Session, mjPlayer *cache.MjPlayer, send *pf.FeedbackSend, abs *pf.AbsMessage) (err error) {
	playerID := mjPlayer.PlayerID
	roomID := mjPlayer.RoomID
	msgID := abs.MsgID

	curTime := time.Now()
	feedback := &dao.Feedback{
		PlayerID:   playerID,
		Channel:    mjPlayer.Channel,
		Version:    mjPlayer.Version,
		ImgURL:     send.ImgUrl,
		Text:       send.Text,
		Status:     0,
		AddTime:    curTime,
		ReviseTime: curTime,
	}
	err = feedback.Insert(db.Pool)
	if err != nil {
		util.LogError(err, "feedback", feedback, playerID, def.ErrInsertFeedback)
	}

	var bytes []byte
	recv := &pf.FeedbackRecv{
		Status: def.StatusOK,
	}

	bytes, err = recv.Marshal()
	if err != nil {
		return
	}

	abs.MsgBody = bytes
	bytes, err = abs.Marshal()
	if err != nil {
		return
	}

	err = s.WriteBinary(bytes)
	if err == nil {
		util.LogRecv(msgID, playerID, roomID, recv, "Feedback")
	}
	return
}

func handleInviteCode(s *melody.Session, mjPlayer *cache.MjPlayer, send *pf.InviteCodeSend, abs *pf.AbsMessage) (err error) {
	playerID := mjPlayer.PlayerID
	roomID := mjPlayer.RoomID
	msgID := abs.MsgID

	var bytes []byte
	recv := &pf.InviteCodeRecv{
		Status: def.StatusOK,
	}

	cache.HandleInviteCode(playerID, send, recv)

	bytes, err = recv.Marshal()
	if err != nil {
		return
	}

	abs.MsgBody = bytes
	bytes, err = abs.Marshal()
	if err != nil {
		return
	}

	err = s.WriteBinary(bytes)
	if err == nil {
		util.LogRecv(msgID, playerID, roomID, recv, "InviteCode")
	}
	return
}

func handleInviteList(s *melody.Session, mjPlayer *cache.MjPlayer, abs *pf.AbsMessage) (err error) {
	playerID := mjPlayer.PlayerID
	roomID := mjPlayer.RoomID
	msgID := abs.MsgID

	var bytes []byte
	recv := &pf.InviteListRecv{
		Status: def.StatusOK,
	}

	cache.HandleInviteList(playerID, recv)

	bytes, err = recv.Marshal()
	if err != nil {
		return
	}

	abs.MsgBody = bytes
	bytes, err = abs.Marshal()
	if err != nil {
		return
	}

	err = s.WriteBinary(bytes)
	if err == nil {
		util.LogRecv(msgID, playerID, roomID, recv, "InviteList")
	}
	return
}

func handleInviteAward(s *melody.Session, mjPlayer *cache.MjPlayer, send *pf.InviteAwardSend, abs *pf.AbsMessage) (err error) {
	playerID := mjPlayer.PlayerID
	roomID := mjPlayer.RoomID
	msgID := abs.MsgID

	var bytes []byte
	recv := &pf.InviteAwardRecv{
		Status: def.StatusOK,
	}

	cache.HandleInviteAward(playerID, send, recv)

	bytes, err = recv.Marshal()
	if err != nil {
		return
	}

	abs.MsgBody = bytes
	bytes, err = abs.Marshal()
	if err != nil {
		return
	}

	err = s.WriteBinary(bytes)
	if err == nil {
		util.LogRecv(msgID, playerID, roomID, recv, "InviteAward")
	}
	return
}

func handleGoodsList(s *melody.Session, mjPlayer *cache.MjPlayer, abs *pf.AbsMessage) (err error) {
	playerID := mjPlayer.PlayerID
	roomID := mjPlayer.RoomID
	msgID := abs.MsgID

	goods := cache.GetGoods(mjPlayer.Channel)

	var bytes []byte
	recv := &pf.GoodsListRecv{
		Status: def.StatusOK,
		Goods:  goods,
	}

	bytes, err = recv.Marshal()
	if err != nil {
		return
	}

	abs.MsgBody = bytes
	bytes, err = abs.Marshal()
	if err != nil {
		return
	}

	err = s.WriteBinary(bytes)
	if err == nil {
		util.LogRecv(msgID, playerID, roomID, recv, "GoodsList")
	}
	return
}

func handleOrderApplyNew(s *melody.Session, mjPlayer *cache.MjPlayer, send *pf.OrderApplyNewSend, abs *pf.AbsMessage) (err error) {
	playerID := mjPlayer.PlayerID
	roomID := mjPlayer.RoomID
	msgID := abs.MsgID

	orderList := make([]*pf.OrderStatus, 0)
	for _, order := range send.OrderList {
		switch order.PayType {
		case def.PayTypeAppStore:
			orderStatus := pay.ValidateAppStore(mjPlayer, order)
			if orderStatus != nil {
				orderList = append(orderList, orderStatus)
			}
		}
	}

	pfPlayer := cache.GetPfPlayer(playerID)
	if pfPlayer == nil {
		return
	}

	var bytes []byte
	recv := &pf.OrderApplyNewRecv{
		Status:    def.StatusOK,
		Coins:     pfPlayer.Coins,
		Cards:     pfPlayer.Cards,
		OrderList: orderList,
	}

	bytes, err = recv.Marshal()
	if err != nil {
		return
	}

	abs.MsgBody = bytes
	bytes, err = abs.Marshal()
	if err != nil {
		return
	}

	err = s.WriteBinary(bytes)
	if err == nil {
		util.LogRecv(msgID, playerID, roomID, recv, "OrderApplyNew")
	}
	return
}

func handleGameRecord(s *melody.Session, mjPlayer *cache.MjPlayer, abs *pf.AbsMessage) (err error) {
	playerID := mjPlayer.PlayerID
	roomID := mjPlayer.RoomID
	msgID := abs.MsgID

	records := redis.GetRecords(playerID)

	var bytes []byte
	recv := &pf.GameRecordRecv{
		Status:  def.StatusOK,
		Records: records,
	}

	bytes, err = recv.Marshal()
	if err != nil {
		return
	}

	abs.MsgBody = bytes
	bytes, err = abs.Marshal()
	if err != nil {
		return
	}

	err = s.WriteBinary(bytes)
	if err == nil {
		util.LogRecv(msgID, playerID, roomID, recv, "GameRecord")
	}
	return
}
