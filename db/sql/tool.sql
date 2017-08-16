SET NAMES utf8mb4;

USE `mj`;

-- ----------------------------
--  Table structure for `account`
-- ----------------------------
CREATE TABLE `account` (
  `index_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '记录索引',
  `email` varchar(16) NOT NULL COMMENT '用户名',
  `password` varchar(64) NOT NULL COMMENT '用户密码',
  `role` varchar(16) NOT NULL COMMENT '用户角色',
  PRIMARY KEY (`index_id`),
  UNIQUE KEY `uidx_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='tool用户表';

-- ----------------------------
--  Table structure for `permission`
-- ----------------------------
CREATE TABLE `permission` (
  `index_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '记录索引',
  `pmsn_type` varchar(16) NOT NULL COMMENT '权限类型',
  `pmsn_content` varchar(32) NOT NULL COMMENT '权限内容',
  `comment` varchar(16) NOT NULL COMMENT '权限描述',
  PRIMARY KEY (`index_id`),
  UNIQUE KEY `uidx_pmsn_type_content` (`pmsn_type`, `pmsn_content`),
  KEY `idx_pmsn_types` (`pmsn_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='tool权限表';

-- ----------------------------
--  Table structure for `role`
-- ----------------------------
CREATE TABLE `role` (
  `index_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '记录索引',
  `role` varchar(32) NOT NULL COMMENT '角色名',
  `comment` varchar(16) NOT NULL COMMENT '角色描述',
  `permissions` json NOT NULL COMMENT '角色权限',
  PRIMARY KEY (`index_id`),
  UNIQUE KEY `uidx_role` (`role`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='tool角色表';
