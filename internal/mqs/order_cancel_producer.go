package mqs

import (
	"context"
	"strconv"

	"CampusTake/internal/svc"

	"github.com/rabbitmq/amqp091-go"
	"github.com/zeromicro/go-zero/core/logx"
)

func PublishDelayCancelOrder(
	ctx *svc.ServiceContext,
	orderID int64,
) error {

	ch, err := ctx.MqConn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// 消息内容：订单ID
	body := []byte(strconv.FormatInt(orderID, 10))

	// 发布延迟消息
	err = ch.PublishWithContext(
		context.Background(),
		ctx.Config.RabbitMQConfig.OrderDelayExchange,   // "campus_order_delay_exchange"
		ctx.Config.RabbitMQConfig.OrderDelayRoutingKey, // "campus.order.delay"
		false,
		false,
		amqp091.Publishing{
			ContentType:  "text/plain",
			Body:         body,
			DeliveryMode: amqp091.Persistent,
		},
	)
	if err != nil {
		logx.Errorf(
			"发送延迟取消订单消息失败，orderID=%d err=%v",
			orderID,
			err,
		)
		return err
	}

	logx.Infof(
		"发送延迟取消订单消息成功，orderID=%d",
		orderID,
	)

	return nil
}
