package main

import "fmt"

func receiptTotal(unitPrice, quantity int) int {
	return unitPrice * quantity
}

func main() {
	fmt.Printf("Tea x3: %d\n", receiptTotal(75, 3))
}
