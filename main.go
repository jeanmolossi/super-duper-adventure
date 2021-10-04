package main

import (
	"github.com/jeanmolossi/super-duper-adventure/services"
	"github.com/joho/godotenv"
	"log"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func main() {
	processor := services.NewProcessorManager()
	processor.Start()
}
