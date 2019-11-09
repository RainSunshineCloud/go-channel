package channel

import (
	"github.com/RainSunshineCloud/go-channel/provider"
	"github.com/RainSunshineCloud/go-channel/receiver"
)

type Channel struct {
	receiver receiver.ReceiverInterface
	InChan chan *InData
}

//传入信息
type InData struct {
	Data *receiver.InMessageDataInterface
	OutChannel chan receiver.OutMessageDataInterface
}

//生成信息传送通道
func NewMessage(receiver receiver.ReceiverInterface) *Channel {
	return &Channel{
		receiver: receiver,
		InChan:make(chan *InData),
	}
}

// 运行
func (this *Channel) run () {
	for {
		res := <- this.InChan; //等待
		res.OutChannel <- this.receiver.Run(*res.Data)
	}
}

//实例化输入数据
func newInData (data receiver.InMessageDataInterface) *InData {
	return &InData{
		&data,
		make(chan receiver.OutMessageDataInterface),
	}
}


//接收
func (this *Channel) Receiver (input_data receiver.InMessageDataInterface) (receiver.OutMessageDataInterface,error) {
		data_handle := newInData(input_data)
		this.InChan <- data_handle
		data := <- data_handle.OutChannel
		return data,nil;
}



//消息管理器
type Manager map[string]*Channel;


//初始化消息处理器
func (this Manager) start() {
	for _,m := range this {
		go m.run();
	}
}

//注册消息接收后的处理程序
func (this Manager) Register (message_key string, receiver receiver.ReceiverInterface) {
	this[message_key]  = NewMessage(receiver);
}

func (this Manager) Run (provider provider.Provider) {
	this.start()
	provider.Run()
}
