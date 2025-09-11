package utils

import (
	"JudgeCore/internal/config"
	"JudgeCore/internal/global"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ConnectRabbitMQ RabbitMQ连接
func ConnectRabbitMQ(rmqConfig config.RabbitMQ) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(rmqConfig.Address)
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, err
	}

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

// PublishJudgeResult 将判题结果发布到消息队列
func PublishJudgeResult(ch *amqp.Channel, result global.JudgeResultMessage) error {
	_, err := ch.QueueDeclare(
		"judgeResults", // 队列名称
		true,           // 持久化
		false,          // 自动删除
		false,          // 排他性
		false,          // 不等待
		nil,            // 参数
	)
	if err != nil {
		return err
	}

	body, err := json.Marshal(result)
	if err != nil {
		return err
	}

	return ch.Publish(
		"",
		"judgeResults",
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	)
}
