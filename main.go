package main

import (
	"net/http"
	"runtime"
)

type Request_Manager struct {
	requests map[uint] *Request
	clients map[uint] *Client
	response map[uint] *Response
	request_to_client map[uint]chan *Request
	client_to_response map[uint]chan interface{}
	response_to_request map[uint]chan interface{}

}

type Request struct {
	req *http.Request
	client_id uint   // 客户端id号
	response_id uint //响应id号
	request_id uint  //请求id号
}

type Client struct {
	Id uint
	http.Client
}

type Response struct {
	Id uint
}


func New () *Request_Manager {
	return &Request_Manager{
		requests:            nil,
		clients:             nil,
		request_to_client:   nil,
		client_to_response:  nil,
		response_to_request: nil,
	}
}

func (this *Request_Manager) RegisterRequest (req *Request) *Request_Manager {
	this.requests[req.request_id] = req;
	//请求发送到客户端
	this.request_to_client[req.request_id] = make(chan *Request);
	//客户端发送响应
	if _,ok := this.client_to_response[req.response_id];!ok {
		this.client_to_response[req.response_id] = make(chan interface{});
	}
	//响应发消息给请求
	if _,ok := this.response_to_request[req.request_id];!ok {
		this.response_to_request[req.request_id] = make(chan interface{});
	}
	return this;
}


func (this *Request_Manager) RegisterClient(my_client *Client) *Request_Manager {
	this.clients[my_client.Id] = my_client;
	return this;
}

func (this *Request_Manager) RegisterResponse (response  *Response) *Request_Manager {
	this.response[response.Id] = response;
	return this;
}

func (this *Request_Manager) RunClient (req *Request) {
	for {
		0 := <-this.request_to_client[req.request_id]
		res,_ := this.clients[request.client_id].Do(request.req);
		this.client_to_response[request.response_id] <- res;
	}
}


func (this *Request_Manager) RunResponse (req_to_res_id uint,res_id uint) {
	for {
		response := <-this.client_to_response[res_id]
		if response == nil {
			 this.response_to_request[req_to_res_id] <- response
		}
	}
}

func (this *Request_Manager) Run () {
	for _,req := range this.requests{
		go this.RunClient(req);
		go this.RunResponse(req);
		go this.SendReq(req);
	}
}


func (this *Request_Manager) SendReq(req_id uint) {
	for {
		select {
		case <-res_to_req:
			this.request_to_client[req_id] <- this.GetReq(req_id);
		default:
			runtime.Gosched()
		}
	}
}


func (this *Request_Manager) GetReq(req_id uint) *Request {
	return nil;
}