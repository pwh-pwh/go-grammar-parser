package grammar

import (
	"grammar_parser_/utils"
	"strings"
)

func (g *Grammar) RemoveLeftFactor() error {
	x := 0
	for x != len(g.NewGrammar) { //循环做左公因子的提取，直到不再改变语法的数量
		x = len(g.NewGrammar)
		err := g.Indirect()
		if err != nil {
			return err
		}
		g.Duplicate()
		g.Direct()
		g.Duplicate()
	}
	g.ShowGrammar()
	return nil
}

func (g *Grammar) Indirect() error {
	remove_index := []int{}
	for _, vnItem := range g.Vn { //遍历所有非终结符号
		tempArr := make([]int, len(*g.VnMapindex[vnItem]))
		copy(tempArr, *g.VnMapindex[vnItem])
		ls := &tempArr
		for _, git := range *ls {
			//------------------------求除了本规则外的所有规则的first集合--------------------------
			already := make(map[string]struct{}) //已经确定的first集合
			utils.ListRemoveOne(ls, git)
			for j := 0; j < len(*ls); j++ {
				firstData := g.GetFirst(string(g.NewGrammar[(*ls)[j]].Right[0]), 1)
				for s := range firstData {
					already[s] = struct{}{}
				}
			}
			if string(g.NewGrammar[git].Right[0]) <= "z" && string(g.NewGrammar[git].Right[0]) >= "a" || (g.NewGrammar[git].Right[0]) == '@' {
				continue
			}
			if utils.SetIntersects(g.GetFirst(string(g.NewGrammar[git].Right[0]), 1), already) {
				inner := utils.SetIntersect(g.GetFirst(string(g.NewGrammar[git].Right[0]), 1), already)
				temp := []string{}
				err := g.LeftFactor2(g.NewGrammar[git].Right, inner, 1, &temp)
				if err != nil {
					return err
				}
				for _, j := range temp {
					g.InsertGrammar(g.NewGrammar[git].Left, j)
				}
				MyInsert(&remove_index, git)
			}
		}
	}
	g.remove_qvector_index(remove_index)
	return nil
}

func (g *Grammar) Direct() {
	remove_index := []int{}
	for _, value := range g.VnMapindex {
		if len(*value) != 1 {
			splits := make(map[string]*[]int)
			for _, item := range *value {
				tempQ, ok := splits[string(g.NewGrammar[item].Right[0])]
				if !ok {
					tempQ = &[]int{}
				}
				*tempQ = append(*tempQ, item)
				splits[string(g.NewGrammar[item].Right[0])] = tempQ
			}
			for _, temp2 := range splits {
				if len(*temp2) == 1 {
					continue
				}
				s := make([]string, 0)
				for _, temp2Item := range *temp2 {
					s = append(s, g.NewGrammar[temp2Item].Right)
				}
				lf := g.LeftFactor(s)
				ql := strings.Split(lf, "(")
				if ql[0] != "" {
					ql[1] = ql[1][:len(ql[1])-1]
					ql2 := strings.Split(ql[1], "|")
					for _, ql2Item := range ql2 {
						g.InsertGrammar(string(int32(g.i2s)), ql2Item)
					}
					g.NewGrammar[(*temp2)[0]].Right = ql[0] + string(int32(g.i2s))
					g.i2s++
					*temp2 = append((*temp2)[1:])
					for _, temp2Item := range *temp2 {
						MyInsert(&remove_index, temp2Item)
					}
				}
			}
		}
	}
	g.remove_qvector_index(remove_index)
}
