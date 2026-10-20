package main

import (
	"fmt"

	"example.org/workshop/catalog"
)

func receiptTotal() int {
	return catalog.TeaPrice() * 3
}

func main() {
	fmt.Printf("Tea x3: %d\n", receiptTotal())
}
