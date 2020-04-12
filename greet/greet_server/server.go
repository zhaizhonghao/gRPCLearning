package main

import (
	"fmt"
	"log"
	"net"

	"github.com/grpcLearning/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

func main() {
	fmt.Println("hello world")

	//create a listener containing the protocol and the port to listen
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}
	//create teh grpc sever
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	//binding the listener and server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server %v", err)
	}

}
