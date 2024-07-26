package client

type Command struct {
	Name   []byte
	Params [][]byte
	Body   []byte
}

func Publish(topic string, body []byte) *Command {
    var params = [][]byte{[]byte(topic)}
    return &Command{[]byte("PUB"), params, body}
}