package main

import "fmt"

func twoSum(nums []int, target int) []int {
	numMap := make(map[int]int)
	for i, num := range nums {
		complement := target - num
		index, found := numMap[complement]

		if found {
			return []int{index, i}
		}

		numMap[num] = i
	}

	return []int{}
}

func main() {
	tests := [][]int{
		{2, 7, 11, 15},
		{3, 2, 4},
		{3, 3},
	}

	targets := []int{9, 6, 6}

	for i, nums := range tests {
		result := twoSum(nums, targets[i])
		fmt.Printf("Input: nums = %v, target = %d\nOutput: %v\n", nums, targets[i], result)
	}
}
