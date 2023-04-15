package grammar

import (
	"grammar_parser_/utils"
	"strings"
)

func (g *Grammar) Invalid() {
	g.CannotArrive()
	g.CannotTerminal()
	g.ShowGrammar()
}

func (g *Grammar) CannotArrive() {
	i := 0
	for i < len(g.NewGrammar) {
		//---------------------如果是文法开始符号------------------------------
		if g.NewGrammar[i].Left == g.NewGrammar[0].Left {
			for _, item := range g.NewGrammar[i].Right {
				//如果规则中的非终止符号已经不可达
				if _, ok := g.Vt[string(item)]; !utils.ListIsContains(g.Vn, string(item)) && !ok {
					g.remove_invalid(i)
					i = -1
					break
				}
			}
			i++
			continue
		}
		//---------------------------如果是其他符号------------------------------
		flag := false
		for j := 0; j < len(g.NewGrammar); j++ {
			//如果该符号在所有的规则的某个右部出现则可达
			if strings.Index(g.NewGrammar[j].Right, g.NewGrammar[i].Left) != -1 && g.NewGrammar[j].Left != g.NewGrammar[i].Left {
				flag = true
				break
			}
		}
		if !flag {
			g.remove_invalid(i)
			i = 0
			continue
		}
		i++
	}
	//-------如果某个非终结符号只在右部出现也不可达-----------
	for vritem := range g.Vr {
		j := 0
		for j < len(g.NewGrammar) {
			if strings.Index(g.NewGrammar[j].Right, vritem) != -1 {
				g.remove_invalid(j)
				g.CannotArrive()
				continue
			}
			j++
		}
	}
}

func (g *Grammar) CannotTerminal() {
	//------------------------寻找不可终止的符号-------------------------------
	untermial := []string{}
	i := 0
	for i < len(g.NewGrammar) {
		if g.NewGrammar[i].Right == g.NewGrammar[i].Left { //形如U->U的有害规则
			g.remove_invalid(i)
			g.CannotArrive()
			continue
		}
		if strings.Index(g.NewGrammar[i].Right, g.NewGrammar[i].Left) != -1 { //存在U->aUb形式
			t := *g.VnMapindex[g.NewGrammar[i].Left]
			j := 0
			for ; j < len(t); j++ {
				if strings.Index(g.NewGrammar[t[j]].Right, g.NewGrammar[i].Left) == -1 { //不止U->aUb形式
					break
				}
			}
			if j == len(t) { //如果仅包含U->aUb形式的符号则为不可终止规则
				MyInsert(&untermial, g.NewGrammar[i].Left)
			}
		}
		i++
	}
	//------------------------移除不可终止的符号-------------------------------
	for i := 0; i < len(untermial); i++ {
		gi := g.VnMapindex[untermial[i]]
		//gLen := 0
		/*if gi == nil {
			//gLen = len(*gi)
			gi = &[]int{}
		}*/
		for j := 0; j < len(*gi); j++ { //把这个不可终止的符号的下标全部删除
			//list := g.VnMapindex[untermial[i]]
			g.remove_invalid((*gi)[j])
			g.CannotArrive()
		}
		for j := 0; j < len(g.NewGrammar); j++ { //把右部含有该不可终止符号的下标删除
			if utils.ListIsContains(g.Vn, g.NewGrammar[j].Left) && strings.Index(g.NewGrammar[j].Right, untermial[i]) != -1 {
				g.remove_invalid(j)
				g.CannotArrive()
			}
		}
		utils.ListRemoveOne(&g.Vn, untermial[i])
	}
}
