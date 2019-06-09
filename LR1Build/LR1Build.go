package LR1Build

import (
	"io/ioutil"
	"fmt"
	"strconv"
	"os"
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
//FIRST集&FOLLOW集
var First []GrammarListStruct
var Follow []GrammarListStruct

//Closure结构体
type ClosureStruct struct {
	G        GrammarListStruct
	Node     int
	End      rune
	NextUnit int
	NextRune rune
}
//ClosureUnit结构体
type ClosureUnitStruct struct {
	Closure  []ClosureStruct
	NextUnit []int
	NextRune []rune
}
//一个ClosureUnit=一个DFA
var ClosureUnit []ClosureUnitStruct

//保存LR1结果结构
type LR1AGStruct struct {
	AG         rune
	End        rune
	ClosureNum int
	Result     int
}
var LR1AG []LR1AGStruct

/*
LR1分析表创建主函数
参数:void
返回:void
*/
func LR1_Build(){
	ReadSourceGrammarFile()
	SaveGrammarList()
	First_Build()
	Follow_Build()
	DFA_Build()
	Table_Build()
	Output()
}

/*
结果输出
FIRST集->FOLLOW集->CLOSURE集->DFA->LR1AG->输出结果到文件
参数:void
返回:void
*/
func Output() {
	for _, v := range First {
		fmt.Print("FIRST(" + string(v.main) + "):")
		for _, vv := range v.next {
			fmt.Print(string(vv) + " ")
		}
		fmt.Println()
	}
	for _, v := range Follow {
		fmt.Print("FOLLOW(" + string(v.main) + "):")
		for _, vv := range v.next {
			fmt.Print(string(vv) + " ")
		}
		fmt.Println()
	}
	fmt.Println("CLOSURE:")
	for k, v := range ClosureUnit {
		//fmt.Print("FOLLOW("+string(v.main)+"):")
		fmt.Println("I" + strconv.Itoa(k) + ":")
		for _, vv := range v.Closure {
			fmt.Print(string(vv.G.main) + "->")
			var temp1 []rune
			temp1 = append(temp1, vv.G.next...) //直接等于由于底层共用数组导致得不出正确结果
			var temp2 []rune
			temp2 = append(temp2, vv.G.next[vv.Node:]...)
			temp1 = append(temp1[0:vv.Node], '.')
			temp1 = append(temp1, temp2...)
			vv.G.next = temp1
			fmt.Print(string(vv.G.next) + "," + string(vv.End) + ",Next:" + string(vv.NextRune) + strconv.Itoa(vv.NextUnit))
			fmt.Println()
		}
	}
	fmt.Println("DFA:")
	for k, v := range ClosureUnit {
		fmt.Println("I" + strconv.Itoa(k) + ":")
		for k, vv := range v.NextUnit {
			fmt.Println(strconv.Itoa(vv) + " " + string(v.NextRune[k]))
		}
	}
	file, err := os.OpenFile("./LR1Table.txt", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0766) //O_TRUNC:如果可能,在打开文件时清空文件
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	//去重复操作
	var FmtLR1AG []LR1AGStruct
	for _, v := range LR1AG {
		if FmtLR1AG == nil {
			FmtLR1AG = append(FmtLR1AG, v)
		} else {
			HaveSame := false
			for _, vv := range FmtLR1AG {
				if v == vv {
					HaveSame = true
					break
				}
			}
			if !HaveSame {
				FmtLR1AG = append(FmtLR1AG, v)
			}
		}
	}
	fmt.Println("LR1AG:")
	for _, v := range FmtLR1AG {
		if v.AG == 'A' {
			fmt.Println("ACTION(" + string(v.End) + "," + strconv.Itoa(v.ClosureNum) + ")=" + strconv.Itoa(v.Result))
		} else {
			fmt.Println("GOTO(" + string(v.End) + "," + strconv.Itoa(v.ClosureNum) + ")=" + strconv.Itoa(v.Result))
		}
		str := string(v.AG) + " " + string(v.End) + " " + strconv.Itoa(v.ClosureNum) + " " + strconv.Itoa(v.Result) + "\n"
		file.WriteString(str)
	}
}

/*
判断当前字符是否是大小写
因为是否是终结符根据大小写判断,大写是非终结符,小写是终结符
常用函数抽出
参数:a rune(当前字符)
返回:bool(大写返回true,false)
*/
func IsUpper(a rune)bool {
	return a >= 'A' && a <= 'Z'
}

/*
判断文法序列中的产生式左侧是否有对应的字符
参数:main rune(要找的产生式左侧字符),GrammarList []GrammarListStruct(要搜索的文法序列)
返回:bool,找到返回true,否则返回false
*/
func IsSameMain(main rune,GrammarList []GrammarListStruct)bool {
	for _, v := range GrammarList {
		if v.main == main {
			return true
		}
	}
	return false
}

/*
找到文法序列中的产生式左侧的字符及对应序号
参数:find rune(要找的产生式左侧字符),GrammarList []GrammarListStruct(要搜索的文法序列)
返回:int (字符在数组中的下标),rune (字符本身的值)
*/
func FindSameMain(find rune,GrammarList []GrammarListStruct)(int,rune) {
	for k, v := range GrammarList {
		if v.main == find {
			return k, v.main
		}
	}
	return -1, ' '
}

/*
判断文法序列中某个文法的产生式右侧是否有对应的字符
参数:next rune(要找产生式右侧的字符),main rune(对应的产生式左侧字符),GrammarList []GrammarListStruct(要搜索的文法序列)
返回:bool,找到返回true,否则返回false
*/
func IsSameNext(next rune,main rune,GrammarList []GrammarListStruct)bool {
	for _, v := range GrammarList {
		if v.main == main {
			for _, vv := range v.next {
				if vv == next {
					return true
				}
			}
		}
	}
	return false
}

/*
保存文法序列
保存文法并把每一个左右两值分别存储进rune和[]rune结构的GrammarListStruct结构体数组中
参数:void
返回:void
*/
func SaveGrammarList() {
	var ListNumber = 0
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
	for k, v := range GrammarList {
		fmt.Println(k, v)
	}
}

/*
读取源文法文件
读取保存在SourceProgram.txt中的源代码
参数:void
返回:void
*/
func ReadSourceGrammarFile() {
	body, err := ioutil.ReadFile("LR1Build/SourceGrammar.txt")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("源文法：")
	fmt.Println(string(body))
	SourceGrammar = []rune(string(body))
}