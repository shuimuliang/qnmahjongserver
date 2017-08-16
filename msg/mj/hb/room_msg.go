package hb

import (
	"qnmahjong/cache"
	"qnmahjong/def"
	"qnmahjong/pf"
	"qnmahjong/util"
	"time"

	"github.com/olahol/melody"
)

func HandleEnterRoom(s *melody.Session, mjPlayer *cache.MjPlayer, send *pf.EnterRoomSend, abs *pf.AbsMessage) (err error) {
	playerID := mjPlayer.PlayerID
	roomID := mjPlayer.RoomID
	msgID := abs.MsgID

	enterRoomRecv := &pf.EnterRoomRecv{
		Status: def.StatusOK,
		MjType: def.MjTypeHb,
	}
	otherJoinRecv := &pf.OtherJoinRecv{
		Status: def.StatusFailed,
	}
	voteCloseRecv := &pf.VoteCloseRecv{
		Status: def.StatusFailed,
	}
	operationNoticeRecv := &pf.OperationNoticeRecv{
		Status: def.StatusFailed,
	}
	discardNoticerecv := &pf.DiscardNoticeRecv{
		Status: def.StatusFailed,
	}
	endRoundRecv := &pf.EndRoundRecv{
		Status: def.StatusFailed,
	}
	otherIDs := make([]int32, 0, def.RoomPlayerCount-1)

	if send.RoomID == 0 {
		// 创建房间
		RoomCache.handleCreateRoom(
			playerID,
			send,
			enterRoomRecv)
	} else if send.Relink == 0 {
		// 进入房间
		RoomCache.handleEnterRoom(
			playerID,
			send,
			enterRoomRecv,
			otherJoinRecv,
			&otherIDs)
	} else {
		// 断线重连
		RoomCache.handleGameRelink(
			playerID,
			send,
			enterRoomRecv,
			voteCloseRecv,
			operationNoticeRecv,
			discardNoticerecv,
			endRoundRecv,
		)
	}

	// 进入房间给自己发消息
	recv, err := enterRoomRecv.Marshal()
	if err == nil {
		abs.MsgBody = recv
		recv, err = abs.Marshal()
		if err == nil {
			err = s.WriteBinary(recv)
			if err == nil {
				util.LogRecv(msgID, playerID, roomID, enterRoomRecv, "EnterRoom")
			}
		}
	}

	// 进入房间给其他人发消息
	if otherJoinRecv.GetStatus() == def.StatusOK {
		msgBody, err := otherJoinRecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.OtherJoin),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				for _, otherID := range otherIDs {
					if cache.SendMessage(otherID, msg) == nil {
						util.LogRecv(absMessage.GetMsgID(), otherID, roomID, otherJoinRecv, "OtherJoin")
					}
				}
			}
		}
	}

	// 给断线人发投票消息
	if voteCloseRecv.GetStatus() == def.StatusOK {
		msgBody, err := voteCloseRecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.VoteClose),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				if cache.SendMessage(playerID, msg) == nil {
					util.LogRecv(absMessage.GetMsgID(), playerID, roomID, voteCloseRecv, "VoteClose")
				}
			}
		}
	}

	// 给断线人发操作提示消息
	if operationNoticeRecv.GetStatus() == def.StatusOK {
		msgBody, err := operationNoticeRecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.OperationNotice),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				if cache.SendMessage(playerID, msg) == nil {
					util.LogRecv(absMessage.GetMsgID(), playerID, roomID, operationNoticeRecv, "OperationNotice")
				}
			}
		}
	}

	// 给断线人发出牌提示消息
	if discardNoticerecv.GetStatus() == def.StatusOK {
		msgBody, err := discardNoticerecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.DiscardNotice),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				if cache.SendMessage(playerID, msg) == nil {
					util.LogRecv(absMessage.GetMsgID(), playerID, roomID, discardNoticerecv, "DiscardNotice")
				}
			}
		}
	}

	// 给断线人发小结算消息
	if endRoundRecv.GetStatus() == def.StatusOK {
		msgBody, err := endRoundRecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.EndRound),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				if cache.SendMessage(playerID, msg) == nil {
					util.LogRecv(absMessage.GetMsgID(), playerID, roomID, endRoundRecv, "EndRound")
				}
			}
		}
	}
	return
}

func HandleExitRoom(s *melody.Session, mjPlayer *cache.MjPlayer, abs *pf.AbsMessage) (err error) {
	playerID := mjPlayer.PlayerID
	roomID := mjPlayer.RoomID
	msgID := abs.MsgID

	exitRoomRecv := &pf.ExitRoomRecv{
		Status: def.StatusOK,
	}
	otherIDs := make([]int32, 0, def.RoomPlayerCount-1)

	RoomCache.handleExitRoom(
		playerID,
		roomID,
		exitRoomRecv,
		&otherIDs)

	// 退出房间给自己发消息
	recv, err := exitRoomRecv.Marshal()
	if err == nil {
		abs.MsgBody = recv
		recv, err = abs.Marshal()
		if err == nil {
			err = s.WriteBinary(recv)
			if err == nil {
				util.LogRecv(msgID, playerID, roomID, exitRoomRecv, "ExitRoom")
			}
		}
	}

	// 退出房间给其他人发消息
	if exitRoomRecv.Status == def.StatusOK {
		msgBody, err := exitRoomRecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.ExitRoom),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				for _, otherID := range otherIDs {
					if cache.SendMessage(otherID, msg) == nil {
						util.LogRecv(absMessage.GetMsgID(), otherID, roomID, exitRoomRecv, "ExitRoom")
					}
				}
			}
		}
	}
	return
}

func HandleCloseRoom(s *melody.Session, mjPlayer *cache.MjPlayer, abs *pf.AbsMessage) (err error) {
	playerID := mjPlayer.PlayerID
	roomID := mjPlayer.RoomID
	msgID := abs.MsgID

	closeRoomRecv := &pf.CloseRoomRecv{
		Status: def.StatusOK,
	}
	otherIDs := make([]int32, 0, def.RoomPlayerCount-1)

	RoomCache.handleCloseRoom(
		playerID,
		roomID,
		closeRoomRecv,
		&otherIDs)

	// 关闭房间给自己发消息
	recv, err := closeRoomRecv.Marshal()
	if err == nil {
		abs.MsgBody = recv
		recv, err = abs.Marshal()
		if err == nil {
			err = s.WriteBinary(recv)
			if err == nil {
				util.LogRecv(msgID, playerID, roomID, closeRoomRecv, "CloseRoom")
			}
		}
	}

	// 关闭房间给其他人发消息
	if closeRoomRecv.Status == def.StatusOK {
		msgBody, err := closeRoomRecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.CloseRoom),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				for _, otherID := range otherIDs {
					if cache.SendMessage(otherID, msg) == nil {
						util.LogRecv(absMessage.GetMsgID(), otherID, roomID, closeRoomRecv, "CloseRoom")
					}
				}
			}
		}
	}
	return
}

func HandleVoteClose(s *melody.Session, mjPlayer *cache.MjPlayer, send *pf.VoteCloseSend, abs *pf.AbsMessage) (err error) {
	playerID := mjPlayer.PlayerID
	roomID := mjPlayer.RoomID
	msgID := abs.MsgID

	voteCloseRecv := &pf.VoteCloseRecv{
		Status: def.StatusOK,
	}
	closeRoomRecv := &pf.CloseRoomRecv{
		Status: def.StatusFailed,
	}
	otherIDs := make([]int32, 0, def.RoomPlayerCount-1)

	RoomCache.handleVoteClose(
		playerID,
		roomID,
		send,
		voteCloseRecv,
		closeRoomRecv,
		&otherIDs)

	// 投票解散房间给自己发消息
	recv, err := voteCloseRecv.Marshal()
	if err == nil {
		abs.MsgBody = recv
		recv, err = abs.Marshal()
		if err == nil {
			err = s.WriteBinary(recv)
			if err == nil {
				util.LogRecv(msgID, playerID, roomID, voteCloseRecv, "VoteClose")
			}
		}
	}

	// 投票解散房间给其他人发消息
	if voteCloseRecv.Status == def.StatusOK {
		msgBody, err := voteCloseRecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.VoteClose),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				for _, otherID := range otherIDs {
					if cache.SendMessage(otherID, msg) == nil {
						util.LogRecv(absMessage.GetMsgID(), otherID, roomID, voteCloseRecv, "VoteClose")
					}
				}
			}
		}
	}

	// 关闭房间给所有人发消息
	if closeRoomRecv.Status == def.StatusOK {
		msgBody, err := closeRoomRecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.CloseRoom),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				deskIDs := append(otherIDs, playerID)
				for _, deskID := range deskIDs {
					if cache.SendMessage(deskID, msg) == nil {
						util.LogRecv(absMessage.GetMsgID(), deskID, roomID, closeRoomRecv, "CloseRoom")
					}
				}
			}
		}
	}
	return
}

func HandleGameChat(s *melody.Session, mjPlayer *cache.MjPlayer, send *pf.GameChatSend, abs *pf.AbsMessage) (err error) {
	playerID := mjPlayer.PlayerID
	roomID := mjPlayer.RoomID
	msgID := abs.MsgID

	gameChatRecv := &pf.GameChatRecv{
		Status: def.StatusOK,
	}
	otherIDs := make([]int32, 0, def.RoomPlayerCount-1)

	RoomCache.handleGameChat(
		playerID,
		roomID,
		send,
		gameChatRecv,
		&otherIDs)

	// 游戏内聊天给自己发消息
	recv, err := gameChatRecv.Marshal()
	if err == nil {
		abs.MsgBody = recv
		recv, err = abs.Marshal()
		if err == nil {
			err = s.WriteBinary(recv)
			if err == nil {
				util.LogRecv(msgID, playerID, roomID, gameChatRecv, "GameChat")
			}
		}
	}

	// 游戏内聊天给其他人发消息
	if gameChatRecv.Status == def.StatusOK {
		msgBody, err := gameChatRecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.GameChat),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				for _, otherID := range otherIDs {
					if cache.SendMessage(otherID, msg) == nil {
						util.LogRecv(absMessage.GetMsgID(), otherID, roomID, gameChatRecv, "GameChat")
					}
				}
			}
		}
	}
	return
}

func HandlePrepareGame(s *melody.Session, mjPlayer *cache.MjPlayer, abs *pf.AbsMessage) (err error) {
	playerID := mjPlayer.PlayerID
	roomID := mjPlayer.RoomID
	msgID := abs.MsgID

	prepareGameRecv := &pf.PrepareGameRecv{
		Status: def.StatusOK,
	}
	startRoundRecv := &pf.StartRoundRecv{
		Status: def.StatusFailed,
	}
	drawCardRecv := &pf.DrawCardRecv{
		Status: def.StatusFailed,
	}
	operationNoticeRecv := &pf.OperationNoticeRecv{
		Status: def.StatusFailed,
	}
	discardNoticerecv := &pf.DiscardNoticeRecv{
		Status: def.StatusFailed,
	}

	bankerID := int32(0)
	otherIDs := make([]int32, 0, def.RoomPlayerCount-1)
	otherCards := make(map[int32][]int32, def.RoomPlayerCount-1)

	RoomCache.handlePrepareGame(
		playerID,
		roomID,
		prepareGameRecv,
		startRoundRecv,
		drawCardRecv,
		operationNoticeRecv,
		discardNoticerecv,
		&bankerID,
		&otherIDs,
		otherCards)

	recv, err := prepareGameRecv.Marshal()
	if err == nil {
		abs.MsgBody = recv
		recv, err = abs.Marshal()
		if err == nil {
			err = s.WriteBinary(recv)
			if err == nil {
				util.LogRecv(msgID, playerID, roomID, prepareGameRecv, "PrepareGame")
			}
		}
	}

	// 准备游戏给其他人发消息
	if prepareGameRecv.Status == def.StatusOK {
		msgBody, err := prepareGameRecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.PrepareGame),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				for _, otherID := range otherIDs {
					if cache.SendMessage(otherID, msg) == nil {
						util.LogRecv(absMessage.GetMsgID(), otherID, roomID, prepareGameRecv, "PrepareGame")
					}
				}
			}
		}
	}

	// 牌局开始给庄家发消息
	if startRoundRecv.Status == def.StatusOK {
		msgBody, err := startRoundRecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.StartRound),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				if cache.SendMessage(playerID, msg) == nil {
					util.LogRecv(absMessage.GetMsgID(), playerID, roomID, startRoundRecv, "StartRound")
				}
			}
		}
	}

	// 牌局开始给其他人发消息
	if startRoundRecv.Status == def.StatusOK {
		for _, otherID := range otherIDs {
			startRoundRecv.MyCards = otherCards[otherID]
			msgBody, err := startRoundRecv.Marshal()
			if err == nil {
				absMessage := &pf.AbsMessage{
					MsgID:   int32(pf.StartRound),
					MsgBody: msgBody,
				}
				msg, err := absMessage.Marshal()
				if err == nil {
					if cache.SendMessage(otherID, msg) == nil {
						util.LogRecv(absMessage.GetMsgID(), otherID, roomID, startRoundRecv, "StartRound")
					}
				}
			}
		}
	}

	// 延迟4s
	time.Sleep(time.Second * 4)

	// 庄家抓牌给庄家发消息
	if drawCardRecv.Status == def.StatusOK {
		msgBody, err := drawCardRecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.DrawCard),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				if cache.SendMessage(bankerID, msg) == nil {
					util.LogRecv(absMessage.GetMsgID(), bankerID, roomID, drawCardRecv, "DrawCard")
				}
			}
		}
	}

	// 庄家抓牌给其他人发消息
	if drawCardRecv.Status == def.StatusOK {
		drawCardRecv.Card = 0
		msgBody, err := drawCardRecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.DrawCard),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				deskIDs := append(otherIDs, playerID)
				for _, deskID := range deskIDs {
					if deskID == bankerID {
						continue
					}
					if cache.SendMessage(deskID, msg) == nil {
						util.LogRecv(absMessage.GetMsgID(), deskID, roomID, drawCardRecv, "DrawCard")
					}
				}
			}
		}
	}

	// 庄家操作提示给庄家发消息
	if operationNoticeRecv.Status == def.StatusOK {
		msgBody, err := operationNoticeRecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.OperationNotice),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				if cache.SendMessage(bankerID, msg) == nil {
					util.LogRecv(absMessage.GetMsgID(), bankerID, roomID, operationNoticeRecv, "OperationNotice")
				}
			}
		}
	}

	// 庄家出牌提示给桌上人发消息
	if discardNoticerecv.Status == def.StatusOK {
		msgBody, err := discardNoticerecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.DiscardNotice),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				deskIDs := append(otherIDs, playerID)
				for _, deskID := range deskIDs {
					if cache.SendMessage(deskID, msg) == nil {
						util.LogRecv(absMessage.GetMsgID(), deskID, roomID, discardNoticerecv, "DiscardNotice")
					}
				}
			}
		}
	}
	return
}

func HandleCancelPrepare(s *melody.Session, mjPlayer *cache.MjPlayer, abs *pf.AbsMessage) (err error) {
	playerID := mjPlayer.PlayerID
	roomID := mjPlayer.RoomID
	msgID := abs.MsgID

	cancelPrepareRecv := &pf.CancelPrepareRecv{
		Status: def.StatusOK,
	}
	otherIDs := make([]int32, 0, def.RoomPlayerCount-1)

	RoomCache.handleCancelPrepare(
		playerID,
		roomID,
		cancelPrepareRecv,
		&otherIDs)

	// 取消准备给自己发消息
	recv, err := cancelPrepareRecv.Marshal()
	if err == nil {
		abs.MsgBody = recv
		recv, err = abs.Marshal()
		if err == nil {
			err = s.WriteBinary(recv)
			if err == nil {
				util.LogRecv(msgID, playerID, roomID, cancelPrepareRecv, "CancelPrepare")
			}
		}
	}

	// 取消准备给其他人发消息
	if cancelPrepareRecv.Status == def.StatusOK {
		msgBody, err := cancelPrepareRecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.CancelPrepare),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				for _, otherID := range otherIDs {
					if cache.SendMessage(otherID, msg) == nil {
						util.LogRecv(absMessage.GetMsgID(), otherID, roomID, cancelPrepareRecv, "CancelPrepare")
					}
				}
			}
		}
	}
	return
}

func HandleDiscard(s *melody.Session, mjPlayer *cache.MjPlayer, send *pf.DiscardSend, abs *pf.AbsMessage) (err error) {
	playerID := mjPlayer.PlayerID
	roomID := mjPlayer.RoomID
	msgID := abs.MsgID

	discardRecv := &pf.DiscardRecv{
		Status: def.StatusOK,
	}
	drawCardRecv := &pf.DrawCardRecv{
		Status: def.StatusFailed,
	}
	operationNoticeRecv := &pf.OperationNoticeRecv{
		Status: def.StatusFailed,
	}
	discardNoticerecv := &pf.DiscardNoticeRecv{
		Status: def.StatusFailed,
	}
	endRoundRecv := &pf.EndRoundRecv{
		Status: def.StatusFailed,
	}
	settlementRecv := &pf.SettlementRecv{
		Status: def.StatusFailed,
	}
	closeRoomRecv := &pf.CloseRoomRecv{
		Status: def.StatusFailed,
	}
	nextID := int32(0)
	otherIDs := make([]int32, 0, def.RoomPlayerCount-1)
	multiHuIDs := make([]int32, 0, def.RoomPlayerCount-1)

	RoomCache.handleDiscard(
		playerID,
		roomID,
		send,
		discardRecv,
		drawCardRecv,
		operationNoticeRecv,
		discardNoticerecv,
		endRoundRecv,
		settlementRecv,
		closeRoomRecv,
		&nextID,
		&otherIDs,
		&multiHuIDs)

	// 出牌给自己发消息
	recv, err := discardRecv.Marshal()
	if err == nil {
		abs.MsgBody = recv
		recv, err = abs.Marshal()
		if err == nil {
			err = s.WriteBinary(recv)
			if err == nil {
				util.LogRecv(msgID, playerID, roomID, discardRecv, "Discard")
			}
		}
	}

	// 出牌给其他人发消息
	if discardRecv.Status == def.StatusOK {
		msgBody, err := discardRecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.Discard),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				for _, otherID := range otherIDs {
					if cache.SendMessage(otherID, msg) == nil {
						util.LogRecv(absMessage.GetMsgID(), otherID, roomID, discardRecv, "Discard")
					}
				}
			}
		}
	}

	// 下家抓牌给下家发消息
	if drawCardRecv.Status == def.StatusOK {
		msgBody, err := drawCardRecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.DrawCard),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				if cache.SendMessage(nextID, msg) == nil {
					util.LogRecv(absMessage.GetMsgID(), nextID, roomID, drawCardRecv, "DrawCard")
				}
			}
		}
	}

	// 下家抓牌给其他人发消息
	if drawCardRecv.Status == def.StatusOK {
		drawCardRecv.Card = 0
		msgBody, err := drawCardRecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.DrawCard),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				deskIDs := append(otherIDs, playerID)
				for _, deskID := range deskIDs {
					if deskID == nextID {
						continue
					}
					if cache.SendMessage(deskID, msg) == nil {
						util.LogRecv(absMessage.GetMsgID(), deskID, roomID, drawCardRecv, "DrawCard")
					}
				}
			}
		}
	}

	// 一炮多响
	if len(multiHuIDs) > 1 {
		// 给接炮人发消息
		if operationNoticeRecv.Status == def.StatusOK {
			operationList := operationNoticeRecv.OperationList
			for i, multiHuID := range multiHuIDs {
				operationNoticeRecv := &pf.OperationNoticeRecv{
					Status: def.StatusOK,
					OperationList: []*pf.OperationItem{
						operationList[2*i],
						operationList[2*i+1],
					},
				}
				msgBody, err := operationNoticeRecv.Marshal()
				if err == nil {
					absMessage := &pf.AbsMessage{
						MsgID:   int32(pf.OperationNotice),
						MsgBody: msgBody,
					}
					msg, err := absMessage.Marshal()
					if err == nil {
						if cache.SendMessage(multiHuID, msg) == nil {
							util.LogRecv(absMessage.GetMsgID(), multiHuID, roomID, operationNoticeRecv, "OperationNotice")
						}
					}
				}
			}
		}

		// 给其他人发等待消息
		if operationNoticeRecv.Status == def.StatusOK {
			watingNoticeRecv := &pf.OperationNoticeRecv{
				Status: def.StatusOK,
				OperationList: []*pf.OperationItem{
					{Type: def.Wating},
				},
			}
			msgBody, err := watingNoticeRecv.Marshal()
			if err == nil {
				absMessage := &pf.AbsMessage{
					MsgID:   int32(pf.OperationNotice),
					MsgBody: msgBody,
				}
				msg, err := absMessage.Marshal()
				if err == nil {
					deskIDs := append(otherIDs, playerID)
					for _, deskID := range deskIDs {
						isMultiHuID := false
						for _, multiHuID := range multiHuIDs {
							if deskID == multiHuID {
								isMultiHuID = true
								break
							}
						}
						if isMultiHuID {
							continue
						}

						if cache.SendMessage(deskID, msg) == nil {
							util.LogRecv(absMessage.GetMsgID(), deskID, roomID, watingNoticeRecv, "OperationNotice")
						}
					}
				}
			}
		}
	} else {
		// 给操作人发消息
		if operationNoticeRecv.Status == def.StatusOK {
			msgBody, err := operationNoticeRecv.Marshal()
			if err == nil {
				absMessage := &pf.AbsMessage{
					MsgID:   int32(pf.OperationNotice),
					MsgBody: msgBody,
				}
				msg, err := absMessage.Marshal()
				if err == nil {
					if cache.SendMessage(nextID, msg) == nil {
						util.LogRecv(absMessage.GetMsgID(), nextID, roomID, operationNoticeRecv, "OperationNotice")
					}
				}
			}
		}

		// 给其他人发等待消息
		if operationNoticeRecv.Status == def.StatusOK {
			watingNoticeRecv := &pf.OperationNoticeRecv{
				Status: def.StatusOK,
				OperationList: []*pf.OperationItem{
					{Type: def.Wating},
				},
			}
			msgBody, err := watingNoticeRecv.Marshal()
			if err == nil {
				absMessage := &pf.AbsMessage{
					MsgID:   int32(pf.OperationNotice),
					MsgBody: msgBody,
				}
				msg, err := absMessage.Marshal()
				if err == nil {
					deskIDs := append(otherIDs, playerID)
					for _, deskID := range deskIDs {
						if deskID == nextID {
							continue
						}
						if cache.SendMessage(deskID, msg) == nil {
							util.LogRecv(absMessage.GetMsgID(), deskID, roomID, watingNoticeRecv, "OperationNotice")
						}
					}
				}
			}
		}
	}

	// 出牌提示给桌上人发消息
	if discardNoticerecv.Status == def.StatusOK {
		msgBody, err := discardNoticerecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.DiscardNotice),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				deskIDs := append(otherIDs, playerID)
				for _, deskID := range deskIDs {
					if cache.SendMessage(deskID, msg) == nil {
						util.LogRecv(absMessage.GetMsgID(), deskID, roomID, discardNoticerecv, "DiscardNotice")
					}
				}
			}
		}
	}

	// 小结算给桌上人发消息
	if endRoundRecv.Status == def.StatusOK {
		msgBody, err := endRoundRecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.EndRound),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				deskIDs := append(otherIDs, playerID)
				for _, deskID := range deskIDs {
					if cache.SendMessage(deskID, msg) == nil {
						util.LogRecv(absMessage.GetMsgID(), deskID, roomID, endRoundRecv, "EndRound")
					}
				}
			}
		}
	}

	// 大结算给桌上人发消息
	if settlementRecv.Status == def.StatusOK {
		msgBody, err := settlementRecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.Settlement),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				deskIDs := append(otherIDs, playerID)
				for _, deskID := range deskIDs {
					if cache.SendMessage(deskID, msg) == nil {
						util.LogRecv(absMessage.GetMsgID(), deskID, roomID, settlementRecv, "Settlement")
					}
				}
			}
		}
	}

	// 房间关闭给桌上人发消息
	if closeRoomRecv.Status == def.StatusOK {
		msgBody, err := closeRoomRecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.CloseRoom),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				deskIDs := append(otherIDs, playerID)
				for _, deskID := range deskIDs {
					if cache.SendMessage(deskID, msg) == nil {
						util.LogRecv(absMessage.GetMsgID(), deskID, roomID, closeRoomRecv, "CloseRoom")
					}
				}
			}
		}
	}
	return
}

func HandleOperation(s *melody.Session, mjPlayer *cache.MjPlayer, send *pf.OperationSend, abs *pf.AbsMessage) (err error) {
	playerID := mjPlayer.PlayerID
	roomID := mjPlayer.RoomID
	msgID := abs.MsgID

	operationRecv := &pf.OperationRecv{
		Status: def.StatusOK,
	}
	drawCardRecv := &pf.DrawCardRecv{
		Status: def.StatusFailed,
	}
	operationNoticeRecv := &pf.OperationNoticeRecv{
		Status: def.StatusFailed,
	}
	discardNoticerecv := &pf.DiscardNoticeRecv{
		Status: def.StatusFailed,
	}
	endRoundRecv := &pf.EndRoundRecv{
		Status: def.StatusFailed,
	}
	settlementRecv := &pf.SettlementRecv{
		Status: def.StatusFailed,
	}
	closeRoomRecv := &pf.CloseRoomRecv{
		Status: def.StatusFailed,
	}
	nextID := int32(0)
	otherIDs := make([]int32, 0, def.RoomPlayerCount-1)

	RoomCache.handleOperation(
		playerID,
		roomID,
		send,
		operationRecv,
		drawCardRecv,
		operationNoticeRecv,
		discardNoticerecv,
		endRoundRecv,
		settlementRecv,
		closeRoomRecv,
		&nextID,
		&otherIDs)

	// 操作给自己发消息
	// 一炮多响特殊处理
	if operationRecv.Status == def.StatusWait {
		watingNoticeRecv := &pf.OperationNoticeRecv{
			Status: def.StatusOK,
			OperationList: []*pf.OperationItem{
				{Type: def.Wating},
			},
		}
		msgBody, err := watingNoticeRecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.OperationNotice),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				err = s.WriteBinary(msg)
				if err == nil {
					util.LogRecv(absMessage.GetMsgID(), playerID, roomID, watingNoticeRecv, "OperationNotice")
				}
			}
		}
	} else {
		recv, err := operationRecv.Marshal()
		if err == nil {
			abs.MsgBody = recv
			recv, err = abs.Marshal()
			if err == nil {
				err = s.WriteBinary(recv)
				if err == nil {
					util.LogRecv(msgID, playerID, roomID, operationRecv, "Operation")
				}
			}
		}
	}

	// 操作给其他人发消息
	if operationRecv.Status == def.StatusOK {
		msgBody, err := operationRecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.Operation),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				for _, otherID := range otherIDs {
					if cache.SendMessage(otherID, msg) == nil {
						util.LogRecv(absMessage.GetMsgID(), otherID, roomID, operationRecv, "Operation")
					}
				}
			}
		}
	}

	// 操作者抓牌给操作者发消息
	if drawCardRecv.Status == def.StatusOK {
		msgBody, err := drawCardRecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.DrawCard),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				if cache.SendMessage(nextID, msg) == nil {
					util.LogRecv(absMessage.GetMsgID(), playerID, roomID, drawCardRecv, "DrawCard")
				}
			}
		}
	}

	// 操作者抓牌给其他人发消息
	if drawCardRecv.Status == def.StatusOK {
		drawCardRecv.Card = 0
		msgBody, err := drawCardRecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.DrawCard),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				deskIDs := append(otherIDs, playerID)
				for _, deskID := range deskIDs {
					if deskID == nextID {
						continue
					}
					if cache.SendMessage(deskID, msg) == nil {
						util.LogRecv(absMessage.GetMsgID(), deskID, roomID, drawCardRecv, "DrawCard")
					}
				}
			}
		}
	}

	// 操作者操作提示给操作者发消息
	if operationNoticeRecv.Status == def.StatusOK {
		msgBody, err := operationNoticeRecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.OperationNotice),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				if cache.SendMessage(nextID, msg) == nil {
					util.LogRecv(absMessage.GetMsgID(), playerID, roomID, operationNoticeRecv, "OperationNotice")
				}
			}
		}
	}

	// 操作者出牌提示给桌上人发消息
	if discardNoticerecv.Status == def.StatusOK {
		msgBody, err := discardNoticerecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.DiscardNotice),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				deskIDs := append(otherIDs, playerID)
				for _, deskID := range deskIDs {
					if cache.SendMessage(deskID, msg) == nil {
						util.LogRecv(absMessage.GetMsgID(), deskID, roomID, discardNoticerecv, "DiscardNotice")
					}
				}
			}
		}
	}

	// 小结算给桌上人发消息
	if endRoundRecv.Status == def.StatusOK {
		msgBody, err := endRoundRecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.EndRound),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				deskIDs := append(otherIDs, playerID)
				for _, deskID := range deskIDs {
					if cache.SendMessage(deskID, msg) == nil {
						util.LogRecv(absMessage.GetMsgID(), deskID, roomID, endRoundRecv, "EndRound")
					}
				}
			}
		}
	}

	// 大结算给桌上人发消息
	if settlementRecv.Status == def.StatusOK {
		msgBody, err := settlementRecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.Settlement),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				deskIDs := append(otherIDs, playerID)
				for _, deskID := range deskIDs {
					if cache.SendMessage(deskID, msg) == nil {
						util.LogRecv(absMessage.GetMsgID(), deskID, roomID, settlementRecv, "Settlement")
					}
				}
			}
		}
	}

	// 房间关闭给桌上人发消息
	if closeRoomRecv.Status == def.StatusOK {
		msgBody, err := closeRoomRecv.Marshal()
		if err == nil {
			absMessage := &pf.AbsMessage{
				MsgID:   int32(pf.CloseRoom),
				MsgBody: msgBody,
			}
			msg, err := absMessage.Marshal()
			if err == nil {
				deskIDs := append(otherIDs, playerID)
				for _, deskID := range deskIDs {
					if cache.SendMessage(deskID, msg) == nil {
						util.LogRecv(absMessage.GetMsgID(), deskID, roomID, closeRoomRecv, "CloseRoom")
					}
				}
			}
		}
	}
	return
}
