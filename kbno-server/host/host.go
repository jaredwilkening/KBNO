package host

import (
	//"openstack"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"

type Host struct {
	ID         string
	InstanceID string
	Status     string
	User       string
	IPAddress  []string
}

type Pool struct {
	Free    chan Host
	flock   sync.Mutex
	Booting chan Host
	Running map[string]Host
	Pending map[string]bool
	Error   map[string]Host
}

func NewPool() (p *Pool) {
	p = &Pool{
		Free:    make(chan Host, 1000),
		Booting: make(chan Host, 1000),
		Running: make(map[string]Host),
		Pending: make(map[string]bool),
		Error:   make(map[string]Host),
	}
	go p.updateStatus()
	return
}

func (p *Pool) updateStatus() {
	for {
		select {
		case h := <-p.Booting:
			// openstack.Update(h)
			fmt.Printf("%s: %#v\n", time.Now(), h)
			if h.Status == "Error" {
				p.Error[h.ID] = h
			} else if h.Status != "Active" {
				go func(p *Pool, h Host) { time.Sleep(30 * time.Second); p.Booting <- h }(p, h)
			}
		}
	}
}

func (p *Pool) Boot(username string) (id string) {
	id = RandString(10)
	p.flock.Lock()
	if len(p.Free) > 0 {
		h := <-p.Free
		h.ID = id
		h.User = username
		p.Running[id] = h
		p.flock.Unlock()
		return
	} else {
		p.flock.Unlock()
		p.Pending[id] = true
		h := Host{}
		//openstack.Boot(h)
		p.Booting <- h
	}
	go func(s string) {
		for {
			select {
			case h := <-p.Free:
				if _, ok := p.Pending[h.ID]; ok {
					delete(p.Pending, id)
					h.ID = id
					h.User = username
					p.Running[id] = h
				} else {
					p.Free <- h
				}
				return
			case <-time.After(1 * time.Minute):
				h := Host{}
				//openstack.Boot(h)
				p.Booting <- h
			}
		}
	}(id)
	return
}

func RandString(l int) (s string) {
	c := make([]byte, l)
	for i := 0; i < l; i++ {
		c[i] = chars[rand.Intn(len(chars))]
	}
	return string(c)
}
