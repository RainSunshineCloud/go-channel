package receiver

//消息数据(需要请自定义)
type InMessageDataInterface interface {
	InHandler([]byte) error
}

type OutMessageDataInterface interface {
	OutHandler(error) []byte
}


type ReceiverInterface interface {
	Run(data InMessageDataInterface) OutMessageDataInterface
}
