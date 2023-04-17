package grammar

import "grammar_parser/utils"

func (g *Grammar) GetAllFirst() (map[string]*map[string]struct{}, error) {
	for _, vnItem := range g.Vn {
		if _, ok := g.First[vnItem]; !ok {
			_, err := g.GetFirstWithErr(vnItem, 1)
			if err != nil {
				return nil, err
			}
		}
	}
	return g.First, nil
}

func (g *Grammar) GetFirstWithErr(nt string, layer int) (map[string]struct{}, error) {
	s := make(map[string]struct{})
	if !utils.ListIsContains(g.Vn, nt) {
		s[nt] = struct{}{}
		return s, nil
	} else {
		vnm := g.VnMapindex[nt]
		for _, item := range *vnm {
			fData, err := g.FirstFun(g.NewGrammar[item].Right, layer)
			if err != nil {
				return nil, err
			}
			for dt := range fData {
				s[dt] = struct{}{}
			}
		}
		g.First[nt] = &s
		return s, nil
	}
}
