package main

import (
	"fmt"
	"math"
)

func square(i float64) float64 {
	value := math.Sqrt(i)
	return value
}

func main() {
	wynik := square(2)
	fmt.Println(wynik)
}
