package tests

import (
	"fmt"
	"testing"
)

func TestChar2Str(t *testing.T) {
	c := 'a'
	s := string(c)
	fmt.Println(s)
	fmt.Println(c)
}

func TestNQ(t *testing.T) {
	var arr []int
	if arr == nil {
		fmt.Println("nil arr")
	} else {
		fmt.Println("0 len arr")
	}
}

func addInt(num int, arr *[]int) {
	*arr = append(*arr, num)
}

func TestAddNum(t *testing.T) {
	num := 1
	var arr []int
	addInt(num, &arr)
	for _, item := range arr {
		fmt.Println(item)
	}
}
