package client

import (
	"io"
	"net"
	"sync"
)

type Conn struct {
	addr string
	conn *net.TCPConn

	mtx sync.Mutex

	r io.Reader // 读缓冲区
	w io.Writer // 写缓冲区

	cmdChan  chan *Command     // 客户端请求指令
	respChan chan *msgResponse // 通过此 channel 接收来自客户端的响应，发往服务端
	exitChan     chan int          // 接收退出指令

	// 记录了有多少消息处于未 ack 状态
	// inFlight 某条消息已经被客户端接收到，但是还未给予服务端 ack 的状态.
	messagesInFlight int64

	config *Config
}

// 处理来自服务端的响应
// 在 Conn.readLoop() 方法中，读取到来自服务端的响应，通过调用 Producer.onConnResponse() 方法，将数据发送到 Producer.responseChan 当中
func (c *Conn) readloop() {

}

// 处理发往服务端的请求
func (c *Conn) writeloop() {

}

func (c *Conn) Connect() {
	dialer := &net.Dialer{
		LocalAddr: c.config.LocalAddr,
		Timeout:   c.config.DialTimeout,
	}
	// 通过 net 包向服务端发起 tcp 连接
	conn, err := dialer.Dial("tcp", c.addr)
	if err != nil {
		panic(err)
	}
	c.conn = conn.(*net.TCPConn)
	// 设置读写缓冲区
	c.r = conn
	c.w = conn

	// 开启两个 goroutine，分别处理来自服务端的响应和发往服务端的请求
	go c.readloop()
	go c.writeloop()
}

func NewConn(addr string, config *Config) *Conn {
	return &Conn{
		addr: addr,

		config: config,

		cmdChan:         make(chan *Command),
		respChan: make(chan *msgResponse),
		exitChan:        make(chan int),
	}
}

func (c *Conn) WriteCommand(cmd *Command) error {
	return nil
}

type msgResponse struct {
}
