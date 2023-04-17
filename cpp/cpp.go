package cpp

// #cgo CXXFLAGS: -std=c++17
// #include "gp.hpp"
import "C"

import (
	"fmt"
)

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
