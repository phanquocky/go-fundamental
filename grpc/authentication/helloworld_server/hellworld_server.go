package main

import (
	"context"
	"crypto/tls"
	"fmt"
	helloworld "go-fundamental/grpc/code_gen"
	"log"
	"net"
	"strings"

	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
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
	// err := alts.ClientAuthorizationCheck(ctx, []string{"foo@iam.gserviceaccount.com"})
	// if err != nil {
	// fmt.Println("cannot authorize client: ", err)
	// return nil, err
	// }
	return &helloworld.HelloReply{Message: "Hello " + in.Name}, nil
}

// main function will create new server and run it
func main() {
	server := newServer()

	cert, err := tls.LoadX509KeyPair("data/server_cert.pem", "data/server_key.pem")
	if err != nil {
		log.Fatalf("failed to load key pair: %s", err)
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(ensureValidToken),
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
	}
	grpcServer := grpc.NewServer(opts...)

	helloworld.RegisterGreeterServer(grpcServer, server)

	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(errors.New("cannot create listener" + err.Error()))
	}

	fmt.Println("Server is running on port :8080")
	grpcServer.Serve(listen)
}

func getAccesstoken(authorization []string) string {
	fmt.Println(authorization)
	token := authorization[0]
	// remove prefix Bearer
	accessToken := strings.TrimPrefix(token, "Bearer ")
	return accessToken
}

func isValidToken(accessToken string) bool {
	return accessToken == "AccessToken"
}

func ensureValidToken(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	fmt.Println("INterceptor: ensureValidToken")
	md, ok := metadata.FromIncomingContext(ctx)
	fmt.Println("metadata: ", md, ok)
	if !ok {
		return nil, errors.New("cannot get metadata")
	}

	authorization := md.Get("authorization")
	if len(authorization) == 0 {
		return nil, errors.New("cannot get authorization token")
	}

	if !isValidToken(getAccesstoken(authorization)) {
		return nil, errors.New("invalid token")
	}

	fmt.Println("Congratulations, you have valid token")
	// check token

	return handler(ctx, req)
}
