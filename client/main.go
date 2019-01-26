package main

import (
	"context"
	"contrib.go.opencensus.io/exporter/stackdriver"
	pb "go-grpc-template/grpc-gen-circleci-template"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"
)

const (
	//address     = "localhost:50051"
	address = "34.73.89.238:50051"
	defaultName = "world"
)

func main() {
	// Create and register a OpenCensus Stackdriver Trace exporter.
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		//ProjectID: os.Getenv("GOOGLE_CLOUD_PROJECT"),
		ProjectID: "sansigma-infra",
	})
	if err != nil {
		log.Fatal(err)
	}
	trace.RegisterExporter(exporter)

	// Configure 100% sample rate, otherwise, few traces will be sampled.
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	//ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	//defer cancel()

	// Create a span with the background context, making this the parent span.
	// A span must be closed.
	bctx := context.Background()
	ctx, span := trace.StartSpan(bctx, "grpc-template.client", trace.WithSampler(trace.AlwaysSample()))
	time.Sleep(80 * time.Millisecond)
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("trace id: %v", ctx)
	log.Printf("Greeting: %s", r.Message)
	span.End()
}
