package cache

import (
	"qnmahjong/db"
	"qnmahjong/db/dao"
	"sync"
)

// RoleMap hold tool role info(name->role)
type RoleMap struct {
	sync.RWMutex
	mmap map[string]*dao.Role
	gmap map[int32]*dao.Role
}

var (
	roleCache RoleMap
)

func init() {
	roleCache = RoleMap{}
}

// InitRole read tool role info from db
func InitRole() {
	roleCache.Lock()
	defer roleCache.Unlock()

	roles, err := dao.RolesFromRole(db.Pool)
	if err != nil {
		return
	}

	roleCache.mmap = make(map[string]*dao.Role, len(roles))
	roleCache.gmap = make(map[int32]*dao.Role, len(roles))
	for _, role := range roles {
		daoRole, err := dao.RoleByRole(db.Pool, role)
		if err != nil {
			continue
		}

		roleCache.mmap[role] = daoRole
		roleCache.gmap[daoRole.IndexID] = daoRole
	}
}

// GetGMTRole get role by index_id
func GetGMTRole(IndexID int32) *dao.Role {
	roleCache.RLock()
	defer roleCache.RUnlock()

	return roleCache.gmap[IndexID]
}
