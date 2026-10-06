//go:build !solution

package first

import "reflect"

func First(chs ...<-chan int) (int, bool) {
	cases := make([]reflect.SelectCase, len(chs))
	for i, ch := range chs {
		cases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		}
	}

	for open := len(cases); open > 0; {
		i, v, ok := reflect.Select(cases)
		if ok {
			return int(v.Int()), true
		}

		cases[i].Chan = reflect.Value{}
		open--
	}

	return 0, false
}
