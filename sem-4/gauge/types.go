package gauge

import "errors"

var ErrNegative = errors.New("negative level")

type Gauge struct{ level int }
