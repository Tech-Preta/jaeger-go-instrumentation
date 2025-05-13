package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func initTracer() func() {
	jaegerEndpoint := os.Getenv("JAEGER_ENDPOINT")
	if jaegerEndpoint == "" {
		log.Fatal("JAEGER_ENDPOINT environment variable is required")
	}

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)))
	if err != nil {
		log.Fatalf("Failed to create Jaeger exporter: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("go-example-service"),
		)),
	)

	otel.SetTracerProvider(tp)
	tracer = tp.Tracer("go-example-service")

	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatalf("Error shutting down tracer provider: %v", err)
		}
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "helloHandler")
	defer span.End()

	// Simulando algum processamento
	time.Sleep(100 * time.Millisecond)

	span.SetAttributes(attribute.String("event", "hello"))
	span.SetAttributes(attribute.String("message", "Hello, World!"))

	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	cleanup := initTracer()
	defer cleanup()

	r := mux.NewRouter()
	r.HandleFunc("/", helloHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
} 