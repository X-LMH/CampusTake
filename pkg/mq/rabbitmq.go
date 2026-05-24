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

// DelayQueueConfig 统一的延迟/死信队列配置
// DelayExchange/DelayQueue/DelayRoutingKey: 延迟队列
// DeadLetterExchange/DeadLetterQueue/DeadLetterRoutingKey: 死信队列
type DelayQueueConfig struct {
	DelayExchange        string
	DelayQueue           string
	DelayRoutingKey      string
	DeadLetterExchange   string
	DeadLetterQueue      string
	DeadLetterRoutingKey string
	TTL                  int64
}

func SetupDelayQueue(conn *amqp.Connection, cfg DelayQueueConfig, scene string) {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("❌ 【%s】获取 Channel 失败: %v", scene, err)
	}
	defer ch.Close()

	// 2. 声明死信交换机
	err = ch.ExchangeDeclare(cfg.DeadLetterExchange, "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("❌ 【%s】声明死信交换机 [%s] 失败: %v", scene, cfg.DeadLetterExchange, err)
	}

	// 3. 声明死信队列
	_, err = ch.QueueDeclare(cfg.DeadLetterQueue, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("❌ 【%s】声明死信队列 [%s] 失败: %v", scene, cfg.DeadLetterQueue, err)
	}

	// 4. 绑定死信队列
	err = ch.QueueBind(cfg.DeadLetterQueue, cfg.DeadLetterRoutingKey, cfg.DeadLetterExchange, false, nil)
	if err != nil {
		log.Fatalf("❌ 【%s】绑定死信队列失败: %v", scene, err)
	}

	// 5. 声明延迟交换机
	err = ch.ExchangeDeclare(cfg.DelayExchange, "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("❌ 【%s】声明延迟交换机 [%s] 失败: %v", scene, cfg.DelayExchange, err)
	}

	// 6. 声明延迟队列
	args := amqp.Table{
		"x-dead-letter-exchange":    cfg.DeadLetterExchange,
		"x-dead-letter-routing-key": cfg.DeadLetterRoutingKey,
		"x-message-ttl":             cfg.TTL * 1000,
	}

	_, err = ch.QueueDeclare(cfg.DelayQueue, true, false, false, false, args)
	if err != nil {
		log.Fatalf("❌ 【%s】声明延迟队列 [%s] 失败! 真正原因: %v", scene, cfg.DelayQueue, err)
	}

	// 7. 绑定延迟队列
	err = ch.QueueBind(cfg.DelayQueue, cfg.DelayRoutingKey, cfg.DelayExchange, false, nil)
	if err != nil {
		log.Fatalf("❌ 【%s】绑定延迟队列失败: %v", scene, err)
	}

	logx.Infof("✅ RabbitMQ %s 延时/死信队列初始化成功", scene)
}
