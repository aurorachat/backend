package packet

import "github.com/aurorachat/backend/chat/packet/types"

type ServerboundChatInitializedPacket struct {
	channels []types.Channel
}

func (p *ServerboundChatInitializedPacket) Type() string {
	return IDServerboundChatInitialized
}
