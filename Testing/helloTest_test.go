package main

import "testing"

func TestHello(t *testing.T) {
	var usecase1 float64
	//var usecase2 float64

	usecase1 = square(22)
	if usecase1 != 2 {
		t.Errorf("Expected: integer got: %v", usecase1)
	}
}
