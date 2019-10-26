package provider

import (
	"github.com/RainSunshineCloud/scrapy/requests"
	"runtime"
	"time"
)

type Base struct {
	Client ClientInterface
	Logger LoggerInterface
	requests map[uint16]requests.RequestInterface
	stopedRequests map[uint16]requests.RequestInterface
	stopedPage chan uint16 //停止的页面
	toStopPage chan uint16 //要停止的页面
	stopClient chan uint16 //要停止的客户端
}

type ClientInterface interface {
	RunPageBefore(requests.RequestInterface) error
	RunPage(requests.RequestInterface) error
	RunPageAfter(requests.RequestInterface) error
	RunBefore() error
	RunAfter() error
}

type LoggerInterface interface {
	Error(interface{})
	Info(interface{})

}


func (this *Base) Run () uint16 {
	//请求前
	err := this.Client.RunBefore()
	if err != nil {
		this.Logger.Error(err);
		return 0;
	}

	for req_id,_ := range this.requests {
		go this.Start(req_id);
	}

	//请求后
	err = this.Client.RunAfter()
	if err != nil {
		this.Logger.Error(err);
		return 0;
	}

	return this.GetStopedPage()
}

func (this *Base) ReRun () {

}

func (this *Base) GetStopedPage () uint16 {
	id := <- this.stopedPage
	if _,ok := this.requests[id];ok {
		return id;
	}
	return 0;
}

func (this *Base) AddStopedPage (id uint16,requests requests.RequestInterface) {
	this.stopedPage <- id;
	if _,ok := this.requests[id];ok {
		if _,ok = this.stopedRequests[id];!ok {
			this.stopedRequests[id] = requests;
		}
	}
}

func (this *Base) Stop(request_id uint16) {
	if request_id == 0 {
		for req_id,_ := range this.requests {
			if _,ok := this.stopedRequests[req_id];!ok {
				this.toStopPage <- req_id;
			}
		}
	} else {
		if _,ok :=  this.requests[request_id]; ok {
			if _,ok := this.stopedRequests[request_id];!ok {
				this.toStopPage <- request_id
			}
		}
	}


	runtime.GC()
}

func (this *Base) Start(request_id uint16) {
	//开始
	if val,ok := this.requests[request_id];ok {

		defer func(request_id uint16,req requests.RequestInterface) {
			if r := recover(); r != nil {
				this.AddStopedPage(request_id,req)
			} else {
				this.AddStopedPage(0,req)
			}

		}(request_id,val)

		if val.GetTimes() == 0 {
			for {
				if this.CheckIfStop(val) {
					continue;
				}
				err := this.Client.RunPageAfter(val)
				if err != nil {
					this.Logger.Error(err);
				}
				err = this.Client.RunPage(val)
				if err != nil {
					this.Logger.Error(err);
				}
				err = this.Client.RunPageAfter(val)
				if err != nil {
					this.Logger.Error(err);
				}
			}
		} else {
			var i uint64 = 0;

			for i < val.GetTimes() {
				if this.CheckIfStop(val) {
					continue;
				}

				err := this.Client.RunPageAfter(val)
				if err != nil {
					this.Logger.Error(err);
				}
				err = this.Client.RunPage(val)
				if err != nil {
					this.Logger.Error(err);
				}
				err = this.Client.RunPageAfter(val)
				if err != nil {
					this.Logger.Error(err);
				}
				i = i + 1;
			}

		}
	}

}

func (this *Base) Restart (request_id uint16) {
	this.Stop(request_id);
	this.Start(request_id)
}




func (this *Base) CheckIfStop(request requests.RequestInterface) bool {
	select {
	case  id := <- this.toStopPage:
		if id == request.GetId() {
			runtime.Goexit()
			return true;
		}
	default:
		if  request.GetStopTime().Sub(time.Now()) >= 0 {
			runtime.Gosched()
			return true;
		}
		return false;
	}
	return false;
}

func NewBase(reqs map[uint16] requests.RequestInterface,client ClientInterface) *Base {
	return &Base{
		requests:       reqs,
		Client:			client,
		stopedRequests: make(map[uint16] requests.RequestInterface,len(reqs)),
		stopedPage:     make(chan uint16,1),
		toStopPage:     make(chan uint16,1),
		stopClient:     make(chan uint16,1),
	}
}