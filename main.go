package main

import (
	"fmt"
	"grammar_parser_/grammar"
)

func main() {
	g, err := grammar.NewGrammar(`S->As|B
A->Sa
B->b`)
	if err != nil {
		fmt.Println(err)
	}
	g.Invalid()
	//g.ShowGrammar()
	/*err = g.RemoveLeftFactor()
	if err != nil {
		fmt.Println(err)
	}*/
	err = g.RemoveLeftRecurse()
	if err != nil {
		fmt.Println(err)
	}
}
