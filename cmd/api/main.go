package main

import (
	"os"
	"os/signal"
	"fmt"
	"log"
	"net"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	proto "github.com/Daaaai0809/grpc-practice/proto/proto"
	protogrpc "github.com/Daaaai0809/grpc-practice/proto/proto/protogrpc"
)

type myServer struct {
	protogrpc.UnimplementedHelloServiceServer
}

func (s *myServer) Hello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error) {
	return &proto.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}

func NewMyServer() *myServer {
	return &myServer{}
}

func main() {
	port := 8081
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	protogrpc.RegisterHelloServiceServer(s, NewMyServer())

	reflection.Register(s)

	go func() {
		log.Printf("starting server on port %d", port)
		s.Serve(listener)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	s.GracefulStop()
}