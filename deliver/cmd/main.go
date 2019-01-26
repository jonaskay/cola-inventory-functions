package main

import (
	"log"
	"net/http"

	"github.com/jonaskay/cola-inventory-functions/deliver"
)

func main() {
	http.HandleFunc("/deliver/", deliver.Latest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
