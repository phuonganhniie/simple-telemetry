package main

import (
	"log"

	"github.com/phuonganhniie/simpleTelemetry/utils"
)

func main() {
	serviceName := "checkout"
	jaegerAddress := utils.EnvString("JAEGER_ADDRESS", "localhost")
	jaegerPort := utils.EnvString("JAEGER_PORT", "6831")
	grpcAddress := utils.EnvString("GRPC_ADDRESS", "localhost:8080")
	amqpUser := utils.EnvString("RABBITMQ_USER", "guest")
	amqpPass := utils.EnvString("RABBITMQ_PASS", "guest")
	amqpHost := utils.EnvString("RABBITMQ_HOST", "localhost")
	amqpPort := utils.EnvString("RABBITMQ_PORT", "5672")

	err := utils.SetGlobalTracer(serviceName, jaegerAddress, jaegerPort)
	if err != nil {
		log.Fatalf("failed to create tracer: %v", err)
	}

	channel, closeConn := utils.ConnectAmqp()
}
