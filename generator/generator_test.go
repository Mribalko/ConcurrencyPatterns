package generator

import (
	"reflect"
	"testing"
)

func Test_makeChannel(t *testing.T) {

	t.Run("Ints", func(t *testing.T) {
		t.Parallel()
		in := []int{1, 2, 3, 4}
		got := readChannel(t, MakeChannel(in))
		if !reflect.DeepEqual(got, in) {
			t.Errorf("makeChannel() = %v, want %v", got, in)
		}
	})

	t.Run("Strings", func(t *testing.T) {
		t.Parallel()
		in := []string{"one", "two", "three", "four"}
		got := readChannel(t, MakeChannel(in))
		if !reflect.DeepEqual(got, in) {
			t.Errorf("makeChannel() = %v, want %v", got, in)
		}
	})

	t.Run("Structs", func(t *testing.T) {
		t.Parallel()

		in := []struct {
			paramInt int
			paramStr string
		}{
			{1, "one"},
			{2, "two"},
		}

		got := readChannel(t, MakeChannel(in))
		if !reflect.DeepEqual(got, in) {
			t.Errorf("makeChannel() = %v, want %v", got, in)
		}
	})

}

func readChannel[T any](t *testing.T, ch <-chan T) []T {
	t.Helper()
	var res []T
	for v := range ch {
		res = append(res, v)
	}
	return res
}
