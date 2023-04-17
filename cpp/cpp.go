package cpp

// #cgo CXXFLAGS: -std=c++17
// #include "wrap_point.hpp"
// #include "tcpp.hpp"
// #include "gp.hpp"
import "C"

import (
	"fmt"
)

func init() {
	fmt.Println("Hi from Go, about to calculate distance in C++ ...")
	distance := C.distance_between(1.0, 1.0, 2.0, 2.0)
	fmt.Printf("Go has result, distance is: %v\n", distance)
}
func P() {

	s := `Program->Stmt_seq
Stmt_seq->Stmt_seq ; Statement|Statement
Statement->If_stmt|Repeat_stmt|Assign_stmt|Read_stmt|Write_stmt
If_stmt->if Exp then Stmt_seq end|if Exp then Stmt_seq else Stmt_seq end
Repeat_stmt->repeat Stmt_seq until Exp
Assign_stmt->id = Exp
Read_stmt->read id
Write_stmt->write Exp
Exp->Smp_exp Cmp_op Smp_exp|Smp_exp
Cmp_op-><|==
Smp_exp->Smp_exp Add Term|Term
Add->+|-
Term->Term MUL Factor|Factor
MUL->*|/
Factor->( Exp )|digit|id`
	s2 := `if digit < id
then id = digit ;
repeat
id = id * id ;
id = id - digit
until id == digit ;
write id
end`

	gStr := C.getTreeF(C.CString(s), C.CString(s2))
	fmt.Println("分割")
	fmt.Println(C.GoString(gStr))

	/*_s := `statement->if_stmt|other
	if_stmt->if ( E ) statement else_part
	else_part->else statement|@
	E->0|1`*/
	//gStr := C.getRR(C.CString(s))
	//_ = C.getRR(gStr)
	//fmt.Printf("rrStr:%v", C.GoString(rrStr))
	//fmt.Printf("data:%v", C.GoString(C.b(C.CString("aa"))))
}
