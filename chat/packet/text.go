package packet

type TextPacket struct {
	SenderID string
	Contents string
}

func (p *TextPacket) Type() string {
	return IDText
}
