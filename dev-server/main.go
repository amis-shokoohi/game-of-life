package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("dev server running on port 8080...")
	log.Fatalln(http.ListenAndServe(":8080", http.FileServer(http.Dir("http"))))
}
