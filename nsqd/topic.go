package minimq

type Topic struct {
	name string
	

	channelMap        map[string]*Channel
	idFactory         *guidFactory // 消息 id 生成器

}

func NewTopic(topicName string ) *Topic {
	return nil
}