package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	grpcserver "github.com/sadwanj/virtuecare-ar-services/sharedlib/grpc"
	//pb "your_project/proto" // your generated protobuf package
)

// Example service implementation
// type MyService struct {
// 	pb.UnimplementedMyServiceServer
// }

// Implement your RPC methods here
// func (s *MyService) SayHello(...) { ... }

// This is the RegisterFunc
func registerServices(s *grpc.Server) {
	//pb.RegisterMyServiceServer(s, &MyService{})
}

func main() {
	config := &grpcserver.Config{
		Host: "0.0.0.0",
		Port: "50051",
	}

	server, err := grpcserver.NewServer(config, registerServices)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	go func() {
		log.Println("Starting gRPC server...")
		if err := server.Start(); err != nil {
			log.Fatalf("server failed: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan

	log.Println("Shutting down...")
	server.Stop()
	log.Println("Stopped.")
}
