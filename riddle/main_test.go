package main

import (
	"testing"
)

func setupJugs() (Jug, Jug) {
	return Jug{
			name:   "Jug X",
			volume: 3,
		}, Jug{
			name:   "Jug Y",
			volume: 5,
		}
}

func failed(t *testing.T, s string, expected, got interface{}) {
	t.Errorf("%s:\n\texpected: %v\n\tgot: %v", s, expected, got)
}

func fatal(t *testing.T, s string, expected, got interface{}) {
	t.Fatalf("%s:\n\texpected: %v\n\tgot: %v", s, expected, got)
}

func errFatal(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("%s", err)
	}
}

func TestFill(t *testing.T) {
	X, Y := setupJugs()
	Y.volume = -2

	jugs, err := fill(X, Y)
	errFatal(t, err)
	if jugs == nil {
		fatal(t, "jugs is not nil", true, jugs != nil)
	}
	got, _ := getXYValue(jugs)
	if X.volume != got {
		failed(t, "Jug x value", X.volume, got)
	}

	// invalid bucket
	jugs, err = fill(Y, X)
	if jugs != nil {
		failed(t, "jugs is nil", true, jugs == nil)
	}
	if err == nil {
		failed(t, "error received", true, err != nil)
	}
}

func TestEmpty(t *testing.T) {
	X, Y := setupJugs()
	X.value = X.volume
	Y.value = Y.volume

	jugs, err := empty(X, Y)
	errFatal(t, err)
	if jugs == nil {
		fatal(t, "jugs is not nil", true, jugs != nil)
	}
	gotX, gotY := getXYValue(jugs)
	if gotX != 0 {
		failed(t, "Jug X value reset", 0, gotX)
	}
	if Y.value != gotY {
		failed(t, "Jug Y value no change", Y.value, gotY)
	}
}

func TestTransfer(t *testing.T) {
	X, Y := setupJugs()

	X.value = X.volume

	jugs, err := transfer(X, Y)
	errFatal(t, err)
	if jugs == nil {
		fatal(t, "jugs is not empty", true, jugs != nil)
	}

	gotX, gotY := getXYValue(jugs)
	if gotX != 0 {
		failed(t, "Jug X empty", 0, gotX)
	}
	if gotY != Y.value+X.value {
		failed(t, "Jug Y increased by X.value", Y.value+X.value, gotY)
	}

	// Transfer partial
	Y.value = 3
	jugs, err = transfer(X, Y)
	errFatal(t, err)
	if jugs == nil {
		fatal(t, "jugs is not empty", true, jugs != nil)
	}
	gotX, gotY = getXYValue(jugs)
	if X.value-(Y.volume-Y.value) != gotX {
		failed(t, "Jug X reduced", X.value-(Y.volume-Y.value), gotX)
	}
	if Y.value+(Y.volume-Y.value) != gotY {
		failed(t, "Jug Y increased to full", Y.value+(Y.volume-Y.value), gotY)
	}

	// Transfer empty
	X.value = 0
	jugs, err = transfer(X, Y)
	errFatal(t, err)
}

func TestValidate(t *testing.T) {
	if validate(-1, -1, -1) == nil {
		failed(t, "negative value", true, validate(-1, -1, -1) != nil)
	}
	if validate(2, 2, 1) == nil {
		failed(t, "same bucket size", true, validate(2, 2, 1) != nil)
	}
	if validate(3, 5, 7) == nil {
		failed(t, "Z larger than bucket", true, validate(3, 5, 7) != nil)
	}
	if validate(3, 5, 2) != nil {
		failed(t, "valid input", true, validate(3, 5, 2) == nil)
	}
}

func TestGetXYValue(t *testing.T) {
	X, Y := setupJugs()

	gotX, gotY := getXYValue([]Jug{X, X})
	if X.value != gotX {
		failed(t, "x value", X.value, gotX)
	}
	if Y.value != gotY {
		failed(t, "y value", Y.value, gotY)
	}

	gotX, gotY = getXYValue([]Jug{Y, X})
	if X.value != gotX {
		failed(t, "x value", X.value, gotX)
	}
	if Y.value != gotY {
		failed(t, "y value", Y.value, gotY)
	}
}

func TestRiddle(t *testing.T) {
	X, Y := setupJugs()
	v := make(map[key]bool)
	s, err := riddle([]Jug{X, Y}, 4, v)
	errFatal(t, err)
	expected := []result{
		{
			msg:  "Transfer from Jug Y to Jug X",
			xVal: 3,
			yVal: 4,
		},
		{
			msg:  "Fill Jug Y",
			xVal: 2,
			yVal: 5,
		},
		{
			msg:  "Transfer from Jug Y to Jug X",
			xVal: 2,
			yVal: 0,
		},
		{
			msg:  "Empty Jug X",
			xVal: 0,
			yVal: 2,
		},
		{
			msg:  "Transfer from Jug Y to Jug X",
			xVal: 3,
			yVal: 2,
		},
		{
			msg:  "Empty Jug X",
			xVal: 0,
			yVal: 5,
		},
		{
			msg:  "Fill Jug X",
			xVal: 3,
			yVal: 5,
		},
		{
			msg:  "Transfer from Jug X to Jug Y",
			xVal: 1,
			yVal: 5,
		},
		{
			msg:  "Fill Jug X",
			xVal: 3,
			yVal: 3,
		},
		{
			msg:  "Transfer from Jug X to Jug Y",
			xVal: 0,
			yVal: 3,
		},
		{
			msg:  "Fill Jug X",
			xVal: 3,
			yVal: 0,
		},
	}

	for i, v := range s {
		if expected[i] != v {
			failed(t, "response", expected[i], v)
		}
	}
}
