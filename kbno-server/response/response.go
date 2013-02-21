package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func WithData(w http.ResponseWriter, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	output, _ := json.Marshal(obj)
	w.Header().Add("Content-Length", strconv.Itoa(len(output)))
	w.Write(output)
	return
}

func OK(w http.ResponseWriter) {
	w.WriteHeader(200)
	fmt.Fprintf(w, "OK")
	return
}

func BadRequest(w http.ResponseWriter) {
	w.WriteHeader(400)
	fmt.Fprintf(w, "Bad Request")
	return
}

func BadRequestWithMessage(w http.ResponseWriter, m string) {
	w.WriteHeader(400)
	fmt.Fprintf(w, "Bad Request: "+m)
	return
}

func Unauthorized(w http.ResponseWriter) {
	w.WriteHeader(401)
	fmt.Fprintf(w, "Unauthorized")
	return
}

func Forbidden(w http.ResponseWriter) {
	w.WriteHeader(403)
	fmt.Fprintf(w, "Forbidden")
	return
}

func NotFound(w http.ResponseWriter) {
	w.WriteHeader(404)
	fmt.Fprintf(w, "Not Found")
	return
}

func InternalServerError(w http.ResponseWriter) {
	w.WriteHeader(500)
	fmt.Fprintf(w, "Internal Server Error")
	return
}

func NotImplemented(w http.ResponseWriter) {
	w.WriteHeader(501)
	fmt.Fprintf(w, "Not Implemented")
	return
}
