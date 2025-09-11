package utils

import (
	"FeasOJ/internal/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ConnectRabbitMQ 建立与 RabbitMQ 的连接
func ConnectRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
	// 连接到 RabbitMQ 服务
	conn, err := amqp.Dial(config.GlobalConfig.RabbitMQ.Address)
	if err != nil {
		return nil, nil, err
	}

	// 创建一个通道
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, err
	}

	// 确保队列存在
	_, err = ch.QueueDeclare(
		"judgeTask", // 队列名称
		true,        // 是否持久化
		false,       // 是否自动删除
		false,       // 是否排他
		false,       // 是否等待消费者
		nil,         // 额外参数
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, nil, err
	}

	return conn, ch, nil
}

// PublishTask 将任务发布到 RabbitMQ 队列
func PublishTask(ch *amqp.Channel, task string) error {
	// 将任务发送到队列
	err := ch.Publish(
		"",          // 默认交换机
		"judgeTask", // 队列名称
		false,       // 是否等待确认
		false,       // 是否是强制的
		amqp.Publishing{
			ContentType: "text/plain", // 消息的类型
			Body:        []byte(task), // 消息内容
		},
	)
	return err
}
