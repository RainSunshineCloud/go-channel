package scrapy

import (
	"io"
	"runtime"
	"time"
)

type Base struct {
	requests map[uint16]RequestInterface
	stopedRequests map[uint16]RequestInterface
	stopedPage chan uint16
	toStopPage chan uint16
	stopClient chan uint16
}

type RequestInterface interface {
	GetUrl() string
	GetMethod() string
	GetBody() io.Reader
	GetTimes() uint64
	GetId() uint16
	GetStopTime() time.Time
}

type ClientInterface interface {
	Run(reqs map[uint16] RequestInterface) RequestInterface
	Stop(request_id uint16)
	Start(request_id uint16)
	ReStart(request_id uint16)
}


func (this *Base) Run ( reqs map[uint16] RequestInterface) RequestInterface {
	//请求前
	this.RunBefore()
	for req_id,_ := range reqs {
		go this.Start(req_id);
	}
	//请求后
	this.RunAfter()
	return this.GetStopedPage()
}

func (this *Base) GetStopedPage () RequestInterface {
	id := <- this.stopedPage
	if val,ok := this.requests[id];ok {
		return val;
	}
	return nil;
}

func (this *Base) AddStopedPage (id uint16,requests RequestInterface) {
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

/*
	启动
*/
func (this *Base) Start(request_id uint16) {
	//开始
	if val,ok := this.requests[request_id];ok {
		defer func(request_id uint16,req RequestInterface) {
			this.AddStopedPage(request_id,req)
		}(request_id,val)

		if val.GetTimes() == 0 {
			for {
				this.CheckIfStop(val)
				this.RunPageBefore(val)
				this.RunPage(val)
				this.RunPageAfter(val)
			}
		} else {
			var i uint64 = 0;
			for i < val.GetTimes() {
				this.CheckIfStop(val)
				this.RunPageBefore(val)
				this.RunPage(val)
				this.RunPageAfter(val)
				i = i + 1;
			}
		}


	}

}


/*
	重新启动
*/
func (this *Base) Restart (request_id uint16) {
	this.Stop(request_id);
	this.Start(request_id)
}



func (this *Base) RunBefore () {

}


func (this *Base) RunAfter () {

}


func (this *Base) RunPageBefore(request RequestInterface) {

}

func (this *Base) RunPage(request RequestInterface) {

}

func (this *Base) RunPageAfter(request RequestInterface) {

}

func (this *Base) CheckIfStop(request RequestInterface) {
	select {
	case  id := <- this.toStopPage:
		if id == request.GetId() {
			runtime.Goexit()
		}
	default:
		if  time.Now().Sub(request.GetStopTime()) <= 0 {
			runtime.Gosched()
			this.CheckIfStop(request)
		}
	}
}