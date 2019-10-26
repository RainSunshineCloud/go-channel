package client_manager

const RUNNING = 1;
const STOPED = 2;
const STOPALL = 0;
type Manager struct{
	clients map[uint] ClientInterface
}

func (this *Manager) Start () {
	for _,client := range this.clients {
		client.Run();
		client.SetStatus(RUNNING);
	}
}

func (this *Manager) Stop (id uint) {
	if id == STOPALL {
		for _,client := range this.clients {
			client.Stop();
			client.SetStatus(STOPED);
		}
	}

	if client,ok := this.clients[id];ok {
		client.Stop();
		client.SetStatus(STOPED);
	}
}

func (this *Manager) Register (client ClientInterface) {
	this.clients[client.GetId()] = client;
}


func New () *Manager {
	return & Manager{clients:make(map[uint]ClientInterface),}
}


type ClientInterface interface {
	Stop();
	Start();
	SetStatus(uint);
	GetId() uint;
	Run();
}

