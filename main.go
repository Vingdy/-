package main

import (
	"Lexical"
	"log"
	"LR1Build"
	"Grammar"
	"Semanteme"
	"Conf"
	"os"
	"fmt"
	"time"
)

//配置
var conf Conf.ConfSturct

/*
运行主程序
参数:void
返回:void
*/
func main() {
	conf.GetConf()
	//fmt.Println(conf)
	LexicalHandle(conf)
	if conf.IsUseLR1Build {
		file, err := os.OpenFile(conf.ProjectPath+conf.LR1TableFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0766) //O_TRUNC:如果可能,在打开文件时清空文件
		if err != nil {
			fmt.Println(err)
		}
		file.Close()
		LR1Build.LR1_Build()
	}
	GrammarHandle(conf)
	SemantemeHandle()
	time.Sleep(10 * time.Second)
}

/*
词法分析处理
参数:void
返回:error
*/
func LexicalHandle(conf Conf.ConfSturct)error{
	Lexical.ReadSourceProgramFile(conf)
	SourceProgram, err := Lexical.PreScan(Lexical.SourceProgram)
	if err != nil {
		log.Panic(err)
		panic(err)
	}
	//fmt.Println(string(SourceProgram))
	err=Lexical.Scan(SourceProgram)
	if err != nil {
		log.Panic(err)
		panic(err)
	}
	return nil
}

/*
语法分析处理
参数:void
返回:error
*/
func GrammarHandle(conf Conf.ConfSturct)error {
	Grammar.ReadLR1TableFile(conf)
	Grammar.SetLR1Table()
	err := Grammar.GetLexicalToAnalysis()
	if err != nil {
		log.Panic(err)
		panic(err)
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
		panic(err)
	}
	return nil
}