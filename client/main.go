package main

import (
	"context"
	"github.com/taybart/log"
	"time"

	// "github.com/google/uuid"
	"github.com/taybart/todo/list"
	pb "github.com/taybart/todo/proto"
	"google.golang.org/grpc"
)

const (
	address   = "localhost:8080"
	defaultID = 1
)

func main() {
	tl, err := list.NewTodo("./todo.db")
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewListClient(conn)

	// Contact the server and print out its response.

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Push(ctx, &pb.String{Contents: "test"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Info("Greeting:", r.Message)
	sy, err := c.Sync(ctx, &pb.NoParams{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Infof("%v\n", sy.Items)
	for _, i := range sy.Items {
		tl.PushItem(&list.Item{
			Contents: i.Contents,
			IsDone:   i.IsDone,
		})
	}
}
