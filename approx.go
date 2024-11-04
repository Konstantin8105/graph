package graph

import (
	"errors"
	"fmt"
)

func Approx(withOutside bool, sort Check, data ...Point) (f func(x float64) (float64, error), err error) {
	defer func() {
		if err != nil {
			err = errors.Join(fmt.Errorf("graph.Approx"), err)
		}
	}()
	// }

	if len(data) < 2 {
		err = ErrorDataset{Id: NotEnougthData}
		return
	}
	// check input data
	if sort == CheckSorted {
		for i := range data {
			if i == 0 {
				continue
			}
			if data[i].X <= data[i-1].X {
				err = ErrorDataset{Id: NotSorted}
				return
			}
		}
	}
	// left border
	leftf := Linear([2]Point{data[0], data[1]})
	// right border
	rightf := Linear([2]Point{data[len(data)-2], data[len(data)-1]})

	return func(x float64) (y float64, err error) {

		// check is X inside graph
		if x < data[0].X {
			if withOutside {
				return leftf(x), nil
				// x0, y0 = data[0].X, data[0].Y
				// x1, y1 = data[1].X, data[1].Y
				// return y0 + (x-x0)/(x1-x0)*(y1-y0), nil
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
				return rightf(x), nil
				// x0, y0 = data[len(data)-2].X, data[len(data)-2].Y
				// x1, y1 = data[len(data)-1].X, data[len(data)-1].Y
				// return y0 + (x-x0)/(x1-x0)*(y1-y0), nil
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

		// old code:
		// var get func(i, j int) (ok bool, y float64)
		// get = func(left, right int) (ok bool, y float64) {
		// 	if right-left < 1 {
		// 		return
		// 	}
		// 	if right-left == 1 {
		// 		x0, y0 := data[left].X, data[left].Y
		// 		x1, y1 := data[right].X, data[right].Y
		// 		return true, y0 + (x-x0)/(x1-x0)*(y1-y0)
		// 	}
		// 	mid := (left + right) / 2
		// 	if x <= data[mid].X {
		// 		right = mid
		// 	} else {
		// 		left = mid
		// 	}
		// 	return get(left, right)
		// }
		// if ok, y := get(0, len(data)-1); ok {
		// 	return y, nil
		// }

		// new algo:
		left, mid, right := 0, 0, len(data)-1
		for left < right {
			if right-left == 1 {
				x0, y0 := data[left].X, data[left].Y
				x1, y1 := data[right].X, data[right].Y
				return y0 + (x-x0)/(x1-x0)*(y1-y0), nil
			}
			mid = (left + right) / 2
			if x <= data[mid].X {
				right = mid
			} else {
				left = mid
			}
		}

		err = ErrorDataset{Id: UndefinedData}
		return

	}, nil

}
