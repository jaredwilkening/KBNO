package host_test

import (
	"fmt"
	. "github.com/MG-RAST/KBNO/kbno-server/host"
	"strconv"
	"testing"
	"time"
)

func TestBoot(t *testing.T) {
	p := NewPool()
	id := p.Boot("jared")
	println(id)
	fmt.Printf("%#v\n", p)
	time.Sleep(90 * time.Second)
	p.Free <- Host{InstanceID: RandString(20), Status: "Active"}
	time.Sleep(2 * time.Second)
	fmt.Printf("%#v\n", p)
	time.Sleep(90 * time.Second)
}

func TestRandString(t *testing.T) {
	for i := 1; i < 10; i++ {
		println(RandString(20))
	}
}

func TestPool(t *testing.T) {
	p := NewPool()
	p.Free <- Host{ID: "free_node", Status: "Active"}
	println("here")
	println(strconv.Itoa(len(p.Free)))
	for i := 0; i <= 10; i++ {
		p.Booting <- Host{ID: strconv.Itoa(i), Status: "Booting"}
		time.Sleep(1 * time.Second)
	}
	//time.Sleep(30 * time.Second)
}
