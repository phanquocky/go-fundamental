// This Package follow the same structure as https://grpc.io/docs/languages/go/basics/
// protoc --go_out=./code-gen --go_opt=paths=source_relative --go-grpc_out=./code-gen --go-grpc_opt=paths=source_relative routeGuide.proto
package main

import "fmt"

func main() {
	fmt.Println("hello world")
}
