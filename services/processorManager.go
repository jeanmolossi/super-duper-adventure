package services

import (
	"fmt"
	"github.com/jeanmolossi/super-duper-adventure/queue"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"log"
	"os"
	"strconv"
)

type ProcessorManager struct {
	MessageChannel  chan amqp.Delivery
	ReturnChannel   chan ProcessedCourseResult
	totalConsumers  int
	channelPrefetch int
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func NewProcessorManager() *ProcessorManager {
	return &ProcessorManager{}
}

func (pm *ProcessorManager) Start() {
	totalConsumers, _ := strconv.Atoi(os.Getenv("TOTAL_OF_CONSUMERS"))
	channelPrefetch, _ := strconv.Atoi(os.Getenv("PREFETCH_BY_CONSUMER"))

	pm.MessageChannel = make(chan amqp.Delivery)
	pm.ReturnChannel = make(chan ProcessedCourseResult)
	pm.totalConsumers = totalConsumers
	pm.channelPrefetch = channelPrefetch

	pm.HeapUpConsumers()

	for result := range pm.ReturnChannel {
		if result.Error != nil {
			log.Printf("%s concluido com erro: %s", result.Course.ID, result.Error.Error())
			continue
		}
	}
}

func (pm *ProcessorManager) HeapUpConsumers() {
	for channelConsumers := 0; channelConsumers < pm.totalConsumers; channelConsumers++ {
		rabbitMQ := queue.NewRabbitMQ()
		rabbitMQ.ConsumerName = fmt.Sprintf("%s-%d", os.Getenv("RABBITMQ_CONSUMER_NAME"), channelConsumers+1)

		channel := rabbitMQ.Connect()
		err := channel.Qos(pm.channelPrefetch, 0, false)
		if err != nil {
			log.Fatalf("Falha ao definir Qos")
		}

		rabbitMQ.Consume(pm.MessageChannel)
		pm.RunProcessors()
	}
}

func (pm *ProcessorManager) RunProcessors() {
	for processors := 0; processors < pm.channelPrefetch; processors++ {
		go func() {
			courseProcessor := NewCourseProcessor(pm.MessageChannel, pm.ReturnChannel)
			courseProcessor.Start()
		}()
	}
}
