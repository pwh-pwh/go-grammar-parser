package grammar

import (
	"errors"
	"fmt"
	"grammar_parser_/utils"
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
	/**
	/*
	     * 将获取的文法文本按行分开
	     * param s: 需要处理的文法文本
	*/
	/*
		s.remove(" "); //去除多余的空格
		int index = 0;
		while(s[index] != '\x0')//如果没达到末尾则继续做
		{
			QString temp = "";
			while(s[index] != '\n')//直到换行符
			{
				temp.append(s[index]);
				index++;
			}
			if(temp != "")//如果不是空白行则加入原始文法向量中
				origin_grammar.push_back(temp);
			index++;
		}
	*/
	s = strings.ReplaceAll(s, " ", "")
	lines := strings.Split(s, "\n")
	for _, i := range lines {
		g.OriginGrammar = append(g.OriginGrammar, i)
	}
}

func (g *Grammar) SplitLeftAndRight() error {
	/*
	 * 将每行的文法区分成左部和右部
	 * param n: 可能出现的错误号
	 */

	//---------------------判断文法的合法性并收集左部符号-------------------
	/*int length = origin_grammar.size();
	for(int i=0;i<length;i++)
	{
	QStringList ql = origin_grammar[i].split("->");
	if(ql.size() == 1)//如果不能构成文法规则的输入则返回出错
	{
	n = 2;
	return;
	}
	else if(ql[1].indexOf("{") != -1 || ql[1].indexOf("[") != -1)
	{
	n = 10;
	return;
	}
	myInsert(vn, ql[0]);//收集非终止符号*/
	//}
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
	/**
	int i = 0;
	    int new_i = 0;
	    while(new_i<length)
	    {
	        QStringList ql = origin_grammar[i].split("->");//区分成左部和右部

	        if(ql[1].indexOf('(') != -1)
	        {
	            simplize(ql[1], n);//调用化简函数
	            if(n == 3 || n == 4)
	                return;
	        }

	        QStringList rarray = ql[1].split("|");//右部所有规则的序列
	        for(int j=0;j<rarray.size();j++)
	        {
	            int k=0;
	            while(rarray[j][k] != '\x0')
	            {
	                if((rarray[j][k] <= 'z' && rarray[j][k] >= 'a') || rarray[j][k] == '@') vt.insert(QString(rarray[j][k]));//收集终止符号
	                else vr.insert(QString(rarray[j][k]));//收集右部的非终止符号
	                k++;
	            }
	        }
	*/
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
	//-------------------所有序列加入文法------------------------
	/*grammarNode* g = new grammarNode(ql[0], rarray[0]);
		new_grammar.append(g);
		vnMapindex[ql[0]].append(new_i);
		for(int j=1;j<rarray.size();j++)
		{
		grammarNode* g = new grammarNode(ql[0], rarray[j]);
		new_grammar.append(g);
		new_i++;
		length++;
		vnMapindex[ql[0]].append(new_i);
		}
		//---------------------------------------------------------
		i++;//origin_grammar的下标
		new_i++;//new_grammar的下标
	}
	//qDebug()<<vn<<endl<<vt<<endl;
	//qDebug()<<vnMapindex<<endl;
	vr.subtract(vn.toList().toSet());//生成只在右部而不在左部出现的符号集，方便后续移除有害规则

	qDebug()<<"start";
	showGrammar();*/
	g.ShowGrammar()
	return nil
}

func (g *Grammar) Simplize(s *string) error {
	/*QStringList ql = s.split("(");//左部除去左括号
	if(ql.size() > 2)
	{
		n = 3;
		return;
	}*/
	ql := strings.Split(*s, "(")
	if len(ql) > 2 {
		return errors.New("有嵌套左公因子的情况存在，请引入新的非终结符号或直接拆开嵌套左公因子，消除嵌套左公因子的现象")
	}
	//--------------------------------------------------------------

	/*QString leftsplit = ql[0];//左部

	//-------------------------检查括号的合法性------------------------
	QStringList ql1 = ql[1].split(")");//右部除去右括号
	if(ql[0] == "" && ql1[0] == "")//如果括号多余直接返回
	{
		s = ql1[1];
		return;
	}
	if(ql1.size() < 2)//如果以右括号分割没有两部分字符串则说明少了右括号
	{
		n = 4;
		return;
	}*/

	leftsplit := ql[0]
	ql1 := strings.Split(ql[1], ")")
	if ql[0] == "" && ql1[0] == "" {
		*s = ql1[1]
		return nil
	}
	if len(ql1) < 2 {
		return errors.New("规则右部缺少右括号')'，请重新设计文法")
	}

	//--------------------------------------------------------------

	//------------拆分括号内的每一条规则和括号外的进行相乘-----------------
	/*QStringList middle = ql1[0].split("|");
	QString rightsplit = ql1[1];
	s = leftsplit + middle[0] + rightsplit;
	for(int i=1;i<middle.size();i++)
	{
	s = s + "|" + leftsplit + middle[i] + rightsplit;
	}*/
	middle := strings.Split(ql1[0], "|")
	rightsplit := ql1[1]
	*s = leftsplit + middle[0] + rightsplit
	for i := 1; i < len(middle); i++ {
		*s = *s + "|" + leftsplit + middle[i] + rightsplit
	}
	return nil
}

/**
void grammar::showGrammar()
{
    /*
     * 显示文法
*/
/*for(int i=0;i<new_grammar.size();i++)
{
qDebug()<<new_grammar[i]->left<<" "<<new_grammar[i]->right;
}
}
*/
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
	/*QVector<int> invalid_index;
	for(int i=0;i<new_grammar.size();i++)
	{
	for(int j=i+1;j<new_grammar.size();j++)
	{

	if(QString::compare(new_grammar[i]->left, new_grammar[j]->left) == 0 &&
	QString::compare(new_grammar[i]->right, new_grammar[j]->right) == 0)//发现重复的
	{
	myInsert(invalid_index, j);
	}
	}
	}

	remove_qvector_index(invalid_index);//对所有重复的行进行删除*/
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
	/*qSort(remove_index.begin(), remove_index.end());//为了方便循环删除需要先排序
	for(int ii=0;ii<remove_index.size();ii++)
	{
	remove_invalid(remove_index[ii] - ii);//相应变化ii
	}*/
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
	/*if(i >= 0 && i < new_grammar.size())
	{
		new_grammar.remove(i);
		reset_vnMapindex();
		reset_v();
		//showGrammar();
	}*/
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
	/*vnMapindex.clear();
	int length = new_grammar.size();
	for(int i=0;i<length;i++)
	{
	myInsert(vnMapindex[new_grammar[i]->left], i);
	}*/
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
	//vt.clear();

	/*QSet<QString> vl;//新的左部非终止符号
	//------------------对非终结符号先进行收集---------------------
	int length = new_grammar.size();
	for(int i=0;i<length;i++)//由于进行间接左递归需要有序，所以需要先收集左部的所有符号
	{
	vl.insert(new_grammar[i]->left);//左部的非终止符号
	}
	//----------------------------------------------------------

	//------------------对终结符号和只在右部出现的非终结符号进行收集---------------------
	for(int i=0;i<length;i++)
	{
	int j = 0;
	while(new_grammar[i]->right[j] != '\x0')
	{
	if((new_grammar[i]->right[j] <= 'z' && new_grammar[i]->right[j] >= 'a') || new_grammar[i]->right[j] == '@')
	vt.insert(QString(new_grammar[i]->right[j]));//收集终止符号
	else
	vr.insert(QString(new_grammar[i]->right[j]));//收集右部的非终止符号
	j++;
	}
	}*/
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

	/*vr.subtract(vn.toList().toSet());//减去终结符号剩下的就是只在右部出现的非终结符号

	//-----------------------更新其非终结符号集合----------------------------------
	QVector<QString> rmvl = vn.toList().toSet().subtract(vl).values().toVector();
	for(int i=0;i<rmvl.size();i++)
	{
	vn.removeOne(rmvl[i]);
	}*/
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

func (g *Grammar) LeftFactor(s []string) string {
	//-------------------------判断是否还有左公因子-------------------------------
	/*QString a = QString(s[0][0]);//获取第一条规则的第一个字符
	QString res = "";
	bool flag = true;
	int length = s.size();
	for(int i=1;i<length;i++)
	{
	if(a != QString(s[i][0]))//如果有一个规则没有了公共因子，则跳出
	{
	flag = false;
	break;
	}
	}*/
	a := string(s[0][0])
	res := ""
	flag := true
	for _, item := range s {
		if a != string(item[0]) {
			flag = false
			break
		}
	}
	//---------------------是否有左公因子执行相应的代码段-------------------------
	/*if(!flag || (flag && a == "@"))//如果已经没有公共因子，以a(b|c)的形式返回
	{
		res = "(" + s[0];
		for(int i=1;i<length;i++)
		{
		res = res + "|" + s[i];
		}
		res += ")";
		return res;
	}
	else if(flag)//如果还有公因子，递归执行
	{
		res = a;
		for(int i=0;i<length;i++)
		{
		s[i] = s[i].mid(1);
		if(s[i] == "")//如果已经空了，则加上@
		s[i] = "@";
		}
		return res + leftFactor(s);
	}
	//------------------------------------------------------------------------
	return res;*/
	if !flag || (flag && a == "@") {
		res = "(" + s[0]
		for i := 1; i < len(s); i++ {
			res = res + "|" + s[i]
		}
		res += ")"
		return res
	} else if flag {
		res = a
		for i := 0; i < len(s); i++ {
			s[i] = s[i][1:]
			if s[i] == "" {
				s[i] = "@"
			}
		}
		return res + g.LeftFactor(s)
	}
	return res
}

func (g *Grammar) InsertGrammar(l, r string) {
	/*
	 * 插入语法的函数
	 * param l: 左部
	 * param r: 右部
	 */
	/*grammarNode* g = new grammarNode(l, r);
	new_grammar.push_back(g);
	myInsert(vn, l);
	myInsert(vnMapindex[l], new_grammar.size()-1);*/
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
	/*
	 * 字符nt的first集合元素
	 * param nt: 需要提取的非终结符nt
	 * param layer: 层数
	 * param n: 可能发生的错误号
	 */
	/*QSet<QString> s;
	if(!vn.contains(nt))//如果是非终结符号，可以直接返回而不需要存储
	{
		s.insert(nt);
		return s;
	}
	else//如果是终结符号
	{
		for(int i=0;i<vnMapindex[nt].size();i++)
		{
		s += First(new_grammar[vnMapindex[nt][i]]->right, layer, n);//提取该规则的first元素
		}
	}
	//qDebug()<<s;
	first[nt] = s;//存储
	return s;*/
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
	/*
	 * 规则right的first集合元素
	 * param right: 一条规则的右部
	 * param layer: 层数
	 * param n: 可能发生的错误号
	 */
	/*if(layer > 3)//大于3层就不做了
	{
		n = 6;
		return {};
	}

	QSet<QString> f;
	int k = 0;
	int len = right.length();
	while(k < len)
	{
		QSet<QString> xk;
		//qDebug()<<right[k]<<endl;

		if(first.find(QString(right[k])) != first.end()) xk = first[QString(right[k])];//如果有存储，则直接使用
		else xk = get_first(QString(right[k]), layer+1, n);//没有则调用函数

		f = (f + xk).subtract({"@"});//减去@
		if(!xk.contains("@")) break;//如果不含有@需要马上退出
		k++;
	}
	if(k == len)//如果x1x2...xn的每个符号xi都有@，则本规则必有@
	{
		f.insert("@");
	}
	return f;*/
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

func (g *Grammar) LeftFactor2(rule string, already map[string]struct{}, layer int, can *[]string) error {
	/*
	 * 用于获取间接的左公因子
	 * param rule: 单条规则的右部
	 * param already: 该条规则的左部已经有的first集合
	 * param layer: 递归的层数
	 * param n: 可能出现的错误号
	 * param can: 所有有间接左公因子的规则
	 */
	/*if(layer > 3)//存在三层以上的左递归则返回提示错误
	{
		n = 6;
		return;
	}
	QString re = QString(rule[0]);//取出首字母
	for(int i=0;i<vnMapindex[re].size();i++)//循环这个首字母的所有规则
	{
	QString re2 = new_grammar[vnMapindex[re][i]]->right;//取出其规则
	if(vt.contains(QString(re2[0])))//如果是终结符号则代入
	{
	can.append(re2 + rule.mid(1));
	}
	else//如果是非终结符号则继续做下一层
	{
	leftFactor2(re2+rule.mid(1), already, layer+1, n, can);
	}
	}*/
	if layer > 3 {
		return errors.New("至少存在以下情况之一: \n1.3步以上的推导才产生第一个字符是终结符号的情况;\n2.存在有害文法;\n3.存在间接左递归但没先消除;")
	}
	re := string(rule[0])
	vnIndex := g.VnMapindex[re]
	for _, item := range *vnIndex {
		re2 := g.NewGrammar[item].Right
		if _, ok := g.Vt[string(re2[0])]; ok {
			*can = append(*can, re2+rule[1:])
		} else {
			err := g.LeftFactor2(re2+rule[1:], already, layer+1, can)
			if err != nil {
				return errors.New("至少存在以下情况之一: \n1.3步以上的推导才产生第一个字符是终结符号的情况;\n2.存在有害文法;\n3.存在间接左递归但没先消除;")
			}
		}
	}
	return nil
}

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
	/*QVector<int> candidate;
	QVector<int> remove_index;*/
	candidate := []int{}
	remove_index := []int{}
	for i, graItem := range g.NewGrammar {
		if utils.ListIndexOf(g.Vn, string(graItem.Right[0])) != -1 &&
			utils.ListIndexOf(g.Vn, string(graItem.Right[0])) > utils.ListIndexOf(g.Vn, graItem.Left) {
			candidate = append(candidate, i)
		}
	}
	//-------------------------把可能发生代入间接左递归的规则先加入容器---------------------------
	/*for(int i=0;i<new_grammar.size();i++)
	{
	if(vn.indexOf(QString(new_grammar[i]->right[0])) != -1 &&
	(vn.indexOf(QString(new_grammar[i]->right[0])) > vn.indexOf(new_grammar[i]->left)))
	//右部首字母是非终结符并且右部首字母在左部字母之后
	{
	candidate.push_back(i);
	}
	}*/
	//qDebug()<<"candidate"<<candidate;
	//-------------------------------------------------------------------------------------
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
	//-------------------------如果有间接左递归，则插入代入之后的语法规则---------------------------
	/*for(int j=0;j<candidate.size();j++)
	{
	QString rule = new_grammar[candidate[j]]->right;
	QVector<QString> temp;
	QString ut = new_grammar[candidate[j]]->left;
	leftRecurse(rule, ut, 1, n, temp);//求间接左递归
	if(n == 6)
	return;
	if(temp.size() != 0)//如果有间接左递归，则插入代入之后的语法规则
	{
	for(int i=0;i<temp.size();i++)
	{
	insertGrammar(ut, temp[i]);
	myInsert(remove_index, candidate[j]);//把这一行加入待删除行
	}
	}
	}*/
	//-------------------------------------------------------------------------------------

	//remove_qvector_index(remove_index);
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
			/*leftRecurse(re2+rule.mid(1), unterminal, layer+1, n, can);
			if(can.size() != 0)//说明深层有左递归
				flag = true;
			else
			can.push_back(re2 + rule.mid(1));*/
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
		//for(int i=0;i<remove_index.size();i++)
		//        {
		//            vnMapindex[vn[j]].removeOne(remove_index[i]);//移除原来的间接左递归的映射
		//        }
		//        i2s++;//同一个非终结符号做完之后再更新新引入符号
		for _, item := range remove_index {
			utils.ListRemoveOne(g.VnMapindex[vnItem], item)
		}
		g.i2s++
	}
	return nil
}
