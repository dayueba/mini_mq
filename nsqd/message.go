package minimq

const (
	MsgIDLength = 16
)

type MessageID [MsgIDLength]byte

type Message struct {
	ID MessageID
}
