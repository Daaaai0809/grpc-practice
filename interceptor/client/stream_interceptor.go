package interceptor

import (
	"io"
	"fmt"
	"log"
	"errors"
	"context"

	grpc "google.golang.org/grpc"
)

func MyStreamClientInterceptor1(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	fmt.Println("[pre] my stream client interceptor 1:", method)
	stream, err := streamer(ctx, desc, cc, method, opts...)
	return &MyStreamClientWrapper1{stream}, err
}

type MyStreamClientWrapper1 struct {
	grpc.ClientStream
}

func (s *MyStreamClientWrapper1) RecvMsg(m interface{}) error {
	err := s.ClientStream.RecvMsg(m)
	if !errors.Is(err, io.EOF) {
		log.Println("[post] my stream client interceptor 1: ", m)
	}

	return err
}

func (s *MyStreamClientWrapper1) SendMsg(m interface{}) error {
	log.Println("[post] my stream client interceptor 1: ", m)
	return s.ClientStream.SendMsg(m)
}

func (s *MyStreamClientWrapper1) CloseSend() error {
	err := s.ClientStream.CloseSend()
	log.Println("[post] my stream client interceptor 1")
	return err
}