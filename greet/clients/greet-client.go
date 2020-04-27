package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	greetpb "github.com/grpcLearning/greet/pb"

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

	//doUnary(c)

	//doServerStreaming(c)

	//doClientStreaming(c)

	doBiDiStream(c)

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

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a server streaming RPC...")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "zhai",
			LastName:  "zhonghao",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)

	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
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
		log.Printf("Resonse from GreetManyTimes %v", msg.GetResult())
	}

}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a client streaming RPC...")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "zhai",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "liu",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "weng",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "zhang",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())

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
		log.Fatalf("error while receiving response from longGreet: %v", err)
	}

	fmt.Printf("LongGreet Response: %v\n", res)
}

func doBiDiStream(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do BiDi streaming RPC...")

	requests := []*greetpb.GreetEveryOneRequest{
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "zhai",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "liu",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "weng",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "zhang",
			},
		},
	}
	//we create a stream by invoking the client
	stream, err := c.GreetEveryOne(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream:%v", err)
		return
	}
	waitc := make(chan struct{})
	//we send a bunch of messages to the client (go routine)
	go func() {
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	//we receive a bunch of messages from the client (go routine)
	go func() {
		//function to receive a bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitc)
	}()
	//block until everything is done
	<-waitc
}
