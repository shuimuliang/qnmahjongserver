package def

// 登录方式
const (
	LoginTypeWX = 1
)

// 运营微信号
const (
	WeixinWchatID = "poker888mm"
)

// 微信相关配置
const (
	WeixinAppIDHB  = "wxe27fe10d1a58989e"
	WeixinSecretHB = "be7a47c92cfcef96202dfc885baa52c4"
)

// 微信相关配置
const (
	WeixinAppIDHN  = "wx123e1af1dca25886"
	WeixinSecretHN = "bbb74aa992c0a30faaab263bc4d6be81"
)

// 登录渠道
const (
	ChannelIOSHB = 100
	ChannelIOSHN = 110

	ChannelAndroidHB = 200
	ChannelAndroidHN = 210
)

// 牌桌人数
const (
	RoomPlayerCount = 4
)

// 开局发牌数目
const (
	DiscardCount = 13
)

// 初始玩家数目
const (
	InitPlayerCount = 1000
)

// 初始房间数目
const (
	InitRoomCount = 1000
)

// 1 碰 2 胡 3 吃 4 明杠 5 暗杠 6 过 7 一炮多响 8 碰牌再杠 9 等别人操作
const (
	Peng      = 1
	Hu        = 2
	Chi       = 3
	MingGang  = 4
	AnGang    = 5
	Pass      = 6
	MultiHu   = 7
	BuGang    = 8
	Wating    = 9
	ZimoHu    = 10
	DianpaoHu = 11
)

// 吃的类型
const (
	ChiType1 = 1 // 3 4 *
	ChiType2 = 2 // 3 * 5
	ChiType3 = 3 // * 3 4
)

// 邀请奖励是否领取
const (
	InviteAwardAvailable   = 0 // 未领
	InviteAwardUnavailable = 1 // 已领
)

// 添加资产
const (
	RegisterAward       = 60  // 新用户注册
	InviteCodeBindAward = 150 // 绑定上线
	InviteCompleteAward = 60  // 绑定上线成功奖励
)

// 增减资产的理由
const (
	Register       = 1 // 新用户注册
	InviteCodeBind = 2 // 绑定上线
	InviteComplete = 3 // 绑定上线成功奖励

	CreateRoom    = 4 // 开房间
	IPayOrder     = 5 // 爱贝支付订单
	AppStoreOrder = 6 // appstore支付订单
)

// 投票状态
const (
	VoteDefault  = 0 // 未投票
	VoteAgree    = 1 // 同意
	VoteDisagree = 2 // 不同意
)

// 投票倒计时
const (
	VoteLeftTime = 60
)

// 开房资产
const (
	CostCoin = 1
	CostCard = 2
)

// 解散房间
const (
	CloseRoomRoundOver  = 0 // 牌局正常结束
	CloseRoomBeforStart = 1 // 开始游戏前房主强制解散房间
	CloseRoomVoteAgree  = 2 // 投票通过解散房间
)

// 牌局状态
const (
	GameStatusIsGaming  = 1 // 牌局进行中
	GameStatusRoundOver = 2 // 一局结束了
	GameStatusGameOver  = 3 // 全部结束
	GameStatusNotStart  = 4 // 牌局还未开始
)

// 聊天类型
const (
	GameChartTypeInput = 1 // 输入文本
	GameChartTypeText  = 2 // 文本
	GameChartTypeEmoji = 3 // 表情
	GameChartTypeVoice = 4 // 语音
)

// ipay交易状态
const (
	IpayStatusNotStart = -1 // 交易未完成
	IpayStatusSuccess  = 0  // 交易成功(付款成功)
	IpayStatusFailed   = 1  // 交易失败
	IpayStatusComplete = 2  // 交易完成(服务器确认并发放资产)
)

// 更新类型
const (
	HotFix         = 1
	OptionalUpdate = 2
	ForceUpdate    = 3
)

// 快速登录
const (
	QuickLoginSession    = "yangjin"
	QuickLoginPlayerID   = 100000
	QuickLoginOpenID     = "test"
	QuickLoginHeadimgurl = "http://7ktu6w.com1.z0.glb.clouddn.com/mj_head_boy.png"
)

// 反馈审阅状态
const ()

// 麻将类型
const (
	MjTypeHb  = 1 // 河北麻将
	MjTypeZz  = 2 // 郑州麻将
	MjTypeTdh = 3 // 推倒胡
	MjTypeKf  = 4 // 开封麻将
)

// 支付类型
const (
	PayTypeAppStore = 1 // appsotre
	PayTypeIPay     = 2 // 爱贝
)

// 版本是否发布
const (
	VersionDisabled = 0 // 关闭
	VersionEnabled  = 1 // 开启
)

// 分销等级
const (
	AgLevelAdmin = 0
	AgLevelOne   = 1
	AgLevelTwo   = 2
)

// 分销等级名称
const (
	AgLevelAdminStr = "管理员"
	AgLevelOneStr   = "大咖一级"
	AgLevelTwoStr   = "小咖二级"
)

// 分销等级提成
const (
	AgLevelOneRate  = 0.45
	AgLevelTwoRate  = 0.35
	AgFirstBuyRate  = 0.1
	AgLowAgentsRate = 0.05
)

// GPS距离限制
const (
	GPSDistanceLimit = 100
)

// 分销系统结算状态
const (
	AgBillStatusWeidakuan = 0 // 未打款
	AgBillStatusYidakuan  = 1 // 已打款
	AgBillStatusYijieyu   = 2 // 已结余
)
