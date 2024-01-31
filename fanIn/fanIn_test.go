package fanin

import (
	"context"
	"reflect"
	"slices"
	"testing"

	gen "github.com/MRibalko/ConcurrencyPatterns/generator"
)

func Test_fanIn(t *testing.T) {
	inputData := [][]int{
		{1, 2, 3, 4},
		{6, 7, 8},
	}
	t.Run("Ints", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		inChan, want := prepareData(t, inputData)

		gotCh := fanIn(ctx, inChan...)

		var got []int
		for v := range gotCh {
			got = append(got, v)
		}
		slices.Sort(want)
		slices.Sort(got)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})

	t.Run("With cancel", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		inChan, _ := prepareData(t, inputData)

		gotCh := fanIn(ctx, inChan...)

		var got, want []int
		for v := range gotCh {
			got = append(got, v)
		}
		slices.Sort(want)
		slices.Sort(got)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want %v", got, want)
		}
	})
}

func prepareData[T any](t *testing.T, in [][]T) ([]<-chan T, []T) {
	t.Helper()
	var (
		resCh []<-chan T
		resSl []T
	)

	for _, v := range in {
		ch := gen.MakeChannel(v)
		resCh = append(resCh, ch)
		resSl = append(resSl, v...)
	}

	return resCh, resSl
}
