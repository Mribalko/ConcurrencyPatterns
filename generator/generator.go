package generator

func MakeChannel[T any](in []T) <-chan T {
	out := make(chan T)

	go func() {
		for _, v := range in {
			out <- v
		}
		close(out)
	}()

	return out
}
