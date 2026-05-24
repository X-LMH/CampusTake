package mqs

import (
	"context"
	"encoding/json"
	"time"

	"CampusTake/internal/enums"
	"CampusTake/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

// StartRiderCheckConsumer 启动监听死信队列
func StartRiderCheckConsumer(svcCtx *svc.ServiceContext) {

	ch, err := svcCtx.MqConn.Channel()
	if err != nil {
		logx.Errorf("消费者获取 Channel 失败: %v", err)
		return
	}
	defer ch.Close()

	// 监听真正死信队列
	msgs, err := ch.Consume(
		svcCtx.Config.RabbitMQConfig.RiderCheck.DeadLetterQueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logx.Errorf("监听骑手接单检查队列失败: %v", err)
		return
	}

	logx.Infof("骑手接单检查消费者启动成功")

	for d := range msgs {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			var payload RiderCheckPayload
			err := json.Unmarshal(d.Body, &payload)
			if err != nil {
				logx.Errorf("骑手接单检查消息解析失败: %v", err)
				d.Ack(false)
				return
			}

			logx.Infof("捕获到延迟接单检查，orderID=%d riderID=%d", payload.OrderID, payload.RiderID)

			// 查询订单
			order, err := svcCtx.Repo.Order.GetByID(ctx, payload.OrderID)
			if err != nil {
				logx.Errorf("查询订单失败，orderID=%d err=%v", payload.OrderID, err)
				d.Ack(false)
				return
			}

			// 订单状态双检
			if order.Status == enums.OrderRiderCancelled || order.Status == enums.OrderUserCancelled || order.RiderID == nil {
				logx.Infof("订单已取消，无需计数，orderID=%d", payload.OrderID)
				d.Ack(false)
				return
			}

			if *order.RiderID != payload.RiderID {
				logx.Infof("骑手不匹配或已转单，orderID=%d riderID=%d", payload.OrderID, payload.RiderID)
				d.Ack(false)
				return
			}

			// 增加接单数量
			err = svcCtx.Repo.Rider.IncrementAcceptedOrderCount(ctx, payload.RiderID)
			if err != nil {
				logx.Errorf("增加接单数失败，riderID=%d err=%v", payload.RiderID, err)
				d.Nack(false, true)
				return
			}

			logx.Infof("增加接单数成功，riderID=%d", payload.RiderID)
			d.Ack(false)
		}()
	}
}
