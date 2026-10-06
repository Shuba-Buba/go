//go:build !solution

package first

func First(chs ...<-chan int) (int, bool) {
	if len(chs) == 0 {
		return 0, false
	}

	type answer struct {
		v  int
		ok bool
	}

	res := make(chan answer)
	done := make(chan struct{})

	for _, c := range chs {
		go func(ch <-chan int) {
			v, ok := <-ch
			select {
			case res <- answer{v, ok}:
			case <-done:
			}
		}(c)
	}

	for i := 0; i < len(chs); i++ {
		a := <-res
		if a.ok {
			close(done)
			return a.v, true
		}
	}

	close(done)
	return 0, false
}
