// JWT authentication server
// The flow of jwt
// 1. The client sends a request to the server with a username and password.
// 2. The server checks the username and password.
// 3. If the username and password are correct, the server creates a token and sends it back to the client.
// 4. The client saves the token and sends it with every request.
// 5. The server checks the token and allows the client to access the resources.

package main

import (
	grpcjwt "Btc/go-fundamental/grpc/authentication-jwt/code-gen"
	"context"
	"fmt"
	"strings"
	"time"

	"net"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// more long key
var key = []byte("my_secret_key")

type server struct {
	grpcjwt.UnimplementedAuthenticationServiceServer
}

func newServer() *server {
	return &server{}
}

type payload struct {
	Username string
	jwt.StandardClaims
}

func (s *server) ValidateToken(ctx context.Context, req *grpcjwt.Token) (*grpcjwt.TokenValidationResponse, error) {
	return &grpcjwt.TokenValidationResponse{
		Valid: true,
	}, nil
}

func (s *server) Login(ctx context.Context, req *grpcjwt.LoginRequest) (*grpcjwt.LoginResponse, error) {
	payload := &payload{
		Username: req.GetUsername(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	fmt.Println("payload: ", payload.Username, payload.ExpiresAt)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	accessToken, err := token.SignedString([]byte(key))
	if err != nil {
		return nil, err
	}

	return &grpcjwt.LoginResponse{
		Token: accessToken,
	}, nil
}

func main() {
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(jwtUnaryInterceptor),
	}
	grpcServer := grpc.NewServer(opts...)

	handler := newServer()
	grpcjwt.RegisterAuthenticationServiceServer(grpcServer, handler)

	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	fmt.Println("server is running on port 8080")
	grpcServer.Serve(listen)
}

func jwtUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info.FullMethod == "/grpc.authentication.jwt.AuthenticationService/Login" {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	authentication := md.Get("authorization")
	if len(authentication) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
	}

	tokenString := authentication[0]
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	payload := &payload{}
	token, err := jwt.ParseWithClaims(tokenString, payload, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "token is not valid")
	}

	if !token.Valid {
		return nil, status.Error(codes.Unauthenticated, "token is not valid")
	}

	fmt.Println("payload: ", payload.Username, payload.ExpiresAt)

	return handler(ctx, req)
}
