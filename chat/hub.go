package chat

import (
	"fmt"
	"github.com/aurorachat/jwt-tokens/tokens"
	"github.com/aurorachat/ws-sessions/sessions"
	"github.com/gorilla/websocket"
	"strconv"
)

type Hub struct {
	store sessions.Store
}

func NewHub() *Hub {
	return &Hub{sessions.NewStore()}
}

func (h *Hub) ConnectClient(conn *websocket.Conn) {
	_, bytesMsg, err := conn.ReadMessage()
	if err != nil {
		_ = conn.Close()
		return
	}

	claimsOrNil, err := tokens.ValidateToken(string(bytesMsg))

	if err != nil {
		_ = conn.Close()
		return
	}

	claims := *claimsOrNil

	userId := strconv.Itoa(int(claims["sub"].(float64)))

	s := h.store.GetSession(userId)

	if s == nil {
		s = sessions.NewSession(userId)
		h.store.SetSession(userId, s)
		go h.StartListening(s)
	}

	s.RegisterConnection(claims["sessionId"].(string), conn)
}

func (h *Hub) StartListening(s *sessions.Session) {
	for {
		connId, data := s.Receive()
		fmt.Println(fmt.Sprintln(connId, " said that ", data))
	}
}
