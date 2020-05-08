package main

import (
	"fmt"
	"log"

	"golang.org/x/net/context"

	"github.com/grpcLearning/blog/blogpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	fmt.Println("Blog Client")
	tls := false

	opts := grpc.WithInsecure()

	if tls {
		//Certificate Authority Trust Certificate
		certFile := "ssl/ca.crt"
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		if sslErr != nil {
			log.Fatalf("Error while loading CA trust certificate:%v", sslErr)
			return
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	//try to connect the server and return the connection without ssl
	//cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	//try to connect the server and return the connection with ssl
	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	//a magic in go, the following line will execute at very very end of the program
	defer cc.Close()

	//binding the connection to the client
	c := blogpb.NewBlogServiceClient(cc)

	fmt.Println("Creating the blog")

	blog := &blogpb.Blog{
		AuthorId: "Zhai",
		Title:    "My second Blog",
		Content:  "Content of the second blog",
	}

	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})

	if err != nil {
		log.Fatalf("Unexpected error : %v", err)
	}
	fmt.Printf("Blog has been created %v", createBlogRes)
}
