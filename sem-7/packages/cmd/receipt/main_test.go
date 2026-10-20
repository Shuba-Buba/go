package main

import "testing"

func TestReceiptTotal(t *testing.T) {
	for _, test := range []struct {
		name     string
		price    int
		quantity int
		want     int
	}{
		{name: "three", price: 75, quantity: 3, want: 225},
		{name: "free", price: 0, quantity: 3, want: 0},
		{name: "empty", price: 75, quantity: 0, want: 0},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := receiptTotal(test.price, test.quantity); got != test.want {
				t.Fatalf(
					"receiptTotal(%d, %d) = %d; want %d",
					test.price, test.quantity, got, test.want,
				)
			}
		})
	}
}
