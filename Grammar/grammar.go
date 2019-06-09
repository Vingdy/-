package Grammar

import (
	"fmt"
	"os"
	"bufio"
	"io"
	"strconv"
	"io/ioutil"
	"Lexical"
	"github.com/pkg/errors"
	"unicode"
	"Conf"
)

//原文法存放处
var SourceGrammar []rune

//文法结构体
type GrammarListStruct struct {
	main rune
	next []rune
}
//原文法序列
var GrammarList []GrammarListStruct

//LR(1)表结果
type LR1AGStruct struct{
	AG string
	End string
	ClosureNum int
	Result int
}
var LR1AG []LR1AGStruct

//LR(1)分析判断结构
type JudgeStruct struct {
	End string
	ClosureNum int
}
//LR(1)表
var LR1ActionTable map[JudgeStruct]int
var LR1GotoTable map[JudgeStruct]int

//语法分析栈
var Status_Stack []int//状态栈
var Symbol_Stack []string//符号栈
var Word_Stack []Lexical.LexicalResultStruct//词法栈

//树节点结构
type TreeNode struct {
	Word   Lexical.LexicalResultStruct
	Next   []*TreeNode
	Parent *TreeNode
}
//语法森林
var Forest []TreeNode

/*
把生成的语法分析树放入森林
参数:NewTree TreeNode
返回:void
*/
func TreeToForest(NewTree TreeNode) {
	Forest = append(Forest, NewTree)
	/*for k, v := range Forest {
		fmt.Println(k, v)
	}*/
}

/*
生成语法分析树
根据当前语法分析规约的结果并生成对应语法树,返回新生成的结点与修改后的树节点数组
生成顺序是由叶子->根
参数:GrammarNum int(文法下标),TreeSave []*TreeNode(保存的树节点数组)
返回:*TreeNode(新生成的树节点),[]*TreeNode
*/
func GrammarTree(GrammarNum int,TreeSave []*TreeNode)(*TreeNode,[]*TreeNode) {
	//fmt.Println(TreeSave)
	var UpperNum= 0
	var NewNode= TreeNode{}
	//大写符号记录
	for k, _ := range GrammarList[GrammarNum].next {
		if unicode.IsUpper(GrammarList[GrammarNum].next[k]) {
			UpperNum++
		}
	}
	//开始生成
	for k, _ := range GrammarList[GrammarNum].next {
		//fmt.Println(unicode.IsUpper(GrammarList[GrammarNum].next[k]))
		/*
		fmt.Println(UpperNum)
		fmt.Println(len(TreeSave))
		*/
		if unicode.IsUpper(GrammarList[GrammarNum].next[k]) {//大写符号把对应树节点进行连接
			TreeSave[len(TreeSave)-UpperNum].Parent = &NewNode
			NewNode.Next = append(NewNode.Next, TreeSave[len(TreeSave)-UpperNum])
			TreeSave = append(TreeSave[:(len(TreeSave) - UpperNum)], TreeSave[len(TreeSave)-UpperNum+1:]...)
			UpperNum--
		} else {//小写就直接新生成一个叶子结点进行连接
			var NewTree = TreeNode{}
			NewTree.Word = Word_Stack[len(Word_Stack)-(len(GrammarList[GrammarNum].next)-k)]
			NewTree.Parent = &NewNode
			NewNode.Next = append(NewNode.Next, &NewTree)
		}
	}
	//词法栈弹出
	Word_Stack = append(Word_Stack[:(len(Word_Stack) - len(GrammarList[GrammarNum].next))], Word_Stack[len(Word_Stack):]...)
	//词法栈移入
	NewWord := Lexical.LexicalResultStruct{}
	NewWord.Character = string(GrammarList[GrammarNum].main)
	NewNode.Word = NewWord
	Word_Stack = append(Word_Stack, NewWord)

	return &NewNode, TreeSave
}

/*
进行语法分析
对一行词法结果列的形式进行语法分析并生成对应语法树放入森林
参数:List []Lexical.LexicalResultStruct(当前语法分析列)
返回:error
*/
func GrammarAnalysis(List []Lexical.LexicalResultStruct)error {
	//初始化全局栈
	Status_Stack = nil
	Symbol_Stack = nil
	Word_Stack = nil
	Status_Stack = append(Status_Stack, 0)
	Symbol_Stack = append(Symbol_Stack, "#")
	Word_Stack = append(Word_Stack, Lexical.LexicalResultStruct{
		"S",
		0,
		0,
	})
	var ListNum int//文法列数组下标
	var NumberToSymbol string//种别编码转符号保存
	var TreeSave []*TreeNode//文法列树节点保存
	for ListNum = 0; ; ListNum++ {
		/*
		fmt.Println(Status_Stack)
		fmt.Println(Symbol_Stack)
		fmt.Println(Word_Stack)
		fmt.Printf("%p", TreeSave)
		fmt.Println(TreeSave)
		*/
		for ; ListNum < len(List); {
			//fmt.Println(List[ListNum])
			//根据种别编码转化为对应文法符号
			switch List[ListNum].Typenumber > 0 {
			case List[ListNum].Typenumber == 1:
				NumberToSymbol = "v"
			case List[ListNum].Typenumber == 2:
				NumberToSymbol = "p"
			case List[ListNum].Typenumber == 8:
				NumberToSymbol = "x"
			case List[ListNum].Typenumber == 9:
				NumberToSymbol = "c"
			case (List[ListNum].Typenumber >= 10) && (List[ListNum].Typenumber <= 23):
				{
					NumberToSymbol = List[ListNum].Character
				}
			case List[ListNum].Typenumber == 24:
				NumberToSymbol = "#"
			}
			break
		}
		//生成对应判断结构
		var FindAG = JudgeStruct{
			NumberToSymbol,
			Status_Stack[len(Status_Stack)-1],
		}
		//分析开始(0->错误,200->Acc,x(100~199)->r(x-100),x(1~99)->sx)
		if LR1ActionTable[FindAG] != 0 {//正确
			if LR1ActionTable[FindAG] == 200 {//Acc
				//fmt.Println("Acc!")
				TreeToForest(*TreeSave[0])
				return nil
			} else if LR1ActionTable[FindAG] >= 100 {//规约
				var GrammarNum = LR1ActionTable[FindAG] - 100//找出文法列
				var GrammarLength = len(GrammarList[GrammarNum].next)//判断文法列是否相同
				var ok = true
				for k, v := range GrammarList[GrammarNum].next {
					if string(v) != Symbol_Stack[len(Symbol_Stack)-(len(GrammarList[GrammarNum].next)-k)] {
						ok = false
					}
				}
				if ok {
					//fmt.Println("ACTION规约")
					//栈弹出
					Status_Stack = append(Status_Stack[:(len(Status_Stack) - GrammarLength)], Status_Stack[len(Status_Stack):]...)
					Symbol_Stack = append(Symbol_Stack[:(len(Symbol_Stack) - GrammarLength)], Symbol_Stack[len(Symbol_Stack):]...)
					//GOTO搜索
					/*
					fmt.Print("GOTO")
					fmt.Println(LR1GotoTable[JudgeStruct{
						string(GrammarList[GrammarNum].main),
						Status_Stack[len(Status_Stack)-1],
					}])
					*/
					//新状态与规约符号移入栈
					Status_Stack = append(Status_Stack, LR1GotoTable[JudgeStruct{
						string(GrammarList[GrammarNum].main),
						Status_Stack[len(Status_Stack)-1],
					}])
					Symbol_Stack = append(Symbol_Stack, string(GrammarList[GrammarNum].main))
					//树节点生成
					var NewTreeNode *TreeNode
					NewTreeNode, TreeSave = GrammarTree(GrammarNum, TreeSave)
					TreeSave = append(TreeSave, NewTreeNode)
					//阻止词法结果列移动
					ListNum--
				}
			} else {
				//普通ACTION移入操作
				//fmt.Println("ACTION移入")
				Status_Stack = append(Status_Stack, LR1ActionTable[FindAG])
				Symbol_Stack = append(Symbol_Stack, NumberToSymbol)
				Word_Stack = append(Word_Stack, List[ListNum])
			}
		} else {
			//错误
			//fmt.Println("failed")
			return errors.New("语法分析错误")
		}
	}
}

/*
语法分析准备
以词法结果列的形式逐个进行语法分析
参数:void
返回:error
*/
func GetLexicalToAnalysis()error {
	for i := 0; i < len(Lexical.LexicalResultList); i++ {
		if (len(Lexical.LexicalResultList[i].LexicalList) == 1) && (Lexical.LexicalResultList[i].LexicalList[0].Typenumber == 24) {
			continue
		}
		//fmt.Println(Lexical.LexicalResultList[i].LexicalList)
		err := GrammarAnalysis(Lexical.LexicalResultList[i].LexicalList)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
读取保存在LR1Table.txt中的LR1分析结果
参数:void
返回:error
*/
func ReadLR1TableFile(conf Conf.ConfSturct)error {
	err:=ReadSourceGrammarFile(conf)
	if err!=nil{
		return err
	}
	SaveGrammarList()
	file, err := os.OpenFile(conf.ProjectPath+conf.LR1TableFile, os.O_RDWR, 0766)
	if err != nil {
		return err
	}
	defer file.Close()
	rd := bufio.NewReader(file)
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行
		if err != nil || io.EOF == err {
			if line == "" {
				break
			}
		}
		var temp []string
		temp = make([]string, 4)
		for i, j := 0, 0; i < len(line); i++ {
			if string(line[i]) == "\n" {
				break
			}
			if string(line[i]) == " " {
				j++
				continue
			}
			temp[j] += string(line[i])
		}
		Num, err := strconv.Atoi(temp[2])
		if err != nil {
			fmt.Println("Grammar.txt第三列出现内容错误")
			return err
		}
		//fmt.Println([]rune(temp[3]))
		Re, err := strconv.Atoi(temp[3])
		if err != nil {
			fmt.Println("Grammar.txt第四列出现内容错误")
			return err
		}
		LR1AG = append(LR1AG, LR1AGStruct{
			temp[0],
			temp[1],
			Num,
			Re,
		})
	}
	return nil
}

/*
把LR(1)结果转化LR(1)表
根据AG值把对应结果分别存入两个map中,用结构体来形成string+int的key值找到对应的value
参数:void
返回:void
*/
func SetLR1Table() {
	LR1ActionTable = make(map[JudgeStruct]int)
	LR1GotoTable = make(map[JudgeStruct]int)
	for i := 0; i < len(LR1AG); i++ {
		var judge = JudgeStruct{
			LR1AG[i].End,
			LR1AG[i].ClosureNum,
		}
		if LR1AG[i].AG == "A" {
			LR1ActionTable[judge] = LR1AG[i].Result
		} else {
			LR1GotoTable[judge] = LR1AG[i].Result
		}
	}
	/*
	for k, v := range LR1ActionTable {
		fmt.Println(k, v)
	}*/
}

/*
保存文法序列
保存文法并把每一个左右两值分别存储进rune和[]rune结构的GrammarListStruct结构体数组中
参数:void
返回:void
*/
func SaveGrammarList() {
	var ListNumber= 0
	for i := 0; i < len(SourceGrammar); i++ {
		var main rune
		var next []rune
		if SourceGrammar[i] == '-' && SourceGrammar[i+1] == '>' {
			main = SourceGrammar[i-1]
			i += 2
			for {
				if SourceGrammar[i] == '\r' || SourceGrammar[i] == '\n' || i >= len(SourceGrammar)-1 {
					if i >= len(SourceGrammar)-1 {
						next = append(next, SourceGrammar[i])
					}
					GrammarList = append(GrammarList, GrammarListStruct{
						main: main,
						next: next,
					})
					ListNumber++
					i++
					break
				}
				next = append(next, SourceGrammar[i])
				i++
			}
		}
	}
}

/*
读取源文法文件
读取保存在SourceProgram.txt中的源代码
参数:void
返回:void
*/
func ReadSourceGrammarFile(conf Conf.ConfSturct)error {
	body, err := ioutil.ReadFile(conf.ProjectPath+conf.GrammarFile)
	if err != nil {
		return err
	}
	SourceGrammar = []rune(string(body))
	return nil
}