package client

import (
	"sync"
	"sync/atomic"
)

type producerConn interface {
	WriteCommand(*Command) error
}

type Producer struct {
	config Config
	mu     sync.Mutex

	id int64

	addr string

	conn producerConn

	// 生产者 router goroutine 接收服务端响应的 channel
	responseChan chan []byte
	// 生产者 router goroutine 接收错误的 channel
	errorChan chan []byte
	// 生产者 router goroutine 接收生产者关闭信号的 channel
	closeChan chan int
	// 生产者 router goroutine 接收 publish 指令的 channel
	transactionChan chan *ProducerTransaction
	transactions    []*ProducerTransaction

	// 生产者的状态
	state int32
	// 当前有多少 publish 指令并发执行
	concurrentProducers int32
	// 生产者 router goroutine 接收退出信号的 channel
	exitChan chan int

	stopFlag int32
}

func (p *Producer) Publish(topic string, data []byte) error {
	return p.sendCommand(Publish(topic, data))
}

func (p *Producer) sendCommand(cmd *Command) error {
	doneChan := make(chan *ProducerTransaction)
	err := p.sendCommandAsync(cmd, doneChan, nil)
	if err != nil {
		close(doneChan)
		return err
	}
	t := <-doneChan
	return t.Error
	// 常见的异步写法
	// doneChan = make(chan *Message)
	// err = fn(doneChan)
	// if err != nil {
	// 	close(doneChan)
	// 	return err
	// }
	// t := <-doneChan
	// return t.Error
}

func (p *Producer) sendCommandAsync(cmd *Command, doneChan chan *ProducerTransaction,
	args []interface{}) error {
	// keep track of how many outstanding producers we're dealing with
	// in order to later ensure that we clean them all up...
	atomic.AddInt32(&p.concurrentProducers, 1)
	defer atomic.AddInt32(&p.concurrentProducers, -1)

	// 确认当前连接状态
	// if atomic.LoadInt32(p.state) != StateConnected {
	// 	err := w.connect()
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	t := &ProducerTransaction{
		cmd:      cmd,
		doneChan: doneChan,
		Args:     args,
	}

	select {
	case p.transactionChan <- t:
	case <-p.exitChan:
		return ErrStopped
	}

	return nil
}

func (p *Producer) connect() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// double check
	if atomic.LoadInt32(&p.stopFlag) == 1 {
		return ErrStopped
	}

	state := atomic.LoadInt32(&p.state)
	switch {
	case state == StateConnected:
		return nil
	// 如果建立连接失败过 就直接报错
	case state != StateInit:
		return ErrNotConnected
	}

	p.conn = NewConn(p.addr, &p.config)

	atomic.StoreInt32(&p.state, StateConnected)
	p.closeChan = make(chan int)

	// 发送消息给nsq
	go p.router()

	return nil
}

func (p *Producer) router() {
	for {
		select {
		case t := <-p.transactionChan:
			p.transactions = append(p.transactions, t)
			err := p.conn.WriteCommand(t.cmd)
			if err != nil {
			}
			// ...
		// ...
		case <-p.closeChan:
			break
		case <-p.exitChan:
			break
		}
	}

	// clean up
}

type ProducerTransaction struct {
	cmd      *Command
	doneChan chan *ProducerTransaction
	Error    error         // the error (or nil) of the publish command
	Args     []interface{} // the slice of variadic arguments passed to PublishAsync or MultiPublishAsync
}
