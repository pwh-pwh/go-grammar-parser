package tests

import (
	"fmt"
	"grammar_parser/cpp"
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

func insert(num int, arr []int) {
	arr = append(arr, num)
}

func TestInsert(t *testing.T) {
	arr := []int{}
	insert(1, arr)
	fmt.Println(arr)
}

func TestMap(t *testing.T) {
	m := make(map[string]*[]int)
	q := m["aa"]
	q = &[]int{1, 2}
	fmt.Println(m)
	fmt.Println(q)
}

func TestListNi(t *testing.T) {
	var arr *[]int = nil
	list := *arr
	fmt.Println(list)
}

func TestCpp(t *testing.T) {
	cpp.P()
	//fmt.Println(sum)
}
