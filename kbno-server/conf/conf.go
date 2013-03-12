package conf

import (
	"flag"
	"fmt"
	"github.com/jaredwilkening/goconfig/config"
	"os"
)

type value struct {
	Str  string
	Int  int
	Bool bool
	Type string
	Req  bool
}

// Setup conf variables
var (
	ConfigFile = ""
	Shock      = map[string]value{
		"url": value{Type: "string", Req: true},
	}
	Openstack = map[string]value{
		"username":    value{Type: "string", Req: true},
		"password":    value{Type: "string", Req: true},
		"auth_url":    value{Type: "string", Req: true},
		"nova_api":    value{Type: "string", Req: true},
		"tenant_name": value{Type: "string", Req: true},
		"tenant_id":   value{Type: "string", Req: true},
		"vm_key_name": value{Type: "string", Req: true},
		"image":       value{Type: "string", Req: true},
		"flavor":      value{Type: "string", Req: true},
		"boot_script": value{Type: "string", Req: false},
	}
	Nginx = map[string]value{
		"conf_dir": value{Type: "string", Req: true},
	}
)

func init() {
	flag.StringVar(&ConfigFile, "conf", "", "path to config file")
	flag.Parse()

	c, err := config.ReadDefault(ConfigFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: error reading conf file: %v\n", err)
		os.Exit(1)
	}

	for k, v := range map[string]*map[string]value{"shock": &Shock, "openstack": &Openstack, "nginx": &Nginx} {
		err := load(c, k, v)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: error reading conf file: %v\n", err)
			os.Exit(1)
		}
	}

}

func load(c *config.Config, section string, s *map[string]value) (err error) {
	for k, v := range *s {
		switch v.Type {
		case "string":
			if t, e := c.String(section, k); v.Req && e != nil {
				return e
			} else {
				v.Str = t
				(*s)[k] = v
			}
		case "int":
			if t, e := c.Int(section, k); v.Req && e != nil {
				return e
			} else {
				v.Int = t
				(*s)[k] = v
			}
		case "bool":
			if t, e := c.Bool(section, k); v.Req && e != nil {
				return e
			} else {
				v.Bool = t
				(*s)[k] = v
			}
		}
	}
	return
}
