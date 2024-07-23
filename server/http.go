package main

import (
	"log"
	"net/http"
)

func startHTTP() {
	http.Handle("/", http.FileServer(http.Dir("../app/dist")))
	log.Println("http server listening from", addressHTTP)
	log.Panic(http.ListenAndServe(addressHTTP, nil))
}
