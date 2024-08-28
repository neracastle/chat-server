package grpc_server

import (
	"github.com/neracastle/auth/pkg/user_v1/auth"
	syserr "github.com/neracastle/go-libs/pkg/sys/error"

	userdesc "github.com/neracastle/chat-server/pkg/chat_v1"
)

// Connect подключает пользователя к чату
func (s *Server) Connect(req *userdesc.ConnectRequest, stream userdesc.ChatV1_ConnectServer) error {
	s.m.RLock()
	existChat, ok := s.connectedChats[req.GetChatId()]
	s.m.RUnlock()
	if !ok {
		return syserr.New("Чат не найден", syserr.NotFound)
	}

	tokenUser := auth.UserFromContext(stream.Context())
	existChat.Connect(tokenUser.ID, stream)
	s.metrics.IncreaseClients()
	for {
		select {
		case <-stream.Context().Done():
			existChat.Disconnect(tokenUser.ID)
			s.metrics.DecreaseClients()
			return stream.Context().Err()
		}
	}
}
