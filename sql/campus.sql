create database if not exists campus;
use campus;
-- =====================
-- 用户表
-- =====================
CREATE TABLE user
(
    id         BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '用户ID',
    phone      VARCHAR(20)  NOT NULL UNIQUE COMMENT '手机号',
    password   VARCHAR(100) NOT NULL COMMENT '密码',
    nickname   VARCHAR(50)           DEFAULT '' COMMENT '昵称',
    avatar     VARCHAR(255)          DEFAULT '' COMMENT '头像',
    gender     TINYINT      NOT NULL DEFAULT '0' COMMENT '性别 0:保密 1:男 2:女',
    role       TINYINT      NOT NULL DEFAULT 1 COMMENT '角色：1普通用户 2代取员 3管理员',
    status     TINYINT      NOT NULL DEFAULT 1 COMMENT '状态：1正常 2禁用',
    created_at DATETIME              DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME              DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at DATETIME     NULL COMMENT '软删除时间',

    INDEX idx_phone (phone),
    INDEX idx_role (role),
    INDEX idx_deleted (deleted_at)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;


-- =====================
-- 代取员信息表
-- =====================
create table rider_profile
(
    id                    bigint auto_increment primary key,
    user_id               bigint                                  not null comment '关联用户ID',
    real_name             varchar(50)                             not null comment '真实姓名',
    student_no            varchar(50)                             not null comment '学号',
    id_card_no            varchar(30)                             not null comment '身份证号',
    dormitory_building    varchar(100)                            not null comment '宿舍楼',
    dormitory_room        varchar(50)                             not null comment '宿舍号',
    campus_card_front     varchar(255)                            not null comment '校园卡正面',
    campus_card_back      varchar(255)                            not null comment '校园卡反面',
    audit_status          tinyint       default 1                 null comment '审核状态：1待审核 2通过 3拒绝 4撤回',
    audit_remark          varchar(255)                            null comment '审核备注',
    -- =====================
    -- 统计字段
    -- =====================
    rating_avg            decimal(3, 2) default 3.00              null comment '平均评分',
    rating_count          int           default 0                 null comment '评分次数',
    completed_order_count int           default 0                 null comment '完成订单数',
    completion_rate       decimal(5, 2) default 0.00              null comment '完成率',
    created_at            datetime      default CURRENT_TIMESTAMP null,
    updated_at            datetime      default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
    deleted_at            datetime                                null,

    UNIQUE KEY uk_user_id (user_id),
    -- 普通索引
    INDEX idx_audit_status (audit_status),
    INDEX idx_deleted (deleted_at)
) charset = utf8mb4;

-- 索引保持不变

create table rider_profile
(
    id                    bigint auto_increment primary key,

    user_id               bigint                                  not null comment '关联用户ID',

    real_name             varchar(50)                             not null comment '真实姓名',

    student_no            varchar(50)                             not null comment '学号',

    id_card_no            varchar(30)                             not null comment '身份证号',

    dormitory_building    varchar(100)                            not null comment '宿舍楼',

    dormitory_room        varchar(50)                             not null comment '宿舍号',

    campus_card_front     varchar(255)                            not null comment '校园卡正面',

    campus_card_back      varchar(255)                            not null comment '校园卡反面',

    audit_status          tinyint       default 1                 null comment '审核状态：1待审核 2通过 3拒绝 4撤回',

    audit_remark          varchar(255)                            null comment '审核备注',

    -- =====================
    -- 统计字段
    -- =====================

    rating_avg            decimal(3, 2) default 5.00              null comment '平均评分',

    rating_count          int           default 0                 null comment '评分次数',

    completed_order_count int           default 0                 null comment '完成订单数',

    completion_rate       decimal(5, 2) default 0.00              null comment '完成率',

    created_at            datetime      default CURRENT_TIMESTAMP null,

    updated_at            datetime      default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,

    deleted_at            datetime                                null,

    constraint uk_user_id unique (user_id)
);

-- =====================
-- 代取员审核日志
-- =====================
CREATE TABLE rider_audit_log
(
    id         BIGINT PRIMARY KEY AUTO_INCREMENT,
    rider_id   BIGINT   NOT NULL COMMENT '代取员ID',
    auditor_id BIGINT   NOT NULL COMMENT '审核人ID（管理员）',
    result     TINYINT  NOT NULL COMMENT '审核结果：1通过 2拒绝',
    remark     VARCHAR(255) COMMENT '审核说明（历史记录）',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '审核时间',
    deleted_at DATETIME NULL COMMENT '软删除时间',

    INDEX idx_rider (rider_id),
    INDEX idx_auditor (auditor_id),
    INDEX idx_deleted (deleted_at)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;


-- =====================
-- 地址表（type区分）
-- =====================
CREATE TABLE address
(
    id            BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id       BIGINT       NOT NULL COMMENT '所属用户',

    type          TINYINT      NOT NULL COMMENT '地址类型：1收货地址 2取件地址',

    contact_name  VARCHAR(50)  NOT NULL COMMENT '联系人姓名',
    contact_phone VARCHAR(20)  NOT NULL COMMENT '联系人电话',

    building      VARCHAR(100) NOT NULL COMMENT '建筑名称（宿舍/教学楼）',
    room          VARCHAR(50) COMMENT '房间号',
    detail        VARCHAR(255) COMMENT '详细描述',

    is_default    TINYINT  DEFAULT 0 COMMENT '是否默认地址',

    created_at    DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at    DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at    DATETIME     NULL COMMENT '软删除时间',

    INDEX idx_user_type (user_id, type),
    INDEX idx_default (user_id, is_default),
    INDEX idx_deleted (deleted_at)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;


-- =====================
-- 订单表（核心）
-- =====================
CREATE TABLE `order`
(
    id                  BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    order_no            VARCHAR(64)    NOT NULL COMMENT '订单号',
    user_id             BIGINT         NOT NULL COMMENT '下单用户ID',
    rider_id            BIGINT         NULL COMMENT '骑手ID',
    order_type          TINYINT        NOT NULL COMMENT '订单类型：1=快递代取 2=外卖代取',
    pickup_address_id   BIGINT         NOT NULL COMMENT '取件地址ID',
    delivery_address_id BIGINT         NOT NULL COMMENT '送达地址ID',
    reward_amount       DECIMAL(10, 2) NOT NULL COMMENT '悬赏金额',
    status              tinyint        not null comment '订单状态：
        1  待支付
        2  待接单
        3  已接单
        4  配送中
        5  已送达
        6  已完成
        7  用户取消
        8  骑手取消
        9  超时关闭
        10 已退款
        11 异常订单
        12 申诉中',
    payment_status      TINYINT        NOT NULL DEFAULT 0 COMMENT '支付状态：0=未支付 1=已支付 2=已退款',
    remark              VARCHAR(255)   NULL COMMENT '订单备注',
    cancel_reason       VARCHAR(255)   NULL COMMENT '取消原因',
    created_at          DATETIME       NULL     DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at          DATETIME       NULL     DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at          DATETIME       NULL COMMENT '软删除时间',
    paid_at             DATETIME       NULL COMMENT '支付时间',
    accepted_at         DATETIME       NULL COMMENT '骑手接单时间',
    picked_up_at        DATETIME       NULL COMMENT '骑手取件时间',
    delivered_at        DATETIME       NULL COMMENT '骑手送达时间',
    completed_at        datetime       null comment '订单完成时间',
    cancelled_at        DATETIME       NULL COMMENT '取消时间',
    refunded_at         DATETIME       NULL COMMENT '退款时间',
    CONSTRAINT uk_order_no UNIQUE (order_no)
) ENGINE = INNODB
  DEFAULT CHARSET = utf8mb4 COMMENT = '订单主表';

-- =========================
-- 索引
-- =========================
CREATE INDEX idx_user ON `order` (user_id);
CREATE INDEX idx_rider ON `order` (rider_id);
CREATE INDEX idx_status ON `order` (status);
CREATE INDEX idx_payment_status ON `order` (payment_status);
CREATE INDEX idx_created_at ON `order` (created_at);
CREATE INDEX idx_deleted_at ON `order` (deleted_at);


-- =====================
-- 订单状态日志
-- =====================
CREATE TABLE order_log
(
    id            BIGINT PRIMARY KEY AUTO_INCREMENT,
    order_id      BIGINT   NOT NULL COMMENT '订单ID',
    from_status   TINYINT COMMENT '原状态',
    to_status     TINYINT COMMENT '新状态',
    operator_type TINYINT COMMENT '操作人类型：1用户 2骑手 3系统 4管理员',
    operator_id   BIGINT COMMENT '操作人ID',
    remark        VARCHAR(255) COMMENT '操作说明',
    created_at    DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    deleted_at    DATETIME NULL COMMENT '软删除时间',

    INDEX idx_order (order_id),
    INDEX idx_deleted (deleted_at)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;


-- =====================
-- 支付表
-- =====================
CREATE TABLE payment
(
    id          BIGINT PRIMARY KEY AUTO_INCREMENT,
    order_id    BIGINT         NOT NULL COMMENT '订单ID',
    pay_no      VARCHAR(64)    NOT NULL UNIQUE COMMENT '支付单号',
    amount      DECIMAL(10, 2) NOT NULL COMMENT '支付金额',
    status      TINYINT        NOT NULL COMMENT '支付状态：0待支付 1已支付 2已退款',
    method      TINYINT  DEFAULT 1 COMMENT '支付方式：1模拟支付',
    paid_at     DATETIME COMMENT '支付时间',
    refunded_at DATETIME COMMENT '退款时间',
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    deleted_at  DATETIME       NULL COMMENT '软删除时间',

    INDEX idx_order (order_id),
    INDEX idx_status (status),
    INDEX idx_deleted (deleted_at)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;


-- =====================
-- 评价表
-- =====================
create table review
(
    id           bigint auto_increment primary key,
    order_id     bigint                             not null comment '订单ID',
    user_id      bigint                             not null comment '用户ID',
    rider_id     bigint                             not null comment '骑手ID',
    score        tinyint                            not null comment '评分（1-5）',
    content      varchar(255)                       null comment '评价内容',
    is_anonymous tinyint  default 0                 null comment '是否匿名',
    created_at   datetime default CURRENT_TIMESTAMP null,
    deleted_at   datetime                           null,

    UNIQUE KEY uk_order (order_id),
    INDEX idx_rider (rider_id),
    INDEX idx_deleted (deleted_at)
);


-- =====================
-- 申诉表
-- =====================
CREATE TABLE appeal
(
    id              BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    order_id        BIGINT            NOT NULL COMMENT '订单ID',
    applicant_id    BIGINT            NOT NULL COMMENT '申诉人ID',
    type            TINYINT           NOT NULL COMMENT '申诉类型：1用户申诉 2骑手申诉',
    content         VARCHAR(500)      NOT NULL COMMENT '申诉内容',
    status          TINYINT DEFAULT 1 NOT NULL COMMENT '1待处理 2已通过 3已驳回 4已撤销',
    handle_remark   VARCHAR(255)      DEFAULT ''                NULL COMMENT '处理备注',
    refund_amount   DECIMAL(10, 2)    DEFAULT 0.00              NULL COMMENT '退款金额',
    punish_rider    TINYINT           DEFAULT 0                 NULL COMMENT '是否处罚骑手：0否 1是',
    terminate_order TINYINT           DEFAULT 0                 NULL COMMENT '是否终止订单：0否 1是',
    handled_by      BIGINT                                         NULL COMMENT '处理管理员ID',
    handled_at      DATETIME                                       NULL COMMENT '处理时间',
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP             NULL COMMENT '创建时间',
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NULL COMMENT '更新时间',
    deleted_at      DATETIME                                       NULL COMMENT '软删除时间',

    INDEX idx_order_id (order_id),
    INDEX idx_status (status),
    INDEX idx_applicant_id (applicant_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
    COMMENT = '订单申诉表';