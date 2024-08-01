package main

import (
	"fmt"
	"log"
	"math"
)

// Viewable interface ensures that any type has a View method that returns a string
type Viewable interface {
	View() string
}

// Window struct represents a window with x, y coordinates and a model
type Window struct {
	x, y  int
	model Viewable
}

// Define a new type Direction
type Direction int

// Define constants for each direction
const (
	Left Direction = iota
	Right
	Up
	Down
)

// Map to hold direction vectors
var directions = map[Direction][2]int{
	Left:  {-1, 0},
	Right: {1, 0},
	Up:    {0, -1},
	Down:  {0, 1},
}

// FindNextWindow finds the closest window in the specified direction
func FindNextWindow(windows []*Window, currentWindow *Window, direction Direction) *Window {
	dirVector, exists := directions[direction]
	if !exists {
		fmt.Printf("Invalid direction: %v\n", direction)
		return nil
	}

	var bestCandidate *Window
	bestDistance := math.MaxFloat64

	for _, window := range windows {
		if window.x == currentWindow.x && window.y == currentWindow.y {
			continue
		}

		if isValidDirection(*currentWindow, *window, dirVector) {
			distance := calculateDistance(*currentWindow, *window, dirVector)
			if distance < bestDistance {
				bestDistance = distance
				bestCandidate = window
			}
		}
	}
	if bestCandidate == nil {
		bestCandidate = currentWindow
	}
	log.Printf("x %s y %s", bestCandidate.x, bestCandidate.y)
	return bestCandidate
}

// isValidDirection checks if the target window is in the specified direction from the current window
func isValidDirection(currentWindow Window, targetWindow Window, dirVector [2]int) bool {
	switch {
	case dirVector == directions[Right]:
		return targetWindow.x > currentWindow.x
	case dirVector == directions[Left]:
		return targetWindow.x < currentWindow.x
	case dirVector == directions[Up]:
		return targetWindow.y < currentWindow.y
	case dirVector == directions[Down]:
		return targetWindow.y > currentWindow.y
	default:
		return false
	}
}

// calculateDistance computes the Manhattan distance between the current and target windows
func calculateDistance(currentWindow Window, targetWindow Window, dirVector [2]int) float64 {
	return math.Abs(float64(targetWindow.x-currentWindow.x)) + math.Abs(float64(targetWindow.y-currentWindow.y))
}

// func main() {
// 	// Example windows with only positions
// 	windows := []Window{
// 		{0, 0},     // Window 0
// 		{300, 0},   // Window 1
// 		{300, 300}, // Window 2
// 	}

// 	currentWindow := windows[0]

// 	// Example navigation
// 	nextWindow := FindNextWindow(windows, currentWindow, "right")
// 	if nextWindow != nil {
// 		fmt.Printf("Next window in the 'right' direction: %+v\n", *nextWindow)
// 	} else {
// 		fmt.Println("No window found in the 'right' direction.")
// 	}
// }
