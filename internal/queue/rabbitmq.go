package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type RabbitMQConfig struct {
	URL         string `env:"RABBITMQ_URL" default:"amqp://guest:guest@localhost:5672/"`
	QueuePrefix string `env:"RABBITMQ_QUEUE_PREFIX" default:"gofxtest"`
}

type rabbitMQ struct {
	config   RabbitMQConfig
	conn     *amqp.Connection
	channel  *amqp.Channel
	mu       sync.Mutex
	handlers map[string]HandlerFunc
	logger   *zap.Logger
}

func NewRabbitMQ(config RabbitMQConfig, logger *zap.Logger) (*rabbitMQ, error) {
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &rabbitMQ{
		config:   config,
		conn:     conn,
		channel:  ch,
		handlers: make(map[string]HandlerFunc),
		logger:   logger,
	}, nil
}

func (r *rabbitMQ) Publish(ctx context.Context, topic string, v any) error {
	if msg, ok := v.(Message); ok {
		return r.publish(ctx, topic, msg)
	}

	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return err
	}

	err = r.publish(ctx, topic, Message{
		Data: jsonBytes,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *rabbitMQ) publish(ctx context.Context, topic string, msg Message) error {
	headers := amqp.Table{}
	for k, v := range msg.Headers {
		headers[k] = v
	}

	return r.channel.PublishWithContext(ctx,
		topic, // exchange
		"",    // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Headers:     headers,
			ContentType: "application/octet-stream",
			Body:        msg.Data,
		},
	)
}

func (r *rabbitMQ) Subscribe(ctx context.Context, topic string, handler HandlerFunc) error {
	err := r.channel.ExchangeDeclare(
		topic,    // exchange name
		"fanout", // exchange type (use "fanout" for pub/sub)
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return err
	}

	r.logger.Debug("rabbitmq created exchange", zap.String("topic", topic))

	r.mu.Lock()
	r.handlers[topic] = handler
	r.mu.Unlock()

	queue, err := r.channel.QueueDeclare(
		fmt.Sprintf("%s-%s", r.config.QueuePrefix, topic), // name
		true,  // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	r.logger.Debug("rabbitmq created queue", zap.String("queue", queue.Name))

	err = r.channel.QueueBind(
		queue.Name, // queue name
		"",         // routing key
		topic,      // exchange
		false,
		nil,
	)
	if err != nil {
		return err
	}

	r.logger.Debug("rabbitmq created queue bind", zap.String("topic", topic), zap.String("queue", queue.Name))

	msgs, err := r.channel.Consume(
		queue.Name,
		"",
		true,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	r.logger.Debug("rabbitmq start consumer", zap.String("topic", topic), zap.String("queue", queue.Name))

	go func() {
		for d := range msgs {
			headers := map[string]string{}
			for k, v := range d.Headers {
				if str, ok := v.(string); ok {
					headers[k] = str
				}
			}

			msg := Message{
				Headers: headers,
				Data:    d.Body,
			}

			handler(msg)
		}
	}()

	return nil
}

func (r *rabbitMQ) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
	return nil
}
