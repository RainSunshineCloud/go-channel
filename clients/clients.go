package clients

import (
	"io"
	"runtime"
)
type Sign uint8

const STARTREQ Sign = 1;
const STOPREQ Sign = 2;
const SLEEP Sign = 3;


type Client struct {
	log LoggerInterface
	requests map[uint]RequestInterface
	real_client ReqClientInterface
	status uint8
	signers chan uint
}

func (this *Client) Stop () {
	for _,request := range this.requests {
		request.SetSigner(STOPREQ);
	}

	for signer := range this.signers {
		this.log.Info().Println("停止请求:",signer);
	}

	runtime.GC();
}

func (this *Client) Run () {

	sign,err := this.real_client.BeforeRun()
	if sign == STOPREQ {
		this.log.Error().Println(err);
		return;
	}

	for _,request := range this.requests {
		go this.runRequest(request)
	}

	sign,err = this.real_client.AfterRun()
	if sign == STOPREQ {
		this.log.Error().Println(err);
		return;
	}
}

func (this *Client) runRequest (request RequestInterface) {
	for {
		switch request.GetSigner() {
		case STOPREQ:
			this.log.Error().Println("收到停止信号");
			runtime.Goexit()
			return;
		case SLEEP:
			this.log.Error().Println("收到睡眠信号");
			runtime.Gosched()
			continue;
		default:

		}

		sign,err := this.real_client.BeforeRunRequest(request)
		if sign == STOPREQ {
			this.log.Error().Println(err);
			runtime.Goexit()
			return;
		}
		sign,err = this.real_client.RunRequest(request)
		if sign == STOPREQ {
			runtime.Goexit()
			this.log.Error().Println(err);
			return;
		}
		sign,err = this.real_client.AfterRunRequest(request)
		if sign == STOPREQ {
			runtime.Goexit()
			this.log.Error().Println(err);
			return;
		}
	}

}

func (this *Client) SetStatus (status uint8) {
	this.status = status;
}

func (this *Client) GetStatus () uint8 {
	return this.status;
}


func New (req_client ReqClientInterface,requests map[uint]RequestInterface) *Client {
	return &Client{
		log:         nil,
		requests:    requests,
		real_client: req_client,
		status:      0,
		signers:     make(chan uint,len(requests)),
	}
}

/*
* 设置停止客户端的操作
*/
func (this *Client) SetRequestedId (id uint) {
	this.signers <- id;
}


/*
*设置停止
*/
type RequestInterface interface {
	SetSigner (sign Sign)
	GetSigner() Sign;
	Run ()
	GetMethod() string
	GetUrl() string
	GetBody() io.Reader
}

type ReqClientInterface interface {
	BeforeRun() (Sign,error)
	BeforeRunRequest(req RequestInterface) (Sign,error)
	RunRequest(req RequestInterface) (Sign,error)
	AfterRunRequest(req RequestInterface) (Sign,error)
	AfterRun() (Sign,error)
}

type LoggerInterface interface {
	Error() LoggerInterface
	Info() LoggerInterface
	Println(...interface{})
}