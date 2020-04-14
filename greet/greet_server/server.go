package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/grpcLearning/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

//step 1 definte the service of server(server is the 'server' defined above)
func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	//to print the req
	fmt.Printf("Greet function was invoked with %v", req)

	//extract the information from the requset
	firstName := req.GetGreeting().GetFirstName()

	//form the response
	result := "hello" + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

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
