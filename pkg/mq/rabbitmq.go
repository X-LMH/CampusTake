package mq

import (
	"fmt"
	"log"
	"net/url"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zeromicro/go-zero/core/logx"
)

// InitRabbitMQ 初始化 RabbitMQ 连接
func InitRabbitMQ(user, password, host, port, vhost string) *amqp.Connection {
	escapedVhost := url.PathEscape(vhost)

	urlStr := fmt.Sprintf("amqp://%s:%s@%s:%s/%s", user, password, host, port, escapedVhost)
	conn, err := amqp.Dial(urlStr)
	if err != nil {
		log.Fatalf("RabbitMQ 连接失败: %v", err)
	}

	return conn
}

// SetupOrderCancelQueues 核心逻辑：自动建交换机和队列，并绑定死信规则
func SetupOrderCancelQueues(
	conn *amqp.Connection,
	delayEx, delayQueue, delayKey string,
	cancelEx, cancelQueue, cancelKey string,
	ttl int32,
) {
	// 1. 创建独立的专属 Channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("❌ 【第1步失败】RabbitMQ 获取 Channel 失败: %v", err)
	}
	defer ch.Close()

	// 2. 声明死信交换机
	err = ch.ExchangeDeclare(cancelEx, "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("❌ 【第2步失败】声明死信交换机 [%s] 失败: %v", cancelEx, err)
	}

	// 3. 声明死信队列
	_, err = ch.QueueDeclare(cancelQueue, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("❌ 【第3步失败】声明死信队列 [%s] 失败: %v", cancelQueue, err)
	}

	// 4. 绑定死信队列
	err = ch.QueueBind(cancelQueue, cancelKey, cancelEx, false, nil)
	if err != nil {
		log.Fatalf("❌ 【第4步失败】绑定死信队列失败: %v", err)
	}

	// 5. 声明延迟交换机
	err = ch.ExchangeDeclare(delayEx, "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("❌ 【第5步失败】声明延迟交换机 [%s] 失败: %v", delayEx, err)
	}

	// 6. 声明延迟队列（最容易崩的地方）
	args := amqp.Table{
		"x-dead-letter-exchange":    cancelEx,
		"x-dead-letter-routing-key": cancelKey,
		"x-message-ttl":             ttl, // 强制转为 int32
	}

	// 💡 注意：这里我们抓取声明那一刻的真正 err
	_, err = ch.QueueDeclare(delayQueue, true, false, false, false, args)
	if err != nil {
		// 🔥 这里会打印出诸如 PRECONDITION_FAILED 的真正原因，而不是无意义的 504！
		log.Fatalf("❌ 【第6步致命失败】声明延迟队列 [%s] 失败! 真正原因: %v", delayQueue, err)
	}

	// 7. 绑定延迟队列
	err = ch.QueueBind(delayQueue, delayKey, delayEx, false, nil)
	if err != nil {
		log.Fatalf("❌ 【第7步失败】绑定延迟队列失败: %v", err)
	}

	logx.Info("✅ RabbitMQ 延时/死信队列基建初始化成功！")
}
