package main

import (
	"parking_lot_simulator/scenes"
)

const (
	numCars = 100
)

func main() {
	scenes.StartArrival(numCars)
}
