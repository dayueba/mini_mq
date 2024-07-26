package minimq

type NSQD struct {
	// 节点内递增的客户端 id 序号，为每个到来的客户端请求分配唯一的标识 id
	clientIDSequence int64

	topicMap map[string]*Topic


	exitChan             chan int
}
