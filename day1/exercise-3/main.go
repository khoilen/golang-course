package main

import (
	"fmt"
	"sort"
)

func main() {
	var n int

	fmt.Print("Enter number of element: ")
	fmt.Scanln(&n)

	numbers := make([]int, n)

	for i := 0; i < n; i++ {
		fmt.Printf("Enter number %d: ", i+1)
		fmt.Scanln(&numbers[i])
	}

	sort.Ints(numbers)
	min := numbers[0]
	max := numbers[n-1]

	var total int
	for _, num := range numbers {
		total += num
	}
	average := float64(total) / float64(n)

	fmt.Printf("Min  %d\n", min)
	fmt.Printf("Max  %d\n", max)
	fmt.Printf("Sum  %d\n", total)
	fmt.Printf("Average %f\n", average)
	fmt.Printf("Sort Slice %d\n", numbers)
}
