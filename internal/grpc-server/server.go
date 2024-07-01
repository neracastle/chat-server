package grpc_server

import (
	"context"
	"github.com/brianvoe/gofakeit"
	userdesc "github.com/neracastle/chat-server/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

type Server struct {
	userdesc.UnimplementedChatV1Server
}

func (s *Server) Create(ctx context.Context, req *userdesc.CreateRequest) (*userdesc.CreateResponse, error) {
	log.Printf("called Create method with req: %v", req)

	return &userdesc.CreateResponse{Id: gofakeit.Int64()}, nil
}

func (s *Server) Delete(ctx context.Context, req *userdesc.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("called Delete method with req: %v", req)

	return &emptypb.Empty{}, nil
}

func (s *Server) SendMessage(ctx context.Context, req *userdesc.SendMessageRequest) (*emptypb.Empty, error) {
	log.Printf("called SendMessage method with req: %v", req)

	return &emptypb.Empty{}, nil
}
