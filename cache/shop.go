package cache

import (
	"qnmahjong/db"
	"qnmahjong/db/dao"
	"qnmahjong/def"
	"qnmahjong/pf"
	"strconv"
	"sync"
	"time"
)

// ShopMap hold mj shop info(channel->waresid->shop)
type ShopMap struct {
	sync.RWMutex
	mmap map[int32]map[string]*dao.Shop
	gmap map[int32]*dao.Shop
}

var (
	shopCache ShopMap
)

func init() {
	shopCache = ShopMap{}
}

// InitShop read mj shop info from db
func InitShop() {
	shopCache.Lock()
	defer shopCache.Unlock()

	channels, err := dao.ChannelsFromShop(db.Pool)
	if err != nil {
		return
	}

	shopCache.mmap = make(map[int32]map[string]*dao.Shop, len(channels))
	shopCache.gmap = make(map[int32]*dao.Shop, len(channels))
	for _, channel := range channels {
		shops, err := dao.ShopsByChannel(db.Pool, channel)
		if err != nil {
			continue
		}

		shopCache.mmap[channel] = make(map[string]*dao.Shop, len(shops))
		for _, shop := range shops {
			shopCache.mmap[channel][shop.WaresID] = shop
			shopCache.gmap[shop.IndexID] = shop
		}
	}
}

// GetGoods get goods by channel
func GetGoods(channel int32) []*pf.Goods {
	shopCache.RLock()
	defer shopCache.RUnlock()

	shops, ok := shopCache.mmap[channel]
	if !ok {
		return nil
	}

	var goods []*pf.Goods
	for _, shop := range shops {
		goods = append(goods, &pf.Goods{
			Id:        shop.GemID,
			Count:     shop.GoodsCount,
			Extra:     shop.ExtraCount,
			Price:     shop.Price,
			IconUrl:   shop.IconURL,
			PayType:   shop.PayType,
			ProductID: shop.WaresID,
		})
	}
	return goods
}

// GetOrder get goods by channel and waresID
func GetOrderWaresID(channel, playerID, payType int32, waresID string) *dao.Order {
	shopCache.RLock()
	defer shopCache.RUnlock()

	shop, ok := shopCache.mmap[channel][waresID]
	if !ok {
		return nil
	}

	curTime := time.Now()
	orderID := strconv.Itoa(int(playerID)) + strconv.Itoa(int(curTime.UnixNano()))
	order := &dao.Order{
		OrderID:    orderID,
		PlayerID:   playerID,
		Channel:    shop.Channel,
		GemID:      shop.GemID,
		PayType:    payType,
		WaresID:    shop.WaresID,
		WaresName:  shop.WaresName,
		GoodsCount: shop.GoodsCount,
		ExtraCount: shop.ExtraCount,
		Price:      shop.Price,
		Status:     def.IpayStatusNotStart,
		AddTime:    curTime,
		ReviseTime: curTime,
	}
	return order
}

// GetOrderByProductID get goods by channel and productID
func GetOrderByProductID(channel, playerID, payType int32, productID string) *dao.Order {
	shopCache.RLock()
	defer shopCache.RUnlock()

	shop, ok := shopCache.mmap[channel][productID]
	if !ok {
		return nil
	}

	curTime := time.Now()
	orderID := strconv.Itoa(int(playerID)) + strconv.Itoa(int(curTime.UnixNano()))
	order := &dao.Order{
		OrderID:    orderID,
		PlayerID:   playerID,
		Channel:    shop.Channel,
		GemID:      shop.GemID,
		PayType:    payType,
		WaresID:    shop.WaresID,
		WaresName:  shop.WaresName,
		GoodsCount: shop.GoodsCount,
		ExtraCount: shop.ExtraCount,
		Price:      shop.Price,
		Status:     def.IpayStatusNotStart,
		AddTime:    curTime,
		ReviseTime: curTime,
	}
	return order
}

// GetGMTShop get shop by index_id
func GetGMTShop(IndexID int32) *dao.Shop {
	shopCache.RLock()
	defer shopCache.RUnlock()

	return shopCache.gmap[IndexID]
}
