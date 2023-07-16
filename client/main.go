package main

import (
	"io"
	"os"
	"fmt"
	"log"
	"bufio"
	"errors"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/credentials/insecure"
	_ "google.golang.org/genproto/googleapis/rpc/errdetails"

	proto "github.com/Daaaai0809/grpc-practice/proto/proto"
	protogrpc "github.com/Daaaai0809/grpc-practice/proto/proto/protogrpc"
	interceptor "github.com/Daaaai0809/grpc-practice/interceptor/client"
)

var (
	scanner *bufio.Scanner
	client protogrpc.HelloServiceClient
)

func main() {
	fmt.Println("Client is running...")

	scanner = bufio.NewScanner(os.Stdin)

	address := "localhost:8081"
	conn, err := grpc.Dial(
		address, 
		grpc.WithUnaryInterceptor(interceptor.MyUnaryClientInterceptor1),
		grpc.WithStreamInterceptor(interceptor.MyStreamClientInterceptor1),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return
	}
	defer conn.Close()

	client = protogrpc.NewHelloServiceClient(conn)

	for {
		fmt.Println("1: send Request")
		fmt.Println("2: HelloServerStream")
		fmt.Println("3: HelloClientStream")
		fmt.Println("4: HelloBiStream")
		fmt.Println("5: exit")
		fmt.Print("> ")
		
		scanner.Scan()
		in := scanner.Text()

		switch in {
			case "1":
				Hello()
			case "2":
				HelloServerStream()
			case "3":
				HelloClientStream()
			case "4":
				HelloBiStream()
			case "5":
				fmt.Println("Bye!")
				return
		}
	}
}

func Hello() {
	fmt.Println("Enter your name: ")
	scanner.Scan()
	name := scanner.Text()

	req := &proto.HelloRequest{
		Name: name,
	}

	res, err := client.Hello(context.Background(), req)
	if err != nil {
		if stat, ok := status.FromError(err); ok {
			fmt.Printf("Error Code: %v \n", stat.Code())
			fmt.Printf("Error Message: %v \n", stat.Message())
			fmt.Printf("Error Details: %v \n", stat.Details())
			return
		} else {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Println(res.GetMessage())
	}
}

func HelloServerStream() {
	fmt.Println("Enter your name: ")
	scanner.Scan()
	name := scanner.Text()

	req := &proto.HelloRequest{
		Name: name,
	}

	stream, err := client.HelloServerStream(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling HelloServerStream RPC: %v", err)
		return
	}

	for {
		res, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("all the responses have already received.")
			return
		}

		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
			return
		}

		fmt.Println(res)
	}
}

func HelloClientStream() {
	stream, err := client.HelloClientStream(context.Background())
	if err != nil {
		log.Fatalf("error while calling HelloClientStream RPC: %v", err)
		return
	}

	sendCount := 5
	fmt.Println("Enter the names: ")
	for i := 0; i < sendCount; i++ {
		scanner.Scan()
		name := scanner.Text()

		if err := stream.Send(&proto.HelloRequest{
			Name: name,
		}); err != nil {
			log.Fatalf("error while sending stream: %v", err)
			return
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while closing stream: %v", err)
		return
	}

	fmt.Println(res.GetMessage())
}

func HelloBiStream() {
	stream, err := client.HelloBiStream(context.Background())
	if err != nil {
		log.Fatalf("error while calling HelloBiStream RPC: %v", err)
		return
	}

	sendNum := 5
	fmt.Println("Enter the names: ")

	var sendEnd, recvEnd bool
	sendCount := 0
	for !(sendEnd && recvEnd) {
		if !sendEnd {
			scanner.Scan()
			name := scanner.Text()

			sendCount++
			if err := stream.Send(&proto.HelloRequest{
				Name: name,
			}); err != nil {
				log.Fatalf("error while sending stream: %v", err)
				sendEnd = true
			}

			if sendCount == sendNum {
				sendEnd = true
				if err := stream.CloseSend(); err != nil {
					log.Fatalf("error while closing stream: %v", err)
				}
			}
		}
		if !recvEnd {
			if res, err := stream.Recv(); err != nil {
				if !errors.Is(err, io.EOF) {
					log.Fatalf("error while reading stream: %v", err)
				}
				recvEnd = true
			} else {
				fmt.Println(res.GetMessage())
			}
		}
	}
}