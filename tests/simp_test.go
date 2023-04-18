package tests

import (
	"fmt"
	"strings"
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

/**

(===)===0===1===else===if===other===#===
statement===1===1===1===1===statement->if_stmt1statement->other1===1
if_stmt===1===1===1===1===if_stmt->if ( E ) statementelse_part1===1===1
else_part===1 ===1 ===1 ===1else_part->@1===1 ===1else_part->@1
E=== ===E->0E->1=== === === ===

(===)===0===1===else===if===other===#===
statement===============statement->if_stmtstatement->other===
if_stmt===============if_stmt->if ( E ) statement else_part======
else_part============else_part->@======else_part->@
E======E->0E->1============

(===)===0===1===else===if===other===#===
statement===============statement->if_stmt###statement->other###===
if_stmt===============if_stmt->if ( E ) statement else_part###======
else_part============else_part->@###======else_part->@###
E======E->0###E->1###============

*/

func TestSSp(t *testing.T) {
	s := "statement===============statement->if_stmt###statement->other###==="
	//split := strings.Split(s, "===")
	index := strings.Index(s, "===")
	fmt.Println(index)
	var str string = s[strings.Index(s, "==="):]
	fmt.Println(str)
	//todo
}

func TestT(t *testing.T) {
	//引入新的分割符号解决
	//todo
	s := "===1===1===1===1===1statement->if_stmt1statement->other1===1"
	split := strings.Split(s, "1")
	fmt.Println(len(split))
	for _, item := range split {
		fmt.Println("item:", item)
	}
}

/**

(===)===0===1===else===if===other===#===
statement===~~~===~~~===~~~===~~~===~~~statement->if_stmt~~~statement->other~~~===~~~
if_stmt===~~~===~~~===~~~===~~~===~~~if_stmt->if ( E ) statement else_part~~~===~~~===~~~
else_part===~~~===~~~===~~~===~~~else_part->@~~~===~~~===~~~else_part->@~~~
E===~~~===~~~E->0~~~E->1~~~===~~~===~~~===~~~===~~~

*/
func TestN2(t *testing.T) {
	s := "E===~~~===~~~E->0~~~E->1~~~===~~~===~~~===~~~===~~~"
	split := strings.Split(s, "~~~")
	fmt.Println(len(split))
	for _, item := range split {
		fmt.Println("item:", item)
	}
}
