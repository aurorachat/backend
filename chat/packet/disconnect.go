package packet

type DisconnectPacket struct {
	Reason string
}

func (pkt *DisconnectPacket) Type() string {
	return IDDisconnect
}
