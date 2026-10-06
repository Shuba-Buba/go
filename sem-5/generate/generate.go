//go:build !solution

package generate

func Generate(n int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 0; i < n; i++ {
			ch <- i
		}
	}()
	return ch
}
