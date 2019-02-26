package main

import (
	"fmt"
	"os"
)

func main () {
	data := map[string]int {
		"Ross": 22,
		"Katie": 20,
	}

	data2 := map[string]int {
		"Ross": 22,
		"Katie": 20,
	}

	match := equal (data, data2) //return true (does match)

	fmt.Fprintf(os.Stdout, "%v\n", match)

	data = map[string]int {
		"Ross": 22,
		"Katie": 20,
	}

	data2 = map[string]int {
		"Ross": 22,
		"Josh": 22,
		"James": 21,
		"Jack": 21,
	}

	match = equal(data, data2) //returns false (the maps do not match)

	fmt.Fprintf(os.Stdout, "%v\n", match)
}

func equal (x, y map[string]int) bool {
	if len(x) != len(y) {
		return false
	}

	for k, xv := range x {
		yv, ok := y[k]
		if !ok || yv != xv {
			return false
		}
	}

	return true
}
