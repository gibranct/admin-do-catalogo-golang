package main

import "fmt"

func main() {
	nums := []int{1, 2}
	rotate(nums, 3)
	// fmt.Println(t)
}

func rotate(nums []int, k int) {
	size := len(nums)
	k = k % size
	if k == 0 {
		return
	}
	temp := make([]int, 0, size)
	temp = append(temp, nums[size-k:]...)
	temp = append(temp, nums[:size-k]...)

	copy(nums, temp)
	fmt.Println(nums)
	// return nums
}
