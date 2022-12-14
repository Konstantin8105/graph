package graph

import (
	"fmt"
	"math"
	"testing"
)

func TestWrong(t *testing.T) {
	tcs := []struct {
		x    float64
		data []Point
	}{
		{0, []Point{}},
		{+5, []Point{{-1, -2}, {1, 2}}},
		{-5, []Point{{-1, -2}, {1, 2}}},
		{0, []Point{{1, -2}, {-1, 2}}},
	}

	for i := range tcs {
		t.Run(fmt.Sprintf("%2d", i), func(t *testing.T) {
			_, err := Find(tcs[i].x, tcs[i].data...)
			if err == nil {
				t.Errorf("haven`t error")
			} else {
				t.Logf("%v", err)
			}
		})
	}
}

func Test(t *testing.T) {
	type testCase struct {
		x       float64
		data    []Point
		yExpect float64
	}
	tcs := []testCase{
		{0, []Point{{0, 2}}, 2},
		{0, []Point{{-1, -2}, {1, 2}}, 0},
		{0, []Point{{-1, 0}, {1, 4}}, 2},
		{1, []Point{{-1, 0}, {1, 4}}, 4},
		{1.1, []Point{{-1, 0}, {1, 4}, {2, 10}}, 4 + 0.1*6},
		{-1, []Point{{-1, 0}, {1, 4}}, 0},
	}

	for iter := 0; iter < 2; iter++ {
		// swap tests
		for i, size := 0, len(tcs); i < size; i++ {
			var t testCase
			t.x, t.yExpect = tcs[i].yExpect, tcs[i].x
			t.data = Swap(tcs[i].data...)
			tcs = append(tcs, t)
		}
	}

	for i := range tcs {
		t.Run(fmt.Sprintf("%2d", i), func(t *testing.T) {
			y, err := Find(tcs[i].x, tcs[i].data...)
			if err != nil {
				t.Errorf("haven`t error")
			}
			if math.Abs((y-tcs[i].yExpect)/y) > 1e-6 {
				t.Errorf("not valid Y: %v", y)
			}
		})
	}
}

func TestErrorDataset(t *testing.T) {
	for i := 0; i < 100; i++ {
		err := ErrorDataset{Id: DatasetErrorValue(i)}
		if len(err.Error()) == 0 {
			t.Errorf("not enought error information for : %d", i)
		}
	}
}

func TestErrorRange(t *testing.T) {
	for _, b := range []bool{false, true} {
		err := ErrorRange{IsUpper: b}
		if len(err.Error()) == 0 {
			t.Errorf("not enought error information for : %v", b)
		}
	}
}
