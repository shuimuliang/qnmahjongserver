package cache

import (
	"qnmahjong/db"
	"qnmahjong/db/dao"
	"qnmahjong/def"
	"qnmahjong/pf"
	"qnmahjong/util"
	"strings"
	"sync"
	"time"

	"github.com/olahol/melody"
)

// MjPlayer hold mjplayer info
type MjPlayer struct {
	PlayerID  int32 // 玩家id
	RoomID    int32 // 房间id
	MJType    int32 // 麻将类型
	Channel   int32 // 渠道
	Version   int32 // 版本
	LoginType int32 // 登录方式
}

// Player hold player info
type Player struct {
	*dao.Player
	*melody.Session
	RoomID    int32   // 房间id
	MJType    int32   // 麻将类型
	Channel   int32   // 渠道
	Version   int32   // 版本
	LoginType int32   // 登录方式
	LowID     []int32 // 下线id
	IP        string  // ip 地址
	Latitude  float64 // 纬度
	Longitude float64 // 经度
}

// PlayerMap hold tool player info(playerid->player)
type PlayerMap struct {
	sync.RWMutex
	mmap map[int32]*Player
}

var (
	playerCache PlayerMap
)

func init() {
	playerCache = PlayerMap{}
}

// InitPlayer read tool player info from db
func InitPlayer() {
	playerCache.Lock()
	defer playerCache.Unlock()

	playerIDs, err := dao.PlayerIDsFromPlayer(db.Pool)
	if err != nil {
		return
	}

	playerCache.mmap = make(map[int32]*Player, len(playerIDs))
	for _, playerID := range playerIDs {
		daoPlayer, err := dao.PlayerByPlayerID(db.Pool, playerID)
		if err != nil {
			continue
		}

		player := &Player{
			Player: daoPlayer,
		}
		playerCache.mmap[playerID] = player
	}

	for _, player := range playerCache.mmap {
		if player.HighID == 0 {
			continue
		}

		highPlayer, ok := playerCache.mmap[player.HighID]
		if !ok {
			continue
		}

		highPlayer.LowID = append(highPlayer.LowID, player.PlayerID)
	}
}

// GetMjPlayer get mjplayer by playerid
func GetMjPlayer(playerID int32) *MjPlayer {
	playerCache.RLock()
	defer playerCache.RUnlock()

	player, ok := playerCache.mmap[playerID]
	if !ok {
		return nil
	}

	return &MjPlayer{
		PlayerID:  player.PlayerID,
		RoomID:    player.RoomID,
		MJType:    player.MJType,
		Channel:   player.Channel,
		Version:   player.Version,
		LoginType: player.LoginType,
	}
}

// GetPfPlayer get pfplayer by playerid
func GetPfPlayer(playerID int32) *pf.Player {
	playerCache.RLock()
	defer playerCache.RUnlock()

	player, ok := playerCache.mmap[playerID]
	if !ok {
		return nil
	}

	return &pf.Player{
		Id:       player.PlayerID,
		Nickname: player.Nickname,
		Avatar:   strings.TrimRight(player.Headimgurl, "0") + "132",
		Gender:   player.Sex,
		Coins:    player.Coins,
		Cards:    player.Cards,
		Ip:       player.IP,
	}
}

func GetRoomMJType(playerID int32) (int32, int32) {
	playerCache.RLock()
	defer playerCache.RUnlock()

	player, ok := playerCache.mmap[playerID]
	if !ok {
		return 0, 0
	}

	return player.RoomID, player.MJType
}

func SetRoomID(playerID, roomID, mjType int32) {
	playerCache.Lock()
	defer playerCache.Unlock()

	player, ok := playerCache.mmap[playerID]
	if !ok {
		return
	}

	player.RoomID = roomID
	player.MJType = mjType
}

func SetLocation(playerID int32, latitude, longitude float64) {
	playerCache.Lock()
	defer playerCache.Unlock()

	player, ok := playerCache.mmap[playerID]
	if !ok {
		return
	}

	player.Latitude = latitude
	player.Longitude = longitude
}

func CheckCards(playerID, cost int32) bool {
	playerCache.RLock()
	defer playerCache.RUnlock()

	player, ok := playerCache.mmap[playerID]
	if !ok {
		return false
	}

	return player.Cards >= cost
}

func CheckCoins(playerID, cost int32) (ok bool) {
	playerCache.RLock()
	defer playerCache.RUnlock()

	player, ok := playerCache.mmap[playerID]
	if !ok {
		return false
	}

	return player.Coins >= cost
}

func SendMessage(playerID int32, msg []byte) error {
	playerCache.RLock()
	defer playerCache.RUnlock()

	player, ok := playerCache.mmap[playerID]
	if !ok {
		return def.ErrHandleLogic
	}

	return player.Session.WriteBinary(msg)
}

func ExistPlayer(playerID int32) bool {
	playerCache.RLock()
	defer playerCache.RUnlock()

	_, ok := playerCache.mmap[playerID]
	return ok
}

func AddCardsAppStore(playerID, addCards int32) bool {
	playerCache.Lock()
	defer playerCache.Unlock()

	player, ok := playerCache.mmap[playerID]
	if !ok {
		return false
	}

	player.Cards += addCards
	err := player.Update(db.Pool)
	if err != nil {
		util.LogError(err, "player", player, playerID, def.ErrUpdatePlayer)
		return false
	}

	t := &dao.Treasure{
		PlayerID:   playerID,
		Reason:     def.AppStoreOrder,
		Coins:      0,
		Cards:      addCards,
		ChangeTime: time.Now(),
	}
	err = t.Insert(db.Pool)
	if err != nil {
		util.LogError(err, "treasure", t, playerID, def.ErrInsertTreasure)
	}

	return true
}

// IsIPConflict 判断2位玩家IP是否冲突
func IsIPConflict(playerID1, playerID2 int32) bool {
	playerCache.RLock()
	defer playerCache.RUnlock()

	ip1 := playerCache.mmap[playerID1].IP
	ip2 := playerCache.mmap[playerID2].IP
	return ip1 == ip2
}

// IsGeoConflict 判断2位玩家IP是否冲突
func IsGeoConflict(playerID1, playerID2 int32) bool {
	playerCache.RLock()
	defer playerCache.RUnlock()

	player1 := playerCache.mmap[playerID1]
	player2 := playerCache.mmap[playerID2]

	geo1 := &util.Point{
		Lat: player1.Latitude,
		Lng: player1.Longitude,
	}

	geo2 := &util.Point{
		Lat: player2.Latitude,
		Lng: player2.Longitude,
	}

	distance := geo1.GreatCircleDistance(geo2)
	return distance <= def.GPSDistanceLimit
}

func GetGPSLocation(playerID int32) (float64, float64) {
	playerCache.RLock()
	defer playerCache.RUnlock()

	player, ok := playerCache.mmap[playerID]
	if !ok {
		return 0, 0
	}

	return player.Latitude, player.Longitude
}
