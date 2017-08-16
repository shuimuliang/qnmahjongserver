SET NAMES utf8mb4;

USE `mj`;

-- ----------------------------
--  Table structure for `ag_bill`
-- ----------------------------
CREATE TABLE `ag_bill` (
  `index_id` int(11) NOT NULL AUTO_INCREMENT,
  `ag_id` int(11) NOT NULL COMMENT '代理商ID，对应游戏id',
  `bill_money` int(11) NOT NULL DEFAULT '0' COMMENT '提现金额，分为单位',
  `last_week_left` int(11) NOT NULL DEFAULT '0' COMMENT '累计结余金额',
  `last_week_balance` int(11) NOT NULL DEFAULT '0' COMMENT '上周提成金额',
  `last_week_dakuan` int(11) NOT NULL DEFAULT '0' COMMENT '上周打款金额',
  `low_agents_award` int(11) NOT NULL DEFAULT '0' COMMENT '下级代理提成',
  `cards_award` int(11) NOT NULL DEFAULT '0' COMMENT '销售玉提成',
  `first_buy_award` int(11) NOT NULL DEFAULT '0' COMMENT '首充赠送',
  `hongbao` int(11) NOT NULL DEFAULT '0' COMMENT '500分销红包',
  `delflag` int(11) NOT NULL DEFAULT '0' COMMENT '已提现标识位',
  `start_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '结算起始时间',
  `end_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '结算终止时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '结算时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '打款时间',
  PRIMARY KEY (`index_id`),
  UNIQUE KEY `uidx_start_time_ag_id` (`start_time`,`ag_id`),
  KEY `idx_start_time` (`start_time`),
  KEY `idx_ag_id` (`ag_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='代理商返现明细表';

CREATE TABLE `ag_bill` (
  `index_id` int(11) NOT NULL AUTO_INCREMENT,
  `ag_id` int(11) NOT NULL COMMENT '代理商ID，对应游戏id',
  `bill_money` int(11) NOT NULL DEFAULT '0' COMMENT '提现金额，分为单位',
  `last_week_left` int(11) NOT NULL DEFAULT '0' COMMENT '累计结余金额',
  `last_week_balance` int(11) NOT NULL DEFAULT '0' COMMENT '上周提成金额',
  `last_week_dakuan` int(11) NOT NULL DEFAULT '0' COMMENT '上周打款金额',
  `low_agents_award` int(11) NOT NULL DEFAULT '0' COMMENT '下级代理提成',
  `cards_award` int(11) NOT NULL DEFAULT '0' COMMENT '销售玉提成',
  `first_buy_award` int(11) NOT NULL DEFAULT '0' COMMENT '首充赠送',
  `hongbao` int(11) NOT NULL DEFAULT '0' COMMENT '500分销红包',
  `delflag` int(11) NOT NULL DEFAULT '0' COMMENT '已提现标识位',
  `start_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '结算起始时间',
  `end_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '结算终止时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '结算时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '打款时间',
  PRIMARY KEY (`index_id`),
  UNIQUE KEY `uidx_start_time_ag_id` (`start_time`,`ag_id`),
  KEY `idx_start_time` (`start_time`),
  KEY `idx_ag_id` (`ag_id`)
) ENGINE=InnoDB AUTO_INCREMENT=70 DEFAULT CHARSET=utf8mb4 COMMENT='代理商返现明细表';

-- ----------------------------
--  Table structure for `ag_pay`
-- ----------------------------
CREATE TABLE `ag_pay` (
  `index_id` int(11) NOT NULL AUTO_INCREMENT,
  `ag_id` int(11) NOT NULL COMMENT '代理商ID，对应游戏id',
  `customer_id` int(11) NOT NULL COMMENT '下线玩家id，对应游戏id',
  `diamond_cnt` int(11) NOT NULL DEFAULT '0' COMMENT '玩家充值钻石数，对应游戏内千胜玉',
  `money_cnt` int(11) NOT NULL DEFAULT '0' COMMENT '玩家充值金额',
  `delflag` int(11) NOT NULL DEFAULT '0' COMMENT '已结算标识位',
  `first_buy_award` int(11) NOT NULL DEFAULT '0' COMMENT '首充奖励,0否 1是',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '下线玩家充值时间',
  PRIMARY KEY (`index_id`),
  KEY `idx_ag_id` (`ag_id`),
  KEY `idx_ag_id_customer_id` (`ag_id`,`customer_id`),
  KEY `idx_ag_id_delflag` (`ag_id`,`delflag`),
  KEY `idx_ag_id_ag_id` (`ag_id`,`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='下线玩家充值明细表';

-- ----------------------------
--  Table structure for `ag_auth`
-- ----------------------------
CREATE TABLE `ag_auth` (
  `index_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '记录索引',
  `ag_upper_id` int(11) NOT NULL COMMENT '上级代理商ID，对应游戏id',    
  `ag_id` int(11) NOT NULL COMMENT '代理商ID，对应游戏id',
  `ag_level` int(11) NOT NULL COMMENT '授权等级',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '上级代理商授权时间',
  PRIMARY KEY (`index_id`),
  UNIQUE KEY `uidx_ag_upper_id_ag_id` (`ag_upper_id`,`ag_id`),
  KEY `idx_ag_upper_id` (`ag_upper_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='代理商授权表';

-- ----------------------------
--  Table structure for `ag_account`
-- ----------------------------
CREATE TABLE `ag_account` (
  `index_id` int(11) NOT NULL AUTO_INCREMENT,
  `ag_upper_id` int(11) NOT NULL COMMENT '上级代理商ID，对应游戏id',
  `ag_id` int(11) NOT NULL COMMENT '代理商ID，对应游戏id',
  `ag_level` int(11) NOT NULL COMMENT '代理商等级',
  `password` varchar(64) NOT NULL COMMENT '密码',
  `telephone` varchar(64) NOT NULL DEFAULT '' COMMENT '手机号',
  `realname` varchar(64) NOT NULL DEFAULT '' COMMENT '真实姓名',
  `weixin` varchar(64) NOT NULL DEFAULT '' COMMENT '微信',
  `alipay` varchar(64) NOT NULL DEFAULT '' COMMENT '支付宝',
  `email` varchar(64) NOT NULL DEFAULT '' COMMENT '电子邮箱',
  `hongbao` int(11) NOT NULL DEFAULT '0' COMMENT '分销红包是否领取(0否,1是)',
  `total_balance` int(11) NOT NULL DEFAULT '0' COMMENT '总打款金额',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '下级代理商注册时间',
  PRIMARY KEY (`index_id`),
  UNIQUE KEY `uidx_ag_id` (`ag_id`),
  KEY `idx_ag_upper_id` (`ag_upper_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='代理商用户表';