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
		{+0, []Point{{0, 2}}},
		{+0, []Point{}},
		{+5, []Point{{-1, -2}, {1, 2}}},
		{-5, []Point{{-1, -2}, {1, 2}}},
		{+0, []Point{{1, -2}, {-1, 2}}},
	}

	for i := range tcs {
		t.Run(fmt.Sprintf("%2d", i), func(t *testing.T) {
			_, err := Find(tcs[i].x, false, tcs[i].data...)
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
		x           float64
		data        []Point
		yExpect     float64
		withOutside bool
	}
	tcs := []testCase{
		{x: +0.0, data: []Point{{-1, -2}, {+1, +2}}, yExpect: 0},
		{x: +0.0, data: []Point{{-1, +0}, {+1, +4}}, yExpect: 2},
		{x: +1.0, data: []Point{{-1, +0}, {+1, +4}}, yExpect: 4},
		{x: -1.0, data: []Point{{-1, +0}, {+1, +4}}, yExpect: 0},
		{x: +1.1, data: []Point{{-1, +0}, {+1, +4}, {2, 10}}, yExpect: 4 + 0.1*6},
		{x: +2.0, data: []Point{{-1, -1}, {+1, +1}}, yExpect: 2, withOutside: true},
		{x: +2.0, data: []Point{{-10, -2}, {-1, -1}, {+1, +1}}, yExpect: 2, withOutside: true},
		{x: -2.0, data: []Point{{-1, -1}, {+1, +1}}, yExpect: -2, withOutside: true},
		{x: -2.0, data: []Point{{-1, -1}, {+1, +1}, {+100, +2}}, yExpect: -2, withOutside: true},
		{x: +4.0, data: []Point{{-2, -1}, {+2, +1}}, yExpect: 2, withOutside: true},
		{x: -4.0, data: []Point{{-2, -1}, {+2, +1}}, yExpect: -2, withOutside: true},
		{x: +4.0, data: []Point{{-2, +0}, {+2, +2}}, yExpect: 3, withOutside: true},
		{x: -4.0, data: []Point{{-2, +0}, {+2, +2}}, yExpect: -1, withOutside: true},
	}

	for iter := 0; iter < 2; iter++ {
		// swap tests
		for i, size := 0, len(tcs); i < size; i++ {
			var t testCase
			t.x, t.yExpect = tcs[i].yExpect, tcs[i].x
			t.data = Swap(tcs[i].data...)
			t.withOutside = tcs[i].withOutside
			tcs = append(tcs, t)
		}
	}

	for i := range tcs {
		t.Run(fmt.Sprintf("%2d", i), func(t *testing.T) {
			y, err := Find(tcs[i].x, tcs[i].withOutside, tcs[i].data...)
			if err != nil {
				t.Fatalf("haven`t error: %v", err)
			}
			if math.Abs((y-tcs[i].yExpect)/y) > 1e-6 {
				t.Errorf("not valid Y: %v != %v", y, tcs[i].yExpect)
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

func TestLinear(t *testing.T){
	ps := []Point {
		{X: 1.1, Y: 2.1},
		{X: 3.1, Y: 3.050101010101},
	}
	f := Linear([2]Point{ps[0], ps[1]})
	for x := -0.5; x< 5; x += 0.2{
		yl := f(x)
		y, err := Find(x, true, ps...)
		if err != nil {
			t.Errorf("%v", err)
		}
		eps := math.Abs((yl-y)/y)
		if 1e-6 < eps {
			t.Fatalf("x=%.1f y: %.3f != %.3f", x, yl, y)
		}
		t.Logf("x=%4.1f y: %4.3f == %4.3f eps: %.3e", x, yl, y, eps)
	}
}
