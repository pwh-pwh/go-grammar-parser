package main

import (
	"fmt"
	"grammar_parser_/grammar"
)

func main() {
	g, err := grammar.NewGrammar(`S->Be
B->Ce
B->Af
A->Ae|Dg
C->Cf
D->f
E->E
E->hk`)
	if err != nil {
		fmt.Println(err)
	}
	g.Invalid()
	g.ShowGrammar()
}
