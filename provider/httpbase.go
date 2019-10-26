package provider

import (
	"errors"
	"github.com/RainSunshineCloud/scrapy/requests"
	"io/ioutil"
	"net/http"
)

type HttpBase struct {
	http.Client
}

func (this *HttpBase) RunPage (my_req requests.RequestInterface) error {
	req,err := http.NewRequest(my_req.GetMethod(),my_req.GetUrl(),my_req.GetBody())
	if err != nil {
		return err;
	}

	res,err := this.Do(req)
	if err != nil || res == nil {
		return err;
	}

	if res == nil && res.Body == nil && res.StatusCode != 200 {
		return errors.New(res.Status)
	}

	data,err := ioutil.ReadAll(res.Body)

	if err != nil {
		return err;
	}

	this.TransData(data);
	return nil;
}

func (this *HttpBase) RunPageAfter (my_req requests.RequestInterface) error {
	return nil;
}

func (this *HttpBase) RunPageBefore (my_req requests.RequestInterface) error {
	return nil;
}


func (this *HttpBase) RunBefore () error {
	return nil;
}

func (this *HttpBase) RunAfter () error {
	return nil;
}

func NewHttpBase (reqs map[uint16] requests.RequestInterface) *Base {
	client := &HttpBase{
		Client: http.Client{},
	}

	return NewBase(reqs,client)
}


func (this *HttpBase) TransData (res []byte) {

}

