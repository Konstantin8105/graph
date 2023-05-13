package graph

import (
	"fmt"
	"math"
)

// Point coordinates
type Point struct {
	X, Y float64
}

// Swap points coordinates
func Swap(ps ...Point) (swap []Point) {
	swap = make([]Point, len(ps))
	copy(swap, ps)
	for i := range swap {
		swap[i].X, swap[i].Y = swap[i].Y, swap[i].X
	}
	return
}

// Linear approximation by `y = a*X + b`
func Linear(ps [2]Point) (f func(x float64) float64) {
	dX := ps[1].X - ps[0].X
	dY := ps[1].Y - ps[0].Y
	if dX == 0 {
		panic("not valid dX is zero")
	}
	return func(x float64) float64 {
		return ps[0].Y + (x-ps[0].X)*dY/dX
	}
}

// LogLog approximation by `y = 10^(10^(a*Log10(x) + b))-1`
func LogLog(ps [2]Point) (f func(x float64) float64) {
	var (
		logX0    = math.Log10(ps[0].X)
		loglogY0 = math.Log10(math.Log10(ps[0].Y + 1.0))
		logX1    = math.Log10(ps[1].X)
		loglogY1 = math.Log10(math.Log10(ps[1].Y + 1.0))
	)

	// valid algoritm for verification
	// fl := Linear([2]Point{
	// 	{X: logX0, Y: loglogY0},
	// 	{X: logX1, Y: loglogY1},
	// })
	// return func(x float64) float64 {
	// 	y := fl(math.Log10(x))
	// 	return math.Pow(10, math.Pow(10, y)) - 1.0
	// }

	// return func(x float64) float64 {
	// 	return math.Pow(10, math.Pow(10, a*math.Log10(x)+b)) - 1.0
	// }

	// for performance
	// y = 10^(10^(a*Log10(x) + b))-1
	// y = 10^(10^(a*Log10(x))*10^b)-1
	// y = 10^(10^(Log10(x)*a)*10^b)-1
	// y = 10^((10^Log10(x))^a*10^b)-1
	// y = 10^(x^a*10^b)-1
	// C1 = 10^b
	// y = 10^(x^a*C1)-1
	// y = 10^(C1*x^a)-1
	// y = (10^C1)^(x^a)-1
	// C2 = 10^C1 = 10^10^b
	// y = C2^(x^a)-1
	//
	// a := (loglogY1 - loglogY0) / (logX1 - logX0)
	// b := loglogY1 - a*logX1
	// fmt.Println(a, b)
	// C2 := math.Pow(10, math.Pow(10, b))
	// return func(x float64) float64 {
	// 	return math.Pow(C2, math.Pow(x, a))-1
	// }

	// for performance
	// y = 10^(10^(a*Log10(x) + b))-1
	// y = 10^D-1
	// D = 10^(a*Log10(x) + b)
	// D = 10^(a*Log10(x)) * 10^b
	// D = 10^(Log10(x)*a) * 10^b
	// D = (10^Log10(x))^a * 10^b
	// D = x^a * 10^b
	// E = 10^b
	// D = x^a * E
	// y = 10^(x^a * E) - 1
	a := (loglogY1 - loglogY0) / (logX1 - logX0)
	b := loglogY1 - a*logX1
	E := math.Pow(10, b)
	return func(x float64) float64 {
		D := math.Pow(x, a) * E
		return math.Pow(10, D) - 1
	}
}

// Check is type of checking datasets
type Check bool

const (
	CheckSorted   Check = true
	NoCheckSorted Check = false
)

// Find position Y by X in grapth dataset data.
// Dataset shall by sorted by x.
func Find(x float64, withOutside bool, sort Check, data ...Point) (y float64, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("graph.Find error: %v", err)
		}
	}()

	// check input data
	if sort == CheckSorted {
		for i := range data {
			if i == 0 {
				continue
			}
			if data[i-1].X >= data[i].X {
				err = ErrorDataset{Id: NotSorted}
				return
			}
		}
	}
	if len(data) < 2 {
		err = ErrorDataset{Id: NotEnougthData}
		return
	}
	// check is X inside graph
	if x < data[0].X {
		if withOutside {
			x0, y0 := data[0].X, data[0].Y
			x1, y1 := data[1].X, data[1].Y
			return y0 + (x-x0)/(x1-x0)*(y1-y0), nil
		}
		err = ErrorRange{
			IsUpper: false,
			X:       x,
		}
		return
	}
	if x == data[0].X {
		return data[0].Y, nil
	}
	if data[len(data)-1].X < x {
		if withOutside {
			x0, y0 := data[len(data)-2].X, data[len(data)-2].Y
			x1, y1 := data[len(data)-1].X, data[len(data)-1].Y
			return y0 + (x-x0)/(x1-x0)*(y1-y0), nil
		}
		err = ErrorRange{
			IsUpper: true,
			X:       x,
		}
		return
	}
	if x == data[len(data)-1].X {
		return data[len(data)-1].Y, nil
	}
	// find
	var get func(i, j int) (ok bool, y float64)
	get = func(left, right int) (ok bool, y float64) {
		if right-left < 1 {
			return
		}
		if right-left == 1 {
			x0, y0 := data[left].X, data[left].Y
			x1, y1 := data[right].X, data[right].Y
			return true, y0 + (x-x0)/(x1-x0)*(y1-y0)
		}
		mid := (left + right) / 2
		if x <= data[mid].X {
			right = mid
		} else {
			left = mid
		}
		return get(left, right)
	}
	if ok, y := get(0, len(data)-1); ok {
		return y, nil
	}

	err = ErrorDataset{Id: UndefinedData}
	return
}

// ErrorRange is error for range error of X
type ErrorRange struct {
	IsUpper bool
	X       float64
}

// Error return string for implementation error type
func (err ErrorRange) Error() string {
	location := "lower"
	if err.IsUpper {
		location = "upper"
	}
	return fmt.Sprintf("Not acceptable X - out of range[%s]: %8.4f",
		location, err.X)
}

// DatasetErrorValue is value of error in dataset
type DatasetErrorValue int

// Constants of DatasetErrorValue`s
const (
	NotSorted DatasetErrorValue = iota + 1
	NotEnougthData
	UndefinedData
)

// ErrorDataset is error of dataset identification
type ErrorDataset struct {
	Id DatasetErrorValue
}

// Error return string for implementation error type
func (err ErrorDataset) Error() string {
	switch err.Id {
	case NotSorted:
		return "dataset is not sorted"
	case NotEnougthData:
		return "not enought data in dataset"
	case UndefinedData:
		return "undefined dataset"
	}
	return fmt.Sprintf("not valid error data Id: %d", int(err.Id))
}
