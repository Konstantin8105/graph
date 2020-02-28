# graph

linear XY graph

```
package graph // import "github.com/Konstantin8105/graph"


FUNCTIONS

func Find(x float64, data ...Point) (y float64, err error)
    Find position Y by X in grapth dataset data. Dataset must by sorted by x.


TYPES

type DatasetErrorValue int
    DatasetErrorValue is value of error in dataset

const (
	NotSorted DatasetErrorValue = iota + 1
	NotEnougthData
)
    Constants of DatasetErrorValue`s

type ErrorDataset struct {
	Id DatasetErrorValue
}
    ErrorDataset is error of dataset identification

func (err ErrorDataset) Error() string
    Error return string for implementation error type

type ErrorRange struct {
	IsUpper bool
	X       float64
}
    ErrorRange is error for range error of X

func (err ErrorRange) Error() string
    Error return string for implementation error type

type Point struct {
	X, Y float64
}
    Point coordinates
```
