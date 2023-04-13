package grammar

import (
	"errors"
	"fmt"
	"strconv"
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
}

func MyInsert(vn *[]string, s string) {
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
	i := 1
	new_i := 0
	for new_i < length {
		ql := strings.Split(g.OriginGrammar[i], "->")
		if strings.Index(ql[1], "(") != -1 {

		}
		rarray := strings.Split(ql[1], "|")
		for i := 0; i < len(rarray); i++ {
			for _, item := range rarray[i] {
				itemS := strconv.Itoa(int(item))
				if item <= 'z' && item >= 'a' || item == '@' {
					g.Vt[itemS] = struct{}{}
				} else {
					g.Vr[itemS] = struct{}{}
				}
			}
		}
		//-------------------所有序列加入文法------------------------
		gNode := newNode(ql[0], rarray[0])
		g.NewGrammar = append(g.NewGrammar, gNode)
		ints := g.VnMapindex[ql[0]]
		*ints = append(*ints, new_i)
		for i := 1; i < len(rarray); i++ {
			gNode := newNode(ql[0], rarray[i])
			g.NewGrammar = append(g.NewGrammar, gNode)
			new_i++
			length++
			ints := g.VnMapindex[ql[0]]
			*ints = append(*ints, new_i)
		}
		new_i++
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
	for _, item := range g.NewGrammar {
		fmt.Printf("%v %v", item.Left, item.Right)
	}
}

//todo duplicate

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
