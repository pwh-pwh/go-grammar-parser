package main

import (
	"fmt"
	"grammar_parser/grammar"
)

func main() {
	g, err := grammar.NewGrammar(`S->AbB|Bc
A->aA|@
B->d|e`)
	if err != nil {
		fmt.Println(err)
	}
	g.Invalid()
	//g.ShowGrammar()
	err = g.RemoveLeftFactor()
	if err != nil {
		fmt.Println(err)
	}
	err = g.RemoveLeftRecurse()
	if err != nil {
		fmt.Println(err)
	}
	allFirst, err := g.GetAllFirst()
	if err != nil {
		fmt.Println(err)
	} else {
		ShowForm(allFirst)
	}
	err = g.FollowFunc()
	if err != nil {
		fmt.Println(err)
	} else {
		ShowForm(g.Follow)
	}
}

func ShowForm(first map[string]*map[string]struct{}) {
	for key, value := range first {
		fmt.Print("非终结符:" + key + "= ")
		if value != nil {
			for key := range *value {
				fmt.Print(key + ",")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
