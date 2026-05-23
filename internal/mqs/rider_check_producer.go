package mqs

import (
	"context"
	"encoding/json"

	"CampusTake/internal/svc"

	"github.com/rabbitmq/amqp091-go"
	"github.com/zeromicro/go-zero/core/logx"
)

type RiderCheckPayload struct {
	OrderID int64 `json:"order_id"`
	RiderID int64 `json:"rider_id"`
}

func PublishDelayRiderCheck(
	ctx *svc.ServiceContext,
	orderID int64,
	riderID int64,
) error {

	ch, err := ctx.MqConn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	payload := RiderCheckPayload{
		OrderID: orderID,
		RiderID: riderID,
	}
	body, _ := json.Marshal(payload)

	// 发布延迟消息
	err = ch.PublishWithContext(
		context.Background(),
		ctx.Config.RabbitMQConfig.RiderCheck.DelayExchange,
		ctx.Config.RabbitMQConfig.RiderCheck.DelayRoutingKey,
		false,
		false,
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp091.Persistent,
		},
	)
	if err != nil {
		logx.Errorf(
			"发送延迟检查骑手接单状态消息失败，orderID=%d riderID=%d err=%v",
			orderID,
			riderID,
			err,
		)
		return err
	}

	logx.Infof(
		"发送延迟检查骑手接单状态消息成功，orderID=%d riderID=%d",
		orderID,
		riderID,
	)

	return nil
}

