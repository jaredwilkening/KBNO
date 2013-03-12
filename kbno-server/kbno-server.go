package main

import (
	"github.com/MG-RAST/KBNO/kbno-server/host"
	"github.com/MG-RAST/KBNO/kbno-server/response"
	"net/http"
	"strings"
)

var (
	pool = host.NewPool()
)

type hostRes struct {
	ID     string `json:"id"`
	Url    string `json:"url"`
	Status string `json:"status"`
}

func hostHander(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		id, err := pool.Boot("Jared")
		if err != nil {
			println(err.Error())
			response.InternalServerError(w)
			return
		}
		if h, ok := pool.Running[id]; ok {
			response.WithData(w, hostRes{ID: id, Url: h.ID, Status: "Active"})
			return
		} else {
			response.WithData(w, hostRes{ID: id, Url: "", Status: "Pending"})
			return
		}
	case "GET":
		path := strings.Split(r.URL.Path[1:], "/")
		if len(path) > 1 {
			id := path[1]
			if h, ok := pool.Running[id]; ok {
				response.WithData(w, hostRes{ID: id, Url: h.ID, Status: "Active"})
				return
			} else if _, ok := pool.Pending[id]; ok {
				response.WithData(w, hostRes{ID: id, Url: "", Status: "Pending"})
				return
			} else {
				response.BadRequestWithMessage(w, "Unable to locate or missing host id")
				return
			}
		}
	case "DELETE":
		path := strings.Split(r.URL.Path[1:], "/")
		if len(path) > 1 {
			id := path[1]
			if pool.Has(id) {
				if success, err := pool.Delete(id); success {
					response.OK(w)
					return
				} else {
					println(err.Error())
					response.InternalServerError(w)
					return
				}
			} else {
				response.BadRequestWithMessage(w, "Unable to locate or missing host id")
				return
			}
		}
	default:
		response.BadRequest(w)
		return
	}
	return
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { response.NotFound(w) })
	http.HandleFunc("/host/", hostHander)
	http.ListenAndServe(":8888", nil)
}
