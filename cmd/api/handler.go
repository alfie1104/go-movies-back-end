package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request){
	var payload = struct {
		Status string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status : "active",
		Message : "Go Movies up and running",
		Version: "1.0.0",
	}

	out, err := json.Marshal(payload)


	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (app *application) AllMovies(w http.ResponseWriter, r *http.Request){
	movies, err := app.DB.AllMovies()
	if err != nil {
		fmt.Println(err)
		return
	}

	out, err := json.Marshal(movies)

	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}