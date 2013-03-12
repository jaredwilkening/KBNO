package conf_test

import (
	"fmt"
	. "github.com/MG-RAST/KBNO/kbno-server/conf"
	"testing"
)

func TestLoad(t *testing.T) {
	fmt.Printf("%v\n", Shock)
	fmt.Printf("%v\n", Openstack)
}
