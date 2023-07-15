package main

import (
	"os"
	"fmt"
	"log"
	"bufio"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	proto "github.com/Daaaai0809/grpc-practice/proto/proto"
	protogrpc "github.com/Daaaai0809/grpc-practice/proto/proto/protogrpc"
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
		fmt.Println("2: exit")
		fmt.Print("> ")
		
		scanner.Scan()
		in := scanner.Text()

		switch in {
			case "1":
				Hello()
			case "2":
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
		log.Fatalf("error while calling Hello RPC: %v", err)
		return
	} else {
		fmt.Println(res.GetMessage())
	}
}