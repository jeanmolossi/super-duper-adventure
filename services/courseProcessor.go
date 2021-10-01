package services

import (
	"github.com/jeanmolossi/super-duper-adventure/domain"
	"github.com/jeanmolossi/super-duper-adventure/queue"
	"github.com/streadway/amqp"
	"log"
)

type CourseProcessor struct {
	MessageChannel  chan amqp.Delivery
	RabbitMQ        *queue.RabbitMQ
	StudentsChannel chan ProcessedCourseResult
}

func NewCourseProcessor(messageChannel chan amqp.Delivery, rabbitMQ *queue.RabbitMQ, studentsChannel chan ProcessedCourseResult) *CourseProcessor {
	return &CourseProcessor{
		MessageChannel:  messageChannel,
		RabbitMQ:        rabbitMQ,
		StudentsChannel: studentsChannel,
	}
}

func (c *CourseProcessor) Start(channel *amqp.Channel) {
	for processes := 0; processes < 10; processes++ {
		go WorkOn(c.MessageChannel, c.StudentsChannel)
	}
}

func WorkOn(messageChannel chan amqp.Delivery, studentsProcessorChannel chan ProcessedCourseResult) {
	for courseMessage := range messageChannel {
		log.Printf("Processando curso: %s", string(courseMessage.Body))

		courseInstance := domain.NewCourse()
		course, err := courseInstance.GetCourse(string(courseMessage.Body))
		if err != nil {
			log.Printf("Mensagem rejeitada: %s", err.Error())
			courseMessage.Reject(false)
			continue
		}

		studentsProcessorChannel <- ProcessedCourseResult{
			Course:  *course,
			Message: courseMessage,
		}
	}
}
