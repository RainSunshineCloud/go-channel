# go-channel
这是一个消息队列框架
### 说明
- channel 通道，Register 注册一个通道，并指定对应的接收者
- provider 提供者，触发通道的程序
- receiver 接收者，通过通道处理程序，同一个通道，同一时刻，只会运行一个
- OutMessageDataInterface 传出消息接口
- InMessageDataInterface 传入消息接口

### 用法
```
var manager channel.Manager;

func main () {
	manager = channel.Manager{}
	model_receiver := &modelReceiver{}
	manager.Register("key1",model_receiver)
	manager.Register("key2",model_receiver)
	go manager.Run(&timerProvider{})
	go manager.Run(&httpProvider{})
	for {
		time.Sleep(100 * time.Second)
	}
}

//接收者运行
type modelReceiver struct {
}

func (this *modelReceiver) Run(data receiver.InMessageDataInterface) receiver.OutMessageDataInterface {
	da,err := data.(*MessageData)
	fmt.Println(da,err);
	res := &OutMessageData{
		NowCode:[]uint{9},
	}
	return res;
}

//提供者运行
type httpProvider struct {
}

func (this *httpProvider) Run () {
	http.HandleFunc("/",handler)
	http.ListenAndServe("0.0.0.0:8000",nil)
}

func handler (w http.ResponseWriter, r *http.Request) {
	byts,_ := ioutil.ReadAll(r.Body)
	data := &MessageData{}
	(*data).InHandler(byts)
	res,err := manager["key1"].Receiver(data)
	by := res.OutHandler(err)
	fmt.Fprint(w,string(by))

}

//提供者timer
type timerProvider struct {

}

func (this timerProvider) Run () {
	for {
		data := &MessageData{}
		(*data).InHandler([]byte(`{"NowCode":[1,2,3],"NowNum":"123123","NowTime":1573116316}`))
		res,err := manager["key1"].Receiver(data)
		by := res.OutHandler(err)
		time.Sleep(10*time.Second)
		fmt.Println(string(by))
	}
}



//传送数据
type MessageData struct {
	NowCode []uint
}

type OutMessageData struct {
	NowCode []uint
}

func (this *MessageData) InHandler(bytes []byte) error {
	err := json.Unmarshal(bytes,this)
	if err != nil {
		return err;
	}
	return nil;
}

func (this *OutMessageData) OutHandler(err error) []byte {
	if err != nil {
		return []byte(err.Error())
	}
	by,err := json.Marshal(this)
	if err != nil {
		return []byte(err.Error())
	}
	return by;
}

```
