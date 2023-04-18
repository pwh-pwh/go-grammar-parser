package cpp

// #cgo CXXFLAGS: -std=c++17
// #include "gp.hpp"
import "C"

func GetRR(str string) string {
	data := C.getRR(C.CString(str))
	return C.GoString(data)
}

func GetRRAl(str string) string {
	data := C.getRRAL(C.CString(str))
	return C.GoString(data)
}

func GetFirst(str string) string {
	data := C.getFirstF(C.CString(str))
	return C.GoString(data)
}

func GetFollow(str string) string {
	data := C.getFollowF(C.CString(str))
	return C.GoString(data)
}

func GetTable(str string) string {
	data := C.getTableF(C.CString(str))
	return C.GoString(data)
}

func GetTree(str string, tok string) string {
	data := C.getTreeF(C.CString(str), C.CString(tok))
	return C.GoString(data)
}
