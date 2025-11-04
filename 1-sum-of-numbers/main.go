package main

import (
	"fmt"
)

func Sum(numbers []int) int {
	if len(numbers) == 0{
		return 0
	}
	sum:=0

	for _,n := range numbers{
		sum+=n
	}
	return sum
}

func main(){
	fmt.Println(Sum([]int{1,2,3,4,5}))
	fmt.Println(Sum([]int{}))
}