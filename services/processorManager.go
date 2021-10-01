package services

import (
	"github.com/jeanmolossi/super-duper-adventure/domain"
	"github.com/jeanmolossi/super-duper-adventure/queue"
	"github.com/streadway/amqp"
)

type ProcessorManager struct {
	MessageChannel chan amqp.Delivery
	RabbitMQ       *queue.RabbitMQ
	Channel        *amqp.Channel
}

type ProcessedCourseResult struct {
	Course  domain.Course
	Message amqp.Delivery
}

func NewProcessorManager(messageChannel chan amqp.Delivery, rabbit *queue.RabbitMQ, channel *amqp.Channel) *ProcessorManager {
	return &ProcessorManager{
		MessageChannel: messageChannel,
		RabbitMQ:       rabbit,
		Channel:        channel,
	}
}

func (p *ProcessorManager) Run() {
	studentsChannel := make(chan ProcessedCourseResult)
	endChannel := make(chan string)

	courseProcessor := NewCourseProcessor(p.MessageChannel, p.RabbitMQ, studentsChannel)
	courseProcessor.Start(p.Channel)

	studentProcessor := NewStudentProcessor(studentsChannel)

	for parallelCourses := 0; parallelCourses < 20; parallelCourses++ {
		go studentProcessor.Start(endChannel)
	}

	for result := range endChannel {
		if result == "all-done" {
			close(endChannel)
		}
	}
}
