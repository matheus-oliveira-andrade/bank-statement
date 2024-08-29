package broker

import (
	"log/slog"

	"github.com/spf13/viper"

	amqp "github.com/rabbitmq/amqp091-go"
)

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
