package Grammar

import (
	"fmt"
	"os"
	"bufio"
	"io"
	"strconv"
	"io/ioutil"
	"Lexical"
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

type LR1AGStruct struct{
	AG string
	End string
	ClosureNum int
	Result int
}
var LR1AG []LR1AGStruct

type JudgeStruct struct {
	End string
	ClosureNum int
}
var LR1ActionTable map[JudgeStruct]int
var LR1GotoTable map[JudgeStruct]int

var Status_Stack []int
var Symbol_Stack []string

func Do(){
	ReadLR1TableFile()
	SetLR1Table()
	GetLexicalToAnalysis()
}

func GrammarAnalysis(List []Lexical.LexicalResultStruct) {
	Status_Stack = nil
	Symbol_Stack = nil
	Status_Stack = append(Status_Stack, 0)
	Symbol_Stack = append(Symbol_Stack, "#")
	var ListNum int
	var NumberToSymbol string
	for ListNum=0;;ListNum++{
		fmt.Println(Status_Stack)
		fmt.Println(Symbol_Stack)
		for ; ListNum < len(List);  {
			fmt.Println(List[ListNum])
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
		//fmt.Println(NumberToSymbol)
		var FindAG = JudgeStruct{
			NumberToSymbol,
			Status_Stack[len(Status_Stack)-1],
		}
		fmt.Println(FindAG)
		fmt.Println(LR1ActionTable[FindAG])
		if LR1ActionTable[FindAG] != 0 {
			if LR1ActionTable[FindAG] == 200 {
				fmt.Println("Acc!")
				return
			} else if LR1ActionTable[FindAG] >= 100 {
				var GrammarNum= LR1ActionTable[FindAG] - 100
				fmt.Println(string(GrammarList[GrammarNum].next))
				var GrammarLength= len(GrammarList[GrammarNum].next)
				var ok= true
				for k, v := range GrammarList[GrammarNum].next {
					fmt.Println(len(Symbol_Stack))
					fmt.Println(len(GrammarList[GrammarNum].next))
					fmt.Println((len(Symbol_Stack))-((len(GrammarList[GrammarNum].next))-k))
					if string(v) != Symbol_Stack[len(Symbol_Stack)-(len(GrammarList[GrammarNum].next)-k)] {
						ok = false
					}
				}
				fmt.Println(ok)
				if ok {
					fmt.Println("ACTION规约")
					Status_Stack = append(Status_Stack[:(len(Status_Stack) - GrammarLength)], Status_Stack[len(Status_Stack):]...)
					Symbol_Stack = append(Symbol_Stack[:(len(Symbol_Stack) - GrammarLength)], Symbol_Stack[len(Symbol_Stack):]...)
					Status_Stack = append(Status_Stack, LR1GotoTable[JudgeStruct{
						string(GrammarList[GrammarNum].main),
						Status_Stack[len(Status_Stack)-1],
					}])
					Symbol_Stack = append(Symbol_Stack, string(GrammarList[GrammarNum].main))
				}
			} else {
				fmt.Println("ACTION移入")
				Status_Stack = append(Status_Stack, LR1ActionTable[FindAG])
				Symbol_Stack = append(Symbol_Stack, NumberToSymbol)
			}
		} else {
			fmt.Println("failed")
			return
		}

	}
}

func GetLexicalToAnalysis(){
	for i:=0;i<len(Lexical.LexicalResultList);i++{
		fmt.Println(Lexical.LexicalResultList[i].LexicalList)
		GrammarAnalysis(Lexical.LexicalResultList[i].LexicalList)
	}
}

/*
读取保存在LR1Table.txt中的LR1分析结果
参数:void
返回:void
*/
func ReadLR1TableFile() {
	ReadSourceGrammarFile()
	SaveGrammarList()
	file,err:=os.OpenFile("./LR1Table.txt",os.O_RDWR,0766)
	//body, err := ioutil.ReadFile("./LR1Table.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	//bufReader := bufio.NewReader(file)
	//var lines [][]byte
	rd := bufio.NewReader(file)
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行

		if err != nil || io.EOF == err {
			if line==""{
				break
			}
		}
		//fmt.Println(line)
		var temp []string
		temp=make([]string,4)
		for i,j:=0,0;i< len(line);i++{
			if string(line[i]) == "\n"{
				break
			}
			if string(line[i]) == " " {
				j++
				continue
			}
			temp[j] += string(line[i])
		}
		/*
		for _,v:=range temp{
			fmt.Println(v)
		}*/
		Num,err:=strconv.Atoi(temp[2])
		if err!=nil{
			fmt.Println("Grammar.txt第三列出现内容错误")
			return
		}
		//fmt.Println([]rune(temp[3]))
		Re,err:=strconv.Atoi(temp[3])
		if err!=nil{
			fmt.Println("Grammar.txt第四列出现内容错误")
			return
		}
		LR1AG=append(LR1AG,LR1AGStruct{
			temp[0],
			temp[1],
			Num,
			Re,
		})
	}
}



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
	for k,v:=range LR1ActionTable{
		fmt.Println(k,v)
	}
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
	SourceGrammar = []rune(string(body))
}