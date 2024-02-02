package parallel_execution

import (
	"errors"
	"fmt"

	errGroup "github.com/MRibalko/ConcurrencyPatterns/errorGroup"
)

var ErrLimitExceeded = errors.New("limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	g := errGroup.New(m)

	if err := g.SetLimit(n); err != nil {
		return fmt.Errorf("set tasks limit %d: %v", n, err)
	}

	for _, t := range tasks {
		if err := g.Go(errGroup.Task(t)); errors.Is(err, errGroup.ErrErrorLimitExceeded) {
			break
		}
	}

	errors := g.Wait()

	if len(errors) >= m {
		return ErrLimitExceeded
	}

	return nil
}
