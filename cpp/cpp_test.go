package cpp

import (
	"fmt"
	"testing"
)

const s3 = "aa"

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

const token = `if digit < id
then id = digit ;
repeat
id = id * id ;
id = id - digit
until id == digit ;
write id
end`

const spToken = "read id"

func TestGetRR(t *testing.T) {
	fmt.Println(GetRR(s1))
}

func TestGetFirst(t *testing.T) {
	fmt.Println(GetFirst(s1))
}

func TestGetFollow(t *testing.T) {
	fmt.Println(GetFollow(s1))
}

func TestGetRRAl(t *testing.T) {
	fmt.Println(GetRRAl(s1))
}

func TestGetTable(t *testing.T) {
	fmt.Println(GetTable(s1))
}

func TestGetTree(t *testing.T) {
	fmt.Println(GetTree(s1, spToken))
}
