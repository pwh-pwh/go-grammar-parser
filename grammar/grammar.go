package grammar

import (
	"errors"
	"fmt"
	"grammar_parser_/utils"
	"sort"
	"strings"
)

/**
QVector<grammarNode*> new_grammar;//分左右部的存储

    QMap<QString, QVector<int>> vnMapindex;//非终结符号到所在行的映射

    QVector<QString> vn;//非终结符号集合

    QSet<QString> vr;//只在右部出现的非终结符号集合

    QSet<QString> vt;//终结符号集合

    QMap<QString, QSet<QString>> first;//每个非终结符号的first集合元素

    QMap<QString, QSet<QString>> follow;//每个非终结符号的follow集合元素

    QMap<QString, QMap<QString, QString>> table;//LL(1)分析表
*/
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
	i2s           int
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

func (g *Grammar) CannotArrive() {
	i := 0
	for i < len(g.NewGrammar) {
		//---------------------如果是文法开始符号------------------------------
		/*if(QString::compare(new_grammar[i]->left, new_grammar[0]->left) == 0)
		//其不可能在其他规则右部出现，但会出现其右部的非终止符号已经不可达
		{
		int p = 0;
		while(new_grammar[i]->right[p] != '\x0')
		{
		if(!vn.contains(QString(new_grammar[i]->right[p])) && !vt.contains(QString(new_grammar[i]->right[p])))
		//如果规则中的非终止符号已经不可达
		{
		//qDebug()<<i<<" "<<new_grammar[i]->right;
		remove_invalid(i);//移除该行
		i = -1;//由于后面会进行i++，因此使其等于-1
		break;
		}
		p++;
		}
		i++;
		continue;
		}*/
		if g.NewGrammar[i].Left == g.NewGrammar[0].Left {
			for _, item := range g.NewGrammar[i].Right {

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
		/*bool flag = false;
		for(int j=0;j<new_grammar.size();j++)
		{
		//如果该符号在所有的规则的某个右部出现则可达
		if((new_grammar[j]->right.indexOf(new_grammar[i]->left) != -1) && (new_grammar[j]->left != new_grammar[i]->left))
		{
		flag = true;
		break;
		}
		}
		if(!flag)
		{
			//invalid_index.insert(i);
			remove_invalid(i);
			i = 0;
			continue;
		}
		i++;*/
		flag := false
		for j := 0; j < len(g.NewGrammar); j++ {
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
	/*QVector<QString> vrv = vr.values().toVector();
	for(int i=0;i<vrv.size();i++)
	{
	int j=0;
	while(j<new_grammar.size())
	{
	if(new_grammar[j]->right.indexOf(vrv[i]) != -1)
	{
	remove_invalid(j);
	cannot_arrive();//再次清除不可达的符号
	continue;
	}
	j++;
	}
	}*/
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
	/*QVector<QString> untermial;
	int i=0;
	while(i<new_grammar.size())
	{
		if(QString::compare(new_grammar[i]->right, new_grammar[i]->left)==0)//形如U->U的有害规则
		{
		remove_invalid(i);
		cannot_arrive();//再次清除不可达的规则
		continue;
		}
		if(new_grammar[i]->right.indexOf(new_grammar[i]->left) != -1)//存在U->aUb形式
		{
			QVector<int> t = vnMapindex[new_grammar[i]->left];
			int j=0;
			for(;j<t.size();j++)
			{
			if(new_grammar[t[j]]->right.indexOf(new_grammar[i]->left) == -1)//不止U->aUb形式
			{
			break;
			}
			}
			if(j==t.size())//如果仅包含U->aUb形式的符号则为不可终止规则
			{
				myInsert(untermial, new_grammar[i]->left);
			}
		}
		i++;
	}*/
	untermial := []string{}
	i := 0
	for i < len(g.NewGrammar) {
		if g.NewGrammar[i].Right == g.NewGrammar[i].Left {
			g.remove_invalid(i)
			g.CannotArrive()
			continue
		}
		if strings.Index(g.NewGrammar[i].Right, g.NewGrammar[i].Left) != -1 {
			t := *g.VnMapindex[g.NewGrammar[i].Left]
			j := 0
			for ; j < len(t); j++ {
				if strings.Index(g.NewGrammar[t[j]].Right, g.NewGrammar[i].Left) == -1 {
					break
				}
			}
			if j == len(t) {
				MyInsert(&untermial, g.NewGrammar[i].Left)
			}
		}
		i++
	}
	//------------------------移除不可终止的符号-------------------------------
	/*for(int i=0;i<untermial.size();i++)
	{
	for(int j=0;j<vnMapindex[untermial[i]].size();j++)//把这个不可终止的符号的下标全部删除
	{
	remove_invalid(vnMapindex[untermial[i]][j]);
	cannot_arrive();
	}
	for(int j=0;j<new_grammar.size();j++)//把右部含有该不可终止符号的下标删除
	{
	if(vn.contains(new_grammar[j]->left) && (new_grammar[j]->right.indexOf(untermial[i])!=-1))
	{
	remove_invalid(j);
	cannot_arrive();
	}
	}
	vn.removeOne(untermial[i]);
	}*/
	for i := 0; i < len(untermial); i++ {
		gi := g.VnMapindex[untermial[i]]
		//gLen := 0
		/*if gi == nil {
			//gLen = len(*gi)
			gi = &[]int{}
		}*/
		for j := 0; j < len(*gi); j++ {
			//list := g.VnMapindex[untermial[i]]
			g.remove_invalid((*gi)[j])
			g.CannotArrive()
		}
		for j := 0; j < len(g.NewGrammar); j++ {
			if utils.ListIsContains(g.Vn, g.NewGrammar[j].Left) && strings.Index(g.NewGrammar[j].Right, untermial[i]) != -1 {
				g.remove_invalid(j)
				g.CannotArrive()
			}
		}
		utils.ListRemoveOne(&g.Vn, untermial[i])
	}
}

func (g *Grammar) Invalid() {
	/*cannot_arrive();

	cannot_terminal();

	qDebug()<<"remove_invalid"<<vnMapindex;
	showGrammar();*/
	g.CannotArrive()
	g.CannotTerminal()
	g.ShowGrammar()
}

func (g *Grammar) RemoveLeftFactor() error {
	/*int x = 0;
	while(x != new_grammar.size())//循环做左公因子的提取，直到不再改变语法的数量
	{
		x = new_grammar.size();

		indirect(n);

		if(n == 6)
		return;

		duplicate();

		direct(n);

		duplicate();
	}

	qDebug()<<"remove_left_factor";
	showGrammar();*/
	x := 0
	for x != len(g.NewGrammar) {
		x = len(g.NewGrammar)
		//todo vmmap 被清空
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
	for _, vnItem := range g.Vn {
		tempArr := make([]int, len(*g.VnMapindex[vnItem]))
		copy(tempArr, *g.VnMapindex[vnItem])
		ls := &tempArr
		for _, git := range *ls {
			already := make(map[string]struct{})
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
	/*for(int k=0;k<vn.size();k++)//遍历所有非终结符号
	{
	for(int i=0;i<vnMapindex[vn[k]].size();i++)//这些非终结符号的行
	{
	//------------------------求除了本规则外的所有规则的first集合--------------------------
	QSet<QString> already;//已经确定的first集合
	QVector<int> ls = vnMapindex[vn[k]];//取出所有规则
	ls.removeOne(vnMapindex[vn[k]][i]);//移除本规则
	//qDebug()<<"ls"<<ls;
	for(int j=0;j<ls.size();j++)//将其他规则的first集合全部求出
	{
	//qDebug()<<QString(new_grammar[ls[j]]->right[0])<<" "<<get_first(QString(new_grammar[ls[j]]->right[0]), 1, n);
	already.unite(get_first(QString(new_grammar[ls[j]]->right[0]), 1, n));
	}

	if((QString(new_grammar[vnMapindex[vn[k]][i]]->right[0]) <= "z" && QString(new_grammar[vnMapindex[vn[k]][i]]->right[0]) >= "a")
	|| new_grammar[vnMapindex[vn[k]][i]]->right[0] == "@")
	//如果本规则是终结符号则不需要代入，会在提取直接左公因子中进行提取
	{
	continue;
	}
	//qDebug()<<"already"<<already;
	//--------------------------------------------------------------------------------

	//--------------------------查看是否和其他规则的集合有相交---------------------------
	if(get_first(QString(new_grammar[vnMapindex[vn[k]][i]]->right[0]), 1, n).intersects(already))
	//如果和其他规则的first集合相交，意味着需要代入
	{
	QSet<QString> inter = get_first(QString(new_grammar[vnMapindex[vn[k]][i]]->right[0]), 1, n).intersect(already);//相交的部分
	QVector<QString> temp;
	leftFactor2(new_grammar[vnMapindex[vn[k]][i]]->right, inter, 1, n, temp);
	for(int j=0;j<temp.size();j++)
	{
	insertGrammar(new_grammar[vnMapindex[vn[k]][i]]->left, temp[j]);//形成新的规则插入
	}
	myInsert(remove_index, vnMapindex[vn[k]][i]);//把本行放入待删除的集合
	}
	//-------------------------------------------------------------------------------
	}
	}

	remove_qvector_index(remove_index);*/
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
				//todo impl leftFactor
				lf := g.LeftFactor(s)
				ql := strings.Split(lf, "(")
				if ql[0] != "" {
					ql[1] = ql[1][:len(ql[1])-1]
					ql2 := strings.Split(ql[1], "|")
					for _, ql2Item := range ql2 {
						g.InsertGrammar(string(int32(g.i2s)), ql2Item)
					}
					//new_grammar[temp2[0]]->right = ql[0] + QString(QChar(i2s));//更改第一条的映射
					//                    i2s++;//新符号的unicode编码+1
					//
					//                    temp2.remove(0);//除了更改了映射的行，全部放入待删除集最后删除
					//                    for(int k=0;k<temp2.size();k++)
					//                    {
					//                        myInsert(remove_index, temp2[k]);
					//                    }
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
	//todo list is nil?
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
			//todo impl First
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

/**
void showGrammar();//显示文法

    void TextEdit2Vector(QString s);//按行分割

    void simplize(QString& s, int& n);//化简文法

    void splitLeftandRight(int& n);//分开左部和右部

    void cannot_arrive();//不可达

    void cannot_terminal();//不可终止

    void invalid();//去有害/多余规则

    void insertGrammar(QString l, QString r);//增加语法

    void remove_invalid(int i);//删除语法

    void remove_qvector_index(QVector<int> remove_index);//删除一批不合法的行

    QString leftFactor(QVector<QString> s);//返回所有规则s的左公因子供直接左公因子调用

    //返回规则rule的间接左公因子的所有规则can，供间接提取左公因子调用
    void leftFactor2(QString rule, QSet<QString> already, int layer, int& n, QVector<QString>& can);

    void remove_left_factor(int& n);//提取左公因子

    void direct(int& n);//直接提取左公因子

    void indirect(int& n);//间接提取左公因子

    //返回规则rule的间接左递归的所有规则can，供间接消除左递归调用
    void leftRecurse(QString rule, QString unterminal, int layer, int& n, QVector<QString>& can);

    void remove_left_recurse(int& n);//消除左递归

    void direct2(int &n);//直接消除左递归

    void indirect2(int &n);//间接消除左递归

    QMap<QString, QSet<QString>> get_all_first(int& n);//获取所有符号的first集合

    QSet<QString> get_first(QString nt, int layer, int& n);//获取符号ut的first集合

    QSet<QString> First(QString right, int layer, int& n);//获取规则right的first元素

    void Follow(int& n);//计算所有非终结符号的follow集合

    void setTable(int& n);//计算LL(1)分析表

    int searchFirstLine(QString ut, QString element, int& n);

    void reset_vnMapindex();//重置映射表

    void reset_v();//重置所有符号集合

    void duplicate();//去重
*/
