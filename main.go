package main

import (
	"io/ioutil"
	"fmt"
	"Lexical"
	"log"
	"LR1Build"
	"Grammar"
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
主程序
*/
func main() {
	ReadSourceProgramFile()
	SourceProgram, err := Lexical.PreScan(SourceProgram)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(SourceProgram))
	err=Lexical.Scan(SourceProgram)
	if err != nil {
		log.Panic(err)
	}
	LR1Build.LR1_Build()
	Grammar.Do()
}
