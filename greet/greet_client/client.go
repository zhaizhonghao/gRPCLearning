package main

import (
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

	fmt.Printf("Created Client:%f", c)
}
