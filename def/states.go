package def

// 返回给客户端的状态码
const (
	StatusWait   = -2 // 等待
	StatusFailed = -1 // 失败(服务器用)
	StatusOK     = 0  // 成功

	StatusErrorLoginType     = 101 // 错误的登录类型
	StatusErrorWeixin        = 102 // 微信登录失败
	StatusErrorWeixinExpires = 103 // 微信登录过期
	StatusErrorLogin         = 104 // 登录服务器异常

	StatusErrorToken  = 301 // token错误
	StatusErrorPlayer = 302 // 玩家不存在
	StatusErrorRoom   = 303 // 房间不存在
	StatusErrorLogic  = 304 // 逻辑服务器异常

	StatusNoPlayer  = 401 // 玩家不存在
	StatusNoRoom    = 402 // 房间不存在
	StatusNotInRoom = 403 // 玩家不在房间

	StatusErrorInviteCode   = 1101 // 邀请码无效
	StatusErrorInviteFailed = 1102 // 绑定邀请码失败
	StatusErrorInviteInfo   = 1103 // 领取信息有误
	StatusErrorInviteAward  = 1104 // 领取奖励失败

	StatusErrorCreateRoom  = 2101 // 创建房间失败
	StatusErrorNoRoomID    = 2102 // 房间号不存在
	StatusErrorRoomIsFull  = 2103 // 房间满员了
	StatusErrorIPConflict  = 2104 // IP冲突
	StatusErrorGeoConflict = 2105 // GPS位置太近
	StatusErrorGPSNotOpen  = 2106 // GPS未打开

	StatusIsStart   = 2201 // 游戏已经开始(退出房间)
	StatusNoStart   = 2202 // 游戏还没开始(投票退出)
	StatusIsGaming  = 2203 // 游戏正在进行(准备游戏)
	StatusNotGaming = 2204 // 游戏不在进行(取消准备)
	StatusNotOwner  = 2205 // 不是房主(关闭房间)
	StatusOperError = 2206 // 操作非法(操作)
	StatusBuyPao    = 2208 // 已经买跑(买跑)
)
