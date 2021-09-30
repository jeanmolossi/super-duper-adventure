package services

import (
	"github.com/jeanmolossi/super-duper-adventure/queue"
	"github.com/streadway/amqp"
	"log"
)

type CourseProcessor struct {
	MessageChannel chan amqp.Delivery
	RabbitMQ       *queue.RabbitMQ
}

func NewCourseProcessor(messageChannel chan amqp.Delivery, rabbitMQ *queue.RabbitMQ) *CourseProcessor {
	return &CourseProcessor{
		MessageChannel: messageChannel,
		RabbitMQ:       rabbitMQ,
	}
}

func (c *CourseProcessor) Start(channel *amqp.Channel) {
	workReturnChannel := make(chan string)

	for processes := 0; processes < 2; processes++ {
		go WorkOn(c.MessageChannel, workReturnChannel)
	}

	for processWork := range workReturnChannel {
		log.Printf("Process status: %s", processWork)
	}

}

func WorkOn(messageChannel chan amqp.Delivery, workReturnChannel chan string) {
	for course := range messageChannel {
		log.Printf("Curso: %s", string(course.Body))

		err := course.Ack(false)
		if err != nil {
			workReturnChannel <- "end"
			return
		}

		workReturnChannel <- "end"
	}
}
