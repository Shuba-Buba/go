package tablecases

import "errors"

// ErrWeight reports a weight outside the supported range of 1 through 20 kg.
var ErrWeight = errors.New("weight must be between 1 and 20 kg")

// ShippingCost returns a delivery cost in roubles for an integer weight in kg.
func ShippingCost(weight int, express bool) (int, error) {
	if weight < 1 || weight > 20 {
		return 0, ErrWeight
	}

	var cost int

	switch {
	case weight <= 5:
		cost = 100
	case weight <= 10:
		cost = 180
	default:
		cost = 300
	}

	if express {
		cost += 150
	}

	return cost, nil
}
