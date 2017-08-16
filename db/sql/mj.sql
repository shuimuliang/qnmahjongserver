SET NAMES utf8mb4;

CREATE DATABASE `mj` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

USE `mj`;

-- ----------------------------
--  Table structure for `chat`
-- ----------------------------
CREATE TABLE `chat` (
  `index_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '记录索引',
  `send_id` int(11) NOT NULL COMMENT '发起玩家id',
  `room_id` int(11) NOT NULL COMMENT '玩家房间id',
  `msg_type` int(11) NOT NULL COMMENT '消息类型 1 输入文本，2 已定义的文本，3 表情，4 语音',
  `mess_id` int(11) NOT NULL COMMENT 'mess_id',
  `msg_text` varchar(256) NOT NULL COMMENT '消息文本',
  `send_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发送时间',
  PRIMARY KEY (`index_id`),
  KEY `idx_send_id` (`send_id`),
  KEY `idx_room_id` (`room_id`),
  KEY `idx_send_time` (`send_time`)
) ENGINE=InnoDB AUTO_INCREMENT=10000001 DEFAULT CHARSET=utf8mb4 COMMENT='游戏聊天表';

-- ----------------------------
--  Table structure for `cost`
-- ----------------------------
CREATE TABLE `cost` (
  `index_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '记录索引',
  `mj_type` int(11) NOT NULL COMMENT '麻将类型',
  `mj_desc` varchar(16) NOT NULL COMMENT '麻将描述',
  `rounds` int(11) NOT NULL COMMENT '开房局数',
  `cards` int(11) NOT NULL COMMENT '开房消耗玉',
  `coins` int(11) NOT NULL COMMENT '开房消耗币',
  PRIMARY KEY (`index_id`),
  UNIQUE KEY `uidx_mj_type_rounds` (`mj_type`,`rounds`),
  KEY `idx_mj_type` (`mj_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='开房消耗表';

-- ----------------------------
--  Table structure for `feedback`
-- ----------------------------
CREATE TABLE `feedback` (
  `index_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '记录索引',
  `player_id` int(11) NOT NULL COMMENT '玩家id',
  `channel` int(11) NOT NULL COMMENT '游戏渠道',
  `version` int(11) NOT NULL COMMENT '游戏版本',
  `img_url` varchar(256) NOT NULL COMMENT '图片链接',
  `text` varchar(256) NOT NULL COMMENT '问题描述',
  `status` int(11) NOT NULL COMMENT '审阅状态，0未审阅，1审阅',
  `add_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  `revise_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`index_id`),
  KEY `idx_player_id` (`player_id`),
  KEY `idx_channel` (`channel`),
  KEY `idx_channel_version` (`channel`,`version`),
  KEY `idx_status` (`status`),
  KEY `idx_add_time` (`add_time`),
  KEY `idx_revise_time` (`revise_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='玩家反馈表';

-- ----------------------------
--  Table structure for `game`
-- ----------------------------
CREATE TABLE `game` (
  `index_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '记录索引',
  `channel` int(11) NOT NULL COMMENT '游戏渠道',
  `version` int(11) NOT NULL COMMENT '游戏版本',
  `size` int(11) NOT NULL COMMENT '包大小',
  `module` json NOT NULL COMMENT '模块配置',
  `mj_types` varchar(256) NOT NULL COMMENT '包含麻将玩法',
  `enabled` int(11) NOT NULL COMMENT '游戏是否打开',
  `update_type` int(11) NOT NULL COMMENT '游戏更新类型 1:动态更新 2:可跳过更新 3:强更',
  `download_url` varchar(256) NOT NULL COMMENT '游戏下载链接',
  `svn_version` int(11) NOT NULL COMMENT '当前包的svn版本',
  PRIMARY KEY (`index_id`),
  UNIQUE KEY `uidx_channel_version` (`channel`,`version`),
  KEY `idx_channel` (`channel`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='游戏开关表';

-- ----------------------------
--  Table structure for `login`
-- ----------------------------
CREATE TABLE `login` (
  `index_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '记录索引',
  `player_id` int(11) NOT NULL COMMENT '玩家id',
  `login_channel` int(11) NOT NULL COMMENT '登录渠道',
  `login_version` int(11) NOT NULL COMMENT '登录版本',
  `login_type` int(11) NOT NULL COMMENT '登录方式',
  `login_ip` varchar(15) NOT NULL COMMENT '登录地址',
  `login_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
  `login_machine` varchar(32) NOT NULL DEFAULT '' COMMENT '登录机器码',
  PRIMARY KEY (`index_id`),
  KEY `idx_player_id` (`player_id`),
  KEY `idx_login_channel` (`login_channel`),
  KEY `idx_login_channel_login_version` (`login_channel`,`login_version`),
  KEY `idx_login_time` (`login_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='登录日志表';

-- ----------------------------
--  Table structure for `module`
-- ----------------------------
CREATE TABLE `module` (
  `index_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '记录索引',
  `module` varchar(16) NOT NULL COMMENT '模块名',
  `comment` varchar(16) NOT NULL COMMENT '模块描述',
  PRIMARY KEY (`index_id`),
  UNIQUE KEY `uidx_module` (`module`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='游戏模块表';

-- ----------------------------
--  Table structure for `order`
-- ----------------------------
CREATE TABLE `order` (
  `index_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '记录索引',
  `order_id` varchar(32) NOT NULL COMMENT '订单id',
  `trans_id` varchar(32) NOT NULL COMMENT '交易流水号',
  `player_id` int(11) NOT NULL COMMENT '玩家id',
  `channel` int(11) NOT NULL COMMENT '玩家渠道',
  `pay_type` int(11) NOT NULL COMMENT '支付方式',
  `gem_id` int(11) NOT NULL COMMENT '商品id',
  `wares_id` varchar(32) NOT NULL COMMENT '商品编号',
  `wares_name` varchar(32) NOT NULL COMMENT '商品名称',
  `goods_count` int(11) NOT NULL COMMENT '商品数量',
  `extra_count` int(11) NOT NULL COMMENT '赠送数量',
  `price` int(11) NOT NULL COMMENT '商品价格(分为单位)',
  `status` int(11) NOT NULL COMMENT 'ipay交易状态',
  `add_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  `revise_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`index_id`),
  UNIQUE KEY `idx_order_id` (`order_id`),
  KEY `idx_trans_id` (`trans_id`),
  KEY `idx_player_id` (`player_id`),
  KEY `idx_channel` (`channel`),
  KEY `idx_add_time` (`add_time`),
  KEY `idx_revise_time` (`revise_time`)
) ENGINE=InnoDB AUTO_INCREMENT=10000001 DEFAULT CHARSET=utf8mb4 COMMENT='游戏订单表';

-- ----------------------------
--  Table structure for `player`
-- ----------------------------
CREATE TABLE `player` (
  `player_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '用户唯一id',
  `high_id` int(11) NOT NULL COMMENT '上线用户唯一id',
  `invite_award` int(11) NOT NULL COMMENT '上线是否领奖，0未领取，1领取',
  `openid` varchar(32) NOT NULL COMMENT '授权用户唯一标识，普通用户的标识，对当前开发者帐号唯一',
  `access_token` varchar(128) NOT NULL COMMENT '接口调用凭证',
  `expires_in` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'access_token接口调用凭证超时时间',
  `refresh_token` varchar(128) NOT NULL COMMENT '用户刷新access_token',
  `nickname` varchar(16) NOT NULL COMMENT '普通用户昵称',
  `sex` int(11) NOT NULL COMMENT '普通用户性别，0为未知，1为男性，2为女性',
  `province` varchar(16) NOT NULL COMMENT '普通用户个人资料填写的省份',
  `city` varchar(16) NOT NULL COMMENT '普通用户个人资料填写的城市',
  `country` varchar(16) NOT NULL COMMENT '国家，如中国为CN',
  `headimgurl` varchar(256) NOT NULL COMMENT '用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空',
  `unionid` varchar(32) NOT NULL COMMENT '用户统一标识。针对一个微信开放平台帐号下的应用，同一用户的unionid是唯一的。',
  `coins` int(11) NOT NULL COMMENT '游戏币',
  `cards` int(11) NOT NULL COMMENT '房卡',
  `first_buy` int(11) NOT NULL COMMENT '首充，0否，1是',
  PRIMARY KEY (`player_id`),
  UNIQUE KEY `uidx_openid` (`openid`),
  KEY `idx_high_id` (`high_id`),
  KEY `idx_coins` (`coins`),
  KEY `idx_cards` (`cards`)
) ENGINE=InnoDB AUTO_INCREMENT=100001 DEFAULT CHARSET=utf8mb4 COMMENT='玩家信息表';

-- ----------------------------
--  Table structure for `record`
-- ----------------------------
CREATE TABLE `record` (
  `record_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '对局唯一id',
  `create_id` int(11) NOT NULL COMMENT '创建者id',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `room_id` int(11) NOT NULL COMMENT '房间id',
  `mj_type` int(11) NOT NULL COMMENT '麻将类型',
  `total_round` int(11) NOT NULL COMMENT '麻将局数',
  `start_time` datetime NOT NULL COMMENT '开局时间',
  `cur_round` int(11) NOT NULL COMMENT '当前局数',
  `east_id` int(11) NOT NULL COMMENT '东id',
  `south_id` int(11) NOT NULL COMMENT '南id',
  `west_id` int(11) NOT NULL COMMENT '西id',
  `north_id` int(11) NOT NULL COMMENT '北id',
  `east_score` int(11) NOT NULL COMMENT '东score',
  `south_score` int(11) NOT NULL COMMENT '南score',
  `west_score` int(11) NOT NULL COMMENT '西score',
  `north_score` int(11) NOT NULL COMMENT '北score',
  PRIMARY KEY (`record_id`),
  KEY `idx_create_id` (`create_id`),
  KEY `idx_create_time` (`create_time`),
  KEY `idx_create_id_create_time` (`create_id`,`create_time`)
) ENGINE=InnoDB AUTO_INCREMENT=10000001 DEFAULT CHARSET=utf8mb4 COMMENT='游戏战绩表';

-- ----------------------------
--  Table structure for `register`
-- ----------------------------
CREATE TABLE `register` (
  `index_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '记录索引',
  `player_id` int(11) NOT NULL COMMENT '玩家id',
  `register_channel` int(11) NOT NULL COMMENT '注册渠道',
  `register_version` int(11) NOT NULL COMMENT '注册版本',
  `register_type` int(11) NOT NULL COMMENT '注册方式',
  `register_ip` varchar(15) NOT NULL COMMENT '注册地址',
  `register_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '注册时间',
  `register_machine` varchar(32) NOT NULL COMMENT '注册机器码',
  PRIMARY KEY (`index_id`),
  KEY `idx_player_id` (`player_id`),
  KEY `idx_register_channel` (`register_channel`),
  KEY `idx_register_channel_register_version` (`register_channel`,`register_version`),
  KEY `idx_register_time` (`register_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='注册日志表';

-- ----------------------------
--  Table structure for `shop`
-- ----------------------------
CREATE TABLE `shop` (
  `index_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '记录索引',
  `channel` int(11) NOT NULL COMMENT '游戏渠道',
  `pay_type` int(11) NOT NULL COMMENT '支付方式',
  `gem_id` int(11) NOT NULL COMMENT '商品id',
  `wares_id` varchar(32) NOT NULL DEFAULT '' COMMENT '商品编号',
  `wares_name` varchar(32) NOT NULL COMMENT '商品名称',
  `goods_count` int(11) NOT NULL COMMENT '商品数量',
  `extra_count` int(11) NOT NULL COMMENT '赠送数量',
  `price` int(11) NOT NULL COMMENT '商品价格(分为单位)',
  `icon_url` varchar(256) NOT NULL COMMENT '图片链接',
  PRIMARY KEY (`index_id`),
  UNIQUE KEY `uidx_channel_gem_id` (`channel`,`gem_id`),
  UNIQUE KEY `uidx_channel_wares_id` (`channel`,`wares_id`),
  KEY `idx_channel` (`channel`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4 COMMENT='商品价目表';

-- ----------------------------
--  Table structure for `treasure`
-- ----------------------------
CREATE TABLE `treasure` (
  `index_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '记录索引',
  `player_id` int(11) NOT NULL COMMENT '玩家id',
  `reason` int(11) NOT NULL COMMENT '变化原因',
  `coins` int(11) NOT NULL COMMENT '游戏币变化',
  `cards` int(11) NOT NULL COMMENT '房卡变化',
  `change_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '交易时间',
  PRIMARY KEY (`index_id`),
  KEY `idx_player_id` (`player_id`),
  KEY `idx_reason` (`reason`),
  KEY `idx_change_time` (`change_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='资产变化日志表';