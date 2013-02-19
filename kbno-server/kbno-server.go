package main

import (
    "fmt"
    "net/http"    
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "KBase Notebook Orchestration Service")
    })    
    http.HandleFunc("/notebook", notebookHander)
    http.ListenAndServe(":8080", nil)
}