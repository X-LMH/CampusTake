CREATE
    DATABASE IF NOT EXISTS `campus` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE
    `campus`;

-- ==========================================
-- 1. 用户表 (User)
-- ==========================================
CREATE TABLE `user`
(
    `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `phone`      VARCHAR(20)     NOT NULL COMMENT '手机号',
    `password`   VARCHAR(255)    NOT NULL COMMENT '密码',
    `nickname`   VARCHAR(50)     NOT NULL DEFAULT '默认昵称' COMMENT '用户昵称',
    `avatar`     VARCHAR(255)    NOT NULL DEFAULT '' COMMENT '头像URL',
    `status`     TINYINT         NOT NULL DEFAULT 1 COMMENT '状态: 1-正常, 2-封禁',
    `role`       TINYINT         NOT NULL DEFAULT 1 COMMENT '角色: 1-普通用户, 2-接单用户, 3-管理员',
    `created_at` TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP       NULL     DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_phone` (`phone`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='用户基本信息表';

-- ==========================================
-- 2. 地址信息表 (Address)
-- ==========================================
CREATE TABLE `address`
(
    `id`            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `user_id`       BIGINT UNSIGNED NOT NULL COMMENT '归属用户ID',
    `contact_name`  VARCHAR(50)     NOT NULL COMMENT '联系人姓名',
    `contact_phone` VARCHAR(20)     NOT NULL COMMENT '联系电话',
    `building`      VARCHAR(50)     NOT NULL COMMENT '宿舍楼/教学楼(如: 南苑1栋)',
    `room`          VARCHAR(20)     NOT NULL COMMENT '具体房间号(如: 404寝室)',
    `is_default`    TINYINT         NOT NULL DEFAULT 0 COMMENT '是否默认地址: 0-否, 1-是',
    `created_at`    TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`    TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at`    TIMESTAMP       NULL     DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`), -- 用户查自己的地址列表时需要走索引
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='用户地址库表';

-- ==========================================
-- 3. 跑腿订单表 (Order)
-- ==========================================
CREATE TABLE `order`
(
    `id`                  BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `order_no`            VARCHAR(64)     NOT NULL COMMENT '业务订单号(雪花算法生成，全局唯一)',
    `publisher_id`        BIGINT UNSIGNED NOT NULL COMMENT '发单人(用户)ID',
    `receiver_id`         BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '接单人(骑手)ID, 0表示暂未接单',
    `order_type`          TINYINT         NOT NULL COMMENT '订单类型: 1-外卖代取, 2-快递代取',
    `status`              TINYINT         NOT NULL DEFAULT 10 COMMENT '状态: 10-待接单, 20-配送中, 30-已完成, 40-已取消',

    `pickup_location`     VARCHAR(255)    NOT NULL COMMENT '取件地址(如: 北门美团外卖柜 / 菜鸟驿站)',
    `delivery_address_id` BIGINT UNSIGNED NOT NULL COMMENT '送达的地址表ID (关联 address 表)',
    `remark`              VARCHAR(255)    NOT NULL DEFAULT '' COMMENT '给接单人的备注(如: 快递取件码 / 外卖柜手机号尾号)',
    `item_image`          VARCHAR(255)    NOT NULL DEFAULT '' COMMENT '物品/外卖截图URL(可选)',
    `reward_amount`       INT             NOT NULL DEFAULT 0 COMMENT '赏金金额(单位: 分，避免浮点数精度丢失)',

    `expire_time`         TIMESTAMP       NOT NULL COMMENT '接单超时时间(用于结合 Redis 延迟队列)',
    `version`             INT             NOT NULL DEFAULT 0 COMMENT '乐观锁版本号(防超卖核心字段)',

    `created_at`          TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`          TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at`          TIMESTAMP       NULL     DEFAULT NULL COMMENT '软删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_order_no` (`order_no`),             -- 订单号必须唯一
    KEY `idx_publisher_id` (`publisher_id`),           -- 发单人查自己发的单
    KEY `idx_receiver_id` (`receiver_id`),             -- 接单人查自己接的单
    KEY `idx_status_created` (`status`, `created_at`), -- 复合索引：用于首页大厅按时间倒序查"待接单"的单子
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='互助跑腿订单表';