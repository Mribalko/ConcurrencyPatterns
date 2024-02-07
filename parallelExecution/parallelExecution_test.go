package parallel_execution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
)

func TestRun(t *testing.T) {

	errTest := errors.New("Err")

	cases := []struct {
		maxTasksLimit  int
		maxErrorsLimit int
		tasksRes       []error
	}{
		{
			2,
			2,
			[]error{nil, nil, nil},
		},
		{
			2,
			2,
			[]error{errTest, nil, errTest, nil},
		},
		{
			2,
			3,
			[]error{errTest, nil, errTest, nil},
		},
		{
			2,
			1,
			[]error{nil, nil, errTest},
		},
	}
	for _, tt := range cases {
		tt := tt
		t.Run(fmt.Sprintf("tasksLimit=%d,errLimit=%d,errors=%v", tt.maxTasksLimit, tt.maxErrorsLimit, tt.tasksRes),
			func(t *testing.T) {
				t.Parallel()

				var (
					tasks    []Task
					executed int32
					taskErr  int
				)

				for _, res := range tt.tasksRes {
					res := res
					if res != nil {
						taskErr++
					}
					tasks = append(tasks, func() error {
						atomic.AddInt32(&executed, 1)
						time.Sleep(time.Duration(rand.Intn(50)) * time.Microsecond)
						return res
					})
				}
				err := Run(tasks, tt.maxTasksLimit, tt.maxErrorsLimit)

				wantErr := taskErr >= tt.maxErrorsLimit

				if (err != nil) != wantErr {
					t.Errorf("error = %v, wantErr %v", err, wantErr)
				}

				if taskErr == 0 && int(executed) != len(tasks) {
					t.Errorf("got executions = %d, want = %d", executed, len(tasks))
				}

			})
	}
}
