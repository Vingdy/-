# CompilerExperiment

为了最终的编译原理实验写的一个，

真的非常简单的一个，功能十分受限的编译器。

基本是按照原理一步一步尝试修改完成的，很艰难而且痛苦

一开始想得挺多额，想做个比较完整的，最后还是输给了时间和精力，如果以后有精力的话再来优化一下代码吧（虽然我觉得不太可能

很多Go的特性没有用到，而且主要在数据结构的考虑比较久，本应该能做到更快更好的，乖乖看书吧

## 指令集

- `var x;`：定义变量`x`（浮点型），初值默认为0。
- `var x = [expression];`：定义变量`x`（浮点型），并将`[expression]`的值赋给`x`。
- `x = [expression];`：将表达式`[expression]`的值赋值给`x`，表达式仅支持带小括号的四则运算。
- `print [expression];`：输出表达式`[expression]`的值。
- 支持`/*` +`*/`与`//`的注释方式

## 使用方法

exe版:在同一目录下创建SourceProgram.txt和SourceGrammar.txt文件分别输入对应程序和文法，然后创建conf.yaml输入

```
projectPath: "./"
isUseLR1Build: true
grammarFile: "SourceGrammar.txt"
programFile: "SourceProgram.txt"
LR1TableFile: "LR1Table.txt"
```

然后运行exe即可看到结果，10s后关闭

程序版:直接下载然后运行main.go

Ps:isUseLR1Build如果设置为true会看到大量的LR1分析表信息生成一次后尽量改为false

## 使用例子

SourceProgram.txt:

```
var x=1+(2*3);/*asfdf
sdafdas
*/
print x;

;
```

SourceGrammar.txt:

```
X->S
S->vx
S->x=E
S->vx=E
S->pE
E->E+T
E->E-T
E->T
T->T*F
T->T/F
T->F
F->(E)
F->x
F->c
```

conf.yaml:


```
projectPath: "./"
isUseLR1Build: false
grammarFile: "SourceGrammar.txt"
programFile: "SourceProgram.txt"
LR1TableFile: "LR1Table.txt"
```

运行结果:

![1560101394370](C:\Users\zhchx\AppData\Roaming\Typora\typora-user-images\1560101394370.png)![1560101382194](C:\Users\zhchx\AppData\Roaming\Typora\typora-user-images\1560101382194.png)


## 进度

- [x] Compiler开发
  - [x] 词法分析部分
    - [x] 文件读取
    - [x] 预扫描处理注释
    - [x] 扫描主体代码
    - [x] DFA构造
    - [x] 各符号判断
    - [x] 各符号处理
  - [x] 语法分析部分
    - [x] LR表数据文件取出
    - [x] 词法分析总控程序
    - [x] 语法树构造
  - [x] 语义分析部分
    - [x] 对应语义分析
    - [x] 结果输出

- [x] LR(1)表产生器开发
  - [x] 文件读取
  - [x] FIRST集合&FOLLOW集生成
  - [x] DFA生成
  - [x] ACTION和GOTO的生成
  - [x] 表结果生成
  - [x] 文件写入


## 开发日志

**2019.06.10**：增加配置文件，删除测试与无用输出，交叉编译Win版本和Linux版的exe文件

**2019.06.09**：完成语义分析部分，修改语法分析中树生成的BUG，完成Comilper开发

**2019.06.08**：修改LR(1)表有一部分重复BUG，完成语法分析生成语法树的步骤，完成语法分析

**2019.06.07**：发现LR(1)表产生器BUG并完成修复，修改词法分析使其以句子形式保存，语法分析总控程序完成

**2019.06.06**：完成ACTION集合GOTO集合的生成以及对应表结果，文件写入，LR(1)表产生器完成

**2019.06.05**：再次修改数据结构，完成FIRST集合&FOLLOW集的生成以及DFA的生成

**2019.06.04**：发现错误，修改数据结构，重写FIRST集合&FOLLOW集的生成

**2019.06.03**：完成FOLLOW集，继续看书

**2019.06.01-2019.06.02**：六一放假

**2019.05.31**：看书上课看书，完成FIRST集

**2019.05.30**：准备LR(1)表产生，看书上课看书

**2019.05.29**：吹逼摸鱼看书

**2019.05.27-2019.05.28**：完成词法分析器部分。