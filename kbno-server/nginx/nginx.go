package nginx

import (
	"bytes"
	"io/ioutil"
	"text/template"
)

type Location struct {
	ID string
	IP string
}

const headerTemplate = `server {
    server_name kbno.kbase.us;    
    access_log  /var/log/nginx/kbno/access.log main;
    error_log   /var/log/nginx/kbno/error.log;
    
`
const locationTemplate = `    location /{{.ID}} {
      proxy_pass  http://{{.IP}}/;
      proxy_set_header        Host            $host;
      proxy_set_header        X-Real-IP       $remote_addr;
      proxy_set_header        X-Forwarded-For $proxy_add_x_forwarded_for;
      proxy_buffering         on;
      proxy_buffer_size       32k;
      proxy_buffers           512 32k;
      client_body_buffer_size 51200k;
    }
    
`

func Generate(loc []Location, path string) (err error) {
	config := bytes.NewBuffer(nil)
	config.WriteString(headerTemplate)
	t := template.Must(template.New("nginx").Parse(locationTemplate))
	for _, l := range loc {
		if err = t.Execute(config, l); err != nil {
			return err
		}
	}
	config.WriteString("}\n")
	err = ioutil.WriteFile(path, config.Bytes(), 0644)
	return
}
