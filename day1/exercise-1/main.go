package main

import (
	"fmt"
)

func rectangleArea(width, height float64) float64 {
	return width * height
}

func rectanglePerimeter(width, height float64) float64 {
	return 2 * (width + height)
}

func getInput() (float64, float64) {
	var width, height float64
	fmt.Print("Enter width: ")
	fmt.Scanln(&width)

	fmt.Print("Enter height: ")
	fmt.Scanln(&height)

	if width < 0 || height < 0 {
		fmt.Print("Width and Height should be greater than 0, please enter again\n")
		return getInput()
	}

	return width, height
}

func main() {
	width, height := getInput()

	area := rectangleArea(width, height)
	perimeter := rectanglePerimeter(width, height)

	fmt.Printf("Area of ​​rectangle: %.2f\n", area)
	fmt.Printf("Perimeter of ​​rectangle:  %.2f\n", perimeter)
}
