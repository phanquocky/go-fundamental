package main

import (
	grpcjwt "Btc/go-fundamental/grpc/authentication-jwt/code-gen"
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {

	// accessToken := oauth.TokenSource{TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: ""})}
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.WithPerRPCCredentials(accessToken),
	}
	client, err := grpc.NewClient("localhost:8080", opts...)
	if err != nil {
		panic(err)
	}

	authenTicationService := grpcjwt.NewAuthenticationServiceClient(client)

	res, err := authenTicationService.Login(context.Background(), &grpcjwt.LoginRequest{
		Username: "admin",
		Password: "admin",
	})
	if err != nil {
		panic(err)
	}

	println("TOKEN: ", res.GetToken())

	token := res.GetToken()
	token = strings.Replace(token, "Bearer ", "", 1)
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", token))

	res1, err := authenTicationService.ValidateToken(ctx, &grpcjwt.Token{
		Token: token,
	})

	if err != nil {
		panic(err)
	}

	println("VALID: ", res1.GetValid())
}
