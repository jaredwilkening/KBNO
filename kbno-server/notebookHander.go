package main

import (
    "fmt"
    "net/http"    
)

func notebookHander(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "POST":
        println("create new")
    case "GET":
        println("return")
    case "DELETE":
        println("delete")
    }
}