package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/neracastle/chat-server/internal/config"
	chatsrv "github.com/neracastle/chat-server/internal/grpc-server"
	chatdesc "github.com/neracastle/chat-server/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {

	cfg := config.MustLoad()

	conn, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))

	if err != nil {
		log.Fatal(color.RedString("failed to serve grpc server: %v", err))
	}

	log.Print(color.GreenString("ChatAPI grpc server listening on: %s", conn.Addr().String()))

	gsrv := grpc.NewServer()
	reflection.Register(gsrv)

	chatdesc.RegisterChatV1Server(gsrv, &chatsrv.Server{})

	if err = gsrv.Serve(conn); err != nil {
		log.Fatal(color.RedString("failed to serve grpc server: %v", err))
	}
}
