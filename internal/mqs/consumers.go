package mqs

import "CampusTake/internal/svc"

// StartAllConsumers 一键启动所有 MQ 消费者
func StartAllConsumers(ctx *svc.ServiceContext) {
	go StartOrderCancelConsumer(ctx)
	go StartOrderConfirmConsumer(ctx)
}
