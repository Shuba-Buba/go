//go:build !solution

package fallback

func Lookup(primary, backup Source, key string) (string, error) { return "", nil }
