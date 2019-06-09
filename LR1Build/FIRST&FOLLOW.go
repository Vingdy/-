package LR1Build


/*
FOLLOW集创建
先遍历所有字符然后根据情况的不同进行分析判断处理
是产生式中出现的所有字符的FOLLOW集,终结符没有对应FOLLOW集
参数:void
返回:void
*/
func Follow_Build() {
	var main rune
	var next []rune
	for i := 0; i < len(GrammarList); i++ {
		if IsUpper(GrammarList[i].main) {
			main = GrammarList[i].main
		}
		if !IsSameMain(main, Follow) {
			if i == 0 {
				next = append(next, '#')
			} else {
				next = nil
			}
			Follow = append(Follow, GrammarListStruct{
				main: main,
				next: next,
			})
		}
		//终结符没有FOLLOW集不用处理
		for j := 0; j < len(GrammarList[i].next); j++ {
			if !IsUpper(GrammarList[i].next[j]) {
				main = GrammarList[i].next[j]
			} else {
				continue
			}
			if !IsSameMain(main, Follow) {
				Follow = append(Follow, GrammarListStruct{
					main: main,
					next: next,
				})
			}
		}
	}
	for i := 0; i < len(Follow); i++ {
		mainnumber, main := FindSameMain(Follow[i].main, Follow) //main记录是那个Follow集
		if !IsUpper(main) { //如果是FOLLOW集是终结符跳过
			continue
		}
		for j := 0; j < len(GrammarList); j++ {
			Snumber, _ := FindSameMain(GrammarList[j].main, Follow) //记录对应FOLLOW集下标
			for k := 0; k < len(GrammarList[j].next); k++ { //遍历文法的
				if main == GrammarList[j].next[k] { //找到对应字符
					if k == len(GrammarList[j].next)-1 { //是右侧最后一位
						if Follow[Snumber].next != nil { //对应左侧FOLLOW集不为空
							for l := 0; l < len(Follow[Snumber].next); l++ { //把文法产生式左侧的FOLLOW集加入当前FOLLOW集
								if !IsSameNext(Follow[Snumber].next[l], main, Follow) { //判断有无相同
									Follow[mainnumber].next = append(Follow[mainnumber].next, Follow[Snumber].next[l])
								}
							}
						}
					} else if IsUpper(GrammarList[j].next[k+1]) { //如果文法的下一位是非终结符
						for l := 0; l < len(First[Snumber].next); l++ { //把对应的FIRST集加入当前FOLLOW集
							if !IsSameNext(First[Snumber].next[l], main, Follow) { //判断有无相同
								Follow[mainnumber].next = append(Follow[mainnumber].next, First[Snumber].next[l])
							}
						}
					} else { //如果符号的下一位是终结符
						if !IsSameNext(GrammarList[j].next[k+1], main, Follow) { //判断有无相同
							Follow[mainnumber].next = append(Follow[mainnumber].next, GrammarList[j].next[k+1]) //直接加入当前FOLLOW集
						}
					}
				}
			}
		}
	}
}

/*
FIRST集创建递归函数
参数:Nowchar rune(当前搜索字符),mainnumber int(当前First集下标记录),main rune(搜索到的产生式左侧符号)
返回:void
*/
func Dfs(Nowchar rune,mainnumber int,main rune) {
	for i := 0; i < len(GrammarList); i++ {
		if Nowchar == GrammarList[i].main && (GrammarList[i].main != GrammarList[i].next[0]) {
			if !IsUpper(GrammarList[i].next[0]) && !IsSameNext(GrammarList[i].next[0], main, First) {
				First[mainnumber].next = append(First[mainnumber].next, GrammarList[i].next[0])
			}
			if IsUpper(GrammarList[i].next[0]) && !IsSameNext(GrammarList[i].next[0], main, First) {
				Dfs(GrammarList[i].next[0], mainnumber, main)
			}
		}
	}
}

/*
FIRST集创建
先遍历所有字符然后遍历字符递归搜索修改
是产生式中出现的所有字符的FIRST集
参数:void
返回:void
*/
func First_Build() {
	var main rune
	var next []rune
	for i := 0; i < len(GrammarList); i++ {
		if IsUpper(GrammarList[i].main) {
			main = GrammarList[i].main
		}
		if !IsSameMain(main, First) {
			First = append(First, GrammarListStruct{
				main: main,
				next: next,
			})
		}
		//终结符的FIRST集就是它本身
		for j := 0; j < len(GrammarList[i].next); j++ {
			if !IsUpper(GrammarList[i].next[j]) {
				main = GrammarList[i].next[j]
			} else {
				continue
			}
			if !IsSameMain(main, First) {
				First = append(First, GrammarListStruct{
					main: main,
					next: next,
				})
			}
			num, same := FindSameMain(main, First)
			if !IsSameNext(main, same, First) {
				First[num].next = append(First[num].next, same)
			}
		}
	}
	//其余终结符按正常情况进行判断
	for i := 0; i < len(GrammarList); i++ {
		mainnumber, main := FindSameMain(GrammarList[i].main, First)
		if !IsUpper(GrammarList[i].next[0]) && !IsSameNext(GrammarList[i].next[0], main, First) {
			First[mainnumber].next = append(First[mainnumber].next, GrammarList[i].next[0])
			continue
		}
		if IsUpper(GrammarList[i].next[0]) {
			Dfs(GrammarList[i].next[0], mainnumber, main)
		}
	}
}