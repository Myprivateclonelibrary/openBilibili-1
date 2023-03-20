# 新增 block_user 表
CREATE TABLE `block_user` (
`id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增id',
`mid` INT(11) NOT NULL COMMENT '用户mid',
`status` TINYINT(4) NOT NULL COMMENT '封禁状态',
`ctime` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
`mtime` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
PRIMARY KEY (`id`),
UNIQUE KEY `uk_mid` (`mid`),
KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='封禁服务用户表';

# 新增 block_user_detail 表 ，用户详情表，用作聚合数据用
CREATE TABLE `block_user_detail` (
`id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增id',
`mid` INT(11) NOT NULL COMMENT '用户mid',
`block_count` INT(11) NOT NULL COMMENT '封禁计次',
`ctime` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
`mtime` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
PRIMARY KEY (`id`),
UNIQUE KEY `uk_mid` (`mid`),
KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='封禁服务用户详情表';
 
# 新增 block_history 表 —— 10张分表！！
CREATE TABLE `block_history_[0-9]` (
`id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增id',
`mid` INT(11) NOT NULL COMMENT '用户mid',
`admin_id` INT(11) NOT NULL COMMENT '管理员id',
`admin_name` VARCHAR(16) NOT NULL COMMENT '管理员name',
`source` TINYINT(4) NOT NULL COMMENT '封禁来源',
`area` TINYINT(4) NOT NULL COMMENT '封禁业务',
`reason` VARCHAR(50) NOT NULL COMMENT '封禁理由',
`comment` VARCHAR(50) NOT NULL COMMENT '封禁备注',
`action` TINYINT(4) NOT NULL COMMENT '操作类型',
`start_time` TIMESTAMP NOT NULL COMMENT '开始生效时间',
`duration` INT(11) NOT NULL COMMENT '生效时长（秒）',
`notify` TINYINT(4) NOT NULL COMMENT '是否通知',
`ctime` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
`mtime` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
PRIMARY KEY (`id`),
KEY `ix_mid` (`mid`),
KEY `ix_action` (`action`),
KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='封禁服务用户历史表';