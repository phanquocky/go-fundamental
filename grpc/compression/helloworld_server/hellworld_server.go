package main

import (
	"context"
	"fmt"
	helloworld "go-fundamental/grpc/code_gen"
	"net"

	"errors"

	_ "google.golang.org/grpc/encoding/gzip" // Install the gzip compressor

	"google.golang.org/grpc"
)

type server struct {
	helloworld.UnimplementedGreeterServer // this is a must add	mustEmbedUnimplementedGreeterServer() function, you can try remove this line and see the error
}

func newServer() *server {
	return &server{}
}

// SayHello implements SayHello method from GreeterServer
// this function will be called when client call SayHello method
// it will reply with `Hello + name`
func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: "Hello " + in.Name}, nil
}

// main function will create new server and run it
func main() {
	server := newServer()

	// opts := []grpc.ServerOption{
	// 	grpc.UseCompressor(gzip.Name),
	// }
	grpcServer := grpc.NewServer()
	helloworld.RegisterGreeterServer(grpcServer, server)

	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(errors.New("cannot create listener" + err.Error()))
	}

	fmt.Println("Server is running on port :8080")
	grpcServer.Serve(listen)
}
