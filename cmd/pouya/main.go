package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./"))

	log.Println("Server running...")
	log.Fatal(http.ListenAndServe(":8585", fs))
}
