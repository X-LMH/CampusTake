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
create table `order`
(
    id                  bigint auto_increment comment '主键ID'
        primary key,
    order_no            varchar(64)                        not null comment '订单号',
    user_id             bigint                             not null comment '下单用户ID',
    rider_id            bigint                             null comment '骑手ID',
    order_type          tinyint                            not null comment '订单类型：1=快递代取 2=外卖代取',
    pickup_address_id   bigint                             not null comment '取件地址ID',
    delivery_address_id bigint                             not null comment '送达地址ID',
    reward_amount       decimal(10, 2)                     not null comment '悬赏金额',
    status              tinyint                            not null comment '订单状态：
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
        11 异常订单',
    payment_status      tinyint  default 0                 not null comment '支付状态：0=未支付 1=已支付 2=已退款',
    appeal_status       tinyint  default 0                 not null comment '0 无申诉 1 申诉中 2 申诉通过 3 申诉驳回',
    can_reassign        tinyint  default 0                 not null comment '骑手取消后是否可重新派单',
    remark              varchar(255)                       null comment '订单备注',
    cancel_reason       varchar(255)                       null comment '取消原因',
    created_at          datetime default CURRENT_TIMESTAMP null comment '创建时间',
    updated_at          datetime default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP comment '更新时间',
    deleted_at          datetime                           null comment '软删除时间',
    paid_at             datetime                           null comment '支付时间',
    accepted_at         datetime                           null comment '骑手接单时间',
    picked_up_at        datetime                           null comment '骑手取件时间',
    delivered_at        datetime                           null comment '骑手送达时间',
    completed_at        datetime                           null comment '订单完成时间',
    cancelled_at        datetime                           null comment '取消时间',
    refunded_at         datetime                           null comment '退款时间',
    constraint uk_order_no
        unique (order_no)
)
    comment '订单主表' charset = utf8mb4;

create index idx_created_at
    on `order` (created_at);

create index idx_deleted_at
    on `order` (deleted_at);

create index idx_payment_status
    on `order` (payment_status);

create index idx_rider
    on `order` (rider_id);

create index idx_status
    on `order` (status);

create index idx_user
    on `order` (user_id);


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
create table appeal
(
    id               bigint auto_increment comment '申诉ID'
        primary key,
    order_id         bigint                             not null comment '订单ID',
    applicant_id     bigint                             not null comment '申诉人ID',
    applicant_role   tinyint                            not null comment '
    1用户
    2骑手
    ',
    appeal_type      tinyint                            not null comment '
    1未收到货
    2骑手态度恶劣
    3骑手提前点送达
    4用户恶意退款
    5用户恶意投诉
    6订单丢失
    7物品损坏
    8其他
    ',
    content          varchar(500)                       not null comment '申诉内容',
    evidence_urls    json                               null comment '证据图片',
    status           tinyint  default 1                 not null comment '
    1待处理
    2已通过
    3已驳回
    4已撤销
    ',
    handled_by       bigint                             null comment '处理管理员ID',
    handled_at       datetime                           null comment '处理时间',
    created_at       datetime default CURRENT_TIMESTAMP null,
    updated_at       datetime default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
    deleted_at       datetime                           null
)
    comment '订单申诉表' charset = utf8mb4;

-- =====================
-- 申诉处理记录表
-- =====================
create table appeal_handle
(
    id              bigint auto_increment comment '处理记录ID'
        primary key,
    appeal_id       bigint                                   not null comment '申诉ID',
    handler_id      bigint                                   not null comment '管理员ID',
    result          tinyint                                  not null comment '
    1通过
    2驳回
    ',
    remark          varchar(255)                             null comment '处理备注',
    refund_amount   decimal(10, 2) default 0.00              null comment '退款金额',
    punish_rider    tinyint        default 0                 null comment '是否处罚骑手',
    punish_user     tinyint        default 0                 null comment '是否处罚用户',
    terminate_order tinyint        default 0                 null comment '是否终止订单',
    created_at      datetime       default CURRENT_TIMESTAMP null
)
    comment '申诉处理记录表' charset = utf8mb4;
