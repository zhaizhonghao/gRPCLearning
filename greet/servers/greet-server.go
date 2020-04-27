package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	greetpb "github.com/grpcLearning/greet/pb"
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

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "hello " + firstName + " number" + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		//use stream to send the client many times
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet function was invoked with streaming ")
	result := "Hello"
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			//we have finished reading the client stream
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream : %v", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		result += firstName + "! "
	}
}

func (*server) GreetEveryOne(stream greetpb.GreetService_GreetEveryOneServer) error {
	fmt.Printf("GreetEveryOne function was invoked with streaming request\n")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client streaming %v", err)
			return nil
		}
		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName
		sendErr := stream.Send(&greetpb.GreetEveryOneResponse{
			Result: result,
		})
		if err != nil {
			log.Fatalf("Error while sending data to client %v", sendErr)
			return sendErr
		}
	}
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
