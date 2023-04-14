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
