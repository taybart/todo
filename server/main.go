package main

//go:generate protoc -I ../proto --go_out=plugins=grpc:../proto ../proto/todo.proto

import (
	"context"
	"fmt"
	"github.com/taybart/log"
	"github.com/taybart/todo/list"
	pb "github.com/taybart/todo/proto"
	"google.golang.org/grpc"
	"net"
)

var tl *list.Todo

const (
	port = ":8080"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) Toggle(ctx context.Context, in *pb.ToggleRequest) (*pb.ToggleReply, error) {
	log.Info("Received:", in.Id)
	msg := fmt.Sprintf("Toggle %d", in.Id)
	return &pb.ToggleReply{Message: msg}, nil
}

func (s *server) Sync(ctx context.Context, in *pb.SyncParams) (*pb.List, error) {
	log.Info("Sync")
	items := []*pb.Item{}
	for _, i := range tl.Items {
		items = append(items, &pb.Item{Id: i.ID, IsDone: i.IsDone, Contents: i.Contents})
	}
	return &pb.List{Items: items}, nil
}

func main() {
	var err error
	tl, err = list.NewTodo()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTodolistServer(s, &server{})
	log.Info("Running on", port, "...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
