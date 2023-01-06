package main

import (
	"fmt"
	"net/http"
)

func Hello(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "Hello world")
}