package fanin

import (
	"context"
	"sync"
)

func fanIn[T any](ctx context.Context, input ...<-chan T) chan T {
	res := make(chan T)
	var wg sync.WaitGroup

	for _, ch := range input {
		ch := ch
		wg.Add(1)

		go func() {
			defer wg.Done()
			for v := range ch {
				select {
				case <-ctx.Done():
					return
				default:
					res <- v
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(res)
	}()

	return res
}
