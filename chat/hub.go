package chat

import (
	"github.com/aurorachat/backend/chat/packet"
	"github.com/aurorachat/jwt-tokens/tokens"
	"github.com/aurorachat/ws-sessions/sessions"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"strconv"
)

type Hub struct {
	store sessions.Store
}

type PacketPayload struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
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
		go h.StartListening(userId, s)
	}

	s.RegisterConnection(claims["sessionId"].(string), conn)
	s.Send(CreatePacketPayload(&packet.ServerboundChatInitializedPacket{}), claims["sessionId"].(string))
}

func (h *Hub) StartListening(userId string, s *sessions.Session) {
	generalC := h.store.GetChannel("general")

	if generalC == nil {
		generalC = sessions.NewChannel("general")
		h.store.SetChannel("general", generalC)
	}

	s.Subscribe(generalC)

	for {
		_, _, msgBytes := s.Receive()

		var pkPayload PacketPayload

		err := json.Unmarshal(msgBytes, &pkPayload)
		if err != nil {
			continue
		}

		if pkPayload.Type == packet.IDText {
			var textPk packet.TextPacket
			err = mapstructure.Decode(pkPayload.Data, &textPk)
			if err != nil {
				continue
			}
			generalC.Broadcast(CreatePacketPayload(&textPk))
		}
	}
}

func CreatePacketPayload(pk packet.Packet) PacketPayload {
	return PacketPayload{
		Type: pk.Type(),
		Data: pk,
	}
}
