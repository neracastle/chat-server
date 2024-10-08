package grpc_server

import (
	"context"
	"time"

	"github.com/neracastle/auth/pkg/user_v1/auth"
	syserr "github.com/neracastle/go-libs/pkg/sys/error"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	chatdesc "github.com/neracastle/chat-server/pkg/chat_v1"
)

// SendMessage отправляет сообщение в чат
func (s *Server) SendMessage(ctx context.Context, req *chatdesc.SendMessageRequest) (*emptypb.Empty, error) {
	s.m.RLock()
	existChat, ok := s.connectedChats[req.GetChatId()]
	s.m.RUnlock()

	if ok {
		tokenUser := auth.UserFromContext(ctx)
		existChat.AddMessage(&chatdesc.Message{
			From:      tokenUser.ID,
			Text:      req.GetText(),
			Timestamp: timestamppb.New(time.Now()),
		})
		return &emptypb.Empty{}, nil
	}

	return nil, syserr.New("Чат не найден", syserr.NotFound)
}
