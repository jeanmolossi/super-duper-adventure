package services

import (
	"github.com/jeanmolossi/super-duper-adventure/domain"
	"github.com/streadway/amqp"
	"log"
	"time"
)

type ProcessedCourseResult struct {
	Course  domain.Course
	Message amqp.Delivery
	Error   error
}

type CourseProcessor struct {
	MessageChannel chan amqp.Delivery
	ReturnChannel  chan ProcessedCourseResult
}

func NewCourseProcessor(messageChannel chan amqp.Delivery, returnChannel chan ProcessedCourseResult) *CourseProcessor {
	return &CourseProcessor{
		MessageChannel: messageChannel,
		ReturnChannel:  returnChannel,
	}
}

func (c *CourseProcessor) Start() {
	for courseMessage := range c.MessageChannel {
		startedProcess := time.Now()
		log.Printf("Processando curso: %s", string(courseMessage.Body))

		course := domain.NewCourse()
		err := course.GetCourse(string(courseMessage.Body))

		processedResult := ProcessedCourseResult{
			Course:  *course,
			Message: courseMessage,
			Error:   err,
		}

		if err != nil {
			log.Printf("Mensagem rejeitada: %s", err.Error())
			err := courseMessage.Reject(false)
			if err != nil {
				return
			}

			processedResult.Error = err
			c.ReturnChannel <- processedResult

			return
		}

		students, err := course.GetStudents()
		if err != nil {
			log.Printf("Mensagem rejeitada: %s", err.Error())
			err := courseMessage.Reject(true)
			if err != nil {
				return
			}

			processedResult.Error = err
			c.ReturnChannel <- processedResult

			return
		}

		log.Printf("Antes de processar students do curso %s: %s", course.ID, time.Since(startedProcess))

		err = c.ProcessStudents(students)
		if err != nil {
			log.Printf("Processamento de estudantes com erro: %s", err.Error())

			mustRequeue := err.Error() != "Curso já existe no redis"
			err := courseMessage.Reject(mustRequeue)
			if err != nil {
				return
			}

			processedResult.Error = err
			c.ReturnChannel <- processedResult

			return
		}

		log.Printf(
			"Processamento %s terminado em %v. (Total %d alunos)\n",
			processedResult.Course.ID,
			time.Since(startedProcess),
			len(students),
		)

		err = courseMessage.Ack(false)
		if err != nil {
			return
		}
		c.ReturnChannel <- processedResult
	}
}

func (c *CourseProcessor) ProcessStudents(students []domain.Student) error {
	studentsProcessor := NewStudentProcessor()
	studentsProcessor.Students = students
	err := studentsProcessor.Start()

	return err
}
