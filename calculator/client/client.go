package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/grpcLearning/calculator/calculatorpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("hello I'm a calculator client")
	//try to connect the server and return the connection
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	//a magic in go, the following line will execute at very very end of the program
	defer cc.Close()

	//binding the connection to the client
	c := calculatorpb.NewCalculatorServiceClient(cc)

	//doUnary(c)
	//doServerStreaming(c)
	doClientStreaming(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	//Step 2: Define the request
	req := &calculatorpb.SumRequest{
		FirstNum:  1,
		SecondNum: 6,
	}

	//Step 1 : To call the Greet the function in the client
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling RPC %v", err)
	}

	//Step 3: print the result of the req
	log.Printf("Response from the Calculator %v", res.SumResult)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a server streaming RPC...")

	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Num: 120,
	}

	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)

	if err != nil {
		log.Fatalf("error while calling PrimeNumberDecomposition RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			//We've reached the end of the stream
			break
		}

		if err != nil {
			log.Fatalf("error while reading stream %v", err)
		}
		log.Printf("Resonse from GreetManyTimes %v", msg.GetPrime())
	}

}
func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a client streaming RPC...")

	requests := []*calculatorpb.ComputeAverageRequest{
		&calculatorpb.ComputeAverageRequest{
			Num: 1,
		},
		&calculatorpb.ComputeAverageRequest{
			Num: 2,
		},
		&calculatorpb.ComputeAverageRequest{
			Num: 3,
		},
		&calculatorpb.ComputeAverageRequest{
			Num: 4,
		},
	}

	stream, err := c.ComputeAverage(context.Background())

	if err != nil {
		log.Fatalf("error while calling LongGreet : %v", err)
	}

	//we iterate over our slice and send each message individually
	for _, req := range requests {
		fmt.Printf("Sending req:%v", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from ComputeAverage: %v", err)
	}

	fmt.Printf("ComputeAverage Response: %v\n", res)

}
