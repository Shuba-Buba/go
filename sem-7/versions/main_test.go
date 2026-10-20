package main

import (
	"testing"

	"example.org/workshop/cart"
	"example.org/workshop/catalog"
	"example.org/workshop/report"
)

func TestConsumersUseCatalog(t *testing.T) {
	if got, want := cart.Total(), catalog.TeaPrice()*2; got != want {
		t.Errorf("cart.Total() = %d; want %d", got, want)
	}

	if got, want := report.Price(), catalog.TeaPrice(); got != want {
		t.Errorf("report.Price() = %d; want %d", got, want)
	}
}
