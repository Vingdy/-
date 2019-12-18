参考:https://blog.csdn.net/jxch____/article/details/78688894

FIRST集的定义：

设G=(VT,VN,P,S)是上下文无关文法

FIRST(a)={a|a=>*ab,a∈VT, a,b∈V*}

若a=>*ε则规定ε∈FIRST (a)

FIRST(α)就是从α可能推导出的所有开头终结符号和可能的ε所构成的集合。



FIRST集的计算：//用通俗的语言讲

1.置FIRST(α)为空

2.遍历所有产生式左侧是α文法的式子如果右侧产生式第一位字符为终结符则把该字符放入FIRST(α)集

例子:A->aB,FIRST(A)={a}

3.如果右侧产生式第一位字符为非终结符则根据把该字符的FIRST放入FIRST(α)集

例子:A->BC,FIRST(A)=FIRST(B)

4.检查3.步骤中非终结符的FIRST集有无空字符(ε)，若有则把下一个非终结符字符的FIRST集也放入FIRST(α)中

5.重复上述步骤直到所有FIRST集都生成且不再变化



例：

   E  →TE’

   E’→+TE’

   E’→ε

   T  →T'F

   T’→*FT’

   T’→ε

   F→(E)|i



FIRST(E)=FIRST(T)//FIRST(T)未计算

FIRST(E‘)={+,ε}//步骤2.

FIRST(T)=FIRST(T')∪FIRST(F)//由于T’集合存在ε所以要加上FIRST(F)

FIRST(T')={*,ε}

FIRST(F)={(,i}  ->  FIRST(E)=FIRST(T)={*,ε,(,i}



FOLLOW集定义：

FOLLOW(A)={a| S=>*mAb 且a∈FIRST(b),m∈V*,b∈V+}

若 S=>*uAb, 且b =>*ε,则#∈FOLLOW( A)。



FOLLOW集的计算：//用通俗的语言讲

1.对于开始的文法的FOLLOW集加入{#}

例:对于G[S]文法置FOLLOW(S)={#}

2.看产生式左侧的符号，如果在产生式右侧找到则根据情况进行判断

如果该符号的右侧是最后一位，则把产生式左侧的FOLLOW集加入FOLLOW(α)

例:A->Bα，FOLLOW(α)=FOLLOW(A)

如果该符号的下一位是终结符，则把该终结符放入FOLLOW(α)

例:A->αbC，FOLLOW(α)={b}

如果该符号的下一位是非终结符，则把该符号除去空符号(ε)的FIRST集放入FOLLOW(α)

例:A->αBC，FOLLOW(α)=FIRST(B)-{ε}

3.检查2.步骤中非终结符的FIRST集有无空字符(ε)，若有则把下一个字符的进行同2.步骤的判断直至遇到终结符或者不含空符号的非终结符

5.重复上述步骤直到所有FOLLOW集都生成且不再变化



例：

   E  →TE’

   E’→+TE’

   E’→ε

   T  →FT'

   T’→*FT’

   T’→ε

   F→(E)|i



FOLLOW(E)={#，)}//步骤1.+最后一个式子

FOLLOW(E')=FOLLOW(E')=FOLLOW(E)={#,)}//第一个式子

FOLLOW(T)=(FIRST(E')-{ε})∪FOLLOW(E')={+,#,)}//第一个式子or第二个式子，由于E'已经是最后一个则加入FOLLOW(E')

FOLLOW(T’)= FOLLOW(T)= {+,) ,#}

FOLLOW(F)={FIRST(T’)-{e}})∪FOLLOW(T’) = {+,*,) ,#}  