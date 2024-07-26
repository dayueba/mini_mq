package client

type Handler interface {
	HandleMessage(message *Message) error
}

type Consumer struct {
	id      int64
	topic   string
	channel string

	conns            map[string]*Conn
	incomingMessages chan *Message
}

func (r *Consumer) SetHandler(handler Handler) {
	go r.handlerLoop(handler)
}

func (r *Consumer) handlerLoop(handler Handler) {
	for {
		select {
		case msg := <-r.incomingMessages:
			err := handler.HandleMessage(msg)
			if err != nil {
				// ...
			}
		}
	}
}

func (r *Consumer) ConnectToNSQD(addr string) error {
	return nil
}
