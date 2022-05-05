package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb_checkout "github.com/phuonganhniie/simpleTelemetry/proto"
	"github.com/phuonganhniie/simpleTelemetry/utils"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
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

	ch, closeConn := utils.ConnectAmqp(amqpUser, amqpPass, amqpHost, amqpPort)
	defer closeConn()

	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()))

	pb_checkout.RegisterCheckoutServer(s, &server{channel: ch})

	log.Printf("GRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type server struct {
	pb_checkout.UnimplementedCheckoutServer
	channel *amqp.Channel
}

func (s *server) DoCheckOut(ctx context.Context, rq *pb_checkout.CheckOutRequest) (*pb_checkout.CheckOutResponse, error) {
	messageName := "checkout.processed"

	// Create a new span (child of the trace id) to inform the publishing of the message
	tr := otel.Tracer("amqp")
	amqpContext, messageSpan := tr.Start(ctx, fmt.Sprintf("AMQP - publish - %s", messageName))
	defer messageSpan.End()

	// Inject the context in the headers
	headers := utils.InjectAMQPHeaders(amqpContext)
	msg := amqp.Publishing{Headers: headers}
	err := s.channel.Publish("exchange", messageName, false, false, msg)
	if err != nil {
		log.Fatal(err)
	}

	response := &pb_checkout.CheckOutResponse{TotalAmount: 1234}

	// Log specific events for a span
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("response: %v", response))

	return response, nil
}
