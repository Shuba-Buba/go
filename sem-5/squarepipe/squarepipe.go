//go:build !solution

package squarepipe

func Square(in <-chan int) <-chan int {
	out := make(chan int)
	close(out)
	return out
}
