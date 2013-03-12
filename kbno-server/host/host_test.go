package host_test

import (
	"fmt"
	. "github.com/MG-RAST/KBNO/kbno-server/host"
	"testing"
	"time"
)

func TestBoot(t *testing.T) {
	p := NewPool()
	id, err := p.Boot("jared")
	if err != nil {
		println(err.Error())
	}
	println(id)
	for {
		fmt.Printf("%#v\n", p)
		time.Sleep(35 * time.Second)
	}
}

func TestRandString(t *testing.T) {
	for i := 1; i < 10; i++ {
		println(RandString(20))
	}
}
