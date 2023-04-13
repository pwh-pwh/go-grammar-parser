package grammar

type GrammarNode struct {
	Left  string
	Right string
}

func newNode(left, right string) *GrammarNode {
	return &GrammarNode{
		left, right,
	}
}
