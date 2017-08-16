SET NAMES utf8mb4;

USE `mj`;

INSERT INTO `account` (`index_id`, `email`, `password`, `role`)
VALUES
	(1, '921046333@qq.com', '36336ed2e9318ffeeef299551fdc8dd870d7cddf', 'admin');

INSERT INTO `cost` (`index_id`, `mj_type`, `mj_desc`, `rounds`, `cards`, `coins`)
VALUES
	(1, 1, '河北麻将', 1, 3, 30),
	(2, 1, '河北麻将', 2, 6, 60),
	(3, 2, '郑州麻将', 1, 3, 30),
	(4, 2, '郑州麻将', 2, 6, 60),
	(5, 3, '推倒胡', 1, 3, 30),
	(6, 3, '推倒胡', 2, 6, 60),
	(7, 4, '开封麻将', 1, 3, 30),
	(8, 4, '开封麻将', 2, 6, 60);

INSERT INTO `game` (`index_id`, `channel`, `version`, `size`, `module`, `mj_types`, `enabled`, `update_type`, `download_url`, `svn_version`)
VALUES
	(1, 100, 10000, 0, '{}', '1', 1, 1, '', 1497),
	(2, 200, 10000, 0, '{}', '1', 1, 1, '', 1497),
	(3, 110, 11000, 0, '{}', '1|2|3|4', 1, 1, '', 1799),
	(4, 210, 11000, 0, '{}', '1|2|3|4', 1, 1, '', 1799),
	(5, 110, 12000, 0, '{}', '1|2|3|4', 1, 1, '', 1799),
	(6, 210, 12000, 0, '{}', '1|2|3|4', 1, 1, '', 1799);

INSERT INTO `module` (`index_id`, `module`, `comment`)
VALUES
	(1, 'quick_login', '快速登录'),
	(2, 'history_record', '历史战绩'),
	(3, 'shop', '商城充值'),
	(4, 'invite_award', '邀请有礼'),
	(5, 'weixin_login', '微信登录');

INSERT INTO `player` (`player_id`, `high_id`, `invite_award`, `openid`, `access_token`, `expires_in`, `refresh_token`, `nickname`, `sex`, `province`, `city`, `country`, `headimgurl`, `unionid`, `coins`, `cards`, `first_buy`)
VALUES
	(100000, 100000, 0, 'test', '', '2017-04-10 16:21:18', '', '', 0, '', '', '', 'http://7ktu6w.com1.z0.glb.clouddn.com/mj_head_boy.png', '', 1000, 0, 0);

INSERT INTO `shop` (`index_id`, `channel`, `pay_type`,`gem_id`, `wares_id`, `wares_name`, `goods_count`, `extra_count`, `price`, `icon_url`)
VALUES
	(1, 100, 2, 100201, '11', '小胜玉18枚', 12, 6, 1200, 'http://7u2l4h.com1.z0.glb.clouddn.com/mj_interface_buy_small.png'),
	(2, 100, 2, 100202, '12', '小胜玉50枚', 30, 20, 3000, 'http://7u2l4h.com1.z0.glb.clouddn.com/mj_interface_buy_medium.png'),
	(3, 100, 2, 100203, '13', '小胜玉120枚', 60, 60, 6000, 'http://7u2l4h.com1.z0.glb.clouddn.com/mj_interface_buy_common.png'),
	(4, 100, 2, 100204, '14', '小胜玉200枚', 90, 110, 9000, 'http://7u2l4h.com1.z0.glb.clouddn.com/mj_interface_buy_big.png'),
	(5, 200, 2, 200201, '17', '小胜玉18枚', 12, 6, 1200, 'http://7u2l4h.com1.z0.glb.clouddn.com/mj_interface_buy_small.png'),
	(6, 200, 2, 200202, '18', '小胜玉50枚', 30, 20, 3000, 'http://7u2l4h.com1.z0.glb.clouddn.com/mj_interface_buy_medium.png'),
	(7, 200, 2, 200203, '19', '小胜玉120枚', 60, 60, 6000, 'http://7u2l4h.com1.z0.glb.clouddn.com/mj_interface_buy_common.png'),
	(8, 200, 2, 200204, '20', '小胜玉200枚', 90, 110, 9000, 'http://7u2l4h.com1.z0.glb.clouddn.com/mj_interface_buy_big.png'),
	(9, 110, 2, 110201, '11', '小胜玉18枚', 12, 6, 1200, 'http://7u2l4h.com1.z0.glb.clouddn.com/mj_interface_buy_small.png'),
	(10, 110, 2, 110202, '12', '小胜玉50枚', 30, 20, 3000, 'http://7u2l4h.com1.z0.glb.clouddn.com/mj_interface_buy_medium.png'),
	(11, 110, 2, 110203, '13', '小胜玉120枚', 60, 60, 6000, 'http://7u2l4h.com1.z0.glb.clouddn.com/mj_interface_buy_common.png'),
	(12, 110, 2, 110204, '14', '小胜玉200枚', 90, 110, 9000, 'http://7u2l4h.com1.z0.glb.clouddn.com/mj_interface_buy_big.png'),
	(13, 210, 2, 210201, '17', '小胜玉18枚', 12, 6, 1200, 'http://7u2l4h.com1.z0.glb.clouddn.com/mj_interface_buy_small.png'),
	(14, 210, 2, 210202, '18', '小胜玉50枚', 30, 20, 3000, 'http://7u2l4h.com1.z0.glb.clouddn.com/mj_interface_buy_medium.png'),
	(15, 210, 2, 210203, '19', '小胜玉120枚', 60, 60, 6000, 'http://7u2l4h.com1.z0.glb.clouddn.com/mj_interface_buy_common.png'),
	(16, 210, 2, 210204, '20', '小胜玉200枚', 90, 110, 9000, 'http://7u2l4h.com1.z0.glb.clouddn.com/mj_interface_buy_big.png'),
	(17, 110, 1, 110101, 'com.snowcat.hnmj.gem12', '小胜玉18枚', 12, 6, 1200, 'http://7u2l4h.com1.z0.glb.clouddn.com/mj_interface_buy_small.png'),
	(18, 110, 1, 110102, 'com.snowcat.hnmj.gem30', '小胜玉50枚', 30, 20, 3000, 'http://7u2l4h.com1.z0.glb.clouddn.com/mj_interface_buy_medium.png'),
	(19, 110, 1, 110103, 'com.snowcat.hnmj.gem60', '小胜玉120枚', 60, 60, 6000, 'http://7u2l4h.com1.z0.glb.clouddn.com/mj_interface_buy_common.png'),
	(20, 110, 1, 110104, 'com.snowcat.hnmj.gem98', '小胜玉208枚', 98, 110, 9800, 'http://7u2l4h.com1.z0.glb.clouddn.com/mj_interface_buy_big.png');
