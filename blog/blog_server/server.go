package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/grpcLearning/blog/blogpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct {
}

func main() {
	//if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Blog Service Started")

	//create a listener containing the protocol and the port to listen
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}
	tls := false
	opts := []grpc.ServerOption{}
	if tls {
		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"
		//to import the certificate(reference from: https://grpc.io/docs/guides/auth/#authentication-api)
		creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
		if sslErr != nil {
			log.Fatalf("Failed loading certificates:%v", sslErr)
			return
		}
		//create teh grpc sever
		opts = append(opts, grpc.Creds(creds))
	}

	//without ssl
	//s := grpc.NewServer()

	//with ssl
	s := grpc.NewServer(opts...)

	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		fmt.Println("Starting the Server...")
		//binding the listener and server
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to server %v", err)
		}
	}()
	//Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	//Block util a signal is received
	<-ch
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Closing the listener")
	lis.Close()
	fmt.Println("End of listener")
}
