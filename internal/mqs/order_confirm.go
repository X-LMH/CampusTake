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

// StartOrderConfirmConsumer 启动确认收货死信队列监听
func StartOrderConfirmConsumer(svcCtx *svc.ServiceContext) {

	ch, err := svcCtx.MqConn.Channel()
	if err != nil {
		logx.Errorf("消费者获取 Channel 失败: %v", err)
		return
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		svcCtx.Config.RabbitMQConfig.OrderConfirm.DeadLetterQueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logx.Errorf("监听确认收货队列失败: %v", err)
		return
	}

	logx.Infof("订单自动确认收货消费者启动成功")

	for d := range msgs {
		// 利用无参闭包，既避开了具体的 RabbitMQ 类型声明，又完美应用了 defer cancel()
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			orderID, err := strconv.ParseInt(string(d.Body), 10, 64)
			if err != nil {
				logx.Errorf("订单ID解析失败: %v", err)
				d.Ack(false)
				return // 触发 defer cancel()，结束本轮处理
			}

			logx.Infof("捕获到待自动确认收货订单，orderID=%d", orderID)

			// =========================
			// 查询订单（使用带超时的 ctx）
			// =========================
			order, err := svcCtx.Repo.Order.GetByID(ctx, orderID)
			if err != nil {
				logx.Errorf("查询订单失败，orderID=%d err=%v", orderID, err)
				d.Ack(false)
				return
			}

			// =========================
			// 只有已送达才自动确认
			// =========================
			if order.Status != enums.OrderDelivered {
				logx.Infof(
					"订单无需自动确认，orderID=%d status=%s",
					orderID,
					order.Status.String(),
				)
				d.Ack(false)
				return
			}

			// =========================
			// 开事务自动确认（使用带超时的 ctx）
			// =========================
			err = svcCtx.Repo.WithTx(ctx, func(tx *repo.RepoTx) error {
				at := time.Now()
				fromStatus := enums.OrderDelivered
				toStatus := enums.OrderCompleted

				err := tx.Order.UpdateStatusAndTime(
					ctx,
					query.OrderStatusUpdateQuery{OrderID: orderID},
					fromStatus,
					toStatus,
					at,
					nil,
				)
				if err != nil {
					return err
				}

				log := &model.OrderLog{
					OrderID:      orderID,
					FromStatus:   fromStatus,
					ToStatus:     toStatus,
					OperatorType: enums.OperatorTypeSystem,
					OperatorID:   0,
					Remark:       "系统自动确认收货",
					CreatedAt:    at,
				}

				if err := tx.Order.CreateLog(ctx, log); err != nil {
					return err
				}

				if order.AppealStatus == enums.OrderAppealStatusApproved {
					return nil
				}
				if order.RiderID == nil {
					return nil
				}

				return tx.Rider.IncrementCompletedOrderCount(ctx, *order.RiderID)
			},
			)

			if err != nil {
				logx.Errorf("自动确认收货失败，orderID=%d err=%v", orderID, err)
				d.Nack(false, true)
				return
			}

			logx.Infof("订单自动确认收货成功，orderID=%d", orderID)
			d.Ack(false)
		}()
	}
}
