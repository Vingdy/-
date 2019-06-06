package Grammar

import (
	"fmt"
	"os"
	"bufio"
	"io"
	"strconv"
)
type LR1AGStruct struct{
	AG string
	End string
	ClosureNum int
	Result int
}
var LR1AG []LR1AGStruct

var LR1Table map[a]int

func Do(){
	ReadLR1TableFile()
}

/*
读取保存在LR1Table.txt中的LR1分析结果
参数:void
返回:void
*/
func ReadLR1TableFile() {
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

type a struct {
	a string
	b int
}

func SetLR1Table(){
	for i:=0;i<len(LR1AG);i++{
		LR1Table[a{
			LR1AG[i].End,
			LR1AG[i].ClosureNum,
		}]=LR1AG[i].Result
	}
}