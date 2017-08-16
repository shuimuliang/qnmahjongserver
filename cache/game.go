package cache

import (
	"qnmahjong/db"
	"qnmahjong/db/dao"
	"qnmahjong/def"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// VersionInfo hold info for current version
type VersionInfo struct {
	URL        string `json:"url"`
	UpdateType int32  `json:"updateType"`
	NumVersion int32  `json:"numVersion"`
	FileSize   int32  `json:"fileSize"`
}

// VersionCheck hold client version check info
type VersionCheck struct {
	Module      string        `json:"module"`
	VersionList []VersionInfo `json:"versionList"`
}

// GameMap hold game info(channel->version->game)
type GameMap struct {
	sync.RWMutex
	mmap map[int32]map[int32]*dao.Game
	gmap map[int32]*dao.Game
}

var (
	gameCache GameMap
)

func init() {
	gameCache = GameMap{}
}

// InitGame read mj game version info from db
func InitGame() {
	gameCache.Lock()
	defer gameCache.Unlock()

	channels, err := dao.ChannelsFromGame(db.Pool)
	if err != nil {
		return
	}

	gameCache.mmap = make(map[int32]map[int32]*dao.Game, len(channels))
	gameCache.gmap = make(map[int32]*dao.Game, len(channels))
	for _, channel := range channels {
		games, err := dao.GamesByChannel(db.Pool, channel)
		if err != nil {
			continue
		}

		gameCache.mmap[channel] = make(map[int32]*dao.Game, len(games))
		for _, game := range games {
			gameCache.mmap[channel][game.Version] = game
			gameCache.gmap[game.IndexID] = game
		}
	}
}

// GetVersionCheck get client check version info
func GetVersionCheck(channel, version int32) VersionCheck {
	gameCache.RLock()
	defer gameCache.RUnlock()

	versionCheck := VersionCheck{
		Module:      getModules(channel, version),
		VersionList: getVersions(channel, version),
	}
	return versionCheck
}

// GetMjTypes get mj types by channel version
func GetMjTypes(channel, version int32) []int32 {
	gameCache.RLock()
	defer gameCache.RUnlock()

	game, err := dao.GameByChannelVersion(db.Pool, channel, version)
	if err != nil {
		return nil
	}

	var mjTypes []int32
	mjTypesStr := strings.Split(game.MjTypes, "|")
	for _, mjTypeStr := range mjTypesStr {
		mjType, err := strconv.Atoi(mjTypeStr)
		if err != nil {
			continue
		}

		mjTypes = append(mjTypes, int32(mjType))
	}
	return mjTypes
}

func getModules(channel, version int32) string {
	games, ok := gameCache.mmap[channel]
	if !ok {
		return ""
	}

	game, ok := games[version]
	if !ok {
		return ""
	}

	return game.Module
}

func getVersions(channel, version int32) []VersionInfo {
	games, ok := gameCache.mmap[channel]
	if !ok {
		return nil
	}

	var versionList []VersionInfo
	for _, game := range games {
		if game.Version > version && game.Enabled == def.VersionEnabled {
			versionInfo := VersionInfo{
				URL:        game.DownloadURL,
				UpdateType: game.UpdateType,
				NumVersion: game.Version,
				FileSize:   game.Size,
			}
			if game.UpdateType == def.ForceUpdate {
				return []VersionInfo{versionInfo}
			}

			versionList = append(versionList, versionInfo)
		}
	}

	sort.Slice(versionList, func(i, j int) bool {
		return versionList[i].NumVersion <= versionList[j].NumVersion
	})
	return versionList
}

// GetGMTGame get game by index_id
func GetGMTGame(IndexID int32) *dao.Game {
	gameCache.RLock()
	defer gameCache.RUnlock()

	return gameCache.gmap[IndexID]
}
