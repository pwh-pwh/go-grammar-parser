package grammar

import (
	"errors"
	"fmt"
	"grammar_parser_/utils"
	"strings"
)

func (g *Grammar) RemoveLeftRecurse() error {
	err := g.Indirect2()
	if err != nil {
		return err
	}
	g.Duplicate()
	err = g.Direct2()
	if err != nil {
		return err
	}
	g.Duplicate()
	g.ShowGrammar()
	return nil
}

func (g *Grammar) Indirect2() error {
	candidate := []int{}
	remove_index := []int{}
	//-------------------------把可能发生代入间接左递归的规则先加入容器---------------------------
	for i, graItem := range g.NewGrammar {
		if utils.ListIndexOf(g.Vn, string(graItem.Right[0])) != -1 &&
			utils.ListIndexOf(g.Vn, string(graItem.Right[0])) > utils.ListIndexOf(g.Vn, graItem.Left) {
			//右部首字母是非终结符并且右部首字母在左部字母之后
			candidate = append(candidate, i)
		}
	}
	//-------------------------如果有间接左递归，则插入代入之后的语法规则---------------------------
	for _, candItem := range candidate {
		rule := g.NewGrammar[candItem].Right
		temp := []string{}
		ut := g.NewGrammar[candItem].Left
		err := g.LeftRecurse(rule, ut, 1, &temp)
		if err != nil {
			return err
		}
		for _, tempItem := range temp {
			g.InsertGrammar(ut, tempItem)
			MyInsert(&remove_index, candItem)
		}
	}
	g.remove_qvector_index(remove_index)
	return nil
}

func (g *Grammar) Direct2() error {
	for _, vnItem := range g.Vn {
		flag := false
		remove_index := []int{}
		vmap := g.VnMapindex[vnItem]
		for _, vmapItem := range *vmap {
			if g.NewGrammar[vmapItem].Right == vnItem {
				return errors.New("至少存在以下情况之一: \n1.3步以上的推导才产生第一个字符是终结符号的情况;\n2.存在有害文法;\n3.存在间接左递归但没先消除;\n")
			}
		}
		if len(*vmap) == 1 &&
			strings.Index(
				g.NewGrammar[(*vmap)[0]].Right, vnItem,
			) != -1 {
			return errors.New("至少存在以下情况之一: \n1.3步以上的推导才产生第一个字符是终结符号的情况;\n2.存在有害文法;\n3.存在间接左递归但没先消除;\n")
		}
		for _, vmapItem := range *vmap {
			if string(g.NewGrammar[vmapItem].Right[0]) == vnItem {
				ut := vnItem
				newright := g.NewGrammar[vmapItem].Right
				//println newright
				fmt.Println("new_right:" + newright)
				g.NewGrammar[vmapItem].Left = string(int32(g.i2s))
				g.NewGrammar[vmapItem].Right = newright[1:] + string(int32(g.i2s))
				//myInsert(vnMapindex[QString(QChar(i2s))], vnMapindex[vn[j]][i]);//更改映射
				//remove_index.push_back(vnMapindex[vn[j]][i]);//不止一个直接左递归的时候，会删除映射
				tI, ok := g.VnMapindex[string(int32(g.i2s))]
				if !ok {
					tI = &[]int{}
				}
				MyInsert(tI, vmapItem)
				g.VnMapindex[string(int32(g.i2s))] = tI
				remove_index = append(remove_index, vmapItem)
				if !flag {
					g.InsertGrammar(string(g.i2s), "@")
					utMap := g.VnMapindex[ut]
					for _, utMapItem := range *utMap {
						if string(g.NewGrammar[utMapItem].Right[0]) != ut &&
							g.NewGrammar[utMapItem].Left != string(g.i2s) {
							g.NewGrammar[utMapItem].Right += string(g.i2s)
						}
					}
					flag = true
				}
			}
		}
		for _, item := range remove_index {
			utils.ListRemoveOne(g.VnMapindex[vnItem], item)
		}
		g.i2s++
	}
	return nil
}

func (g *Grammar) LeftRecurse(rule, unterminal string, layer int, can *[]string) error {
	if layer > 3 {
		return errors.New("至少存在以下情况之一: \n1.3步以上的推导才产生第一个字符是终结符号的情况;\n2.存在有害文法;\n3.存在间接左递归但没先消除;\n")
	}
	re := string(rule[0])
	//find_index := []int{}
	flag := false
	vnMap := g.VnMapindex[re]
	for _, vmpItem := range *vnMap {
		re2 := g.NewGrammar[vmpItem].Right
		if string(re2[0]) == unterminal {
			*can = append(*can, re2+rule[1:])
			flag = true
		} else if utils.ListIsContains(g.Vn, string(re2[0])) &&
			utils.ListIndexOf(g.Vn, string(re2[0])) > utils.ListIndexOf(g.Vn, re) {
			g.LeftRecurse(re2+rule[1:], unterminal, layer+1, can)
			if len(*can) != 0 {
				flag = true
			} else {
				*can = append(*can, re2+rule[1:])
			}
		} else {
			//can.push_back(re2 + rule.mid(1));
			*can = append(*can, re2+rule[1:])
		}
	}
	if !flag {
		*can = []string{}
	}
	return nil
}
