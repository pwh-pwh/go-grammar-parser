package service

import "grammar_parser/grammar"

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
	data := []GramTuple{}
	for _, item := range s.NewGrammar {
		data = append(data, GramTuple{
			Left:  item.Left,
			Right: item.Right,
		})
	}
	return data, nil
}
