package cache

import (
	"qnmahjong/db"
	"qnmahjong/db/dao"
	"sync"
)

// ModuleMap hold tool module info(name->module)
type ModuleMap struct {
	sync.RWMutex
	mmap map[string]*dao.Module
	gmap map[int32]*dao.Module
}

var (
	moduleCache ModuleMap
)

func init() {
	moduleCache = ModuleMap{}
}

// InitModule read tool module info from db
func InitModule() {
	moduleCache.Lock()
	defer moduleCache.Unlock()

	modules, err := dao.ModulesFromModule(db.Pool)
	if err != nil {
		return
	}

	moduleCache.mmap = make(map[string]*dao.Module, len(modules))
	moduleCache.gmap = make(map[int32]*dao.Module, len(modules))
	for _, module := range modules {
		daoModule, err := dao.ModuleByModule(db.Pool, module)
		if err != nil {
			continue
		}

		moduleCache.mmap[module] = daoModule
		moduleCache.gmap[daoModule.IndexID] = daoModule
	}
}

// GetGMTModule get module by index_id
func GetGMTModule(IndexID int32) *dao.Module {
	moduleCache.RLock()
	defer moduleCache.RUnlock()

	return moduleCache.gmap[IndexID]
}
