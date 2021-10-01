package main

import (
	"github.com/jeanmolossi/super-duper-adventure/queue"
	"github.com/jeanmolossi/super-duper-adventure/services"
	"github.com/streadway/amqp"
)

func main() {
	messageChannel := make(chan amqp.Delivery)

	rabbitMQ := queue.NewRabbitMQ()
	channel := rabbitMQ.Connect()
	defer channel.Close()

	rabbitMQ.Consume(messageChannel)

	processor := services.NewProcessorManager(messageChannel, rabbitMQ, channel)
	processor.Run()
}
