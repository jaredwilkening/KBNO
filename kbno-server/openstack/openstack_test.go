package openstack_test

import (
	"fmt"
	. "github.com/MG-RAST/KBNO/kbno-server/openstack"
	"testing"
)

func TestBoot(t *testing.T) {
	res, err := Boot("jared", "test_boot_node")
	if err != nil {
		println(err.Error())
	} else {
		fmt.Printf("Booted: %#v\n", res)
	}
}

/*
func TestStatus(t *testing.T) {
	if id, err := Boot("test_status_node"); err != nil {
		println(err.Error())
	} else {
		println("Booted: " + id)
		if s, err := Status(id); err != nil {
			println(err.Error())
		} else {
			fmt.Printf("Status: %#v\n", s)
		}
	}

}

func TestStatusAll(t *testing.T) {
	if s, err := StatusAll(); err != nil {
		println(err.Error())
	} else {
		fmt.Printf("Status: %#v\n", s)
	}
}

func TestBootAndDelete(t *testing.T) {
	id, err := Boot("test_delete_node")
	if err != nil {
		println(err.Error())
	} else {
		println("Booted: " + id)
	}
	status, err := Delete(id)
	if err != nil {
		println(err.Error())
	} else if status {
		println("Deleted: " + id)
	} else {
		println("Failed delete: " + id)
	}
}

func TestBootScript(t *testing.T) {
	buf, _ := BootScript("test_user")
	println(buf.String())
}

*/
