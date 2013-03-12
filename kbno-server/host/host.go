package host

import (
	"github.com/MG-RAST/KBNO/kbno-server/openstack"
	"math/rand"
	"time"
)

const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"

var (
	errorStates = map[string]int{"ERROR": 1, "UNKNOWN": 1, "RESCUE": 1, "SHUTOFF": 1, "SUSPENDED": 1}
)

type Host struct {
	ID         string
	InstanceID string
	User       string
	Server     openstack.Server
}

type Pool struct {
	Booting     chan Host
	Running     map[string]Host
	Pending     map[string]Host
	Error       map[string]Host
	del         chan Host
	updateProxy chan bool
}

func NewPool() (p *Pool) {
	p = &Pool{
		Booting:     make(chan Host, 1000),
		Running:     make(map[string]Host),
		Pending:     make(map[string]Host),
		Error:       make(map[string]Host),
		del:         make(chan Host, 1000),
		updateProxy: make(chan bool, 1000),
	}
	go p.update()
	return
}

func (p *Pool) update() {
	for {
		select {
		case h := <-p.Booting:
			// need to add print error to log
			s, _ := openstack.Status(h.Server.ID)
			h.Server = s.Server
			if _, in := errorStates[h.Server.Status]; in {
				p.Error[h.ID] = h
			} else if h.Server.Status != "ACTIVE" {
				go func(p *Pool, h Host) { time.Sleep(30 * time.Second); p.Booting <- h }(p, h)
			} else {
				delete(p.Pending, h.ID)
				p.Running[h.ID] = h
				p.updateProxy <- true
			}
		case h := <-p.del:
			delete(p.Running, h.ID)
			p.updateProxy <- true
		case <-p.updateProxy:

		}
	}
	return
}

func (p *Pool) Has(id string) bool {
	_, running := p.Running[id]
	_, pending := p.Pending[id]
	return running || pending
}

func (p *Pool) Delete(id string) (success bool, err error) {
	if _, running := p.Running[id]; running {
		p.del <- p.Running[id]
		return openstack.Delete(p.Running[id].Server.ID)
	}
	defer delete(p.Pending, id)
	return openstack.Delete(p.Pending[id].Server.ID)
}

func (p *Pool) Boot(username string) (id string, err error) {
	id = RandString(10)
	s, err := openstack.Boot(username, "ipy_kbnb_"+id)
	if err != nil {
		return "", err
	}
	h := Host{ID: id, User: username, Server: s.Server}
	p.Pending[id] = h
	p.Booting <- h
	go func(u, id string) {
		for {
			select {
			case <-time.After(10 * time.Minute):
				openstack.Delete(p.Pending[id].Server.ID)
				delete(p.Pending, id)
				id = RandString(10)
				s, _ := openstack.Boot(username, "ipy_kbnb_"+id)
				// log error
				h := Host{ID: id, User: username, Server: s.Server}
				p.Pending[id] = h
				p.Booting <- h
			}
		}
	}(username, id)
	return
}

func RandString(l int) (s string) {
	rand.Seed(time.Now().UTC().UnixNano())
	c := make([]byte, l)
	for i := 0; i < l; i++ {
		c[i] = chars[rand.Intn(len(chars))]
	}
	return string(c)
}
