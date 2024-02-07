package error_group

import (
	"errors"
	"sync"
)

var (
	ErrTasksLimitUnchanable = errors.New("modifing limit while goroutines are active")
	ErrErrorLimitExceeded   = errors.New("errors limit exceeded")
)

type (
	token struct{}
	Task  func() error
	Group struct {
		wg       sync.WaitGroup
		limiter  chan token
		errLimit int
		mu       sync.Mutex
		errors   []error
	}
)

// Creates new instance of group and sets the number of errors limit.
// A negative or zero value sets the limit to 1
func New(errLimit int) *Group {
	return &Group{
		errLimit: max(errLimit, 1),
	}
}

// Sets the number of tasks running simultaneously. A negative or zero value indicates no limit.
// The limit must not be modified while any tasks are executed.
func (g *Group) SetLimit(n int) error {

	if len(g.limiter) != 0 {
		return ErrTasksLimitUnchanable
	}

	if n <= 0 {
		g.limiter = nil
		return nil
	}

	g.limiter = make(chan token, n)
	return nil
}

func (g *Group) done() {
	if g.limiter != nil {
		<-g.limiter
	}
	g.wg.Done()
}

func (g *Group) exceeded() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	return len(g.errors) >= g.errLimit
}

// Runs function in a new goroutine. Blocks until the new goroutine can be added
// (not exceeding the configured limit).
// Return ErrErrorLimitExceeded if the number of occured errros exceeds the limit.
func (g *Group) Go(f Task) error {

	if g.exceeded() {
		return ErrErrorLimitExceeded
	}

	if g.limiter != nil {
		g.limiter <- token{}
	}

	g.wg.Add(1)

	go func() {
		defer g.done()

		if err := f(); err != nil {
			g.mu.Lock()
			g.errors = append(g.errors, err)
			g.mu.Unlock()
		}
	}()

	return nil
}

// Blocks until all tasks return from the Go method, returns all occured errors.
func (g *Group) Wait() []error {
	g.wg.Wait()
	if g.limiter != nil {
		close(g.limiter)
	}
	return g.errors
}
