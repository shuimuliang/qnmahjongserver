package pay

import (
	"encoding/json"
	"qnmahjong/cache"
	"qnmahjong/db"
	"qnmahjong/def"
	"qnmahjong/pf"
	"qnmahjong/util"

	log "github.com/Sirupsen/logrus"
)

type AppStoreSend struct {
	TransactionId string `json:"transaction_id,omitempty"`
	ReceiptData   string `json:"receipt-data,omitempty"`
}

func ValidateAppStore(mjPlayer *cache.MjPlayer, orderData *pf.OrderData) *pf.OrderStatus {
	send := AppStoreSend{}
	err := json.Unmarshal([]byte(orderData.JsonStr), &send)
	if err != nil {
		return nil
	}

	receipt, err := VerifyReceipt(send.ReceiptData, false)
	goiapErr, ok := err.(ErrorWithCode)

	if ok && goiapErr.Code() == SandboxReceiptOnProd {
		receipt, err = VerifyReceipt(send.ReceiptData, true)
	}

	if err != nil {
		return nil
	}

	log.WithFields(log.Fields{
		"receipt":   *receipt,
		"mjPlayer":  *mjPlayer,
		"orderData": *orderData,
	}).Info("ValidateAppStore")

	order := cache.GetOrderByProductID(mjPlayer.Channel, mjPlayer.PlayerID, def.PayTypeAppStore, receipt.ProductId)
	if order == nil {
		util.LogError(err, "order", order, mjPlayer.PlayerID, def.ErrCreateDBOrder)
		return nil
	}

	order.TransID = receipt.TransactionId
	order.Status = def.IpayStatusSuccess
	err = order.Insert(db.Pool)
	if err != nil {
		util.LogError(err, "order", order, mjPlayer.PlayerID, def.ErrInsertOrder)
		return nil
	}

	ok = cache.AddCardsAppStore(mjPlayer.PlayerID, order.GoodsCount+order.ExtraCount)
	if !ok {
		return nil
	}

	order.Status = def.IpayStatusComplete
	err = order.Update(db.Pool)
	if err != nil {
		util.LogError(err, "order", order, mjPlayer.PlayerID, def.ErrUpdateOrder)
		return nil
	}

	return &pf.OrderStatus{
		OrderID:   order.TransID,
		Status:    0,
		ErrorDesc: "",
		JsonStr:   "",
		GemID:     order.GemID,
	}
}
