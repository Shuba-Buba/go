//go:build !solution

package merge

func Merge(chs ...<-chan int) <-chan int {
	out := make(chan int)
	close(out)
	return out
}
