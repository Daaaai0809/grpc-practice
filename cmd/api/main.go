package main

import (
	"io"
	"os"
	"os/signal"
	"fmt"
	"log"
	"net"
	"time"
	"errors"
	"strings"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"
	// "google.golang.org/genproto/googleapis/rpc/errdetails"

	interceptor "github.com/Daaaai0809/grpc-practice/interceptor/server"
	proto "github.com/Daaaai0809/grpc-practice/proto/proto"
	protogrpc "github.com/Daaaai0809/grpc-practice/proto/proto/protogrpc"
)

type myServer struct {
	protogrpc.UnimplementedHelloServiceServer
}

func (s *myServer) Hello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error) {
	// stat := status.New(codes.Unknown, "unknown error has ocurred")
	// stat, _ = stat.WithDetails(&errdetails.DebugInfo{
	// 	Detail: "detail reason of err",
	// })
	// err := stat.Err()
	return &proto.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}

func (s *myServer) HelloServerStream(req *proto.HelloRequest, stream protogrpc.HelloService_HelloServerStreamServer) error {
	resCount := 10
	for i := 0; i < resCount; i++ {
		if err := stream.Send(&proto.HelloResponse{
			Message: fmt.Sprintf("Hello, %s! [%d]", req.GetName(), i),
		}); err != nil {
			return err
		}

		time.Sleep(1 * time.Second)
	}
	return nil
}

func (s *myServer) HelloClientStream(stream protogrpc.HelloService_HelloClientStreamServer) error {
	nameList := make([]string, 0)
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			message := fmt.Sprintf("Hello, %s!", strings.Join(nameList, ", "))
			return stream.SendAndClose(&proto.HelloResponse{
				Message: message,
			})
		}
		if err != nil {
			return err
		}
		nameList = append(nameList, req.GetName())
	}
}

func (s *myServer) HelloBiStream(stream protogrpc.HelloService_HelloBiStreamServer) error {
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}

		message := fmt.Sprintf("Hello, %s!", req.GetName())

		if err := stream.Send(&proto.HelloResponse{
			Message: message,
		}); err != nil {
			return err
		}
	}
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

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptor.MyUnaryServerInterceptor1, interceptor.MyUnaryServerInterceptor2),
		grpc.ChainStreamInterceptor(interceptor.MyStreamServerInterceptor1, interceptor.MyStreamServerInterceptor2),
	)

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