package main

import (
	"context"
	"contrib.go.opencensus.io/exporter/stackdriver"
	pb "go-grpc-template/grpc-gen-circleci-template"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"time"
)

const (
	port = "0.0.0.0:50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("ctx: %v", ctx)
	ctx, span := trace.StartSpan(ctx, "grpc-template.server")//, trace.WithSampler(trace.AlwaysSample()))
	defer span.End()
	log.Printf("Received: %v", in.Name)
	time.Sleep(80 * time.Millisecond)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	// Create and register a OpenCensus Stackdriver Trace exporter.
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: os.Getenv("GOOGLE_CLOUD_PROJECT"),
	})
	if err != nil {
		log.Fatal(err)
	}
	trace.RegisterExporter(exporter)

	/*
	// Register the views to collect server request count.
	if err := view.Register(ocgrpc.DefaultServerViews...); err != nil {
		log.Fatal(err)
	}
	*/
	// Configure 100% sample rate, otherwise, few traces will be sampled.
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Set up a new server with the OpenCensus
	// stats handler to enable stats and tracing.
	s := grpc.NewServer(grpc.StatsHandler(new(ocgrpc.ServerHandler)))
	pb.RegisterGreeterServer(s, &server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
