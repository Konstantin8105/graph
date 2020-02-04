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
	tcs := []struct {
		x       float64
		data    []Point
		yExpect float64
	}{
		{0, []Point{{0, 2}}, 2},
		{0, []Point{{-1, -2}, {1, 2}}, 0},
		{0, []Point{{-1, 0}, {1, 4}}, 2},
		{1, []Point{{-1, 0}, {1, 4}}, 4},
		{-1, []Point{{-1, 0}, {1, 4}}, 0},
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
