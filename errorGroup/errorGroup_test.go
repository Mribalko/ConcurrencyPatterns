package error_group

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
)

func TestGo(t *testing.T) {
	t.Run("errors limit", func(t *testing.T) {
		t.Parallel()
		errTest := errors.New("Err")

		cases := []struct {
			maxTasksLimit  int
			maxErrorsLimit int
			tasksRes       []error
		}{
			{
				0,
				0,
				[]error{nil, errTest, nil},
			},
			{
				2,
				0,
				[]error{nil, errTest, nil},
			},
			{
				2,
				2,
				[]error{nil, errTest, nil},
			},
			{
				2,
				1,
				[]error{errTest, errTest, nil, nil, nil},
			},
			{
				1,
				1,
				[]error{errTest, errTest, errTest},
			},
		}

		for _, tt := range cases {
			tt := tt
			t.Run(fmt.Sprintf("tasksLimit=%d,errLimit=%d,errors=%v", tt.maxTasksLimit, tt.maxErrorsLimit, tt.tasksRes),
				func(t *testing.T) {
					t.Parallel()
					g := New(tt.maxErrorsLimit)
					g.SetLimit(tt.maxTasksLimit)

					var errIteration int
					for i, v := range tt.tasksRes {
						v := v
						err := g.Go(func() error {
							time.Sleep(time.Duration(rand.Intn(50)) * time.Microsecond)
							return v
						})

						if errors.Is(err, ErrErrorLimitExceeded) {
							t.Log("breaking the cycle")
							errIteration = i
							break
						}
					}
					t.Logf("error occured on iteration = %d", errIteration)
					if errIteration > tt.maxErrorsLimit+tt.maxTasksLimit {
						t.Errorf("error occured on task = %d, wanted less or equal %d",
							errIteration, tt.maxErrorsLimit+tt.maxTasksLimit)
					}

				})
		}
	})
	t.Run("simultanious tasks", func(t *testing.T) {
		t.Parallel()
		cases := []struct {
			maxTasksLimit int
			testTasksNum  int
		}{
			{
				3,
				100,
			},
			{
				5,
				50,
			},
			{
				1,
				20,
			},
		}

		for _, tt := range cases {
			tt := tt
			t.Run(fmt.Sprintf("maxLimit=%d,tasks=%d", tt.maxTasksLimit, tt.testTasksNum),
				func(t *testing.T) {
					t.Parallel()
					g := New(1)
					g.SetLimit(tt.maxTasksLimit)

					var active int32
					samples := make(chan int32, tt.testTasksNum)

					for i := 0; i < tt.testTasksNum; i++ {
						g.Go(func() error {
							samples <- atomic.AddInt32(&active, 1)
							time.Sleep(time.Duration(rand.Intn(50)) * time.Microsecond)
							atomic.AddInt32(&active, -1)
							return nil
						})
					}

					g.Wait()
					close(samples)
					var maxTasksActive int32
					for v := range samples {
						maxTasksActive = max(maxTasksActive, v)
					}

					if int(maxTasksActive) > tt.maxTasksLimit {
						t.Errorf("saw %d active gorounines; want <= %d", maxTasksActive, tt.maxTasksLimit)
					}

				})
		}
	})
}
