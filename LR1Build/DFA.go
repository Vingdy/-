package LR1Build

/*
LR1表构造
参数:void
返回:void
*/
func Table_Build() {
	for i := 0; i < len(ClosureUnit); i++ { //遍历项目集
		for j := 0; j < len(ClosureUnit[i].Closure); j++ {
			if ClosureUnit[i].Closure[j].Node == len(ClosureUnit[i].Closure[j].G.next) {
				if ClosureUnit[i].Closure[j].G.main == 'X' && ClosureUnit[i].Closure[j].End == '#' { //如果是X->S.,#
					LR1AG = append(LR1AG, LR1AGStruct{
						AG:         'A',
						End:        '#',
						ClosureNum: i,
						Result:     200,
					})
				} else { //规约判断
					var SameNum int //找到是第几个产生式
					for m := 0; m < len(GrammarList); m++ {
						flag := true//去重判断
						for n := 0; n < len(GrammarList[m].next); n++ {
							if (len(ClosureUnit[i].Closure[j].G.next) != len(GrammarList[m].next)) || (ClosureUnit[i].Closure[j].G.next[n] != GrammarList[m].next[n]) {
								flag = false
								break
							}
						}
						if flag && (ClosureUnit[i].Closure[j].G.main == GrammarList[m].main) { //找到了跳出循环
							SameNum = m
							break
						}
					}
					LR1AG = append(LR1AG, LR1AGStruct{
						AG:         'A',
						End:        ClosureUnit[i].Closure[j].End,
						ClosureNum: i,
						Result:     100 + SameNum,
					})
				}
			} else if IsUpper(ClosureUnit[i].Closure[j].G.next[ClosureUnit[i].Closure[j].Node]) { //非终结符
				LR1AG = append(LR1AG, LR1AGStruct{
					AG:         'G',
					End:        ClosureUnit[i].Closure[j].G.next[ClosureUnit[i].Closure[j].Node],
					ClosureNum: i,
					Result:     ClosureUnit[i].Closure[j].NextUnit,
				})
			} else if !IsUpper(ClosureUnit[i].Closure[j].G.next[ClosureUnit[i].Closure[j].Node]) { //终结符
				LR1AG = append(LR1AG, LR1AGStruct{
					AG:         'A',
					End:        ClosureUnit[i].Closure[j].G.next[ClosureUnit[i].Closure[j].Node],
					ClosureNum: i,
					Result:     ClosureUnit[i].Closure[j].NextUnit,
				})
			}
		}
	}
}

/*
Closure项目集初始化
构造第一个ClosureUnit和第一个Closure便于后面进行DFA_Build
参数:void
返回:void
*/
func ClosureInit() {
	ClosureUnit = append(ClosureUnit, ClosureUnitStruct{})
	ClosureUnit[0].Closure = append(ClosureUnit[0].Closure, ClosureStruct{
		G:    GrammarList[0],
		Node: 0,
		End:  '#',
	})
}
//错误留档
/*
for i := 0; i < len(GrammarList); i++ {
	for j := 0; j <= len(GrammarList[i].next); j++ {
		var ClosureNewNext []rune
		ClosureNewNext = append(ClosureNewNext, GrammarList[i].next...) //直接等于由于底层共用数组导致得不出正确结果
		var temp []rune
		temp = append(temp, ClosureNewNext[j:]...)
		ClosureNewNext = append(ClosureNewNext[0:j], '.')
		ClosureNewNext = append(ClosureNewNext, temp...)
		Closure = append(Closure, ClosureStruct{
			G: GrammarListStruct{
				S:    GrammarList[i].main,
				next: ClosureNewNext,
			},
			Node: j,
			End:  '#',
		})
	}
}*/

/*
DFA的构建
最重要的函数也比较复杂,总体思路跟书上区别不大
但是在操作上是先遍历构建完一个产生式的闭包直接把这个产生式扔到原理上的GO(),遍历完一个Closure再走下一个Closure,不回头
比较复杂,很多函数可以抽出来,懒了懒了
参数:void
返回:void
*/
func DFA_Build() {
	ClosureInit()
	for i := 0; i < len(ClosureUnit); i++ {
		for j := 0; j < len(ClosureUnit[i].Closure); j++ { //GO
			if ClosureUnit[i].Closure[j].Node == len(ClosureUnit[i].Closure[j].G.next) { //点到了最后跳过
				continue
			} else {
				//Clouse集
				for k := 0; k < len(GrammarList); k++ { //遍历原文法
					if ClosureUnit[i].Closure[j].G.next[ClosureUnit[i].Closure[j].Node] == GrammarList[k].main { //找到点后一个符号与原文法产生式左侧相同的
						var l int
						l = ClosureUnit[i].Closure[j].Node //用l表示下一个结点
						if len(ClosureUnit[i].Closure[j].G.next) == (ClosureUnit[i].Closure[j].Node + 1) { //如果只有一个,直接把当前终结符放入,因为β=End

							ClosureUnit[i].Closure = append(ClosureUnit[i].Closure, ClosureStruct{
								G:    GrammarList[k],
								Node: 0,
								End:  ClosureUnit[i].Closure[j].End,
							})
						} else if IsUpper(ClosureUnit[i].Closure[j].G.next[l]) { //如果是非终结符,遍历对应First并放入
							for m := 0; m < len(First); m++ {
								if First[m].main == ClosureUnit[i].Closure[j].G.next[l+1] {
									for n := 0; n < len(First[m].next); n++ {
										var IsSame bool //这一部位是去除Closure构造出现的相同产生式
										IsSame = false
										for o := 0; o < len(ClosureUnit[i].Closure); o++ {
											var flag bool
											flag = true
											for p := 0; p < len(ClosureUnit[i].Closure[o].G.next); p++ {
												if (len(ClosureUnit[i].Closure[o].G.next) != len(GrammarList[k].next)) || (ClosureUnit[i].Closure[o].G.next[p] != GrammarList[k].next[p]) {
													flag = false
													break
												}
											}
											if flag && (ClosureUnit[i].Closure[o].G.main == GrammarList[k].main) && (ClosureUnit[i].Closure[o].Node == 0) && (ClosureUnit[i].Closure[o].End == First[m].next[n]) { //找到了跳出循环
												IsSame = true
												break
											}
										}
										if !IsSame {
											ClosureUnit[i].Closure = append(ClosureUnit[i].Closure, ClosureStruct{
												G:    GrammarList[k],
												Node: 0,
												End:  First[m].next[n],
											})
										}
									}
								}
							}
						} else { //如果遇到终结符,是自己的闭包,不处理
							continue
						}
					}
				}

			}
			//GO
			//新Closure
			var NewClosure ClosureStruct
			NewClosure = ClosureUnit[i].Closure[j]
			NewClosure.Node++ //当前Node+1说明已经是下一个Struct了
			var IsSame bool
			IsSame = false
			var SameInt int
			for k := 0; k < len(ClosureUnit); k++ { //遍历DFA找有无相等
				for l := 0; l < len(ClosureUnit[k].Closure); l++ {
					var flag bool//去重判断
					flag = true
					for m := 0; m < len(ClosureUnit[k].Closure[l].G.next); m++ {
						if (len(ClosureUnit[k].Closure[l].G.next) != len(NewClosure.G.next)) || (ClosureUnit[k].Closure[l].G.next[m] != NewClosure.G.next[m]) {
							flag = false
							break
						}
					}
					if flag && (ClosureUnit[k].Closure[l].G.main == NewClosure.G.main) && (ClosureUnit[k].Closure[l].Node == NewClosure.Node) && (ClosureUnit[k].Closure[l].End == NewClosure.End) { //找到了跳出循环
						IsSame = true
						SameInt = k
						break
					}
				}
				if IsSame { //找到了继续跳出循环
					break
				}
			}
			if len(ClosureUnit[i].NextRune) == 0 { //指向为空直接加入
				ClosureUnit[i].NextRune = append(ClosureUnit[i].NextRune, ClosureUnit[i].Closure[j].G.next[NewClosure.Node-1]) //Rune加入DFA边
				if IsSame { //之前找到的处理
					ClosureUnit[i].NextUnit = append(ClosureUnit[i].NextUnit, SameInt) //Node放入DFA边
					//ClosureUnit[SameInt].Closure=append(ClosureUnit[SameInt].Closure,NewClosure)
					//Rune和Node放入单Closure
					ClosureUnit[i].Closure[j].NextRune = ClosureUnit[i].Closure[j].G.next[NewClosure.Node-1]
					ClosureUnit[i].Closure[j].NextUnit = SameInt
				} else { //没找到开新集合
					//开一个新的ClosureUnit
					ClosureUnit = append(ClosureUnit, ClosureUnitStruct{})
					ClosureUnit[len(ClosureUnit)-1].Closure = append(ClosureUnit[len(ClosureUnit)-1].Closure, NewClosure) //把当前Closure放入ClosureUnit
					ClosureUnit[i].NextUnit = append(ClosureUnit[i].NextUnit, len(ClosureUnit)-1)                         //开一个新的ClosureUnit
					//Rune和Node放入单Closure
					ClosureUnit[i].Closure[j].NextRune = ClosureUnit[i].Closure[j].G.next[NewClosure.Node-1]
					ClosureUnit[i].Closure[j].NextUnit = len(ClosureUnit) - 1
				}
			} else { //判断是不是已经DFA边
				IsFindSameInDFA := false
				var DFASameUnitInt int
				for k := 0; k < len(ClosureUnit[i].NextRune); k++ { //遍历项目集的指向
					if NewClosure.G.next[NewClosure.Node-1] == ClosureUnit[i].NextRune[k] { //没有相同的
						IsFindSameInDFA = true
						DFASameUnitInt = ClosureUnit[i].NextUnit[k]
					}
				}
				if IsFindSameInDFA { //找到相同->加入
					if IsSame {
						var HaveSame bool //这一部位是去除Closure移入时出现的相同产生式
						HaveSame = false
						for o := 0; o < len(ClosureUnit[DFASameUnitInt].Closure); o++ {
							var flag bool
							flag = true
							for p := 0; p < len(ClosureUnit[DFASameUnitInt].Closure[o].G.next); p++ {
								if (len(ClosureUnit[DFASameUnitInt].Closure[o].G.next) != len(NewClosure.G.next)) || (ClosureUnit[DFASameUnitInt].Closure[o].G.next[p] != NewClosure.G.next[p]) {
									flag = false
									break
								}
							}
							if flag && (ClosureUnit[DFASameUnitInt].Closure[o].G.main == NewClosure.G.main) && (ClosureUnit[DFASameUnitInt].Closure[o].Node == NewClosure.Node) && (ClosureUnit[DFASameUnitInt].Closure[o].End == NewClosure.End) { //找到了跳出循环
								HaveSame = true
								break
							}
						}
						if !HaveSame {

							ClosureUnit[DFASameUnitInt].Closure = append(ClosureUnit[DFASameUnitInt].Closure, NewClosure)
						}

						ClosureUnit[i].NextRune = append(ClosureUnit[i].NextRune, NewClosure.G.next[NewClosure.Node-1])
						ClosureUnit[i].NextUnit = append(ClosureUnit[i].NextUnit, DFASameUnitInt)
						ClosureUnit[i].Closure[j].NextRune = NewClosure.G.next[NewClosure.Node-1]
						ClosureUnit[i].Closure[j].NextUnit = DFASameUnitInt
					} else {
						ClosureUnit[DFASameUnitInt].Closure = append(ClosureUnit[DFASameUnitInt].Closure, NewClosure)
						ClosureUnit[i].NextRune = append(ClosureUnit[i].NextRune, NewClosure.G.next[NewClosure.Node-1])
						ClosureUnit[i].NextUnit = append(ClosureUnit[i].NextUnit, DFASameUnitInt)
						ClosureUnit[i].Closure[j].NextRune = NewClosure.G.next[NewClosure.Node-1]
						ClosureUnit[i].Closure[j].NextUnit = DFASameUnitInt
					}
				} else { //找不到相同可能返回
					ClosureUnit[i].NextRune = append(ClosureUnit[i].NextRune, NewClosure.G.next[NewClosure.Node-1])
					if IsSame { //找不到相同DFA边但是在前面DFA找到
						ClosureUnit[i].NextUnit = append(ClosureUnit[i].NextUnit, SameInt)
						ClosureUnit[i].Closure[j].NextRune = NewClosure.G.next[NewClosure.Node-1]
						ClosureUnit[i].Closure[j].NextUnit = SameInt
					} else {
						ClosureUnit = append(ClosureUnit, ClosureUnitStruct{})
						ClosureUnit[len(ClosureUnit)-1].Closure = append(ClosureUnit[len(ClosureUnit)-1].Closure, NewClosure)
						ClosureUnit[i].NextUnit = append(ClosureUnit[i].NextUnit, len(ClosureUnit)-1)
						//Rune和Node放入单Closure
						ClosureUnit[i].Closure[j].NextRune = ClosureUnit[i].Closure[j].G.next[NewClosure.Node-1]
						ClosureUnit[i].Closure[j].NextUnit = len(ClosureUnit) - 1
					}
				}
			}
		}
	}
}
