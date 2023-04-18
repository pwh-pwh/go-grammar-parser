package service

import (
	"grammar_parser/grammar"
)

type SpGramService struct {
	*grammar.Grammar
}

func NewSpGramService(gramStr string) (*SpGramService, error) {
	newGrammar, err := grammar.NewGrammar(gramStr)
	if err != nil {
		return nil, err
	}
	return &SpGramService{newGrammar}, nil
}

type GramTuple struct {
	Left  string
	Right string
}

func (s *SpGramService) GetInvalid() ([]GramTuple, error) {
	s.Invalid()
	return s.getGrammarData(), nil
}

func (s *SpGramService) GetRemoveLeftFactor() ([]GramTuple, error) {
	err := s.RemoveLeftFactor()
	if err != nil {
		return nil, err
	}
	return s.getGrammarData(), nil
}

func (s *SpGramService) GetRemoveLeftRecurse() ([]GramTuple, error) {
	err := s.RemoveLeftRecurse()
	if err != nil {
		return nil, err
	}
	return s.getGrammarData(), nil
}

func (s *SpGramService) GetFirst() ([]GramTuple, error) {
	first, err := s.GetAllFirst()
	if err != nil {
		return nil, err
	}
	return s.TrafTableData(first), nil
}

func (s *SpGramService) GetFollow() ([]GramTuple, error) {
	err := s.FollowFunc()
	if err != nil {
		return nil, err
	}
	return s.TrafTableData(s.Follow), nil
}

func (s *SpGramService) TrafTableData(data map[string]*map[string]struct{}) []GramTuple {
	var res []GramTuple
	for key, value := range data {
		var gt GramTuple
		gt.Left = key
		for vIt := range *value {
			gt.Right = gt.Right + vIt + ","
		}
		if len(gt.Right) != 0 {
			gt.Right = gt.Right[:len(gt.Right)-1]
		}
		res = append(res, gt)
	}
	return res
}

func (s *SpGramService) getGrammarData() []GramTuple {
	data := []GramTuple{}
	m := make(map[string]string)
	for _, item := range s.NewGrammar {
		if _, ok := m[item.Left]; ok {
			m[item.Left] = m[item.Left] + "|" + item.Right
		} else {
			m[item.Left] = item.Right
		}
	}
	for key, value := range m {
		data = append(data, GramTuple{
			Left:  key,
			Right: value,
		})
	}
	return data
}
