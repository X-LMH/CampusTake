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
    gender     TINYINT   NOT NULL DEFAULT '0' COMMENT '性别 0:保密 1:男 2:女',
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
CREATE TABLE rider_profile
(
    id                 BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id            BIGINT       NOT NULL COMMENT '关联用户ID',
    real_name          VARCHAR(50)  NOT NULL COMMENT '真实姓名',
    student_no         VARCHAR(50)  NOT NULL COMMENT '学号',
    id_card_no         VARCHAR(30)  NOT NULL COMMENT '身份证号',
    dormitory_building VARCHAR(100) NOT NULL COMMENT '宿舍楼',
    dormitory_room     VARCHAR(50)  NOT NULL COMMENT '宿舍号',
    campus_card_photo  VARCHAR(255) COMMENT '校园卡照片',

    audit_status       TINYINT       DEFAULT 0 COMMENT '审核状态：0待审核 1通过 2拒绝',
    audit_remark       VARCHAR(255) COMMENT '审核备注（当前结果说明）',

    rating_avg         DECIMAL(3, 2) DEFAULT 3.00 COMMENT '平均评分',
    rating_count       INT           DEFAULT 0 COMMENT '评价次数',
    completion_rate    DECIMAL(5, 2) DEFAULT 0 COMMENT '完成率',
    punctual_rate      DECIMAL(5, 2) DEFAULT 0 COMMENT '准时率',

    created_at         DATETIME      DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at         DATETIME      DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at         DATETIME     NULL COMMENT '软删除时间',

    UNIQUE KEY uk_user_id (user_id),
    INDEX idx_audit_status (audit_status),
    INDEX idx_deleted (deleted_at)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;


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
    id                  BIGINT PRIMARY KEY AUTO_INCREMENT,
    order_no            VARCHAR(64)    NOT NULL UNIQUE COMMENT '订单号',

    user_id             BIGINT         NOT NULL COMMENT '用户ID',
    rider_id            BIGINT   DEFAULT NULL COMMENT '代取员ID',

    order_type          TINYINT        NOT NULL COMMENT '订单类型：1快递 2外卖',

    pickup_address_id   BIGINT         NOT NULL COMMENT '取件地址ID',
    delivery_address_id BIGINT         NOT NULL COMMENT '收货地址ID',

    reward_amount       DECIMAL(10, 2) NOT NULL COMMENT '悬赏金额',

    status              TINYINT        NOT NULL COMMENT '
    1待支付 2待接单 3已接单 4取件中 5配送中
    6已完成 7已取消 8已退款 9超时 10异常 11申诉中 12已撤单',

    payment_status      TINYINT  DEFAULT 0 COMMENT '支付状态：0未支付 1已支付 2已退款',

    remark              VARCHAR(255) COMMENT '订单备注',
    cancel_reason       VARCHAR(255) COMMENT '取消原因',

    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at          DATETIME       NULL COMMENT '软删除时间',

    paid_at             DATETIME COMMENT '支付时间',
    accepted_at         DATETIME COMMENT '接单时间',
    finished_at         DATETIME COMMENT '完成时间',

    INDEX idx_user (user_id),
    INDEX idx_rider (rider_id),
    INDEX idx_status (status),
    INDEX idx_deleted (deleted_at)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;


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
CREATE TABLE review
(
    id         BIGINT PRIMARY KEY AUTO_INCREMENT,
    order_id   BIGINT   NOT NULL COMMENT '订单ID',
    user_id    BIGINT   NOT NULL COMMENT '用户ID',
    rider_id   BIGINT   NOT NULL COMMENT '代取员ID',
    score      TINYINT  NOT NULL COMMENT '评分（1-5）',
    content    VARCHAR(255) COMMENT '评价内容',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    deleted_at DATETIME NULL COMMENT '软删除时间',

    UNIQUE KEY uk_order (order_id),
    INDEX idx_rider (rider_id),
    INDEX idx_deleted (deleted_at)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;


-- =====================
-- 申诉表
-- =====================
CREATE TABLE appeal
(
    id           BIGINT PRIMARY KEY AUTO_INCREMENT,
    order_id     BIGINT   NOT NULL COMMENT '订单ID',
    applicant_id BIGINT   NOT NULL COMMENT '申诉人ID',
    type         TINYINT COMMENT '类型：1用户 2骑手',
    content      VARCHAR(255) COMMENT '申诉内容',
    status       TINYINT  DEFAULT 0 COMMENT '状态：0待处理 1已处理',
    handled_by   BIGINT COMMENT '处理人ID',
    handled_at   DATETIME COMMENT '处理时间',
    created_at   DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    deleted_at   DATETIME NULL COMMENT '软删除时间',

    INDEX idx_order (order_id),
    INDEX idx_status (status),
    INDEX idx_deleted (deleted_at)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;