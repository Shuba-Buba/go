package catalog

import "testing"

func TestTeaPrice(t *testing.T) {
	if TeaPrice() <= 0 {
		t.Fatal("tea must have a positive price")
	}
}
