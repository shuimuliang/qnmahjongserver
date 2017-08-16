package hb

import (
	"qnmahjong/cache"
	"qnmahjong/db"
	"qnmahjong/db/dao"
	"qnmahjong/def"
	"qnmahjong/pf"
	"qnmahjong/util"
	"math/rand"
	"time"
)

func (m RoomMap) handleCreateRoom(
	id int32,
	send *pf.EnterRoomSend,
	recv *pf.EnterRoomRecv) {

	// 检查是否是进入房间输入了0
	if send.MjType == 0 {
		recv.Status = def.StatusErrorCreateRoom
		return
	}

	// 检查配置和房卡金币
	config, ok := checkConfig(id, send)
	if !ok {
		recv.Status = def.StatusErrorCreateRoom
		return
	}

	m.Lock()
	defer m.Unlock()

	cliConfig := send.GetConfigs()
	mjType := send.GetMjType()
	cheat := send.GetCheat()
	latitude := send.GetLatitude()
	longitude := send.GetLongitude()
	cache.SetLocation(id, latitude, longitude)

	// 判断是否拿到gps位置
	if cheat && int(latitude) == 0 && int(longitude) == 0 {
		recv.Status = def.StatusErrorGPSNotOpen
		return
	}

	// 生成房间id
	var roomID int32
	for {
		roomID = rand.Int31n(9000) + 1000 + 10000
		if _, ok := m.mmap[roomID]; !ok {
			break
		}
	}

	// 初始化玩家列表
	roomSeats := make([]*MJSeat, def.RoomPlayerCount)
	roomSeats[0] = &MJSeat{
		ID:  id,
		Pos: 0,
	}

	// 初始化麻将信息
	var roomTiles *MJTiles
	if config.Daifengpai {
		roomTiles = &MJTiles{
			Tiles: TileWithFeng,
		}
	} else {
		roomTiles = &MJTiles{
			Tiles: TileWithoutFeng,
		}
	}

	// 初始化胡牌信息
	var roomHupai = &MJHupai{}

	// 初始化状态信息
	var roomStatus = &MJStatus{
		GameStatus: def.GameStatusNotStart,
	}

	// 初始化配置信息
	var roomConfig = &MJConfig{
		MJType:     mjType,
		CreateTime: time.Now(),
		Config:     config,
		CliConfig:  cliConfig,
		Cheat:      cheat,
	}

	// 初始化房间信息
	room := &MJRoom{
		RoomID:     roomID,
		RoomSeats:  roomSeats,
		RoomTiles:  roomTiles,
		RoomHupai:  roomHupai,
		RoomStatus: roomStatus,
		RoomConfig: roomConfig,
	}
	m.mmap[roomID] = room
	cache.SetRoomID(id, roomID, mjType)

	// 获取玩家列表
	var playerList []*pf.RoomPlayer
	for _, seat := range room.RoomSeats {
		if seat == nil {
			continue
		}
		player := cache.GetPfPlayer(seat.ID)
		lat, long := cache.GetGPSLocation(seat.ID)
		if player != nil {
			playerList = append(playerList, &pf.RoomPlayer{
				Player:    player,
				Score:     seat.Score,
				Pos:       seat.Pos + 1,
				Latitude:  lat,
				Longitude: long,
			})
		}
	}

	// 获取返回值
	recv.GameStatus = room.RoomStatus.GameStatus
	recv.RoomID = room.RoomID
	recv.TotalRound = room.RoomConfig.Config.Jushu
	recv.CurRound = room.RoomStatus.CurRound
	recv.LeftCards = room.RoomTiles.LeftCnt
	recv.BankerPos = room.RoomStatus.BankerPos + 1
	recv.CurPos = room.RoomStatus.CurPos + 1
	recv.Configs = room.RoomConfig.CliConfig
	recv.PlayerList = playerList
	recv.Cheat = cheat
	return
}

func (m RoomMap) handleEnterRoom(
	id int32,
	send *pf.EnterRoomSend,
	recv *pf.EnterRoomRecv,
	otherJoinRecv *pf.OtherJoinRecv,
	otherIDs *[]int32) {

	m.Lock()
	defer m.Unlock()

	roomID := send.GetRoomID()
	latitude := send.GetLatitude()
	longitude := send.GetLongitude()
	cache.SetLocation(id, latitude, longitude)

	// 判断房间是否存在
	room, ok := RoomCache.mmap[roomID]
	if !ok {
		recv.Status = def.StatusErrorNoRoomID
		return
	}

	// 判断是否防作弊
	if room.RoomConfig.Cheat {
		// 判断是否拿到gps位置
		if int(latitude) == 0 && int(longitude) == 0 {
			recv.Status = def.StatusErrorGPSNotOpen
			return
		}
		for _, seat := range room.RoomSeats {
			if seat != nil && cache.IsIPConflict(id, seat.ID) {
				recv.Status = def.StatusErrorIPConflict
				return
			}

			if seat != nil && cache.IsGeoConflict(id, seat.ID) {
				recv.Status = def.StatusErrorGeoConflict
				return
			}
		}
	}

	// 判断是否能进入房间
	var success bool
	for pos, seat := range room.RoomSeats {
		if seat == nil {
			room.RoomSeats[pos] = &MJSeat{
				ID:  id,
				Pos: int32(pos),
			}
			cache.SetRoomID(id, roomID, def.MjTypeHb)
			success = true
			break
		}
	}
	if !success {
		recv.Status = def.StatusErrorRoomIsFull
		return
	}

	// 获取玩家列表
	var playerList []*pf.RoomPlayer
	var roomPlayer *pf.RoomPlayer
	for _, seat := range room.RoomSeats {
		if seat == nil {
			continue
		}
		player := cache.GetPfPlayer(seat.ID)
		if player != nil {
			lat, long := cache.GetGPSLocation(seat.ID)
			playerList = append(playerList, &pf.RoomPlayer{
				Player:    player,
				Score:     seat.Score,
				Pos:       seat.Pos + 1,
				Latitude:  lat,
				Longitude: long,
			})
			if seat.ID == id {
				roomPlayer = &pf.RoomPlayer{
					Player:    player,
					Score:     seat.Score,
					Pos:       seat.Pos + 1,
					Latitude:  lat,
					Longitude: long,
				}
			} else {
				*otherIDs = append(*otherIDs, seat.ID)
			}
		}
	}

	// 获取返回值
	recv.GameStatus = room.RoomStatus.GameStatus
	recv.RoomID = room.RoomID
	recv.TotalRound = room.RoomConfig.Config.Jushu
	recv.CurRound = room.RoomStatus.CurRound
	recv.LeftCards = room.RoomTiles.LeftCnt
	recv.BankerPos = room.RoomStatus.BankerPos + 1
	recv.CurPos = room.RoomStatus.CurPos + 1
	recv.Configs = room.RoomConfig.CliConfig
	recv.PlayerList = playerList
	recv.Cheat = room.RoomConfig.Cheat
	otherJoinRecv.Status = def.StatusOK
	otherJoinRecv.Player = roomPlayer
	return
}

func (m RoomMap) handleGameRelink(
	id int32,
	send *pf.EnterRoomSend,
	recv *pf.EnterRoomRecv,
	voteCloseRecv *pf.VoteCloseRecv,
	operationNoticeRecv *pf.OperationNoticeRecv,
	discardNoticerecv *pf.DiscardNoticeRecv,
	endRoundRecv *pf.EndRoundRecv) {

	m.Lock()
	defer m.Unlock()

	roomID := send.GetRoomID()

	room, ok := m.mmap[roomID]
	if !ok {
		recv.Status = def.StatusNoRoom
		return
	}

	pos, ok := room.InRoom(id, roomID, nil)
	if !ok {
		recv.Status = def.StatusNotInRoom
		return
	}

	// 获取玩家列表
	// 获取玩家麻将
	var playerList []*pf.RoomPlayer
	var cardsList []*pf.CardsInfo
	for _, seat := range room.RoomSeats {
		if seat == nil {
			continue
		}
		player := cache.GetPfPlayer(seat.ID)
		if player != nil {
			lat, long := cache.GetGPSLocation(seat.ID)
			playerList = append(playerList, &pf.RoomPlayer{
				Player:    player,
				Score:     seat.Score,
				Pos:       seat.Pos + 1,
				Latitude:  lat,
				Longitude: long,
			})
			cardsList = append(cardsList, seat.GetCardsInfo())
		}
	}

	// 获取返回值
	recv.GameStatus = room.RoomStatus.GameStatus
	recv.RoomID = room.RoomID
	recv.TotalRound = room.RoomConfig.Config.Jushu
	recv.CurRound = room.RoomStatus.CurRound
	recv.LeftCards = room.RoomTiles.LeftCnt
	recv.BankerPos = room.RoomStatus.BankerPos + 1
	recv.CurPos = room.RoomStatus.CurPos + 1
	recv.LastPos = room.RoomStatus.LastPos + 1
	recv.Configs = room.RoomConfig.CliConfig
	recv.PlayerList = playerList
	recv.CardsInfoList = cardsList
	recv.Cheat = room.RoomConfig.Cheat
	if room.RoomConfig.HunTile != nil {
		recv.HunCard = room.RoomConfig.HunTile.Front
	}

	if room.RoomStatus.IsVoting {
		var voteList []*pf.VoteInfo
		for _, seat := range room.RoomSeats {
			voteList = append(voteList, &pf.VoteInfo{
				Pos:    seat.Pos + 1,
				Action: seat.VoteStatus,
			})
		}
		voteCloseRecv.Status = def.StatusOK
		voteCloseRecv.LeftTime = room.RoomStatus.VoteTime
		voteCloseRecv.FirstPos = room.RoomStatus.VotePos + 1
		voteCloseRecv.VoteList = voteList
	}

	if room.RoomStatus.GameStatus == def.GameStatusIsGaming {
		if room.RoomHupai.IsMultiDianPao {

			// pos
			isJiePaoPos := false
			index := 0
			for i, jiePaoPos := range room.RoomHupai.MultiJiePaoPos {
				if jiePaoPos == pos {
					index = i
					isJiePaoPos = true
					break
				}
			}

			if isJiePaoPos {
				isHu := false
				for _, huPos := range room.RoomHupai.MultiJiePaoHupaiPos {
					if huPos == pos {
						isHu = true
						break
					}
				}

				isPass := false
				for _, passPos := range room.RoomHupai.MultiJiePaoPassPos {
					if passPos == pos {
						isPass = true
						break
					}
				}

				if !isHu && !isPass {
					operationNoticeRecv.Status = room.RoomStatus.operationNoticeRecv.GetStatus()
					operationNoticeRecv.OperationList = room.RoomStatus.operationNoticeRecv.GetOperationList()[2*index : 2*index+2]
				} else {
					operationNoticeRecv.Status = room.RoomStatus.operationNoticeRecv.GetStatus()
					operationNoticeRecv.OperationList = []*pf.OperationItem{
						{Type: def.Wating},
					}
				}
			} else {
				operationNoticeRecv.Status = room.RoomStatus.operationNoticeRecv.GetStatus()
				operationNoticeRecv.OperationList = []*pf.OperationItem{
					{Type: def.Wating},
				}
			}
		} else {
			if room.RoomStatus.CurPos == pos {
				if room.RoomStatus.operationNoticeRecv != nil {
					operationNoticeRecv.Status = room.RoomStatus.operationNoticeRecv.GetStatus()
					operationNoticeRecv.OperationList = room.RoomStatus.operationNoticeRecv.GetOperationList()
				}
				if room.RoomStatus.discardNoticerecv != nil {
					discardNoticerecv.Status = room.RoomStatus.discardNoticerecv.GetStatus()
					discardNoticerecv.Pos = room.RoomStatus.discardNoticerecv.GetPos()
				}
			}
		}
	}

	if room.RoomStatus.GameStatus == def.GameStatusRoundOver {
		if room.RoomStatus.endRoundRecv != nil {
			endRoundRecv.Status = room.RoomStatus.endRoundRecv.GetStatus()
			endRoundRecv.ItemList = room.RoomStatus.endRoundRecv.GetItemList()
		}
	}
	return
}

func (m RoomMap) handleExitRoom(
	id, roomID int32,
	recv *pf.ExitRoomRecv,
	otherIDs *[]int32) {

	m.RLock()
	defer m.RUnlock()

	room, ok := RoomCache.mmap[roomID]
	if !ok {
		recv.Status = def.StatusNoRoom
		return
	}

	pos, ok := room.InRoom(id, roomID, otherIDs)
	if !ok {
		recv.Status = def.StatusNotInRoom
		return
	}

	if room.RoomStatus.IsStart {
		recv.Status = def.StatusIsStart
		return
	}

	room.RoomSeats[pos] = nil
	cache.SetRoomID(id, 0, 0)
	recv.Pos = pos + 1
	return
}

func (m RoomMap) handleCloseRoom(
	id, roomID int32,
	recv *pf.CloseRoomRecv,
	otherIDs *[]int32) {

	m.Lock()
	defer m.Unlock()

	room, ok := RoomCache.mmap[roomID]
	if !ok {
		recv.Status = def.StatusNoRoom
		return
	}

	pos, ok := room.InRoom(id, roomID, otherIDs)
	if !ok {
		recv.Status = def.StatusNotInRoom
		return
	}

	if room.RoomStatus.IsStart {
		recv.Status = def.StatusIsStart
		return
	}

	if pos != 0 {
		recv.Status = def.StatusNotOwner
		return
	}

	for _, seat := range room.RoomSeats {
		if seat == nil {
			continue
		}
		cache.SetRoomID(seat.ID, 0, 0)
	}

	recv.CloseType = def.CloseRoomBeforStart
	delete(m.mmap, roomID)
	return
}

func (m RoomMap) handleVoteClose(
	id, roomID int32,
	send *pf.VoteCloseSend,
	recv *pf.VoteCloseRecv,
	closeRoomRecv *pf.CloseRoomRecv,
	otherIDs *[]int32) {

	m.Lock()
	defer m.Unlock()

	action := send.GetAction()

	room, ok := RoomCache.mmap[roomID]
	if !ok {
		recv.Status = def.StatusNoRoom
		return
	}

	pos, ok := room.InRoom(id, roomID, otherIDs)
	if !ok {
		recv.Status = def.StatusNotInRoom
		return
	}

	if !room.RoomStatus.IsStart {
		recv.Status = def.StatusNoStart
		return
	}

	if !room.RoomStatus.IsVoting {
		room.RoomStatus.IsVoting = true
		room.RoomStatus.VotePos = pos
		room.RoomStatus.VoteTime = def.VoteLeftTime
		go RoomCache.handleVoteCountdown(roomID)
	}
	room.RoomSeats[pos].VoteStatus = action

	var agreeCnt int32
	var voteList []*pf.VoteInfo
	for _, seat := range room.RoomSeats {
		voteList = append(voteList, &pf.VoteInfo{
			Pos:    seat.Pos + 1,
			Action: seat.VoteStatus,
		})
		if seat.VoteStatus == def.VoteAgree {
			agreeCnt++
		}
	}
	recv.LeftTime = room.RoomStatus.VoteTime
	recv.FirstPos = room.RoomStatus.VotePos + 1
	recv.VoteList = voteList

	// 有人拒绝继续游戏
	if action == def.VoteDisagree {
		room.RoomStatus.IsVoting = false
		room.RoomStatus.VotePos = 0
		room.RoomStatus.VoteTime = 0
		for _, seat := range room.RoomSeats {
			seat.VoteStatus = def.VoteDefault
		}
		return
	}

	// 全部同意关闭房间
	if agreeCnt >= def.RoomPlayerCount {
		closeRoomRecv.Status = def.StatusOK
		closeRoomRecv.CloseType = def.CloseRoomVoteAgree
		for _, seat := range room.RoomSeats {
			cache.SetRoomID(seat.ID, 0, 0)
			// 战绩
			room.RoomRecord.ScoreList[seat.Pos] = seat.Score
			room.DaoRecord.CurRound = 0
			switch seat.Pos + 1 {
			case East:
				room.DaoRecord.EastID = seat.ID
				room.DaoRecord.EastScore = seat.Score
			case South:
				room.DaoRecord.SouthID = seat.ID
				room.DaoRecord.SouthScore = seat.Score
			case West:
				room.DaoRecord.WestID = seat.ID
				room.DaoRecord.WestScore = seat.Score
			case North:
				room.DaoRecord.NorthID = seat.ID
				room.DaoRecord.NorthScore = seat.Score
			}
		}

		// 战绩记录
		if room.RoomStatus.IsSubCard {
			err := room.DaoRecord.Insert(db.Pool)
			if err != nil {
				util.LogError(err, "record", room.DaoRecord, 0, def.ErrInsertRecord)
			}
			room.RecordToRedis()
		}
		delete(m.mmap, roomID)
	}
	return
}

func (m RoomMap) handleVoteCountdown(
	roomID int32) {

	defer util.Stack()

	for {
		time.Sleep(time.Second)
		m.Lock()

		room, ok := RoomCache.mmap[roomID]
		if !ok {
			m.Unlock()
			break
		}

		if !room.RoomStatus.IsVoting {
			m.Unlock()
			break
		}

		room.RoomStatus.VoteTime--
		if room.RoomStatus.VoteTime > 0 {
			m.Unlock()
			continue
		}

		// 大结算
		settlementRecv := &pf.SettlementRecv{
			Status:      def.StatusOK,
			SettleList:  room.GetSettleItems(),
			IsVoteClose: 1,
		}
		msgBody, err := settlementRecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.Settlement),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				for _, seat := range room.RoomSeats {
					if cache.SendMessage(seat.ID, msg) == nil {
						util.LogRecv(int32(pf.Settlement), seat.ID, roomID, settlementRecv, "Settlement")
					}
				}
			}
		}

		// 关闭房间
		closeRoomRecv := &pf.CloseRoomRecv{
			Status:    def.StatusOK,
			CloseType: def.CloseRoomVoteAgree,
		}
		msgBody, err = closeRoomRecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.CloseRoom),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				for _, seat := range room.RoomSeats {
					if cache.SendMessage(seat.ID, msg) == nil {
						util.LogRecv(int32(pf.CloseRoom), seat.ID, roomID, closeRoomRecv, "CloseRoom")
					}
				}
			}
		}

		// 清理房间
		if closeRoomRecv.GetStatus() == def.StatusOK {
			for _, seat := range room.RoomSeats {
				cache.SetRoomID(seat.ID, 0, 0)
			}
			delete(m.mmap, roomID)
		}

		m.Unlock()
	}
}

func (m RoomMap) handleGameChat(
	id, roomID int32,
	send *pf.GameChatSend,
	recv *pf.GameChatRecv,
	otherIDs *[]int32) {

	m.RLock()
	defer m.RUnlock()

	types := send.GetTypes()
	messID := send.GetMessID()
	text := send.GetText()

	room, ok := RoomCache.mmap[roomID]
	if !ok {
		recv.Status = def.StatusNoRoom
		return
	}

	_, ok = room.InRoom(id, roomID, otherIDs)
	if !ok {
		recv.Status = def.StatusNotInRoom
		return
	}

	player := cache.GetPfPlayer(id)
	if player == nil {
		recv.Status = def.StatusNoPlayer
		return
	}

	chat := &dao.Chat{
		SendID:   id,
		RoomID:   roomID,
		MsgType:  types,
		MessID:   messID,
		MsgText:  text,
		SendTime: time.Now(),
	}

	// 聊天记录
	err := chat.Insert(db.Pool)
	if err != nil {
		util.LogError(err, "chat", chat, id, def.ErrInsertChat)
	}

	recv.PlayerID = id
	recv.NickName = player.GetNickname()
	recv.Types = types
	recv.MessID = messID
	recv.Text = text
	recv.Avatar = player.GetAvatar()
	return
}

func (m RoomMap) handlePrepareGame(
	id, roomID int32,
	recv *pf.PrepareGameRecv,
	startRoundRecv *pf.StartRoundRecv,
	drawCardRecv *pf.DrawCardRecv,
	operationNoticeRecv *pf.OperationNoticeRecv,
	discardNoticerecv *pf.DiscardNoticeRecv,
	bankerID *int32,
	otherIDs *[]int32,
	otherCards map[int32][]int32) {

	m.Lock()
	defer m.Unlock()

	room, ok := RoomCache.mmap[roomID]
	if !ok {
		recv.Status = def.StatusNoRoom
		return
	}

	pos, ok := room.InRoom(id, roomID, otherIDs)
	if !ok {
		recv.Status = def.StatusNotInRoom
		return
	}

	if room.RoomStatus.IsGaming {
		recv.Status = def.StatusIsGaming
		return
	}

	room.RoomSeats[pos].Prepared = true
	recv.Pos = pos + 1

	var preparedCnt int32
	for _, seat := range room.RoomSeats {
		if seat == nil {
			continue
		}
		if seat.Prepared {
			preparedCnt++
		}
	}

	if preparedCnt < def.RoomPlayerCount {
		return
	}

	// 初始化牌桌
	room.ClearRoom()

	// 开局
	for i := range room.RoomSeats {
		j := (int32(i) + room.RoomStatus.BankerPos) % def.RoomPlayerCount
		seat := room.RoomSeats[j]
		seat.HandTiles = room.RoomTiles.DealTiles()
		startRoundRecv.CardsCount = append(startRoundRecv.CardsCount, int32(len(seat.HandTiles)))
		if seat.ID == id {
			startRoundRecv.MyCards = append(startRoundRecv.MyCards, seat.HandTiles.GetFrontCards()...)
		} else {
			otherCards[seat.ID] = seat.HandTiles.GetFrontCards()
		}
		startRoundRecv.ScoreList = append(startRoundRecv.ScoreList, seat.Score)
	}

	startRoundRecv.Status = def.StatusOK
	startRoundRecv.CurRound = room.RoomStatus.CurRound
	startRoundRecv.LeftCards = room.RoomTiles.LeftCnt
	startRoundRecv.BankerPos = room.RoomStatus.BankerPos + 1
	startRoundRecv.TotalRound = room.RoomConfig.Config.Jushu

	// 抓牌
	bankerPos := room.RoomStatus.BankerPos
	bankerSeat := room.RoomSeats[bankerPos]
	*bankerID = bankerSeat.ID
	bankerSeat.DrawTile = room.RoomTiles.DealOneTile()
	drawCardRecv.Status = def.StatusOK
	drawCardRecv.Pos = bankerPos + 1
	drawCardRecv.Card = bankerSeat.DrawTile.Front

	// 混牌
	if room.RoomConfig.Config.Daihuner {
		preHunTile := room.RoomTiles.DealOneTile()
		hunCard := HunTilesMap[preHunTile.Unicode]
		hunTile := room.RoomTiles.Tiles.GetEndTile(hunCard)
		room.RoomConfig.HunTile = &hunTile
		startRoundRecv.HunCard = hunTile.Front
		startRoundRecv.PreHunCard = preHunTile.Front
	}

	// 操作提示
	operationNoticeRecv.OperationList = bankerSeat.GetSelfHupaiList(room, false)
	if len(operationNoticeRecv.OperationList) > 0 {
		operationNoticeRecv.Status = def.StatusOK
	} else {
		operationNoticeRecv.OperationList = bankerSeat.GetSelfGangList()
		if len(operationNoticeRecv.OperationList) > 0 {
			operationNoticeRecv.Status = def.StatusOK
		}
	}

	// 出牌提示
	if operationNoticeRecv.GetStatus() != def.StatusOK {
		discardNoticerecv.Status = def.StatusOK
		discardNoticerecv.Pos = bankerPos + 1
	}

	// 断线重连数据
	if operationNoticeRecv.GetStatus() == def.StatusOK {
		room.RoomStatus.operationNoticeRecv = operationNoticeRecv
		room.RoomStatus.discardNoticerecv = nil
		room.RoomStatus.endRoundRecv = nil
	} else if discardNoticerecv.GetStatus() == def.StatusOK {
		room.RoomStatus.operationNoticeRecv = nil
		room.RoomStatus.discardNoticerecv = discardNoticerecv
		room.RoomStatus.endRoundRecv = nil
	}
	return
}

func (m RoomMap) handleCancelPrepare(
	id, roomID int32,
	recv *pf.CancelPrepareRecv,
	otherIDs *[]int32) {

	m.Lock()
	defer m.Unlock()

	room, ok := RoomCache.mmap[roomID]
	if !ok {
		recv.Status = def.StatusNoRoom
		return
	}

	pos, ok := room.InRoom(id, roomID, otherIDs)
	if !ok {
		recv.Status = def.StatusNotInRoom
		return
	}

	if room.RoomStatus.IsGaming {
		recv.Status = def.StatusIsGaming
		return
	}

	room.RoomSeats[pos].Prepared = false
	recv.Pos = pos + 1
	return
}

func (m RoomMap) handleDiscard(
	id, roomID int32,
	send *pf.DiscardSend,
	recv *pf.DiscardRecv,
	drawCardRecv *pf.DrawCardRecv,
	operationNoticeRecv *pf.OperationNoticeRecv,
	discardNoticerecv *pf.DiscardNoticeRecv,
	endRoundRecv *pf.EndRoundRecv,
	settlementRecv *pf.SettlementRecv,
	closeRoomRecv *pf.CloseRoomRecv,
	nextID *int32,
	otherIDs *[]int32,
	multiHuIDs *[]int32) {

	m.Lock()
	defer m.Unlock()

	card := send.GetCard()

	room, ok := RoomCache.mmap[roomID]
	if !ok {
		recv.Status = def.StatusNoRoom
		return
	}

	pos, ok := room.InRoom(id, roomID, otherIDs)
	if !ok {
		recv.Status = def.StatusNotInRoom
		return
	}

	if !room.RoomStatus.IsGaming {
		recv.Status = def.StatusNotGaming
		return
	}

	RoomCache.mmap[roomID].PrintTile("出牌前")

	discardPos := pos
	discardSeat := room.RoomSeats[discardPos]
	discardTile := room.RoomTiles.Tiles.GetEndTile(card)
	discardSeat.DiscardTile = &discardTile
	discardSeat.DeskTiles = append(discardSeat.DeskTiles, discardTile)
	if discardSeat.DrawTile == nil {
		discardSeat.HandTiles = discardSeat.HandTiles.DelFrontCard(card)
	} else {
		discardSeat.HandTiles = append(discardSeat.HandTiles, *discardSeat.DrawTile).DelFrontCard(card)
	}
	discardSeat.DrawTile = nil
	room.RoomStatus.LastPos = discardPos
	recv.Card = card
	recv.Pos = pos + 1
	recv.LeftCards = discardSeat.HandTiles.GetFrontCards()

	RoomCache.mmap[roomID].PrintTile("出牌后")

	// 操作提示
	for i := int32(1); i < def.RoomPlayerCount; i++ {
		operationPos := (discardPos + i) % def.RoomPlayerCount
		operationSeat := room.RoomSeats[operationPos]
		operationList := operationSeat.GetOtherHupaiList(room, discardTile, discardPos)
		if len(operationList) > 0 {
			operationNoticeRecv.OperationList = append(operationNoticeRecv.OperationList, operationList...)

			operationNoticeRecv.Status = def.StatusOK
			room.RoomStatus.CurPos = operationPos
			*multiHuIDs = append(*multiHuIDs, operationSeat.ID)

			huTiles := make(Tiles, len(room.RoomHupai.DianPaoHuTiles))
			copy(huTiles, room.RoomHupai.DianPaoHuTiles)
			room.RoomHupai.MultiJiePaoPos = append(room.RoomHupai.MultiJiePaoPos, operationPos)
			room.RoomHupai.MultiDianPaoHuTiles = append(room.RoomHupai.MultiDianPaoHuTiles, huTiles)
			room.RoomHupai.MultiDianPaoHuFanshu = append(room.RoomHupai.MultiDianPaoHuFanshu, room.RoomHupai.DianPaoHuFanshu)
		}
	}

	// 一炮多响
	if len(*multiHuIDs) > 1 {
		room.RoomHupai.IsDianPao = false
		room.RoomHupai.IsMultiDianPao = true
		room.RoomHupai.IsMultiDianPaoHupai = false
		room.RoomHupai.MultiDianPaoPos = discardPos
		room.RoomHupai.MultiDianPaoTile = &discardTile
	}

	// 混牌不能碰杠吃
	if !room.RoomConfig.Config.Daihuner || room.RoomConfig.HunTile.Unicode != card {
		if len(operationNoticeRecv.OperationList) == 0 {
			for i := int32(1); i < def.RoomPlayerCount; i++ {
				operationPos := (discardPos + i) % def.RoomPlayerCount
				operationSeat := room.RoomSeats[operationPos]
				operationNoticeRecv.OperationList = operationSeat.GetOtherGangPengList(discardTile, discardPos)
				if len(operationNoticeRecv.OperationList) > 0 {
					operationNoticeRecv.Status = def.StatusOK
					room.RoomStatus.CurPos = operationPos
					break
				}
			}
		}
		if len(operationNoticeRecv.OperationList) == 0 && room.RoomConfig.Config.Kechipai {
			operationPos := (discardPos + 1) % def.RoomPlayerCount
			operationSeat := room.RoomSeats[operationPos]
			operationNoticeRecv.OperationList = operationSeat.GetOtherChiList(discardTile, discardPos)
			if len(operationNoticeRecv.OperationList) > 0 {
				operationNoticeRecv.Status = def.StatusOK
				room.RoomStatus.CurPos = operationPos
			}
		}
	}

	// 抓牌
	if operationNoticeRecv.GetStatus() != def.StatusOK {
		// 判断是否流局
		if room.RoomTiles.IsLiuju() {
			endRoundRecv.Status = def.StatusOK
			endRoundRecv.ItemList = room.GetResultItems()
			room.RoomStatus.GameStatus = def.GameStatusRoundOver
			room.RoomStatus.IsGaming = false
			// 判断是否大结算
			if room.RoomStatus.CurRound >= room.RoomConfig.Config.Jushu {
				settlementRecv.Status = def.StatusOK
				settlementRecv.SettleList = room.GetSettleItems()
				// 关闭房间
				closeRoomRecv.Status = def.StatusOK
				closeRoomRecv.CloseType = def.CloseRoomRoundOver
			}
		} else {
			room.RoomStatus.CurPos = (room.RoomStatus.CurPos + 1) % def.RoomPlayerCount
			drawPos := room.RoomStatus.CurPos
			drawSeat := room.RoomSeats[drawPos]
			drawSeat.DrawTile = room.RoomTiles.DealOneTile()

			drawCardRecv.Status = def.StatusOK
			drawCardRecv.Pos = drawPos + 1
			drawCardRecv.Card = drawSeat.DrawTile.Front
			room.RoomStatus.CurPos = drawPos

			// 操作提示
			operationNoticeRecv.OperationList = drawSeat.GetSelfHupaiList(room, false)
			if len(operationNoticeRecv.OperationList) > 0 {
				operationNoticeRecv.Status = def.StatusOK
			} else {
				operationNoticeRecv.OperationList = drawSeat.GetSelfGangList()
				if len(operationNoticeRecv.OperationList) > 0 {
					operationNoticeRecv.Status = def.StatusOK
				}
			}

			// 出牌提示
			if operationNoticeRecv.GetStatus() != def.StatusOK {
				discardNoticerecv.Status = def.StatusOK
				discardNoticerecv.Pos = drawPos + 1
			}
		}
	}

	if drawCardRecv.GetStatus() == def.StatusOK {
		RoomCache.mmap[roomID].PrintTile("出牌后抓牌")
	}
	*nextID = room.RoomSeats[room.RoomStatus.CurPos].ID

	// 断线重连数据
	if operationNoticeRecv.GetStatus() == def.StatusOK {
		room.RoomStatus.operationNoticeRecv = operationNoticeRecv
		room.RoomStatus.discardNoticerecv = nil
		room.RoomStatus.endRoundRecv = nil
	} else if discardNoticerecv.GetStatus() == def.StatusOK {
		room.RoomStatus.operationNoticeRecv = nil
		room.RoomStatus.discardNoticerecv = discardNoticerecv
		room.RoomStatus.endRoundRecv = nil
	} else if endRoundRecv.GetStatus() == def.StatusOK {
		room.RoomStatus.operationNoticeRecv = nil
		room.RoomStatus.discardNoticerecv = nil
		room.RoomStatus.endRoundRecv = endRoundRecv
	}

	// 清理房间
	if closeRoomRecv.GetStatus() == def.StatusOK {
		for _, seat := range room.RoomSeats {
			cache.SetRoomID(seat.ID, 0, 0)
		}
		delete(m.mmap, roomID)
	}
	return
}

func (m RoomMap) handleOperation(
	id, roomID int32,
	send *pf.OperationSend,
	recv *pf.OperationRecv,
	drawCardRecv *pf.DrawCardRecv,
	operationNoticeRecv *pf.OperationNoticeRecv,
	discardNoticerecv *pf.DiscardNoticeRecv,
	endRoundRecv *pf.EndRoundRecv,
	settlementRecv *pf.SettlementRecv,
	closeRoomRecv *pf.CloseRoomRecv,
	nextID *int32,
	otherIDs *[]int32) {

	m.Lock()
	defer m.Unlock()

	operation := send.GetOperation()

	room, ok := RoomCache.mmap[roomID]
	if !ok {
		recv.Status = def.StatusNoRoom
		return
	}

	_, ok = room.InRoom(id, roomID, otherIDs)
	if !ok {
		recv.Status = def.StatusNotInRoom
		return
	}

	if !room.RoomStatus.IsGaming {
		recv.Status = def.StatusNotGaming
		return
	}

	// 检查操作
	if room.RoomHupai.IsMultiDianPao {
		// 一炮多响检查
		if operation.Type == def.DianpaoHu {
			operation.Type = def.MultiHu
		}
	} else {
		if !room.RoomStatus.CheckOperation(operation) {
			recv.Status = def.StatusOperError
			return
		}
	}

	opType := operation.GetType()
	card := operation.GetKeycard()
	cards := operation.GetCardsList()
	srcPos := operation.GetSrcPos() - 1
	srcSeat := room.RoomSeats[srcPos]
	desPos := operation.GetDesPos() - 1
	desSeat := room.RoomSeats[desPos]

	recv.Operation = operation
	RoomCache.mmap[roomID].PrintTile("操作前")

	switch opType {
	case def.Peng:
		tile := room.RoomTiles.Tiles.GetEndTile(card)
		srcSeat.DeskTiles = srcSeat.DeskTiles.DelFrontCard(card)
		srcSeat.IsPengTiles = append(srcSeat.IsPengTiles, tile)
		desSeat.HandTiles = desSeat.HandTiles.DelPengCard(card)
		desSeat.PengTiles = append(desSeat.PengTiles, Tiles{tile, tile, tile}...)

		discardNoticerecv.Status = def.StatusOK
		discardNoticerecv.Pos = desPos + 1
	case def.Hu:
	case def.Chi:
		tile := room.RoomTiles.Tiles.GetEndTile(card)
		tiles := room.RoomTiles.Tiles.GetEndTiles(cards)
		chiTiles := tiles.DelFrontCard(card)
		srcSeat.DeskTiles = srcSeat.DeskTiles.DelFrontCard(card)
		srcSeat.IsChiTiles = append(srcSeat.IsChiTiles, tile)
		desSeat.HandTiles = desSeat.HandTiles.DelChiCard(chiTiles)
		desSeat.ChiTiles = append(desSeat.ChiTiles, tiles...)

		discardNoticerecv.Status = def.StatusOK
		discardNoticerecv.Pos = desPos + 1
	case def.MingGang:
		tile := room.RoomTiles.Tiles.GetEndTile(card)
		srcSeat.DeskTiles = srcSeat.DeskTiles.DelFrontCard(card)
		srcSeat.IsMingGangTiles = append(srcSeat.IsMingGangTiles, tile)
		desSeat.HandTiles = desSeat.HandTiles.DelGangCard(card)
		desSeat.MingGangTiles = append(srcSeat.MingGangTiles, Tiles{tile, tile, tile, tile}...)
		room.AddMingGangScore(srcPos, desPos)

		if room.RoomTiles.IsLiuju() {
			endRoundRecv.Status = def.StatusOK
			endRoundRecv.ItemList = room.GetResultItems()
			room.RoomStatus.GameStatus = def.GameStatusRoundOver
			room.RoomStatus.IsGaming = false
			// 判断是否大结算
			if room.RoomStatus.CurRound >= room.RoomConfig.Config.Jushu {
				settlementRecv.Status = def.StatusOK
				settlementRecv.SettleList = room.GetSettleItems()
				// 关闭房间
				closeRoomRecv.Status = def.StatusOK
				closeRoomRecv.CloseType = def.CloseRoomRoundOver
			}
		} else {
			drawPos := desPos
			drawSeat := room.RoomSeats[drawPos]
			drawSeat.DrawTile = room.RoomTiles.DealOneTile()
			room.RoomTiles.GangCnt++

			drawCardRecv.Status = def.StatusOK
			drawCardRecv.Pos = drawPos + 1
			drawCardRecv.Card = drawSeat.DrawTile.Front

			operationNoticeRecv.OperationList = drawSeat.GetSelfHupaiList(room, true)
			if len(operationNoticeRecv.OperationList) > 0 {
				operationNoticeRecv.Status = def.StatusOK
			} else {
				operationNoticeRecv.OperationList = drawSeat.GetSelfGangList()
				if len(operationNoticeRecv.OperationList) > 0 {
					operationNoticeRecv.Status = def.StatusOK
				}
			}

			if operationNoticeRecv.GetStatus() != def.StatusOK {
				discardNoticerecv.Status = def.StatusOK
				discardNoticerecv.Pos = drawPos + 1
			}
		}
	case def.AnGang:
		tile := room.RoomTiles.Tiles.GetEndTile(card)
		desSeat.HandTiles = append(desSeat.HandTiles, *desSeat.DrawTile)
		desSeat.HandTiles = desSeat.HandTiles.DelGangCard(card)
		desSeat.AnGangTiles = append(srcSeat.AnGangTiles, Tiles{tile, tile, tile, tile}...)
		room.AddAnGangScore(desPos)

		if room.RoomTiles.IsLiuju() {
			endRoundRecv.Status = def.StatusOK
			endRoundRecv.ItemList = room.GetResultItems()
			room.RoomStatus.GameStatus = def.GameStatusRoundOver
			room.RoomStatus.IsGaming = false
			// 判断是否大结算
			if room.RoomStatus.CurRound >= room.RoomConfig.Config.Jushu {
				settlementRecv.Status = def.StatusOK
				settlementRecv.SettleList = room.GetSettleItems()
				// 关闭房间
				closeRoomRecv.Status = def.StatusOK
				closeRoomRecv.CloseType = def.CloseRoomRoundOver
			}
		} else {
			drawPos := desPos
			drawSeat := room.RoomSeats[drawPos]
			drawSeat.DrawTile = room.RoomTiles.DealOneTile()
			room.RoomTiles.GangCnt++

			drawCardRecv.Status = def.StatusOK
			drawCardRecv.Pos = drawPos + 1
			drawCardRecv.Card = drawSeat.DrawTile.Front

			operationNoticeRecv.OperationList = drawSeat.GetSelfHupaiList(room, true)
			if len(operationNoticeRecv.OperationList) > 0 {
				operationNoticeRecv.Status = def.StatusOK
			} else {
				operationNoticeRecv.OperationList = drawSeat.GetSelfGangList()
				if len(operationNoticeRecv.OperationList) > 0 {
					operationNoticeRecv.Status = def.StatusOK
				}
			}

			if operationNoticeRecv.GetStatus() != def.StatusOK {
				discardNoticerecv.Status = def.StatusOK
				discardNoticerecv.Pos = drawPos + 1
			}
		}
	case def.Pass:
		if srcPos == desPos {
			if room.RoomStatus.GetOperation() == def.Hu {
				if room.RoomHupai.IsZiMo {
					room.RoomHupai = &MJHupai{}
				}
				operationNoticeRecv.OperationList = desSeat.GetSelfGangList()
				if len(operationNoticeRecv.OperationList) > 0 {
					operationNoticeRecv.Status = def.StatusOK
				}
			}

			if operationNoticeRecv.GetStatus() != def.StatusOK {
				discardNoticerecv.Status = def.StatusOK
				discardNoticerecv.Pos = desPos + 1
			}
		} else {
			discardTile := room.RoomTiles.Tiles.GetEndTile(card)
			tempType := room.RoomStatus.GetOperation()
			if room.RoomHupai.IsMultiDianPao {
				if tempType == def.DianpaoHu {
					room.RoomHupai.MultiJiePaoPassPos = append(room.RoomHupai.MultiJiePaoPassPos, desPos)
					if len(room.RoomHupai.MultiJiePaoPos) == len(room.RoomHupai.MultiJiePaoHupaiPos)+len(room.RoomHupai.MultiJiePaoPassPos) {
						if len(room.RoomHupai.MultiJiePaoHupaiPos) > 0 {
							room.SuanFen()
							for _, jiePaoPos := range room.RoomHupai.MultiJiePaoHupaiPos {
								for i, pos := range room.RoomHupai.MultiJiePaoPos {
									if pos == jiePaoPos {
										operation.HuInfoList = append(operation.HuInfoList, &pf.MutiHu{
											Pos:      pos + 1,
											CardList: room.RoomHupai.MultiDianPaoHuTiles[i].DelHuCard().GetFrontCards(),
										})
										break
									}
								}
							}

							// 只有一人变成点炮胡
							if len(room.RoomHupai.MultiJiePaoHupaiPos) == 1 {
								operation.Type = def.DianpaoHu
								operation.DesPos = operation.HuInfoList[0].Pos
								operation.CardsList = operation.HuInfoList[0].CardList
								operation.HuInfoList = nil
							}

							// 小结算
							endRoundRecv.Status = def.StatusOK
							endRoundRecv.ItemList = room.GetResultItems()
							room.RoomStatus.GameStatus = def.GameStatusRoundOver
							room.RoomStatus.IsGaming = false

							// 大结算
							if room.RoomStatus.CurRound >= room.RoomConfig.Config.Jushu {
								settlementRecv.Status = def.StatusOK
								settlementRecv.SettleList = room.GetSettleItems()
								// 关闭房间
								closeRoomRecv.Status = def.StatusOK
								closeRoomRecv.CloseType = def.CloseRoomRoundOver
							}
						} else {
							room.RoomHupai = &MJHupai{}

							for i := int32(1); i < def.RoomPlayerCount; i++ {
								operationPos := (srcPos + i) % def.RoomPlayerCount
								operationSeat := room.RoomSeats[operationPos]
								operationNoticeRecv.OperationList = operationSeat.GetOtherGangPengList(discardTile, srcPos)
								if len(operationNoticeRecv.OperationList) > 0 {
									operationNoticeRecv.Status = def.StatusOK
									room.RoomStatus.CurPos = operationPos
									break
								}
							}
							if operationNoticeRecv.GetStatus() != def.StatusOK {
								if room.RoomConfig.Config.Kechipai {
									operationPos := (srcPos + 1) % def.RoomPlayerCount
									operationSeat := room.RoomSeats[operationPos]
									operationNoticeRecv.OperationList = operationSeat.GetOtherChiList(discardTile, srcPos)
									if len(operationNoticeRecv.OperationList) > 0 {
										operationNoticeRecv.Status = def.StatusOK
										room.RoomStatus.CurPos = operationPos
									}
								}
							}

							if operationNoticeRecv.GetStatus() != def.StatusOK {
								if room.RoomTiles.IsLiuju() {
									endRoundRecv.Status = def.StatusOK
									endRoundRecv.ItemList = room.GetResultItems()
									room.RoomStatus.GameStatus = def.GameStatusRoundOver
									room.RoomStatus.IsGaming = false
									// 判断是否大结算
									if room.RoomStatus.CurRound >= room.RoomConfig.Config.Jushu {
										settlementRecv.Status = def.StatusOK
										settlementRecv.SettleList = room.GetSettleItems()
										// 关闭房间
										closeRoomRecv.Status = def.StatusOK
										closeRoomRecv.CloseType = def.CloseRoomRoundOver
									}
								} else {
									room.RoomStatus.CurPos = (srcPos + 1) % def.RoomPlayerCount
									drawPos := room.RoomStatus.CurPos
									drawSeat := room.RoomSeats[drawPos]
									drawSeat.DrawTile = room.RoomTiles.DealOneTile()

									drawCardRecv.Status = def.StatusOK
									drawCardRecv.Pos = drawPos + 1
									drawCardRecv.Card = drawSeat.DrawTile.Front

									operationNoticeRecv.OperationList = drawSeat.GetSelfHupaiList(room, false)
									if len(operationNoticeRecv.OperationList) > 0 {
										operationNoticeRecv.Status = def.StatusOK
									} else {
										operationNoticeRecv.OperationList = drawSeat.GetSelfGangList()
										if len(operationNoticeRecv.OperationList) > 0 {
											operationNoticeRecv.Status = def.StatusOK
										}
									}

									if operationNoticeRecv.GetStatus() != def.StatusOK {
										discardNoticerecv.Status = def.StatusOK
										discardNoticerecv.Pos = drawPos + 1
									}
								}
							}
						}
					} else {
						recv.Status = def.StatusWait
					}
				} else {

				}
			} else {
				if tempType == def.DianpaoHu {
					if room.RoomHupai.IsDianPao {
						room.RoomHupai = &MJHupai{}
					}
					for i := int32(1); i < def.RoomPlayerCount; i++ {
						operationPos := (srcPos + i) % def.RoomPlayerCount
						operationSeat := room.RoomSeats[operationPos]
						operationNoticeRecv.OperationList = operationSeat.GetOtherGangPengList(discardTile, srcPos)
						if len(operationNoticeRecv.OperationList) > 0 {
							operationNoticeRecv.Status = def.StatusOK
							room.RoomStatus.CurPos = operationPos
							break
						}
					}
					if operationNoticeRecv.GetStatus() != def.StatusOK {
						if room.RoomConfig.Config.Kechipai {
							operationPos := (srcPos + 1) % def.RoomPlayerCount
							operationSeat := room.RoomSeats[operationPos]
							operationNoticeRecv.OperationList = operationSeat.GetOtherChiList(discardTile, srcPos)
							if len(operationNoticeRecv.OperationList) > 0 {
								operationNoticeRecv.Status = def.StatusOK
								room.RoomStatus.CurPos = operationPos
							}
						}
					}
				} else if tempType == def.MingGang || tempType == def.Peng {
					if room.RoomConfig.Config.Kechipai {
						operationPos := (srcPos + 1) % def.RoomPlayerCount
						operationSeat := room.RoomSeats[operationPos]
						operationNoticeRecv.OperationList = operationSeat.GetOtherChiList(discardTile, srcPos)
						if len(operationNoticeRecv.OperationList) > 0 {
							operationNoticeRecv.Status = def.StatusOK
							room.RoomStatus.CurPos = operationPos
						}
					}
				}

				if operationNoticeRecv.GetStatus() != def.StatusOK {
					if room.RoomTiles.IsLiuju() {
						endRoundRecv.Status = def.StatusOK
						endRoundRecv.ItemList = room.GetResultItems()
						room.RoomStatus.GameStatus = def.GameStatusRoundOver
						room.RoomStatus.IsGaming = false
						// 判断是否大结算
						if room.RoomStatus.CurRound >= room.RoomConfig.Config.Jushu {
							settlementRecv.Status = def.StatusOK
							settlementRecv.SettleList = room.GetSettleItems()
							// 关闭房间
							closeRoomRecv.Status = def.StatusOK
							closeRoomRecv.CloseType = def.CloseRoomRoundOver
						}
					} else {
						room.RoomStatus.CurPos = (srcPos + 1) % def.RoomPlayerCount
						drawPos := room.RoomStatus.CurPos
						drawSeat := room.RoomSeats[drawPos]
						drawSeat.DrawTile = room.RoomTiles.DealOneTile()
						drawCardRecv.Status = def.StatusOK
						drawCardRecv.Pos = drawPos + 1
						drawCardRecv.Card = drawSeat.DrawTile.Front
						operationNoticeRecv.OperationList = drawSeat.GetSelfHupaiList(room, false)
						if len(operationNoticeRecv.OperationList) > 0 {
							operationNoticeRecv.Status = def.StatusOK
						} else {
							operationNoticeRecv.OperationList = drawSeat.GetSelfGangList()
							if len(operationNoticeRecv.OperationList) > 0 {
								operationNoticeRecv.Status = def.StatusOK
							}
						}
						if operationNoticeRecv.GetStatus() != def.StatusOK {
							discardNoticerecv.Status = def.StatusOK
							discardNoticerecv.Pos = drawPos + 1
						}
					}
				}
			}
		}
	case def.MultiHu:
		if room.RoomHupai.IsMultiDianPao {
			room.RoomHupai.IsMultiDianPaoHupai = true
			room.RoomHupai.MultiJiePaoHupaiPos = append(room.RoomHupai.MultiJiePaoHupaiPos, desPos)
			if len(room.RoomHupai.MultiJiePaoPos) == len(room.RoomHupai.MultiJiePaoHupaiPos)+len(room.RoomHupai.MultiJiePaoPassPos) {
				room.SuanFen()
				for _, jiePaoPos := range room.RoomHupai.MultiJiePaoHupaiPos {
					for i, pos := range room.RoomHupai.MultiJiePaoPos {
						if pos == jiePaoPos {
							operation.HuInfoList = append(operation.HuInfoList, &pf.MutiHu{
								Pos:      pos + 1,
								CardList: room.RoomHupai.MultiDianPaoHuTiles[i].DelHuCard().GetFrontCards(),
							})
							break
						}
					}
				}

				// 只有一人变成点炮胡
				if len(room.RoomHupai.MultiJiePaoHupaiPos) == 1 {
					operation.Type = def.DianpaoHu
					operation.DesPos = operation.HuInfoList[0].Pos
					operation.CardsList = operation.HuInfoList[0].CardList
					operation.HuInfoList = nil
				}

				// 小结算
				endRoundRecv.Status = def.StatusOK
				endRoundRecv.ItemList = room.GetResultItems()
				room.RoomStatus.GameStatus = def.GameStatusRoundOver
				room.RoomStatus.IsGaming = false

				// 大结算
				if room.RoomStatus.CurRound >= room.RoomConfig.Config.Jushu {
					settlementRecv.Status = def.StatusOK
					settlementRecv.SettleList = room.GetSettleItems()
					// 关闭房间
					closeRoomRecv.Status = def.StatusOK
					closeRoomRecv.CloseType = def.CloseRoomRoundOver
				}
			} else {
				recv.Status = def.StatusWait
			}
		}
	case def.BuGang:
		tile := room.RoomTiles.Tiles.GetEndTile(card)
		for i := int32(1); i < def.RoomPlayerCount; i++ {
			srcPos = (desPos + i) % def.RoomPlayerCount
			srcSeat := room.RoomSeats[srcPos]
			if srcSeat.IsPengTiles.InTiles(tile) {
				srcSeat.IsPengTiles = srcSeat.IsPengTiles.DelFrontCard(card)
				srcSeat.IsBuGangTiles = append(srcSeat.IsBuGangTiles, tile)
				break
			}
		}

		// 补杠先把抓的牌拿到手上
		if desSeat.DrawTile != nil {
			desSeat.HandTiles = append(desSeat.HandTiles, *desSeat.DrawTile)
		}
		desSeat.HandTiles = desSeat.HandTiles.DelGangCard(card)
		desSeat.PengTiles = desSeat.PengTiles.DelGangCard(card)
		desSeat.BuGangTiles = append(srcSeat.BuGangTiles, Tiles{tile, tile, tile, tile}...)
		room.AddBuGangScore(desPos)

		if room.RoomTiles.IsLiuju() {
			endRoundRecv.Status = def.StatusOK
			endRoundRecv.ItemList = room.GetResultItems()
			room.RoomStatus.GameStatus = def.GameStatusRoundOver
			room.RoomStatus.IsGaming = false
			// 判断是否大结算
			if room.RoomStatus.CurRound >= room.RoomConfig.Config.Jushu {
				settlementRecv.Status = def.StatusOK
				settlementRecv.SettleList = room.GetSettleItems()
				// 关闭房间
				closeRoomRecv.Status = def.StatusOK
				closeRoomRecv.CloseType = def.CloseRoomRoundOver
			}
		} else {
			drawPos := desPos
			drawSeat := room.RoomSeats[drawPos]
			drawSeat.DrawTile = room.RoomTiles.DealOneTile()
			room.RoomTiles.GangCnt++

			drawCardRecv.Status = def.StatusOK
			drawCardRecv.Pos = drawPos + 1
			drawCardRecv.Card = drawSeat.DrawTile.Front

			operationNoticeRecv.OperationList = drawSeat.GetSelfHupaiList(room, true)
			if len(operationNoticeRecv.OperationList) > 0 {
				operationNoticeRecv.Status = def.StatusOK
			} else {
				operationNoticeRecv.OperationList = drawSeat.GetSelfGangList()
				if len(operationNoticeRecv.OperationList) > 0 {
					operationNoticeRecv.Status = def.StatusOK
				}
			}

			if operationNoticeRecv.GetStatus() != def.StatusOK {
				discardNoticerecv.Status = def.StatusOK
				discardNoticerecv.Pos = drawPos + 1
			}
		}
	case def.Wating:
	case def.ZimoHu:
		if room.RoomHupai.IsZiMo {
			room.RoomHupai.IsZiMoHupai = true
			room.SuanFen()
			operation.CardsList = room.RoomHupai.ZiMoHuTiles.GetFrontCards()

			// 小结算
			endRoundRecv.Status = def.StatusOK
			endRoundRecv.ItemList = room.GetResultItems()
			room.RoomStatus.GameStatus = def.GameStatusRoundOver
			room.RoomStatus.IsGaming = false

			// 大结算
			if room.RoomStatus.CurRound >= room.RoomConfig.Config.Jushu {
				settlementRecv.Status = def.StatusOK
				settlementRecv.SettleList = room.GetSettleItems()
				// 关闭房间
				closeRoomRecv.Status = def.StatusOK
				closeRoomRecv.CloseType = def.CloseRoomRoundOver
			}
		}
	case def.DianpaoHu:
		if room.RoomHupai.IsDianPao {
			room.RoomHupai.IsDianPaoHupai = true
			room.SuanFen()
			operation.CardsList = room.RoomHupai.ZiMoHuTiles.GetFrontCards()

			// 小结算
			endRoundRecv.Status = def.StatusOK
			endRoundRecv.ItemList = room.GetResultItems()
			room.RoomStatus.GameStatus = def.GameStatusRoundOver
			room.RoomStatus.IsGaming = false

			// 大结算
			if room.RoomStatus.CurRound >= room.RoomConfig.Config.Jushu {
				settlementRecv.Status = def.StatusOK
				settlementRecv.SettleList = room.GetSettleItems()
				// 关闭房间
				closeRoomRecv.Status = def.StatusOK
				closeRoomRecv.CloseType = def.CloseRoomRoundOver
			}
		}
	}
	RoomCache.mmap[roomID].PrintTile("操作后")

	if drawCardRecv.GetStatus() == def.StatusOK {
		RoomCache.mmap[roomID].PrintTile("操作后抓牌")
	}
	*nextID = room.RoomSeats[room.RoomStatus.CurPos].ID

	// 断线重连数据
	if operationNoticeRecv.GetStatus() == def.StatusOK {
		room.RoomStatus.operationNoticeRecv = operationNoticeRecv
		room.RoomStatus.discardNoticerecv = nil
		room.RoomStatus.endRoundRecv = nil
	} else if discardNoticerecv.GetStatus() == def.StatusOK {
		room.RoomStatus.operationNoticeRecv = nil
		room.RoomStatus.discardNoticerecv = discardNoticerecv
		room.RoomStatus.endRoundRecv = nil
	} else if endRoundRecv.GetStatus() == def.StatusOK {
		room.RoomStatus.operationNoticeRecv = nil
		room.RoomStatus.discardNoticerecv = nil
		room.RoomStatus.endRoundRecv = endRoundRecv
	}

	// 清理房间
	if closeRoomRecv.GetStatus() == def.StatusOK {
		for _, seat := range room.RoomSeats {
			cache.SetRoomID(seat.ID, 0, 0)
		}
		delete(m.mmap, roomID)
	}
	return
}
