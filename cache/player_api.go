package cache

import (
	"qnmahjong/db"
	"qnmahjong/db/dao"
	"qnmahjong/def"
	"qnmahjong/pf"
	"qnmahjong/util"
	"strings"
	"time"

	"github.com/olahol/melody"
)

func HandleDirty(playerID int32, s *melody.Session, c *util.Claims, recv *pf.DirtyRecv) {
	playerCache.Lock()
	defer playerCache.Unlock()

	// 检查玩家
	daoPlayer, err := dao.PlayerByPlayerID(db.Pool, playerID)
	if err != nil {
		recv.Status = def.StatusErrorLogic
		return
	}

	player, ok := playerCache.mmap[playerID]
	if !ok {
		player = &Player{}
		playerCache.mmap[playerID] = player
	}
	playerCache.mmap[playerID].Player = daoPlayer

	// 刷新玩家信息
	player.Session = s
	player.Channel = c.Channel
	player.Version = c.Version
	player.LoginType = c.LoginType
	player.IP = util.RealIP(player.Request)

	pfPlayer := &pf.Player{
		Id:       player.PlayerID,
		Nickname: player.Nickname,
		Avatar:   player.Headimgurl,
		Gender:   player.Sex,
		Coins:    player.Coins,
		Cards:    player.Cards,
		Ip:       player.IP,
	}
	mjTypes := GetMjTypes(c.Channel, c.Version)
	jushus := GetJushus(mjTypes)
	module := getModules(c.Channel, c.Version)

	// 获取返回值
	recv.RoomID = player.RoomID
	recv.HighID = player.HighID
	recv.Player = pfPlayer
	recv.Jushus = jushus
	recv.WchatID = def.WeixinWchatID
	recv.Module = module
	recv.RefreshToken = player.RefreshToken
	return
}

func HandleInviteCode(playerID int32, send *pf.InviteCodeSend, recv *pf.InviteCodeRecv) {
	playerCache.Lock()
	defer playerCache.Unlock()

	highID := send.Code

	// 检查是否自己
	if playerID == highID {
		recv.Status = def.StatusErrorInviteCode
		return
	}

	// 快速登录阻止
	if highID == def.QuickLoginPlayerID {
		recv.Status = def.StatusErrorInviteCode
		return
	}

	// 检查玩家、检查上线玩家
	player, ok := playerCache.mmap[playerID]
	highPlayer, highOK := playerCache.mmap[highID]
	if !ok || !highOK {
		recv.Status = def.StatusErrorInviteCode
		return
	}

	// 上线代理限制
	// if highPlayer.HighID == 0 {
	// 	recv.Status = def.StatusErrorInviteCode
	// 	return
	// }

	// 检查上线id
	if player.HighID != 0 {
		recv.Status = def.StatusErrorInviteFailed
		return
	}

	// 绑定上线、添加下线、添加奖励
	player.HighID = highID
	player.Coins += def.InviteCodeBindAward
	highPlayer.LowID = append(highPlayer.LowID, playerID)
	err := player.Update(db.Pool)
	if err != nil {
		util.LogError(err, "player", player, playerID, def.ErrUpdatePlayer)
		recv.Status = def.StatusErrorInviteFailed
		return
	}

	// 奖励记录
	treasure := &dao.Treasure{
		PlayerID:   player.PlayerID,
		Reason:     def.InviteCodeBind,
		Coins:      def.InviteCodeBindAward,
		Cards:      0,
		ChangeTime: time.Now(),
	}
	err = treasure.Insert(db.Pool)
	if err != nil {
		util.LogError(err, "treasure", treasure, playerID, def.ErrInsertTreasure)
	}

	// 获取返回值
	recv.HighID = player.HighID
	recv.Coins = player.Coins
	return
}

func HandleInviteList(playerID int32, recv *pf.InviteListRecv) {
	playerCache.RLock()
	defer playerCache.RUnlock()

	// 检查玩家
	player, ok := playerCache.mmap[playerID]
	if !ok {
		return
	}

	// 获取下线玩家列表
	inviteList := make([]*pf.InvitePlayer, 0, len(player.LowID))
	for _, lowID := range player.LowID {
		lowPlayer, lowOK := playerCache.mmap[lowID]
		if !lowOK {
			continue
		}

		inviteList = append(inviteList, &pf.InvitePlayer{
			Id:       lowPlayer.PlayerID,
			Nickname: lowPlayer.Nickname,
			Avatar:   strings.TrimRight(lowPlayer.Headimgurl, "0") + "132",
			State:    lowPlayer.InviteAward,
		})
	}

	// 获取返回值
	recv.Players = inviteList
	return
}

func HandleInviteAward(id int32, send *pf.InviteAwardSend, recv *pf.InviteAwardRecv) {
	playerCache.Lock()
	defer playerCache.Unlock()

	lowID := send.GetId()

	// 检查玩家、检查下线玩家
	player, ok := playerCache.mmap[id]
	lowPlayer, lowOK := playerCache.mmap[lowID]
	if !ok || !lowOK {
		recv.Status = def.StatusErrorInviteInfo
		return
	}

	// 判断是否为下线
	exist := false
	for _, tempID := range player.LowID {
		if tempID == lowID {
			exist = true
			break
		}
	}
	if !exist {
		recv.Status = def.StatusErrorInviteInfo
		return
	}

	// 检查下线玩家状态
	if lowPlayer.InviteAward != def.InviteAwardAvailable {
		recv.Status = def.StatusErrorInviteInfo
		return
	}

	// 更新下线状态
	lowPlayer.InviteAward = def.InviteAwardUnavailable
	err := lowPlayer.Update(db.Pool)
	if err != nil {
		util.LogError(err, "player", lowPlayer, lowPlayer.PlayerID, def.ErrUpdatePlayer)
		recv.Status = def.StatusErrorInviteAward
		return
	}

	// 添加奖励
	player.Coins += def.InviteCompleteAward
	err = player.Update(db.Pool)
	if err != nil {
		util.LogError(err, "player", player, id, def.ErrUpdatePlayer)
		recv.Status = def.StatusErrorInviteFailed
		return
	}

	// 奖励记录
	treasure := &dao.Treasure{
		PlayerID:   player.PlayerID,
		Reason:     def.InviteComplete,
		Coins:      def.InviteCompleteAward,
		Cards:      0,
		ChangeTime: time.Now(),
	}
	err = treasure.Insert(db.Pool)
	if err != nil {
		util.LogError(err, "treasure", treasure, id, def.ErrInsertTreasure)
	}

	// 获取返回值
	recv.Coins = player.Coins
	recv.Id = lowID
	return
}

func HandleOrderApply(orderID string) {
	playerCache.Lock()
	defer playerCache.Unlock()

	order, err := dao.OrderByOrderID(db.Pool, orderID)
	if err != nil {
		return
	}

	if order.Status != def.IpayStatusSuccess {
		return
	}

	playerID := order.PlayerID
	player, ok := playerCache.mmap[playerID]
	if !ok {
		return
	}

	player.Cards += order.GoodsCount + order.ExtraCount
	// 是否首充
	firstBuyAward := int32(0)
	if player.FirstBuy == 0 {
		player.FirstBuy = 1
		firstBuyAward = 1
	}
	err = player.Update(db.Pool)
	if err != nil {
		util.LogError(err, "player", player, playerID, def.ErrUpdatePlayer)
		return
	}

	t := &dao.Treasure{
		PlayerID:   playerID,
		Reason:     def.IPayOrder,
		Coins:      0,
		Cards:      order.GoodsCount + order.ExtraCount,
		ChangeTime: time.Now(),
	}
	err = t.Insert(db.Pool)
	if err != nil {
		util.LogError(err, "treasure", t, playerID, def.ErrInsertTreasure)
	}

	order.Status = def.IpayStatusComplete
	order.ReviseTime = time.Now()
	err = order.Update(db.Pool)
	if err != nil {
		util.LogError(err, "order", order, playerID, def.ErrUpdateOrder)
		return
	}

	// 往分销系统写记录
	agPay := dao.AgPay{
		AgID:          player.HighID,
		CustomerID:    player.PlayerID,
		DiamondCnt:    order.GoodsCount + order.ExtraCount,
		MoneyCnt:      order.Price,
		Delflag:       0,
		FirstBuyAward: firstBuyAward,
		CreateTime:    time.Now(),
	}
	err = agPay.Insert(db.Pool)
	if err != nil {
		util.LogError(err, "agPay", agPay, playerID, def.ErrInsertAgPay)
	}

	if player.Session == nil {
		return
	}

	roomID := player.RoomID
	msgID := int32(pf.OrderApply)
	s := player.Session
	abs := &pf.AbsMessage{
		MsgID: msgID,
	}

	var bytes []byte
	recv := &pf.OrderApplyRecv{
		Status:    def.StatusOK,
		Coins:     player.Coins,
		Cards:     player.Cards,
		OrderID:   order.OrderID,
		Price:     order.Price / 100,
		PayType:   def.PayTypeIPay,
		Count:     order.GoodsCount + order.ExtraCount,
		ProductID: order.WaresID,
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
		util.LogRecv(msgID, playerID, roomID, recv, "OrderApply")
	}
}

func HandleSubCard(playerID, subCards, subCoins int32) {
	playerCache.Lock()
	defer playerCache.Unlock()

	player, ok := playerCache.mmap[playerID]
	if !ok {
		return
	}

	player.Coins -= subCoins
	player.Cards -= subCards
	err := player.Update(db.Pool)
	if err != nil {
		util.LogError(err, "player", player, playerID, def.ErrUpdatePlayer)
		return
	}

	if player.Session == nil {
		return
	}

	roomID := player.RoomID
	msgID := int32(pf.ResourceChange)
	s := player.Session
	abs := &pf.AbsMessage{
		MsgID: msgID,
	}

	var bytes []byte
	recv := &pf.ResourceChangeRecv{
		Status: def.StatusOK,
		Coins:  player.Coins,
		Cards:  player.Cards,
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
		util.LogRecv(msgID, playerID, roomID, recv, "ResourceChange")
	}
}
