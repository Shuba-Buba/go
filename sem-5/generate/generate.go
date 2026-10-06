//go:build !solution

package generate

func Generate(n int) <-chan int {
	ch := make(chan int)
	close(ch)
	return ch
}
