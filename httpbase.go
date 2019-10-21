package scrapy

import (
	"errors"
	"io/ioutil"
	"net/http"
)

type HttpBase struct {
	Base
	http.Client
}

func (this *HttpBase) RunPage (my_req RequestInterface) error {
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

func (this *HttpBase) TransData (res []byte) {

}
