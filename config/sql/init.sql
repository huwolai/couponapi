-- +migrate Up
-- 优惠券
CREATE TABLE IF NOT EXISTS coupon(
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  v_code VARCHAR(100) UNIQUE  COMMENT '优惠券唯一标记',
  app_id VARCHAR(255) DEFAULT '' COMMENT 'appID',
  title VARCHAR(100) DEFAULT '' COMMENT '优惠券标题',
  remark VARCHAR(1000) DEFAULT '' COMMENT '优惠券备注',
  amount NUMERIC(10,2) COMMENT '优惠券金额',
  publish_num int COMMENT '发放数量 0.表示没限制',
  published_num int COMMENT '已发放数量',
  is_one int COMMENT '是否只能使用一次 0.否 1.是',
  status int COMMENT '优惠券状态 1.正常 0.禁止',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳'
) CHARACTER SET utf8mb4;

-- 用户领取的优惠券
CREATE TABLE IF NOT EXISTS coupon_user(
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  app_id VARCHAR(255) DEFAULT '' COMMENT 'appID',
  open_id VARCHAR(100) DEFAULT '' COMMENT '用户ID',
  coupon_code VARCHAR(100) DEFAULT '' COMMENT '券唯一代号',
  title VARCHAR(100) DEFAULT '' COMMENT '优惠券标题',
  remark VARCHAR(1000) DEFAULT '' COMMENT '优惠券备注',
  amount NUMERIC(10,2) COMMENT '优惠券金额',
  balance NUMERIC(10,2) COMMENT '余额',
  flag VARCHAR(100) COMMENT '标识 ACCOUNT_RECHARGE',
  is_one int COMMENT '是否只能使用一次 0.否 1.是',
  use_status int COMMENT '使用状态 0.未激活 1.未使用或未使用完 2.已使用',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳',
  KEY open_id (open_id),
  KEY coupon_code (coupon_code)
) CHARACTER SET utf8mb4;


-- 优惠券追踪
CREATE TABLE IF NOT EXISTS coupon_track(
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  app_id VARCHAR(255) COMMENT 'app_id',
  trade_no VARCHAR(100) COMMENT '交易号',
  trade_type int COMMENT '交易类型',
  track_code VARCHAR(100) DEFAULT '' COMMENT '追踪编号',
  open_id VARCHAR(100) DEFAULT '' COMMENT '用户ID',
  coupon_code VARCHAR(100) DEFAULT '' COMMENT '第三方券代号',
  title VARCHAR(100) COMMENT '标题',
  remark VARCHAR(255) COMMENT '备注',
  amount NUMERIC(14,2)  COMMENT '实际价格',
  track_type int COMMENT '记录类型 1. 券使用 2.券生成',
  coupon_amount NUMERIC(10,2) COMMENT '优惠掉的金额',
  status int COMMENT '0.待使用 1.已使用',
  create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间戳',
  KEY open_id (open_id),
  KEY coupon_code (coupon_code),
  KEY track_code (track_code)
) CHARACTER SET utf8mb4;


