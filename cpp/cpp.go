package cpp

// #cgo CXXFLAGS: -std=c++11
// #include "wrap_point.hpp"
// #include "tcpp.hpp"
// #include "gp.hpp"
import "C"

import "fmt"

func init() {
	fmt.Println("Hi from Go, about to calculate distance in C++ ...")
	distance := C.distance_between(1.0, 1.0, 2.0, 2.0)
	fmt.Printf("Go has result, distance is: %v\n", distance)
}
func P() {
	distance := C.distance_between(1.0, 1.0, 2.0, 2.0)
	fmt.Println(distance)
	msg := C.a(C.CString("你好"))
	fmt.Printf("msg:%v", C.GoString(msg))
}
