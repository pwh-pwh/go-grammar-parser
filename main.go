package main

import (
	"fmt"
	"grammar_parser_/grammar"
)

func main() {
	g, err := grammar.NewGrammar(`S->AbB|Bc
A->aA|@
B->d|e`)
	if err != nil {
		fmt.Println(err)
	}
	g.Invalid()
	g.ShowGrammar()
}
