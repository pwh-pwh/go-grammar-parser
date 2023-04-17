package tests

import (
	"fmt"
	"grammar_parser/grammar"
	"testing"
)

func TestNewGrammar(t *testing.T) {
	newGrammar, err := grammar.NewGrammar(`S->aBe
B->a(b|c)`)
	if err != nil {
		fmt.Println(err)
	}
	if newGrammar != nil {
		//fmt.Println(*newGrammar)
	}
}
