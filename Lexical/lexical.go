package Lexical

import (
	"github.com/pkg/errors"
	"fmt"
	"strconv"
)

//保留字
var ReserveWords=map[string]int{
	"if": 8, "else": 9, "while": 10, "var": 11, "true": 20, "false": 21,
}
//界符
var DecollatorWords=map[string]int{
	"{": 5, "}": 6, "(": 15, ")": 16, ";": 17,
}
//操作符
var OperatorWords= map[string]int{
	"+": 1, "-": 2, "*": 3, "/": 4, "=": 7, "!=": 12, ">": 13, "<": 14, ">=": 18, "<=": 19,
}
//变量
var VariableWords []string
var Variable=20
//常量
var Const=21
//结果存放
var SaveNumber=0//二元式数量
type LexicalResultStruct struct{
	character      string//内容
	classification string//符号类别
	value          int//值
}
var LexicalResult []LexicalResultStruct

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
	return SourceProgram,nil
}

/*
输出词法分析结果
输出词法分析结果查看用于编译阶段测试用
参数:void
返回:void
*/
func OutputLexicalResult() {
	if (SaveNumber == 0) {
		return
	}
	for _, Result := range LexicalResult {
		fmt.Println(Result.character + "->(" + Result.classification + "," + strconv.Itoa(Result.value) + ")")
	}
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

	for i := 0; i < len(SourceProgram); i++ {
		if IsSpace(string(SourceProgram[i])) {
			token = ""
			continue
		} else if DecollatorWordsHandle(token) {
			token = ""
			i--
			continue
		} else if OperatorWordsHandle(token) {
			token = ""
			i--
			continue
		} else if IsEnglish(string(SourceProgram[i])) {
			for ; ; i++ {
				if ReserveWordsHandle(token) {
					token = ""
					break
				} else if IsSpecialCharacter(string(SourceProgram[i])) || IsSpace(string(SourceProgram[i])) {
					VariableWordsHandle(token)
					i--
					token = ""
					break
				} else if i > len(SourceProgram)-1 {
					return errors.New("词法分析扫描错误")
				} else {
					token += string(SourceProgram[i])
				}
			}
		} else if IsMath(string(SourceProgram[i])) {
			for ; ; i++ {
				if IsEnglish(string(SourceProgram[i])) {
					i--
					break
				} else if IsSpecialCharacter(string(SourceProgram[i])) || IsSpace(string(SourceProgram[i])) {
					ConstWordsHandle(token)
					token = ""
					i--
					break
				} else if i > len(SourceProgram) {
					return errors.New("词法分析扫面错误")
				} else {
					token += string(SourceProgram[i])
				}
			}
		} else {
			token += string(SourceProgram[i])
		}
	}
	OutputLexicalResult()
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
	for char, value := range ReserveWords {
		if (token == char) {
			//fmt.Println(1)
			Result := LexicalResultStruct{
				character:      char,
				classification: "保留字",
				value:          value,
			}
			LexicalResult = append(LexicalResult, Result)
			SaveNumber++
			//fmt.Println(LexicalResult)
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
	for char, value := range DecollatorWords {
		if (token == char) {
			//fmt.Println(1)
			Result := LexicalResultStruct{
				character:      char,
				classification: "界符",
				value:          value,
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
	for char, value := range OperatorWords {
		if (token == char) {
			//fmt.Println(1)
			Result := LexicalResultStruct{
				character:      char,
				classification: "操作符",
				value:          value,
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
	Result := LexicalResultStruct{
		character:      token,
		classification: "变量",
		value:          Variable,
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
	Result := LexicalResultStruct{
		character:      token,
		classification: "常量",
		value:          Const,
	}
	LexicalResult = append(LexicalResult, Result)
	SaveNumber++
}