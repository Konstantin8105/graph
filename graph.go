package graph

import "fmt"

// Point coordinates
type Point struct {
	X, Y float64
}

// Find position Y by X in grapth dataset data.
// Dataset must by sorted by x.
func Find(x float64, data ...Point) (y float64, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("graph.Find error: %v", err)
		}
	}()

	// check input data
	for i := range data {
		if i == 0 {
			continue
		}
		if data[i-1].X >= data[i].X {
			err = fmt.Errorf("input data is not sorted by x. Index [%d,%d]: %8.4f >= %8.4f",
				i-1, i, data[i-1].X, data[i].X)
			return
		}
	}
	if len(data) == 0 {
		err = fmt.Errorf("Not enought dataset")
		return
	}
	if x < data[0].X {
		err = fmt.Errorf("Not acceptable X - out of range[lower]: %8.4f", x)
		return
	}
	if x > data[len(data)-1].X {
		err = fmt.Errorf("Not acceptable X - out of range[upper]: %8.4f", x)
		return
	}

	// find
	for i := range data {
		if x == data[i].X {
			return data[i].Y, nil
		}
		if i == 0 {
			continue
		}
		if data[i-1].X <= x && x <= data[i].X {
			x0, y0 := data[i-1].X, data[i-1].Y
			x1, y1 := data[i].X, data[i].Y
			return y0 + (x-x0)/(x1-x0)*(y1-y0), nil
		}
	}
	return -42, fmt.Errorf("Out of grapth")
}
