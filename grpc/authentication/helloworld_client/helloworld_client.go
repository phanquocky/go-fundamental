package main

import (
	"context"
	helloworld "go-fundamental/grpc/code_gen"

	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

func main() {

	perRPC := oauth.TokenSource{TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "AccessToken"})}
	cred, err := credentials.NewClientTLSFromFile("data/server_cert.pem", "*.test.example.com")
	if err != nil {
		panic("cannot load ca_cert.pem" + err.Error())
	}
	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(perRPC),
		grpc.WithTransportCredentials(cred),
	}

	client, err := grpc.NewClient("localhost:8080", opts...)
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
