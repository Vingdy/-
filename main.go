package main

import (
	"io/ioutil"
	"fmt"
	"Lexical"
	"log"
	"LR1Build"
	"Grammar"
	"Semanteme"
)

//源程序存放处
var SourceProgram []rune

/*
读取源程序文件
读取保存在SourceProgram.txt中的源代码
参数:void
返回:void
*/
func ReadSourceProgramFile() {
	body, err := ioutil.ReadFile("SourceProgram.txt")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("源程序：")
	fmt.Println(string(body))
	SourceProgram = []rune(string(body))
}

/*
运行主程序
参数:void
返回:void
*/
func main() {
	ReadSourceProgramFile()
	LexicalHandle()
	LR1Build.LR1_Build()
	GrammarHadle()
	SemantemeHandle()
}

/*
词法分析处理
参数:void
返回:error
*/
func LexicalHandle()error{
	SourceProgram, err := Lexical.PreScan(SourceProgram)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(SourceProgram))
	err=Lexical.Scan(SourceProgram)
	if err != nil {
		log.Panic(err)
	}
	return nil
}

/*
语法分析处理
参数:void
返回:error
*/
func GrammarHadle()error {
	Grammar.ReadLR1TableFile()
	Grammar.SetLR1Table()
	err := Grammar.GetLexicalToAnalysis()
	if err != nil {
		log.Panic(err)
	}
	return nil
}

/*
语义分析处理
参数:void
返回:error
*/
func SemantemeHandle()error {
	err := Semanteme.ForestAnalysis()
	if err != nil {
		log.Panic(err)
	}
	return nil
}
