package mqs

import (
	"context"
	"strconv"
	"time"

	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	"CampusTake/internal/repo"
	"CampusTake/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

// StartOrderCancelConsumer 启动监听死信队列
func StartOrderCancelConsumer(ctx *svc.ServiceContext) {

	ch, err := ctx.MqConn.Channel()
	if err != nil {
		logx.Errorf("消费者获取 Channel 失败: %v", err)
		return
	}
	defer ch.Close()

	// 监听真正死信队列
	msgs, err := ch.Consume(
		ctx.Config.RabbitMQConfig.OrderCancelQueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logx.Errorf("监听取消订单队列失败: %v", err)
		return
	}

	logx.Infof("订单超时取消消费者启动成功")

	for d := range msgs {

		orderID, err := strconv.ParseInt(string(d.Body), 10, 64)
		if err != nil {
			logx.Errorf("订单ID解析失败: %v", err)
			d.Ack(false)
			continue
		}

		logx.Infof("捕获到超时订单，orderID=%d", orderID)

		// =========================
		// 查询订单
		// =========================

		order, err := ctx.Repo.Order.GetByID(
			context.Background(),
			orderID,
		)
		if err != nil {
			logx.Errorf("查询订单失败，orderID=%d err=%v", orderID, err)
			d.Ack(false)
			continue
		}

		// =========================
		// 如果已经不是待支付
		// 说明用户已经支付
		// =========================

		if order.Status != enums.OrderPendingPay {

			logx.Infof(
				"订单已支付，无需自动取消，orderID=%d status=%s",
				orderID,
				order.Status.String(),
			)

			d.Ack(false)
			continue
		}

		// =========================
		// 开事务自动取消
		// =========================

		err = ctx.Repo.WithTx(
			context.Background(),
			func(tx *repo.RepoTx) error {

				at := time.Now()

				fromStatus := enums.OrderPendingPay
				toStatus := enums.OrderTimeoutClosed

				// 更新订单状态
				err := tx.Order.SystemUpdateStatusAndTime(
					context.Background(),
					orderID,
					fromStatus,
					toStatus,
					at,
					map[string]interface{}{
						"cancel_reason": "订单超时未支付，系统自动关闭",
					},
				)
				if err != nil {
					return err
				}

				// 写订单日志
				log := &model.OrderLog{
					OrderID:      orderID,
					FromStatus:   fromStatus,
					ToStatus:     toStatus,
					OperatorType: enums.OperatorTypeSystem,
					OperatorID:   0,
					Remark:       "订单超时未支付，系统自动关闭",
					CreatedAt:    at,
				}

				return tx.Order.CreateLog(
					context.Background(),
					log,
				)
			},
		)

		if err != nil {

			logx.Errorf(
				"自动取消订单失败，orderID=%d err=%v",
				orderID,
				err,
			)

			continue
		}

		logx.Infof("订单自动取消成功，orderID=%d", orderID)

		// ACK
		d.Ack(false)
	}
}
