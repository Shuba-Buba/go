//go:build !solution

package parallel

func Parallel(fns ...func() error) error {
	return nil
}
