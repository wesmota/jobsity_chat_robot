package rabbitmq

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
	"github.com/wesmota/go-jobsity-chat-robot/csvprocessor"
)

type Broker struct {
	ReceiverQueue  amqp.Queue
	PublisherQueue amqp.Queue
	Channel        *amqp.Channel
}

type ChatMessage struct {
	Type        int    `json:"type"`
	ChatMessage string `json:"chatmessage"`
	ChatUser    string `json:"chatuser"`
	ChatRoomId  uint   `json:"chatroomId"`
}

type MessageResponse struct {
	RoomId  uint   `json:"RoomId"`
	Message string `json:"Message"`
}

// Setup creates(or connects if not existing) the reciever and publisher queues
func (b *Broker) Setup(ch *amqp.Channel) {
	//based on https://www.rabbitmq.com/tutorials/tutorial-one-go.html

	receiverQueue := "JOBSITY_RECEIVER"
	publisherQueue := "JOBSITY_PUBLISHER"

	qR, err := ch.QueueDeclare(
		receiverQueue, // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		log.Info().Msg("Receiver queue already exists")
		return
	}

	qP, err := ch.QueueDeclare(
		publisherQueue, // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		log.Info().Msg("Publisher queue already exists")
		return
	}

	b.ReceiverQueue = qR
	b.PublisherQueue = qP
	b.Channel = ch
}

// Publish sends messages to receiver queue
func (b *Broker) Publish(message []byte) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := b.Channel.PublishWithContext(ctx,
		"",                    // exchange
		b.PublisherQueue.Name, // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})

	if err != nil {
		log.Err(err).Msg("Failed to publish a message")
	}
	log.Info().Msg("Published a message")
}

// Read reads messages from receiver queue
func (b *Broker) Read() {
	entries, err := b.Channel.Consume(
		b.ReceiverQueue.Name, // queue
		"",                   // consumer
		true,                 // auto-ack
		false,                // exclusive
		false,                // no-local
		false,                // no-wait
		nil,                  // args
	)
	if err != nil {
		log.Printf("ReadMessages Error occured %s\n", err)
		return
	}
	log.Info().Msg("Reading messages")
	receivedMessages := make(chan ChatMessage)
	go toMsgResponse(entries, receivedMessages)
	go processAndPublish(receivedMessages, b)

}

func toMsgResponse(entries <-chan amqp.Delivery, receivedMessages chan ChatMessage) {
	var msg ChatMessage
	for d := range entries {
		log.Info().Msgf("Received a message: %s", d.Body)
		err := json.Unmarshal([]byte(d.Body), &msg)
		if err != nil {
			log.Printf("Error on received request : %s ", err)
			continue
		}
		log.Info().Msgf("Received a message: %+v", msg)
		receivedMessages <- msg
	}
}

func processAndPublish(msgs <-chan ChatMessage, b *Broker) {
	for m := range msgs {
		log.Info().Msgf("Processing message %s for room %d", m.ChatMessage, m.ChatRoomId)
		items := strings.Split(m.ChatMessage, "=")
		if len(items) != 2 {
			log.Info().Msg("Invalid message format")
			continue
		}
		skey := items[1]
		quoteValue := csvprocessor.ProcessCSVStockFile(skey)
		log.Info().Msgf("quoteValue: %s", quoteValue)
		// send message to publisher queue
		body, err := json.Marshal(ChatMessage{
			Type:        1,
			ChatMessage: quoteValue,
			ChatUser:    "robot",
			ChatRoomId:  m.ChatRoomId,
		})

		if err != nil {
			log.Err(err).Msg("Failed to marshal message")
			continue
		}
		go b.Publish(body)
		log.Info().Msg("Message sent to publisher queue")
	}
}
