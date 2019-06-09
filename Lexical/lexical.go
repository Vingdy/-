package Lexical

import (
	"github.com/pkg/errors"
	"strconv"
	"io/ioutil"
	"fmt"
	"Conf"
)

//源程序存放处
var SourceProgram []rune

//保留字
var ReserveWords=map[string]int{
	"var": 1, "print": 2, "if": 3, "else": 4, "while": 5, "true": 6, "false": 7,
}
//界符
var DecollatorWords=map[string]int{
	"{": 22, "}": 23, "(": 20, ")": 21, ";": 24,
}
//操作符
var OperatorWords= map[string]int{
	"+": 10, "-": 11, "*": 12, "/": 13, "=": 14, "!=": 15, "<": 16, ">": 17, "<=": 18, ">=": 19,
}
//变量
var VariableWords=make(map[int]int)//初始化
var VariableTempSave []string//用于检查变量重复命名
var VariableWordNum=0//用于map中的key值
var Variable=8//种别编码
//常量
var Const=9//种别编码
//结果存放
var SaveNumber=0//二元式数量

//文法结果
type LexicalResultStruct struct {
	Character  string //内容
	Codevalue  int    //符号类别
	Typenumber int    //值
}
var LexicalResult []LexicalResultStruct

//文法结果列
type LexicalResultListStruct struct {
	LexicalList []LexicalResultStruct
	ListNum     int
}
var LexicalResultList []LexicalResultListStruct

/*
预扫描
扫描文本并清除其中的注释语句,有错误SourceProgram为nil,否则err为空
参数:SourceProgram []rune
返回:SourceProgram []rune,err error
*/
func PreScan(SourceProgram []rune)([]rune,error) {
	//fmt.Println(SourceProgram)
	for i := 0; i < len(SourceProgram); i++ {
		if i == len(SourceProgram)-1 {
			continue
		} else if SourceProgram[i] == '/' && SourceProgram[i+1] == '/' {
			for ; i < len(SourceProgram); i++ {
				if SourceProgram[i] == '\n' {
					break
				}
				SourceProgram[i] = ' '
			}
		} else if SourceProgram[i] == '/' && SourceProgram[i+1] == '*' {
			for ; i < len(SourceProgram); i++ {
				if SourceProgram[i] == '*' && SourceProgram[i+1] == '/' {
					SourceProgram[i] = ' '
					SourceProgram[i+1] = ' '
					break
				}
				if i >= len(SourceProgram)-1 {
					return nil, errors.New("注释错误")
				}
				SourceProgram[i] = ' '
			}
		}
	}
	return SourceProgram, nil
}

/*
输出词法分析结果
输出词法分析结果并把结果转化为词法分析结果列查看用于编译阶段测试用
参数:void
返回:void
*/
func OutputLexicalResult()error {
	if SaveNumber == 0 {
		return nil
	}
	var ListNum= 0
	LexicalResultList = append(LexicalResultList, LexicalResultListStruct{})
	for Num, Result := range LexicalResult {
		//fmt.Println(Result.Character + "->(" + strconv.Itoa(Result.Codevalue) + "," + strconv.Itoa(Result.Typenumber) + ")")
		LexicalResultList[ListNum].LexicalList = append(LexicalResultList[ListNum].LexicalList, Result)
		if Result.Character == ";" && Num != (len(LexicalResult)-1) { //以;为一列
			ListNum++
			LexicalResultList = append(LexicalResultList, LexicalResultListStruct{})
		}
	}
	return nil
}

/*
扫描
扫描预扫描后的源程序并生成二元式
扫描优先级:空格>界符>空格>英文>数字
参数:SourceProgram []rune
返回:err error
*/
func Scan(SourceProgram []rune)(error) {

	var token = "" //保存当前字符串

	//遍历源程序
	for i := 0; i <= len(SourceProgram); i++ {
		//fmt.Println(token)
		if (i == len(SourceProgram)) || IsSpace(string(SourceProgram[i])) { //判断空格
			if DecollatorWordsHandle(token) {
				token = ""
				continue
			} else if OperatorWordsHandle(token) {
				token = ""
				continue
			}
		} else if DecollatorWordsHandle(token) { //判断界符
			token = ""
			i--
			continue
		} else if OperatorWordsHandle(token) { //判断操作符
			token = ""
			i--
			continue
		} else if IsEnglish(string(SourceProgram[i])) { //判断英文
			for ; ; i++ {
				if IsSpecialCharacter(string(SourceProgram[i])) || IsSpace(string(SourceProgram[i])) { //检测到特殊字符或者空格
					if ReserveWordsHandle(token) { //判断当前字符串是否是关键字
						i--
						token = ""
						break
					} else {
						VariableWordsHandle(token) //跳出并认定为标识符
						i--
						token = ""
						break
					}
				} else if i > len(SourceProgram)-1 { //扫到最后一个仍然为标识符
					//return errors.New("源程序错误，词法分析出错")
					VariableWordsHandle(token) //跳出并认定为标识符
					i--
					token = ""
					break
				} else {
					token += string(SourceProgram[i]) //加入下一个字符
				}
			}
		} else if IsMath(string(SourceProgram[i])) { //判断数字
			for ; ; i++ {
				if IsEnglish(string(SourceProgram[i])) { //出现英语返回上一个进入标识符判断
					i--
					break
				} else if IsSpecialCharacter(string(SourceProgram[i])) || IsSpace(string(SourceProgram[i])) { //判断特殊字符和空格
					ConstWordsHandle(token) //认定为常数处理
					token = ""
					i--
					break
				} else if i > len(SourceProgram)-1 { //最后处理
					//return errors.New("源程序错误，词法分析出错")
					ConstWordsHandle(token)
					token = ""
					i--
					break
				} else {
					token += string(SourceProgram[i]) //加入下一个字符
				}
			}
		} else {
			token += string(SourceProgram[i]) //加入下一个字符
		}
	}
	err := OutputLexicalResult() //输出词法分析结果
	if err != nil {
		return err
	}
	return nil
}

/*
空符号判断
判断当前字符是否是\n,\r,\t, ,四个空符号
参数:token string
返回:ok bool
*/
func IsSpace(token string)bool {
	return token == " " || token == "\r" || token == "\t" || token == "\n"
}

/*
英文字符判断
判断当前字符是否是英文字符
参数:token string
返回:ok bool
*/
func IsEnglish(token string)bool {
	return (token >= "A" && token <= "Z") || (token >= "a" && token <= "z")
}

/*
数字字符判断
判断当前字符是否是数字符号
参数:token string
返回:ok bool
*/
func IsMath(token string)bool {
	return token >= "0" && token <= "9"
}

/*
特殊字符判断
判断当前字符是否是特殊符号
参数:token string
返回:ok bool
*/
func IsSpecialCharacter(token string)bool {
	 for char, _ := range OperatorWords {
		 if token == char {
			 return true
		 }
	 }
	 for char, _:=range DecollatorWords {
		 if token == char {
			 return true
		 }
	 }
	 if token=="!"||token==">"||token=="<"{
	 	return true
	 }
	 return false
 }

/*
保留字处理
判断是否是保留字并保存到结果中
参数:token string
返回:ok bool
*/
func ReserveWordsHandle(token string)(bool) {
	for char, Typenumber := range ReserveWords {
		if token == char {
			Result := LexicalResultStruct{
				Character:  char,
				Codevalue:  0,
				Typenumber: Typenumber,
			}
			LexicalResult = append(LexicalResult, Result)
			SaveNumber++
			return true
		}
	}
	return false
}

/*
界符处理
判断是否是界符并保存到结果中
参数:token string
返回:ok bool
*/
func DecollatorWordsHandle(token string)(bool) {
	for char, Typenumber := range DecollatorWords {
		if token == char {
			Result := LexicalResultStruct{
				Character:  char,
				Codevalue:  0,
				Typenumber: Typenumber,
			}
			LexicalResult = append(LexicalResult, Result)
			SaveNumber++
			return true
		}
	}
	return false
}

/*
操作符处理
判断是否是操作符并保存到结果中
参数:token string
返回:ok bool
*/
func OperatorWordsHandle(token string)(bool) {
	for char, Typenumber := range OperatorWords {
		if token == char {
			Result := LexicalResultStruct{
				Character:  char,
				Codevalue:  0,
				Typenumber: Typenumber,
			}
			LexicalResult = append(LexicalResult, Result)
			SaveNumber++
			return true
		}
	}
	return false
}

/*
变量处理
判断是否是变量并保存到结果中
参数:token string
返回:void
*/
func VariableWordsHandle(token string) {
	VariableWords[VariableWordNum] = VariableWordNum
	Result := LexicalResultStruct{
		Character:  token,
		Codevalue:  VariableWordNum,
		Typenumber: Variable,
	}
	VariableWordNum++
	VariableTempSave = append(VariableTempSave, token)
	for _, v := range VariableTempSave {
		if v == token {
			VariableWordNum--
			VariableTempSave = append(VariableTempSave[:(len(VariableTempSave) - 1)], VariableTempSave[len(VariableTempSave):]...)
		}
	}
	LexicalResult = append(LexicalResult, Result)
	SaveNumber++
}

/*
常量处理
判断是否是常量并保存到结果中
参数:token string
返回:void
*/
func ConstWordsHandle(token string) {
	tokennum,_:=strconv.Atoi(token)
	Result := LexicalResultStruct{
		Character:  token,
		Codevalue:  tokennum,
		Typenumber: Const,
	}
	LexicalResult = append(LexicalResult, Result)
	SaveNumber++
}

/*
读取源程序文件
读取保存在SourceProgram.txt中的源代码
参数:void
返回:void
*/
func ReadSourceProgramFile(conf Conf.ConfSturct) {
	body, err := ioutil.ReadFile(conf.ProjectPath+conf.ProgramFile)
	if err != nil {
		fmt.Println(err)
	}
	SourceProgram = []rune(string(body))
}
