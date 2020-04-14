package main

import (
	"context"
	"fmt"
	"log"

	"github.com/grpcLearning/greet/greetpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("hello I'm a client")
	//try to connect the server and return the connection
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	//a magic in go, the following line will execute at very very end of the program
	defer cc.Close()

	//binding the connection to the client
	c := greetpb.NewGreetServiceClient(cc)

	doUnary(c)

}

func doUnary(c greetpb.GreetServiceClient) {
	//Step 2: Define the request
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "zhai",
			LastName:  "zhonghao",
		},
	}

	//Step 1 : To call the Greet the function in the client
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling RPC %v", err)
	}

	//Step 3: print the result of the req
	log.Printf("Response from the Greeting %v", res.Result)
}
