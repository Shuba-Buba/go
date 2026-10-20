package main

import (
	"fmt"

	"example.org/workshop/cart"
	"example.org/workshop/catalog"
	"example.org/workshop/report"
)

func main() {
	fmt.Printf(
		"Catalog %s; cart: %d; report: %d\n",
		catalog.Version(), cart.Total(), report.Price(),
	)
}
