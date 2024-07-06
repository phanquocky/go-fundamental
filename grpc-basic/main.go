// Package help you understand all of about grpc using Go
// You first install tools
// go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
// go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
// brew install protobuf

/*
	 Prerequisites
	 1. First create helloworld.proto file for create format to communicate between client and server
	 2. Second using protoc to generate code for client and server
	  ```
	 		#the source_relative flag is used to relative path for generated code
	    $ protoc --go_out=./code_gen --go_opt=paths=source_relative --go-grpc_out=./code_gen --go-grpc_opt=paths=source_relative helloworld.proto
		```
*/

// `go run helloworld_server/hellworld_server.go` to run grpc server
// `go run helloworld_client/helloworld_client.go` to run grpc client
// See the result: `Message:  Hello John`

// Congratulations! You have just created a simple gRPC server and client using Go.
package main

import "fmt"

func main() {
	fmt.Println("hello world")
}
