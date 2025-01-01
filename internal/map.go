package internal

func MapValues[T any, R any](values []T, fn func(v T) R) []R {
	var out []R
	for _, value := range values {
		out = append(out, fn(value))
	}

	return out
}
