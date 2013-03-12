package nginx_test

import (
	//"fmt"
	. "github.com/MG-RAST/KBNO/kbno-server/nginx"
	"testing"
)

func TestLoad(t *testing.T) {
	Generate([]Location{Location{ID: "test1", IP: "10.10.10.20"}, Location{ID: "test2", IP: "10.10.10.12"}}, "/tmp/nginx.conf")
}
