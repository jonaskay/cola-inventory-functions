package main

import (
	"log"
	"net/http"

	"github.com/jonaskay/cola-inventory-functions/fetch"
)

func main() {
	http.HandleFunc("/fetch/", fetch.Latest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
