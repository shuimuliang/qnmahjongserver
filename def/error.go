package def

import (
	"errors"
)

// login error
var (
	ErrStartLogin    = errors.New("start login error")
	ErrShutdownLogin = errors.New("shutdown login error")

	ErrHandleLogin = errors.New("handle login error")
)

// logic error
var (
	ErrStartLogic    = errors.New("start logic error")
	ErrShutdownLogic = errors.New("shutdown logic error")

	ErrHandleLogic = errors.New("handle logic error")
)

// tool error
var (
	ErrStartTool    = errors.New("start tool error")
	ErrShutdownTool = errors.New("shutdown tool error")
)

// pay error
var (
	ErrStartPay    = errors.New("start pay error")
	ErrShutdownPay = errors.New("shutdown pay error")
)

// sale error
var (
	ErrStartSale    = errors.New("start sale error")
	ErrShutdownSale = errors.New("shutdown sale error")
)

// notice error
var (
	ErrStartLoginNotice = errors.New("start login notice error")
	ErrStartLogicNotice = errors.New("start logic notice error")
	ErrStartToolNotice   = errors.New("start tool notice error")
	ErrStartPayNotice   = errors.New("start pay notice error")
	ErrStartSaleNotice = errors.New("start sale notice error")

	ErrIpayOrderParm       = errors.New("ipay order parm error")
	ErrCreateDBOrder       = errors.New("create db order error")
	ErrCreateIpayOrder     = errors.New("create ipay order error")
	ErrGetHTML5RedirectURL = errors.New("get html5 redirect url error")
	ErrParseIpayNotifyInfo = errors.New("parse ipay notify info error")
	ErrQueryIpayResult     = errors.New("query ipay result error")

	ErrDialNoticeRPC = errors.New("dial rpc notice error")
	ErrCallNoticeRPC = errors.New("call rpc notice error")
)

// db error
var (
	ErrOpenDb  = errors.New("open db error")
	ErrCloseDb = errors.New("close db error")

	ErrInsertAccount    = errors.New("insert account error")
	ErrInsertChat       = errors.New("insert chat error")
	ErrInsertCost       = errors.New("insert cost error")
	ErrInsertFeedback   = errors.New("insert feedback error")
	ErrInsertGame       = errors.New("insert game error")
	ErrInsertLogin      = errors.New("insert login error")
	ErrInsertModule     = errors.New("insert module error")
	ErrInsertOrder      = errors.New("insert order error")
	ErrInsertPermission = errors.New("insert permission error")
	ErrInsertPlayer     = errors.New("insert player error")
	ErrInsertRecord     = errors.New("insert record error")
	ErrInsertRegister   = errors.New("insert register error")
	ErrInsertRole       = errors.New("insert role error")
	ErrInsertShop       = errors.New("insert shop error")
	ErrInsertTreasure   = errors.New("insert treasure error")

	ErrUpdateAccount    = errors.New("update account error")
	ErrUpdateChat       = errors.New("update chat error")
	ErrUpdateCost       = errors.New("update cost error")
	ErrUpdateFeedback   = errors.New("update feedback error")
	ErrUpdateGame       = errors.New("update game error")
	ErrUpdateLogin      = errors.New("update login error")
	ErrUpdateModule     = errors.New("update module error")
	ErrUpdateOrder      = errors.New("update order error")
	ErrUpdatePermission = errors.New("update permission error")
	ErrUpdatePlayer     = errors.New("update player error")
	ErrUpdateRecord     = errors.New("update record error")
	ErrUpdateRegister   = errors.New("update register error")
	ErrUpdateRole       = errors.New("update role error")
	ErrUpdateShop       = errors.New("update shop error")
	ErrUpdateTreasure   = errors.New("update treasure error")

	ErrDeleteAccount    = errors.New("delete account error")
	ErrDeleteChat       = errors.New("delete chat error")
	ErrDeleteCost       = errors.New("delete cost error")
	ErrDeleteFeedback   = errors.New("delete feedback error")
	ErrDeleteGame       = errors.New("delete game error")
	ErrDeleteLogin      = errors.New("delete login error")
	ErrDeleteModule     = errors.New("delete module error")
	ErrDeleteOrder      = errors.New("delete order error")
	ErrDeletePermission = errors.New("delete permission error")
	ErrDeletePlayer     = errors.New("delete player error")
	ErrDeleteRecord     = errors.New("delete record error")
	ErrDeleteRegister   = errors.New("delete register error")
	ErrDeleteRole       = errors.New("delete role error")
	ErrDeleteShop       = errors.New("delete shop error")
	ErrDeleteTreasure   = errors.New("delete treasure error")

	ErrInsertAgAuth    = errors.New("insert ag_auth error")
	ErrInsertAgAccount = errors.New("insert ag_aaacount error")
	ErrInsertAgBill    = errors.New("insert ag_bill error")
	ErrInsertAgPay     = errors.New("insert ag_pay error")

	ErrUpdateAgAuth    = errors.New("update ag_auth error")
	ErrUpdateAgAccount = errors.New("update ag_aaacount error")
	ErrUpdateAgBill    = errors.New("update ag_bill error")
	ErrUpdateAgPay     = errors.New("update ag_pay error")

	ErrDeleteAgAuth    = errors.New("delete ag_auth error")
	ErrDeleteAgAccount = errors.New("delete ag_aaacount error")
	ErrDeleteAgBill    = errors.New("delete ag_bill error")
	ErrDeleteAgPay     = errors.New("delete ag_pay error")
)

// log error
var (
	ErrLogToFile = errors.New("create/open log file error")
)

// conf error
var (
	ErrReadConfFile = errors.New("read conf file error")
)

// cron error
var (
	ErrAddCronFunc = errors.New("add cron func error")
)

// redis error
var (
	ErrDialRedis  = errors.New("dial redis error")
	ErrCloseRedis = errors.New("close redis error")

	ErrPutRecordToRedis    = errors.New("put record to redis error")
	ErrGetRecordsFromRedis = errors.New("get records from redis error")

	ErrMarshalRedisRecord   = errors.New("marshal redis record error")
	ErrUnmarshalRedisRecord = errors.New("unmarshal redis record error")

	ErrParseRecordCreateTime = errors.New("parse record createtime error")
	ErrTrimRecordToRedis     = errors.New("trim record to redis error")
)

// token error
var (
	ErrCreateToken   = errors.New("create token error")
	ErrParseToken    = errors.New("parse token error")
	ErrValidateToken = errors.New("validate token error")
)

// weixin error
var (
	ErrWeixinLogin = errors.New("login with weixin error")
)
