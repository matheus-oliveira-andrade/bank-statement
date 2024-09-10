package broker

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/matheus-oliveira-andrade/bank-statement/statement-service/shared/events"
	"github.com/spf13/viper"

	amqp "github.com/rabbitmq/amqp091-go"
)

type BrokerInterface interface {
	Produce(eventPublish *events.EventPublish, configs *ProduceConfigs) error
}

type RabbitMQBroker struct {
	url string
}

func NewBroker(url string) BrokerInterface {
	return &RabbitMQBroker{
		url: url,
	}
}

func (b *RabbitMQBroker) Produce(eventPublish *events.EventPublish, configs *ProduceConfigs) error {
	conn, err := amqp.Dial(b.url)
	if err != nil {
		slog.Error("Failed to connect to RabbitMQ", "error", err)
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		slog.Error("Failed to open a channel", "error", err)
		return err
	}
	defer ch.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(eventPublish)
	if err != nil {
		slog.Error("Error marshaling event", "error", err)
		return err
	}

	var exchange string
	if configs != nil {
		exchange = configs.Topic
	}

	err = ch.PublishWithContext(ctx,
		exchange,
		eventPublish.Type,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	if err != nil {
		slog.Error("Failed to publish a message", "error", err)
		return err
	}

	return nil
}

type ProduceConfigs struct {
	Topic string
}

func NewConnection(url string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		slog.Error("Failed to connect to RabbitMQ", "error", err)
		return nil, err
	}

	return conn, nil
}

func BuildConnectionUrl() string {
	user := viper.GetString("broker.user")
	password := viper.GetString("broker.password")
	host := viper.GetString("broker.host")
	port := viper.GetString("broker.port")
	protocol := viper.GetString("broker.protocol")

	return protocol + "://" + user + ":" + password + "@" + host + ":" + port + "/"
}
