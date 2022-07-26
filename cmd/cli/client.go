package cli

import (
	"context"
	"log"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=client-gen.conf  ../../books.yaml
////////go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=client-types-gen.conf  ../../books.yaml

func main() {

	initTracer()

	clientOA3 := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	client, err := NewClientWithResponses("http://0.0.0.0:9234", WithHTTPClient(&clientOA3))
	if err != nil {
		log.Fatalf("Failed instantiate client: %s", err)
	}

	for i := 0; i < 10; i++ {
		_, err := client.SearchBooksByAuthor(context.Background(),
			&SearchBooksByAuthorParams{
				AuthorId: 1,
			})
		if err != nil {
			log.Fatalf("Couldn't create task %s", err)
		}
	}
}

const jaegerEndpoint = "http://localhost:14268/api/traces"

func initTracer() {

	jaegerExporter, err := jaeger.New(
		jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)),
	)
	if err != nil {
		log.Fatalln("Failed initialize exporter", err)
	}

	// Create stdout exporter to be able to retrieve the collected spans.
	_, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatalln("Couldn't initialize exporter", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(jaegerExporter),
		trace.WithResource(resource.NewSchemaless(attribute.KeyValue{
			Key:   semconv.ServiceNameKey,
			Value: attribute.StringValue("rest-server"),
		})),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
}
