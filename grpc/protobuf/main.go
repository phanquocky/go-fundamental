package main

import (
	protobuf "Btc/go-fundamental/grpc/protobuf/code-gen"
	"encoding/json"
	"fmt"
)

func main() {
	req := &protobuf.Request{
		Message:     "Hello",
		Code:        200,
		Status:      "OK",
		Description: "This is a description",
		Metadata:    "This is a metadata",
		Data:        "This is a data, ",
	}

	grpcBytes := []byte(req.String())
	fmt.Println("protobuf len", len(grpcBytes))

	jsonBytes, _ := json.Marshal(req)
	fmt.Println("json len", len(jsonBytes))

	
}
