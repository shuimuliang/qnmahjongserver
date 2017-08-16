package cache

import (
	"qnmahjong/db"
	"qnmahjong/db/dao"
	"sync"
)

// PermissionMap hold tool permission info(name->permission)
type PermissionMap struct {
	sync.RWMutex
	mmap map[string]map[string]*dao.Permission
	gmap map[int32]*dao.Permission
}

var (
	permissionCache PermissionMap
)

func init() {
	permissionCache = PermissionMap{}
}

// InitPermission read tool permission info from db
func InitPermission() {
	permissionCache.Lock()
	defer permissionCache.Unlock()

	pmsnTypes, err := dao.PmsnTypesFromPermission(db.Pool)
	if err != nil {
		return
	}

	permissionCache.mmap = make(map[string]map[string]*dao.Permission, len(pmsnTypes))
	permissionCache.gmap = make(map[int32]*dao.Permission, len(pmsnTypes))
	for _, pmsnType := range pmsnTypes {
		permissions, err := dao.PermissionsByPmsnType(db.Pool, pmsnType)
		if err != nil {
			continue
		}

		permissionCache.mmap[pmsnType] = make(map[string]*dao.Permission, len(permissions))
		for _, permission := range permissions {
			permissionCache.mmap[pmsnType][permission.PmsnContent] = permission
			permissionCache.gmap[permission.IndexID] = permission
		}
	}
}

// GetGMTPermission get permission by index_id
func GetGMTPermission(IndexID int32) *dao.Permission {
	permissionCache.RLock()
	defer permissionCache.RUnlock()

	return permissionCache.gmap[IndexID]
}
