package main

//go:generate protoc -I ../proto --go_out=plugins=grpc:../proto ../proto/todo.proto

import (
	"fmt"
	"github.com/taybart/log"
	"github.com/taybart/todo/proto"
	"google.golang.org/grpc"
	"net"
)

const (
	port = ":8080"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 14586))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := todo.Server{}
	s.Init("./test.db")
	grpcServer := grpc.NewServer()
	// attach the Ping service to the server
	todo.RegisterListServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	} else {
		log.Printf("Server started successfully")
	}
}
