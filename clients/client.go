package clients

import (
	"errors"
	"io/ioutil"
	"net/http"
)

type HttpClient struct {
	http.Client
}

func (this *HttpClient) BeforeRun () (Sign,error) {
	return STARTREQ,nil;
}
func (this *HttpClient) BeforeRunRequest(req RequestInterface) (Sign,error) {
	return STARTREQ,nil;
}
func (this *HttpClient) RunRequest(my_req RequestInterface) (Sign,error) {
	defer func () {
		if r := recover();r != nil {

		}
	}()

	req,err := http.NewRequest(my_req.GetMethod(),my_req.GetUrl(),my_req.GetBody())
	if err != nil {
		return STARTREQ,err;
	}

	res,err := this.Do(req)
	if err != nil || res == nil {
		return STARTREQ,err;
	}

	if res == nil && res.Body == nil && res.StatusCode != 200 {
		return STARTREQ,errors.New(res.Status)
	}

	data,err := ioutil.ReadAll(res.Body)

	if err != nil {
		return STARTREQ,err;
	}

	this.TransData(data);
	return STARTREQ,nil;
}


func  (this *HttpClient)  AfterRunRequest(req RequestInterface) (Sign,error) {
	return STARTREQ,nil;
}
func  (this *HttpClient)  AfterRun() (Sign,error) {
	return STARTREQ,nil;
}

func (this *HttpClient) TransData (res []byte) {

}



func NewHttpClient() *HttpClient {
	return &HttpClient{}
}