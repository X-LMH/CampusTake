package mqs

import (
	"CampusTake/internal/repo/query"
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
func StartOrderCancelConsumer(svcCtx *svc.ServiceContext) {

	ch, err := svcCtx.MqConn.Channel()
	if err != nil {
		logx.Errorf("消费者获取 Channel 失败: %v", err)
		return
	}
	defer ch.Close()

	// 监听真正死信队列
	msgs, err := ch.Consume(
		svcCtx.Config.RabbitMQConfig.OrderCancel.DeadLetterQueue,
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
		// 绝招：不传参，直接在闭包内部使用外层的 d
		// 这样不管你用什么 rabbitmq 库，都绝对不会报类型错误！
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			orderID, err := strconv.ParseInt(string(d.Body), 10, 64)
			if err != nil {
				logx.Errorf("订单ID解析失败: %v", err)
				d.Ack(false)
				return // 触发 defer cancel()
			}

			logx.Infof("捕获到超时订单，orderID=%d", orderID)

			// 查询订单
			order, err := svcCtx.Repo.Order.GetByID(ctx, orderID)
			if err != nil {
				logx.Errorf("查询订单失败，orderID=%d err=%v", orderID, err)
				d.Ack(false)
				return
			}

			if order.Status != enums.OrderPendingPay {
				logx.Infof("订单已支付，无需自动取消，orderID=%d status=%s", orderID, order.Status.String())
				d.Ack(false)
				return
			}

			// 开事务自动取消
			err = svcCtx.Repo.WithTx(ctx, func(tx *repo.RepoTx) error {
				at := time.Now()
				fromStatus := enums.OrderPendingPay
				toStatus := enums.OrderTimeoutClosed

				err := tx.Order.UpdateStatusAndTime(
					ctx,
					query.OrderStatusUpdateQuery{OrderID: orderID},
					fromStatus, toStatus, at,
					map[string]interface{}{"cancel_reason": "订单超时未支付，系统自动关闭"},
				)
				if err != nil {
					return err
				}

				log := &model.OrderLog{
					OrderID:      orderID,
					FromStatus:   fromStatus,
					ToStatus:     toStatus,
					OperatorType: enums.OperatorTypeSystem,
					Remark:       "订单超时未支付，系统自动关闭",
					CreatedAt:    at,
				}
				return tx.Order.CreateLog(ctx, log)
			},
			)

			if err != nil {
				logx.Errorf("自动取消订单失败，orderID=%d err=%v", orderID, err)
				d.Nack(false, true)
				return
			}

			logx.Infof("订单自动取消成功，orderID=%d", orderID)
			d.Ack(false)
		}() // 注意：这里不需要传 d 了
	}
}
