package requests

import (
	"io"
	"time"
)

type RequestInterface interface {
	GetUrl() string
	GetMethod() string
	GetBody() io.Reader
	GetTimes() uint64
	GetId() uint16
	GetStopTime() time.Time
}
type Base struct {
	url string
	method string
	body io.ReadSeeker
	max_req uint64
	id uint16
	stop_time time.Time
}

func (this *Base) GetUrl() string {
	return this.url
}

func (this *Base) GetMethod() string {
	return this.method
}

func (this *Base) GetBody() io.Reader {
	return this.body
}

func (this *Base) GetTimes() uint64 {
	return this.max_req
}

func (this *Base) GetId() uint16 {
	return this.id
}

func (this *Base) GetStopTime() time.Time {
	return this.stop_time
}

func New(url string , method string ,max_req uint64, id uint16) RequestInterface {
	return &Base{url:url,method:method,max_req:max_req,id:id,stop_time:time.Now().Add(1 * time.Second)}
}