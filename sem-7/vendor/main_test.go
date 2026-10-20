package main

import (
	"testing"

	"example.org/workshop/catalog"
)

func TestReceiptTotal(t *testing.T) {
	if got, want := receiptTotal(), catalog.TeaPrice()*3; got != want {
		t.Fatalf("receiptTotal() = %d; want %d", got, want)
	}
}
