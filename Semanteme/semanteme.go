package Semanteme

import (
	"Grammar"
	"github.com/pkg/errors"
	"fmt"
	"Lexical"
)

func DFS(Tree *Grammar.TreeNode)error{


	for _,v:=range Tree.Next{
		//fmt.Println(v)
		//fmt.Println(v.Word)
		DFS(v)
	}
	//fmt.Println(Tree.Word)
	switch len(Tree.Next) {
	case 0:
		return nil
	case 1:
		if Tree.Next[0].Word.Typenumber == 9 {
			Tree.Word.Typenumber = Tree.Next[0].Word.Typenumber
			Tree.Word.Codevalue = Tree.Next[0].Word.Codevalue
			return nil
		}
		if Tree.Next[0].Word.Typenumber == 8 {
			//fmt.Println(Lexical.VariableWords[Tree.Next[0].Word.Codevalue])
			Tree.Word.Typenumber = Tree.Next[0].Word.Typenumber
			Tree.Word.Codevalue=Lexical.VariableWords[Tree.Next[0].Word.Codevalue]
			return nil
		}
	case 2:
		if Tree.Next[0].Word.Typenumber == 2 {
			fmt.Println(Tree.Next[1].Word.Codevalue)
		}
		if Tree.Next[0].Word.Typenumber == 1 {
			Lexical.VariableWords[Tree.Next[0].Word.Codevalue]=Tree.Next[1].Word.Codevalue
			return nil //未处理
		}
	case 3:
		if Tree.Next[0].Word.Typenumber == 20 && Tree.Next[2].Word.Typenumber == 21 {
			Tree.Word.Typenumber = Tree.Next[1].Word.Typenumber
			Tree.Word.Codevalue = Tree.Next[1].Word.Codevalue
			return nil
		}
		switch Tree.Next[1].Word.Typenumber {
		case 10:
			{
				Tree.Word.Codevalue = Tree.Next[0].Word.Codevalue + Tree.Next[2].Word.Codevalue
			}
		case 11:
			{
				Tree.Word.Codevalue = Tree.Next[0].Word.Codevalue - Tree.Next[2].Word.Codevalue
			}
		case 12:
			{
				Tree.Word.Codevalue = Tree.Next[0].Word.Codevalue * Tree.Next[2].Word.Codevalue
			}
		case 13:
			{
				Tree.Word.Codevalue = Tree.Next[0].Word.Codevalue / Tree.Next[2].Word.Codevalue
			}
		case 14:
			{
				Lexical.VariableWords[Tree.Next[0].Word.Codevalue]=Tree.Next[2].Word.Codevalue
			}
		}
		Tree.Word.Typenumber = Tree.Next[0].Word.Typenumber
		return nil
	case 4:
		if Tree.Next[0].Word.Typenumber==1{
			Lexical.VariableWords[Tree.Next[1].Word.Codevalue]=Tree.Next[3].Word.Codevalue
			return nil
		}
	}
	return errors.New("语义分析错误")
}

func ForestAnalysis()error{

	fmt.Println(Grammar.Forest[0].Next[0].Word)
	fmt.Println(Grammar.Forest[0].Next[3].Word)
	for _,Tree:=range Grammar.Forest{
		err:=DFS(&Tree)
		if err!=nil{
			return errors.New("语义分析错误")
		}
	}
	return nil
}
