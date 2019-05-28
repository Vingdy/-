package main

import (
	"io/ioutil"
	"fmt"
	"Lexical"
	"log"
)

//源程序存放处
var proText []rune

/*
读取文件
读取保存在proText.txt中的源代码
参数:void
返回:void
 */
func ReadFile() {
	body, err := ioutil.ReadFile("proText.txt")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("源程序：")
	fmt.Println(string(body))
	proText = []rune(string(body))
}

func main() {
	ReadFile()
	proText, err := Lexical.PreScan(proText)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(proText))
	err=Lexical.Scan(proText)
	if err != nil {
		log.Panic(err)
	}
	//fmt.Println(err)
}
