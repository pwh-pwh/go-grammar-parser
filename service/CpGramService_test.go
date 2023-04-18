package service

import (
	"fmt"
	_ "reflect"
	"testing"
)

const s1 = `Program->Stmt_seq
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
const s2 = `statement->if_stmt|other
if_stmt->if ( E ) statement else_part
else_part->else statement|@
E->0|1`

func TestCpGramService_CpGetRR(t *testing.T) {
	rr, _ := CpGetRR(s1)
	for _, item := range rr {
		fmt.Printf("item:%v\n", item)
	}
}

func TestCpGetRRAl(t *testing.T) {
	al, err := CpGetRRAl("aa")
	if err != nil {
		fmt.Println(err)
	}
	for _, item := range al {
		fmt.Printf("item:%v\n", item)
	}
}

func TestCpGetFirst(t *testing.T) {
	first, err := CpGetFirst(s1)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, item := range first {
		fmt.Printf("left:%v right:%v\n", item.Left, item.Right)
	}
}

func TestCpGetFollow(t *testing.T) {
	follow, err := CpGetFollow(s1)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, item := range follow {
		fmt.Printf("left:%v right:%v\n", item.Left, item.Right)
	}
}

//(☉)☉*☉+☉-☉/☉;☉<☉=☉==☉digit☉else☉end☉id☉if☉read☉repeat☉then☉until☉write☉#☉
//(☉)☉*☉+☉-☉/☉;☉<☉=☉==☉digit☉else☉end☉id☉if☉read☉repeat☉then☉until☉write☉#☉
func TestCpGetTable(t *testing.T) {
	table, err := CpGetTable(s1)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, item := range table {
		for _, stItem := range item {
			fmt.Printf("str:%v  ", stItem)
		}
		fmt.Println()
	}
}
