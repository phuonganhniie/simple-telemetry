package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/phuonganhniie/simpleTelemetry/gateway/handler"
	pb_checkout "github.com/phuonganhniie/simpleTelemetry/proto"
	"github.com/phuonganhniie/simpleTelemetry/utils"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func main() {
	serviceName := "gateway"
	jaegerAdd := utils.EnvString("JAEGER_ADDRESS", "localhost")
	jaegerPort := utils.EnvString("JAEGER_PORT", "6831")
	checkOutAdd := utils.EnvString("CHECKOUT_SERVICE_ADDRESS", "localhost:8080")
	httpAdd := utils.EnvString("HTTP_ADDRESS", ":8081")

	err := utils.SetGlobalTracer(serviceName, jaegerAdd, jaegerPort)
	if err != nil {
		log.Fatalf("failed to create tracer: %v", err)
	}

	// add grpc
	// propagate the trace via gRPC to the checkout service
	conn, err := grpc.Dial(
		checkOutAdd,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb_checkout.NewCheckoutClient(conn)

	// http config
	router := http.NewServeMux()
	router.HandleFunc("/api/checkout", handler.CheckOutHandler(c))

	fmt.Println("HTTP server listening at port", httpAdd)
	log.Fatal(http.ListenAndServe(httpAdd, router))
}
