package main

import (
	"context"
	helloworld "go-fundamental/grpc/code_gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cred := insecure.NewCredentials() // this cred is used for testing purpose, it will not check the server's certificate
	client, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(cred))
	if err != nil {
		panic("cannot create conn to localhost:8080" + err.Error())
	}

	greeterService := helloworld.NewGreeterClient(client)

	res, err := greeterService.SayHello(context.Background(), &helloworld.HelloRequest{Name: "John"})
	if err != nil {
		panic("cannot call SayHello" + err.Error())
	}

	println("Message: ", res.Message)
}
