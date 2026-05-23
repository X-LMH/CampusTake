package mqs

import (
	"context"
	"encoding/json"

	"CampusTake/internal/enums"
	"CampusTake/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

// StartRiderCheckConsumer 启动监听死信队列
func StartRiderCheckConsumer(ctx *svc.ServiceContext) {

	ch, err := ctx.MqConn.Channel()
	if err != nil {
		logx.Errorf("消费者获取 Channel 失败: %v", err)
		return
	}
	defer ch.Close()

	// 监听真正死信队列
	msgs, err := ch.Consume(
		ctx.Config.RabbitMQConfig.RiderCheck.DeadLetterQueue,
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

		var payload RiderCheckPayload
		err := json.Unmarshal(d.Body, &payload)
		if err != nil {
			logx.Errorf("骑手接单检查消息解析失败: %v", err)
			d.Ack(false)
			continue
		}

		logx.Infof("捕获到延迟接单检查，orderID=%d riderID=%d", payload.OrderID, payload.RiderID)

		// =========================
		// 查询订单
		// =========================

		order, err := ctx.Repo.Order.GetByID(
			context.Background(),
			payload.OrderID,
		)
		if err != nil {
			logx.Errorf("查询订单失败，orderID=%d err=%v", payload.OrderID, err)
			d.Ack(false)
			continue
		}

		// =========================
		// 订单状态双检
		// 如果订单不再属于该骑手，或者已取消等，跳过
		// =========================

		// 判断订单是不是被转单、取消了，还是由当前骑手处理
		if order.RiderID != payload.RiderID {
			logx.Infof("骑手不匹配或已转单，orderID=%d riderID=%d", payload.OrderID, payload.RiderID)
			d.Ack(false)
			continue
		}

		// 检查状态，如果处于接单之后的状态（已接单、配送中、已送达、已完成等等）
		// 如果是已取消、退款等异常状态，也需要看具体业务。正常来说，只要骑手接过且目前还是这个骑手，就可以计数
		if order.Status == enums.OrderCanceled || order.Status == enums.OrderTimeoutClosed {
			logx.Infof("订单已取消，无需计数，orderID=%d", payload.OrderID)
			d.Ack(false)
			continue
		}

		// =========================
		// 增加接单数量
		// =========================

		err = ctx.Repo.Rider.IncrementAcceptedOrderCount(
			context.Background(),
			payload.RiderID,
		)

		if err != nil {
			logx.Errorf(
				"增加接单数失败，riderID=%d err=%v",
				payload.RiderID,
				err,
			)
			d.Nack(false, true)
			continue
		}

		logx.Infof("增加接单数成功，riderID=%d", payload.RiderID)

		// ACK
		d.Ack(false)
	}
}

