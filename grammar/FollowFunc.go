package grammar

import (
	"grammar_parser_/utils"
)

func (g *Grammar) FollowFunc() error {
	g.Follow[g.Vn[0]] = &map[string]struct{}{"$": struct{}{}}
	for i := 1; i < len(g.Vn); i++ {
		g.Follow[g.Vn[i]] = &map[string]struct{}{}
	}
	has_change := false
	for {
		has_change = false
		for _, ngItem := range g.NewGrammar {
			for j, rightItem := range ngItem.Right {
				if utils.ListIsContains(g.Vn, string(rightItem)) {
					if j == len(ngItem.Right)-1 {
						//如果有改变则置标志
						/*if(follow[QString(new_grammar[i]->right[j])] != follow[QString(new_grammar[i]->right[j])] + follow[new_grammar[i]->left])
							has_change = true;

						follow[QString(new_grammar[i]->right[j])] += follow[new_grammar[i]->left];*/
						fItem := g.Follow[string(rightItem)]
						addFol := utils.SetAdd(*fItem, *g.Follow[ngItem.Left])
						if !utils.IsSetEq(*fItem, addFol) {
							has_change = true
						}
						g.Follow[string(rightItem)] = &addFol
					}
					temp, err := g.FirstFun(ngItem.Right[j+1:], 1)
					if err != nil {
						return err
					}
					//如果有改变则置标志
					/*if(follow[QString(new_grammar[i]->right[j])] != (follow[QString(new_grammar[i]->right[j])] + temp).subtract({"@"}))
					has_change = true;

					follow[QString(new_grammar[i]->right[j])] = (follow[QString(new_grammar[i]->right[j])] + temp).subtract({"@"});*/
					fItem := g.Follow[string(rightItem)]
					addFol := utils.SetAdd(*fItem, temp)
					delete(addFol, "@")
					if !utils.IsSetEq(*fItem, addFol) {
						has_change = true
					}
					g.Follow[string(rightItem)] = &addFol
					if _, ok := temp["@"]; ok {
						//如果有改变则置标志
						/*if(follow[QString(new_grammar[i]->right[j])] != follow[QString(new_grammar[i]->right[j])] + follow[new_grammar[i]->left])
							has_change = true;

						follow[QString(new_grammar[i]->right[j])] += follow[new_grammar[i]->left];*/
						fItem := g.Follow[string(rightItem)]
						addFol := utils.SetAdd(*fItem, *g.Follow[ngItem.Left])
						if !utils.IsSetEq(*fItem, addFol) {
							has_change = true
						}
						g.Follow[string(rightItem)] = &addFol
					}
				}
			}
		}
		if !has_change {
			break
		}
	}
	return nil
}
