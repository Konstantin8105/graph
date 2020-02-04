# graph

linear XY graph

```
package graph // import "."

type Point struct {
	X, Y float64
}
    Point coordinates

func Find(x float64, data ...Point) (y float64, err error)
    Find position Y by X in grapth dataset data. Dataset must by sorted by x.

```
