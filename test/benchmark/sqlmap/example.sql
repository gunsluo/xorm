CREATE TABLE `t_apply_detail` (
  `apply_detail_id` int(11) NOT NULL AUTO_INCREMENT,
  `apply_object_id` int(11) DEFAULT NULL,
  `user_id` varchar(255) DEFAULT NULL COMMENT '申请人',
  `user_name` varchar(255) DEFAULT NULL COMMENT ' 申请人姓名',
  `auth_type` varchar(255) DEFAULT NULL COMMENT '申请类型',
  `key` varchar(255) DEFAULT NULL COMMENT 'auth key',
  `apply_flag` varchar(255) DEFAULT NULL COMMENT '申请or共享',
  `apply_status` varchar(255) DEFAULT NULL COMMENT '申请状态',
  `created_time` datetime DEFAULT NULL COMMENT '用户创建时间',
  `updated_time` datetime DEFAULT NULL COMMENT '用户修改时间',
  `deleted_time` datetime DEFAULT NULL COMMENT '用户删除时间',
  `deleted` tinyint(4) DEFAULT NULL COMMENT '用户状态（0表示正常，1表示已删除，2表示禁用）',
  PRIMARY KEY (`apply_detail_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `t_apply_object` (
  `apply_object_id` int(11) NOT NULL AUTO_INCREMENT,
  `apply_object_name` varchar(255) DEFAULT NULL COMMENT '名称',
  `environment_id` int(11) DEFAULT NULL,
  `environment` varchar(255) DEFAULT NULL COMMENT '运行环境',
  `tps` varchar(255) DEFAULT NULL COMMENT '预估tps',
  `peak_value` varchar(255) DEFAULT NULL COMMENT '预估峰值',
  `object_type` varchar(255) DEFAULT NULL COMMENT '申请对象类型',
  `owner_id` int(11) DEFAULT NULL,
  `owner_name` varchar(255) DEFAULT NULL,
  `department_id` int(255) DEFAULT NULL COMMENT '申请部门ID',
  `department_name` varchar(255) DEFAULT NULL COMMENT '申请部门名称',
  `describe` varchar(255) DEFAULT NULL COMMENT '场景描述',
  `created_time` datetime DEFAULT NULL COMMENT '用户创建时间',
  `updated_time` datetime DEFAULT NULL COMMENT '用户修改时间',
  `deleted_time` datetime DEFAULT NULL COMMENT '用户删除时间',
  `deleted` tinyint(4) DEFAULT NULL COMMENT '用户状态（0表示正常，1表示已删除，2表示禁用）',
  PRIMARY KEY (`apply_object_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
