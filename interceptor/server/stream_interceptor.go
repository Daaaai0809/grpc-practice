package interceptor

import (
	"log"
	"io"
	"errors"

	"google.golang.org/grpc"
)

func MyStreamServerInterceptor1(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Println("[pre] my stream server interceptor 1: ", info.FullMethod)
	err := handler(srv, &myServerStreamWrapper1{ss})
	log.Println("[post] my stream server interceptor 1: ", err)
	return err
}

func MyStreamServerInterceptor2(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Println("[pre] my stream server interceptor 2: ", info.FullMethod)
	err := handler(srv, &myServerStreamWrapper2{ss})
	log.Println("[post] my stream server interceptor 2: ", err)
	return err
}

type myServerStreamWrapper1 struct {
	grpc.ServerStream
}

func (s *myServerStreamWrapper1) RecvMsg(m interface{}) error {
	err := s.ServerStream.RecvMsg(m)
	if errors.Is(err, io.EOF) {
		log.Println("[my server stream wrapper 1] recv EOF")
	}
	return err
}

func (s *myServerStreamWrapper1) SendMsg(m interface{}) error {
	log.Println("[my server stream wrapper 1] send message: ", m)
	return s.ServerStream.SendMsg(m)
}

type myServerStreamWrapper2 struct {
	grpc.ServerStream
}

func (s *myServerStreamWrapper2) RecvMsg(m interface{}) error {
	err := s.ServerStream.RecvMsg(m)
	if errors.Is(err, io.EOF) {
		log.Println("[my server stream wrapper 2] recv EOF")
	}
	return err
}

func (s *myServerStreamWrapper2) SendMsg(m interface{}) error {
	log.Println("[my server stream wrapper 2] send message: ", m)
	return s.ServerStream.SendMsg(m)
}