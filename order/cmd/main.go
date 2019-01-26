package main

import (
	"log"
	"net/http"

	"github.com/jonaskay/cola-inventory-functions/order"
)

func main() {
	http.HandleFunc("/order/", order.Create)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
