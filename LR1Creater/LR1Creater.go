package LR1Creater

import (
	"io/ioutil"
	"fmt"
)

//原文法存放处
var SourceGrammar []rune

//var First map[[]rune]rune

type GrammarListStruct struct {
	S rune
	next []rune
}
var GrammarList []GrammarListStruct


type GroupStruct struct {
	main rune
	group []rune
}
var First []GroupStruct
var Follow []GroupStruct

//var First map[rune]rune

/*
LR1分析表创建主函数
*/
func LR1Create(){
	ReadSourceGrammarFile()
	SaveGrammarList()
	First_Create()
}

func IsUpper(a rune)bool{
	return a>='A'&&a<='Z'
}

/*
dfs
*/
func Dfs(Nowchar rune,mainnumber int,main rune,Group []GroupStruct){
	fmt.Println(string(Nowchar))
	for i:=0;i<len(GrammarList);i++{
		if Nowchar==GrammarList[i].S&&GrammarList[i].S!=GrammarList[i].next[0]{
			if !IsUpper(GrammarList[i].next[0])&&!IsSameNext(GrammarList[i].next[0],main,First){
				First[mainnumber].group=append(First[mainnumber].group,GrammarList[i].next[0])
			}
			if IsUpper(GrammarList[i].next[0])&&!IsSameNext(GrammarList[i].next[0],main,First){
				Dfs(GrammarList[i].next[0],mainnumber,main,Group)
			}
		}
	}
}

/*
FIRST集创建
*/
func First_Create() {
	var main rune
	var group []rune
	for i := 0; i < len(GrammarList); i++ {
		main = GrammarList[i].S
		if !IsSameMain(main,First) {
			First = append(First, GroupStruct{
				main:  main,
				group: group,
			})
		}
	}

	for k, v := range First {
		fmt.Println(k, v)
	}

	for i:=0;i<len(GrammarList);i++{
		mainnumber,main:=FindSameMain(GrammarList[i].S,First)
		if !IsUpper(GrammarList[i].next[0])&&!IsSameNext(GrammarList[i].next[0],main,First){
			First[mainnumber].group=append(First[mainnumber].group,GrammarList[i].next[0])
			continue
		}
		if IsUpper(GrammarList[i].next[0]){
			Dfs(GrammarList[i].next[0],mainnumber,main,First)
		}
	}
	for k, v := range First {
		fmt.Println(k, v)
	}
}

/*
判断FIRST是否有相同的左值
*/
func IsSameMain(main rune,Group []GroupStruct)bool {
	for _, v := range Group {
		if v.main == main{
			return true
		}
	}
	return false
}

/*
判断FIRST是否有相同的左值
*/
func FindSameMain(find rune,Group []GroupStruct)(int,rune) {
	for k, v := range Group {
		if v.main == find{
			return k,v.main
		}
	}
	return -1,' '
}

/*
判断FIRST是否有相同的子集
*/
func IsSameNext(next rune,main rune,Group []GroupStruct)bool {
	for _, v := range Group {
		if v.main == main {
			for _, vv := range v.group {
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
		var S rune
		var next []rune
		if SourceGrammar[i] == '-' && SourceGrammar[i+1] == '>' {
			S = SourceGrammar[i-1]
			i += 2
			for {
				if SourceGrammar[i] == '\r' || SourceGrammar[i] == '\n' || i >= len(SourceGrammar)-1 {
					if i >= len(SourceGrammar)-1 {
						next = append(next, SourceGrammar[i])
					}
					GrammarList = append(GrammarList, GrammarListStruct{
						S:    S,
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
	body, err := ioutil.ReadFile("LR1Creater/SourceGrammar.txt")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("源文法：")
	fmt.Println(string(body))
	SourceGrammar = []rune(string(body))
}