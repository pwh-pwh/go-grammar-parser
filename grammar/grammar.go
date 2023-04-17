package grammar

import (
	"errors"
	"fmt"
	"grammar_parser/utils"
	"sort"
	"strings"
)

type Grammar struct {
	OriginGrammar []string
	NewGrammar    []*GrammarNode
	VnMapindex    map[string]*[]int
	Vn            []string
	Vr            map[string]struct{}
	Vt            map[string]struct{}
	First         map[string]*map[string]struct{}
	Follow        map[string]*map[string]struct{}
	Table         map[string]*map[string]string
	i2s           int32
}

func NewGrammar(s string) (*Grammar, error) {
	g := new(Grammar)
	g.VnMapindex = make(map[string]*[]int)
	g.Vr = make(map[string]struct{})
	g.Vt = make(map[string]struct{})
	g.First = make(map[string]*map[string]struct{})
	g.Follow = make(map[string]*map[string]struct{})
	g.Table = make(map[string]*map[string]string)
	g.i2s = 945
	g.TextEdit2Vec(s)
	err := g.SplitLeftAndRight()
	if err != nil {
		return nil, err
	}
	return g, nil
}

func MyInsert[T int | string](vn *[]T, s T) {
	flag := true
	for _, v := range *vn {
		if s == v {
			flag = false
		}
	}
	if flag {
		*vn = append(*vn, s)
	}
}

func (g *Grammar) TextEdit2Vec(s string) {
	/*
	 * 将获取的文法文本按行分开
	 * param s: 需要处理的文法文本
	 */
	s = strings.ReplaceAll(s, " ", "")
	lines := strings.Split(s, "\n")
	for _, i := range lines {
		g.OriginGrammar = append(g.OriginGrammar, i)
	}
}

func (g *Grammar) SplitLeftAndRight() error {
	//-----------------------------------------------------------------
	length := len(g.OriginGrammar)
	for i := 0; i < length; i++ {
		ql := strings.Split(g.OriginGrammar[i], "->")
		if len(ql) == 1 {
			return errors.New("某个文法规则不能准确区分左部和右部，请重新设计文法")
		} else if strings.Index(ql[1], "{") != -1 || strings.Index(ql[1], "[") != -1 {
			return errors.New("输入不支持扩充BNF的形式")
		}
		//
		MyInsert(&g.Vn, ql[0])
	}

	i := 0
	new_i := 0
	for new_i < length {
		ql := strings.Split(g.OriginGrammar[i], "->")
		if strings.Index(ql[1], "(") != -1 {
			err := g.Simplize(&ql[1])
			if err != nil {
				return err
			}
		}
		rarray := strings.Split(ql[1], "|")
		for i := 0; i < len(rarray); i++ {
			for _, item := range rarray[i] {
				if item <= 'z' && item >= 'a' || item == '@' {
					g.Vt[string(item)] = struct{}{}
				} else {
					g.Vr[string(item)] = struct{}{}
				}
			}
		}
		//-------------------所有序列加入文法------------------------
		gNode := newNode(ql[0], rarray[0])
		g.NewGrammar = append(g.NewGrammar, gNode)
		ints, ok := g.VnMapindex[ql[0]]
		if !ok {
			ints = &[]int{}
		}
		*ints = append(*ints, new_i)
		g.VnMapindex[ql[0]] = ints
		for i := 1; i < len(rarray); i++ {
			gNode := newNode(ql[0], rarray[i])
			g.NewGrammar = append(g.NewGrammar, gNode)
			new_i++
			length++
			ints, ok := g.VnMapindex[ql[0]]
			if !ok {
				ints = &[]int{}
			}
			*ints = append(*ints, new_i)
		}
		new_i++
		i++
	}
	for _, it := range g.Vn {
		delete(g.Vr, it)
	}
	g.ShowGrammar()
	return nil
}

func (g *Grammar) Simplize(s *string) error {
	ql := strings.Split(*s, "(")
	if len(ql) > 2 {
		return errors.New("有嵌套左公因子的情况存在，请引入新的非终结符号或直接拆开嵌套左公因子，消除嵌套左公因子的现象")
	}
	leftsplit := ql[0]
	ql1 := strings.Split(ql[1], ")")
	if ql[0] == "" && ql1[0] == "" {
		*s = ql1[1]
		return nil
	}
	if len(ql1) < 2 {
		return errors.New("规则右部缺少右括号')'，请重新设计文法")
	}
	middle := strings.Split(ql1[0], "|")
	rightsplit := ql1[1]
	*s = leftsplit + middle[0] + rightsplit
	for i := 1; i < len(middle); i++ {
		*s = *s + "|" + leftsplit + middle[i] + rightsplit
	}
	return nil
}

func (g Grammar) ShowGrammar() {
	fmt.Println("show grammar")
	for _, item := range g.NewGrammar {
		fmt.Printf("right:%v left:%v  -", item.Left, item.Right)
	}
	fmt.Println()
}

func (g *Grammar) Duplicate() {
	/*
	 * 去重函数，进行双重循环把重复的语法删除
	 */
	var invalid_index = []int{}
	for i := 0; i < len(g.NewGrammar); i++ {
		for j := i + 1; j < len(g.NewGrammar); j++ {
			if g.NewGrammar[i].Left == g.NewGrammar[j].Left && g.NewGrammar[i].Right == g.NewGrammar[j].Right {
				invalid_index = append(invalid_index, j)
			}
		}
	}

}

//
func (g *Grammar) remove_qvector_index(removeIndex []int) {
	sort.Ints(removeIndex)
	for i := 0; i < len(removeIndex); i++ {
		g.remove_invalid(removeIndex[i] - i)
	}
}

func (g *Grammar) remove_invalid(i int) {
	/*
	 * 删除一条规则
	 * param i: 需要删除的语法规则行号
	 */
	if i >= 0 && i < len(g.NewGrammar) {
		fmt.Println("inv remove_invalid ngLen:%v", len(g.NewGrammar))
		g.NewGrammar = append(g.NewGrammar[0:i], g.NewGrammar[i+1:]...)
		g.reset_vnMapindex()
		g.reset_v()
	}
}

func (g *Grammar) reset_vnMapindex() {
	/*
	 * 更新<非终结符号，规则行>映射表
	 */
	g.VnMapindex = make(map[string]*[]int)
	length := len(g.NewGrammar)
	for i := 0; i < length; i++ {
		list, ok := g.VnMapindex[g.NewGrammar[i].Left]
		if !ok {
			list = &[]int{}
		}
		MyInsert(list, i)
		g.VnMapindex[g.NewGrammar[i].Left] = list
	}
	fmt.Println("update vmap:%v", len(g.VnMapindex))
	fmt.Println("newgrammar len:%v", length)
}

func (g *Grammar) reset_v() {
	/*
	 * 更新非终结符号和终结符号的集合，在删除某一行的时候会用到
	 */
	g.Vt = make(map[string]struct{})
	vl := make(map[string]struct{})
	length := len(g.NewGrammar)
	for i := 0; i < length; i++ {
		vl[g.NewGrammar[i].Left] = struct{}{}
	}
	for i := 0; i < length; i++ {
		for _, item := range g.NewGrammar[i].Right {
			if (item <= 'z' && item >= 'a') || item == '@' {
				g.Vt[string(item)] = struct{}{}
			} else {
				g.Vr[string(item)] = struct{}{}
			}
		}
	}

	for _, item := range g.Vn {
		delete(g.Vr, item)
	}
	rmvl := utils.Set2List(
		utils.SetSubList(
			utils.List2Set(g.Vn), utils.Set2List(vl),
		),
	)
	for _, item := range rmvl {
		for index, it := range g.Vn {
			if item == it {
				g.Vn = append(g.Vn[:index], g.Vn[index+1:]...)
				break
			}
		}
	}
}

func (g *Grammar) InsertGrammar(l, r string) {
	/*
	 * 插入语法的函数
	 * param l: 左部
	 * param r: 右部
	 */
	gNode := newNode(l, r)
	g.NewGrammar = append(g.NewGrammar, gNode)
	MyInsert(&g.Vn, l)
	list, ok := g.VnMapindex[l]
	if !ok {
		list = &[]int{}
	}
	MyInsert(list, len(g.NewGrammar)-1)
	g.VnMapindex[l] = list
}

func (g *Grammar) GetFirst(nt string, layer int) map[string]struct{} {
	s := make(map[string]struct{})
	if !utils.ListIsContains(g.Vn, nt) {
		s[nt] = struct{}{}
		return s
	} else {
		vnm := g.VnMapindex[nt]
		for _, item := range *vnm {
			fData, _ := g.FirstFun(g.NewGrammar[item].Right, layer)
			for dt := range fData {
				s[dt] = struct{}{}
			}
		}
		g.First[nt] = &s
		return s
	}
}

func (g *Grammar) FirstFun(right string, layer int) (map[string]struct{}, error) {
	if layer > 3 {
		return nil, errors.New("至少存在以下情况之一: \n1.3步以上的推导才产生第一个字符是终结符号的情况;\n2.存在有害文法;\n3.存在间接左递归但没先消除;\n")
	}
	f := make(map[string]struct{})
	k := 0
	for _, item := range right {
		xk := make(map[string]struct{})
		if _, ok := g.First[string(item)]; ok {
			xk = *g.First[string(item)]
		} else {
			xk = g.GetFirst(string(item), layer+1)
		}
		for key := range xk {
			f[key] = struct{}{}
		}
		delete(f, "@")
		if _, ok := xk["@"]; !ok {
			break
		}
		k++
	}
	if k == len(right) {
		f["@"] = struct{}{}
	}
	return f, nil
}
