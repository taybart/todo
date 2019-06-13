package todo

//go:generate protoc -I . --go_out=plugins=grpc:. ./todo.proto

import (
	"context"
	"fmt"
	"github.com/taybart/log"
	"github.com/taybart/todo/list"
)

// Server is used to implement protocol
type Server struct {
	tl *list.Todo
}

// Init asdf
func (s *Server) Init(db ...string) error {
	var err error
	name := ""
	if len(db) > 0 && db[0] != "" {
		name = db[0]
	}
	s.tl, err = list.NewTodo(name)
	return err
}

// Toggle asdf
func (s *Server) Toggle(ctx context.Context, in *IDRequest) (*SimpleReply, error) {
	log.Info("Toggle:", in.Id)
	msg := fmt.Sprintf("Toggle %s", in.Id)
	return &SimpleReply{Message: msg}, nil
}

// Sync sadf
func (s *Server) Sync(ctx context.Context, in *NoParams) (*List, error) {
	log.Info("Sync")
	items := []*Item{}
	for _, i := range s.tl.Items {
		items = append(items, &Item{Id: i.ID.String(), IsDone: i.IsDone, Contents: i.Contents})
	}
	return &List{Items: items}, nil
}

// Push asdf
func (s *Server) Push(ctx context.Context, in *String) (*SimpleReply, error) {
	s.tl.Push(in.Contents)
	return &SimpleReply{Message: "done"}, nil
}

// PushItem asdf
func (s *Server) PushItem(ctx context.Context, in *Item) (*SimpleReply, error) {
	s.tl.Push(in.Contents)
	return &SimpleReply{Message: "done"}, nil
}

// Del asdf
func (s *Server) Del(ctx context.Context, in *IDRequest) (*SimpleReply, error) {
	return &SimpleReply{Message: "done"}, nil
}

// Edit asdf
func (s *Server) Edit(ctx context.Context, in *IDRequest) (*SimpleReply, error) {
	return &SimpleReply{Message: "done"}, nil
}

// Current asdf
func (s *Server) Current(ctx context.Context, in *NoParams) (*Item, error) {
	i := s.tl.Current()
	return &Item{Id: i.ID.String(), IsDone: i.IsDone, Contents: i.Contents}, nil
}
