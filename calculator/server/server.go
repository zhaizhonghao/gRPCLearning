package main

import (
	"context"
	"fmt"
	"io"
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

/**
k = 2
N = 210
while N > 1:
	if N % k == 0:   // if k evenly divides into N
		print k      // this is a factor
		N = N / k    // divide N by k so that we have the rest of the number left.
	else:
		k = k + 1
*/
func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("PrimeNumberDecomposition function was invoked with %v\n", req)
	num := req.GetNum()

	k := 2
	for {
		if num <= 1 {
			break
		}
		if num%int32(k) == 0 {
			res := &calculatorpb.PrimeNumberDecompositionResponse{
				Prime: int32(k),
			}
			stream.Send(res)
			num = num / int32(k)
		} else {
			k = k + 1
		}
	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("ComputeAverage function was invoked with streaming ")
	result := float32(0)
	counter := float32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			//we have finished reading the client stream
			fmt.Printf("counter %v", counter)
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Result: result / counter,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream : %v", err)
		}
		result += req.GetNum()
		counter += float32(1)
	}
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
