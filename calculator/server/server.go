package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/grpcLearning/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

//step 1 definte the service of server(server is the 'server' defined above)
func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	//to print the req
	fmt.Printf("Sum function was invoked with %v", req)

	//extract the information from the requset
	firstNum := req.GetFirstNum()
	secondNum := req.GetSecondNum()

	//form the response
	result := firstNum + secondNum
	res := &calculatorpb.SumResponse{
		SumResult: result,
	}
	return res, nil
}

func main() {
	fmt.Println("Calculator server start!")

	//create a listener containing the protocol and the port to listen
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}
	//create teh grpc sever
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	//binding the listener and server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server %v", err)
	}
}
