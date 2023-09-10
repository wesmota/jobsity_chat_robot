package main

import (
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
	"github.com/wesmota/go-jobsity-chat-robot/rabbitmq"
)

func main() {
	// RabbitMQ
	rmqHost := os.Getenv("RMQ_HOST")
	rmqUserName := os.Getenv("RMQ_USERNAME")
	rmqPassword := os.Getenv("RMQ_PASSWORD")
	rmqPort := os.Getenv("RMQ_PORT")
	dsn := "amqp://" + rmqUserName + ":" + rmqPassword + "@" + rmqHost + ":" + rmqPort + "/"

	conn, err := amqp.Dial(dsn)
	if err != nil {
		panic(err)
	}
	log.Info().Msg("Connected to RabbitMQ")
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	log.Info().Msg("Connected to RabbitMQ Channel")
	br := &rabbitmq.Broker{}
	br.Setup(ch)
	go br.Read()
	select {}
}
