package main

import (
	"context"
	"fmt"
	"log"

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

	doUnary(c)

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
