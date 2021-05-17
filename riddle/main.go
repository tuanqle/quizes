package main

import (
	"fmt"
	"os"
)

type Action int64

const (
	Fill Action = iota
	Transfer
	Empty
)

var ActionToString = map[Action]string{
	Fill:     "Fill",
	Transfer: "Transfer",
	Empty:    "Empty",
}

var Actions = []func(Jug, Jug) ([]Jug, error){
	fill,
	transfer,
	empty,
}

type Jug struct {
	name   string
	volume int64
	value  int64
}

func fill(x, y Jug) ([]Jug, error) {
	if x.volume < 1 {
		return nil, fmt.Errorf("invalid volume")
	}
	x.value = x.volume
	return []Jug{x, y}, nil
}

func empty(x, y Jug) ([]Jug, error) {
	x.value = 0
	return []Jug{x, y}, nil
}

func transfer(x, y Jug) ([]Jug, error) {
	dif := y.volume - y.value
	if x.value == 0 || dif == 0 {
		return nil, nil
	}

	if dif < 0 {
		return nil, fmt.Errorf("invalid received bucket state")
	}
	if x.value >= dif {
		x.value = x.value - dif
		y.value = y.volume
	} else {
		y.value += x.value
		x.value = 0
	}
	return []Jug{x, y}, nil
}

func validate(x, y, z int64) error {
	if x <= 0 || y <= 0 || z <= 0 {
		return fmt.Errorf("volume can't be less than 1")
	}
	if x == y && z != x {
		return fmt.Errorf("Jug X and Y are the same")
	}

	largest := x
	if x < y {
		largest = y
	}

	if z > largest {
		return fmt.Errorf("Z amount is larger than the largest Jug can hold")
	}
	return nil
}

func getXYValue(jugs []Jug) (int64, int64) {
	if jugs[0].name == "Jug X" {
		return jugs[0].value, jugs[1].value
	}
	return jugs[1].value, jugs[0].value
}

func printStage(act Action, x Jug, y Jug) string {
	var str string
	switch act {
	case Fill:
		str = fmt.Sprintf("Fill %s", x.name)
	case Transfer:
		str = fmt.Sprintf("Transfer from %s to %s", x.name, y.name)
	case Empty:
		str = fmt.Sprintf("Empty %s", x.name)
	}
	return str
}

type result struct {
	msg  string
	xVal int64
	yVal int64
}

func riddle(jugs []Jug, z int64, visited map[key]bool) ([]result, error) {
	var (
		x, y = getXYValue(jugs)
		k    = key{
			xVal: x,
			yVal: y,
		}
	)

	if visited[k] {
		return nil, nil
	}

	visited[k] = true
	for _, j := range jugs {
		if j.value == z {
			return []result{}, nil
		}
	}

	for i := range jugs {
		for act, f := range Actions {
			if f == nil {
				return nil, fmt.Errorf("action func not found")
			}
			updated, err := f(jugs[i], jugs[(i+1)%2])
			if err != nil {
				return nil, err
			}
			if updated == nil {
				continue
			}
			results, err := riddle(updated, z, visited)
			if err != nil {
				return nil, err
			}
			if results != nil {
				x, y := getXYValue(updated)
				str := printStage(Action(act), jugs[i], jugs[(i+1)%2])
				results = append(results, result{
					msg:  str,
					xVal: x,
					yVal: y,
				})
				return results, nil
			}
		}
	}
	return nil, fmt.Errorf("No Solution")
}

type key struct {
	xVal, yVal int64
}

func bailIf(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func main() {
	fmt.Println("Water Jug Riddle")
	fmt.Printf("Enter value for X, Y, and Z (where X,Y,Z is number of gallon): ")

	var x, y, z int64
	fmt.Scanf("%d %d %d", &x, &y, &z)

	err := validate(x, y, z)
	bailIf(err)

	jugs := []Jug{
		{
			name:   "Jug X",
			volume: x,
		},
		{
			name:   "Jug Y",
			volume: y,
		},
	}

	visited := make(map[key]bool)
	results, err := riddle(jugs, z, visited)
	bailIf(err)

	fmt.Println("\n\t\t\t\t\tJug X\t\tJug Y")
	for i := len(results) - 1; i > -1; i-- {
		r := results[i]
		fmt.Printf("%30s\t\t%5d\t\t%5d\n", r.msg, r.xVal, r.yVal)
	}
}
