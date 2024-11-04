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
			_, err := Find(tcs[i].x, false, CheckSorted, tcs[i].data...)
			if err == nil {
				t.Errorf("haven`t error")
			} else {
				t.Logf("%v", err)
			}
		})
	}
}

func TestData(t *testing.T) {
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
		t.Run(fmt.Sprintf("%04d", i), func(t *testing.T) {
			y, err := Find(tcs[i].x, tcs[i].withOutside, CheckSorted, tcs[i].data...)
			if err != nil {
				t.Fatalf("haven`t error: %v", err)
			}
			if math.Abs((y-tcs[i].yExpect)/y) > 1e-6 {
				t.Fatalf("not valid Y: %v != %v", y, tcs[i].yExpect)
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

func TestLinear(t *testing.T) {
	ps := []Point{
		{X: 1.1, Y: 2.1},
		{X: 3.1, Y: 3.050101010101},
	}
	f := Linear([2]Point{ps[0], ps[1]})
	for x := -0.5; x < 5; x += 0.2 {
		yl := f(x)
		y, err := Find(x, true, CheckSorted, ps...)
		if err != nil {
			t.Errorf("%v", err)
		}
		eps := math.Abs((yl - y) / y)
		if 1e-6 < eps {
			t.Fatalf("x=%.1f y: %.3f != %.3f", x, yl, y)
		}
		t.Logf("x=%4.1f y: %4.3f == %4.3f eps: %.3e", x, yl, y, eps)
	}
}

func TestLogLog(t *testing.T) {
	expect := func(x float64) float64 {
		arg := -2.6181*math.Log10(x) + 6.333
		y := math.Pow(10, math.Pow(10, arg)) - 1
		return y
	}
	appr := LogLog([2]Point{
		{X: 160 + 273, Y: expect(160 + 273)},
		{X: 260 + 273, Y: expect(260 + 273)},
	})
	for x := 160 + 273.0; x < 260+273.0; x += 0.02 {
		//t.Run(fmt.Sprintf("%.1f", x), func(t *testing.T) {
		e := expect(x)
		a := appr(x)
		eps := math.Abs((a - e) / e)
		if 1e-6 < eps {
			t.Fatalf("x=%.1f y: %.3f != %.3f", x, a, e)
		}
		if math.IsNaN(a) {
			t.Fatalf("x=%.1f y: %.3f != %.3f", x, a, e)
		}
		if math.IsInf(a, 0) {
			t.Fatalf("x=%.1f y: %.3f != %.3f", x, a, e)
		}
		t.Logf("x=%4.1f y: %4.3f == %4.3f eps: %.3e", x, a, e, eps)
		//})
	}
}

func expectF(x float64) float64 {
	return x*x + 1.1
}

func dataset() (ps []Point) {
	size := 10000
	for x := -10.0; x < 10.0; x += (20.0 / float64(size)) {
		ps = append(ps, Point{X: x, Y: expectF(x)})
	}
	return
}

func TestBigDataset(t *testing.T) {
	ps := dataset()
	for x := -10.0; x < 10.0; x += 0.0012 {
		// t.Run(fmt.Sprintf("%06.2f", x), func(t *testing.T) {
		y, err := Find(x, true, CheckSorted, ps...)
		if err != nil {
			t.Fatalf("x=%e. %v", x, err)
		}
		if eps := math.Abs((y - expectF(x)) / y); 1e-6 < eps {
			t.Errorf("precision x = %e y = [%e != %e]. eps = %e",
				x, y, expectF(x), eps)
		}
		// })
	}
}

// cpu: Intel(R) Xeon(R) CPU           X5550  @ 2.67GHz
// Benchmark/-1/5-8         	   63352	     19874 ns/op	       0 B/op	       0 allocs/op
// Benchmark/+1/5-8         	   48730	     23438 ns/op	       0 B/op	       0 allocs/op
// Benchmark/+1/2-8         	   38374	     31590 ns/op	       0 B/op	       0 allocs/op
// Benchmark/+4/5-8         	   32684	     36771 ns/op	       0 B/op	       0 allocs/op
// Benchmark/+6/5-8         	   29103	     41493 ns/op	       0 B/op	       0 allocs/op
//
// cpu: Intel(R) Xeon(R) CPU           X5550  @ 2.67GHz
// Benchmark/-1/5-8         	   74947	     15209 ns/op	       0 B/op	       0 allocs/op
// Benchmark/+1/5-8         	   60488	     19031 ns/op	       0 B/op	       0 allocs/op
// Benchmark/+1/2-8         	   63114	     19702 ns/op	       0 B/op	       0 allocs/op
// Benchmark/+4/5-8         	   60438	     19063 ns/op	       0 B/op	       0 allocs/op
// Benchmark/+6/5-8         	   77912	     15249 ns/op	       0 B/op	       0 allocs/op
//
// cpu: Intel(R) Xeon(R) CPU           X5550  @ 2.67GHz
// Benchmark/-1/5-8         	99403454	        12.10 ns/op	       0 B/op	       0 allocs/op
// Benchmark/+1/5-8         	11367764	       103.6 ns/op	       0 B/op	       0 allocs/op
// Benchmark/+1/2-8         	12548338	        95.78 ns/op	       0 B/op	       0 allocs/op
// Benchmark/+4/5-8         	11909890	       105.2 ns/op	       0 B/op	       0 allocs/op
// Benchmark/+6/5-8         	81059797	        15.51 ns/op	       0 B/op	       0 allocs/op
//
// cpu: Intel(R) Xeon(R) CPU           X5550  @ 2.67GHz
// Benchmark/-1/5-8         	96414884	        12.68 ns/op	       0 B/op	       0 allocs/op
// Benchmark/+1/5-8         	11853046	       105.3 ns/op	       0 B/op	       0 allocs/op
// Benchmark/+1/2-8         	13053266	        93.78 ns/op	       0 B/op	       0 allocs/op
// Benchmark/+4/5-8         	11500910	       103.5 ns/op	       0 B/op	       0 allocs/op
// Benchmark/+6/5-8         	78748124	        15.08 ns/op	       0 B/op	       0 allocs/op
//
// cpu: Intel(R) Xeon(R) CPU           X5550  @ 2.67GHz
// Benchmark/-1/5-16     	100000000	        11.30 ns/op	       0 B/op	       0 allocs/op
// Benchmark/+1/5-16     	13036677	        91.24 ns/op	       0 B/op	       0 allocs/op
// Benchmark/+1/2-16     	14085202	        83.60 ns/op	       0 B/op	       0 allocs/op
// Benchmark/+4/5-16     	13023190	        94.49 ns/op	       0 B/op	       0 allocs/op
// Benchmark/+6/5-16     	86749096	        13.74 ns/op	       0 B/op	       0 allocs/op
//
// Benchmark/-1/5-8    	103623026	        11.51 ns/op	       0 B/op	       0 allocs/op
// Benchmark/+1/5-8    	11721333	        101.5 ns/op	       0 B/op	       0 allocs/op
// Benchmark/+1/2-8    	12663447	        94.75 ns/op	       0 B/op	       0 allocs/op
// Benchmark/+4/5-8    	10692764	        98.84 ns/op	       0 B/op	       0 allocs/op
// Benchmark/+6/5-8    	96120549	        12.45 ns/op	       0 B/op	       0 allocs/op
//
// Benchmark/-1/5-8    	90896162	        13.02 ns/op	       0 B/op	       0 allocs/op
// Benchmark/+1/5-8    	25119486	        47.82 ns/op	       0 B/op	       0 allocs/op
// Benchmark/+1/2-8    	25630789	        46.47 ns/op	       0 B/op	       0 allocs/op
// Benchmark/+4/5-8    	23895999	        50.30 ns/op	       0 B/op	       0 allocs/op
// Benchmark/+6/5-8    	93084987	        12.30 ns/op	       0 B/op	       0 allocs/op
//
// cpu: Intel(R) Xeon(R) CPU           X5550  @ 2.67GHz
// Benchmark/Find-1/5-8         	41240526	        28.16 ns/op	       0 B/op	       0 allocs/op
// Benchmark/Find+1/5-8         	15268278	        69.50 ns/op	       0 B/op	       0 allocs/op
// Benchmark/Find+1/2-8         	17245395	        68.02 ns/op	       0 B/op	       0 allocs/op
// Benchmark/Find+4/5-8         	16446465	        71.86 ns/op	       0 B/op	       0 allocs/op
// Benchmark/Find+6/5-8         	26873926	        37.24 ns/op	       0 B/op	       0 allocs/op
// Benchmark/Appox-1/5-8        	131612073	         9.006 ns/op	       0 B/op	       0 allocs/op
// Benchmark/Appox+1/5-8        	24080991	        49.36 ns/op	       0 B/op	       0 allocs/op
// Benchmark/Appox+1/2-8        	22773793	        57.43 ns/op	       0 B/op	       0 allocs/op
// Benchmark/Appox+4/5-8        	24829455	        48.12 ns/op	       0 B/op	       0 allocs/op
// Benchmark/Appox+6/5-8        	132755846	         8.735 ns/op	       0 B/op	       0 allocs/op
func Benchmark(b *testing.B) {
	ps := dataset()
	bcs := []struct {
		name string
		x    float64
	}{
		{"-1/5", -12.0},
		{"+1/5", -6.0},
		{"+1/2", 1e-6},
		{"+4/5", +6.0},
		{"+6/5", +12.0},
	}
	for i := range bcs {
		b.Run("Find"+bcs[i].name, func(b *testing.B) {
			var err error
			b.ReportAllocs()
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				_, err = Find(bcs[i].x, true, NoCheckSorted, ps...)
				if err != nil {
					panic(err)
				}
			}
			_ = err
		})
	}
	for i := range bcs {
		b.Run("Appox"+bcs[i].name, func(b *testing.B) {
			f, err := Approx(true, NoCheckSorted, ps...)
			if err != nil {
				b.Fatal(err)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				_, err = f(bcs[i].x)
				if err != nil {
					panic(err)
				}
			}
			_ = err
		})
	}
}

// cpu: Intel(R) Xeon(R) CPU           X5550  @ 2.67GHz
// BenchmarkLogLog-16    	 4256602	       284.9 ns/op	       0 B/op	       0 allocs/op
// BenchmarkLogLog-16    	 4799954	       239.5 ns/op	       0 B/op	       0 allocs/op
// BenchmarkLogLog-8   	     4822491	       243.5 ns/op	       0 B/op	       0 allocs/op
// BenchmarkLogLog-8   	     4280361	       245.0 ns/op	       0 B/op	       0 allocs/op
func BenchmarkLogLog(b *testing.B) {
	expect := func(x float64) float64 {
		arg := -2.6181*math.Log10(x) + 6.333
		y := math.Pow(10, math.Pow(10, arg)) - 1
		return y
	}
	appr := LogLog([2]Point{
		{X: 160 + 273, Y: expect(160 + 273)},
		{X: 260 + 273, Y: expect(260 + 273)},
	})
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = appr(160 + 273.15)
	}
}
