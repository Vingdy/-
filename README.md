# CompilerExperiment

为了最终的编译原理大作业写的一个，
非常奶衣服的一个，功能十分受限的编译器。

## 指令集（暂定）

- `var x;`：定义变量`x`（浮点型），初值默认为0。
- `var x = [expression];`：定义变量`x`（浮点型），并将`[expression]`的值赋给`x`。
- `x = [expression];`：将表达式`[expression]`的值赋值给`x`，表达式仅支持带小括号的四则运算。
- `print [expression];`：输出表达式`[expression]`的值。

## 进度

- [ ] CompilerExperiment开发
  - [x] 词法分析部分
    - [x] 文件读取
    - [x] 预扫描处理注释
    - [x] 扫描主体代码
    - [x] DFA构造
    - [x] 各符号判断
    - [x] 各符号处理
  - [ ] 语法分析部分
    - [ ] LR表数据文件取出
    - [ ] AST构造
  - [ ] 语义分析部分
    - [ ]  

- [x] LR(1)表产生器开发
  - [x] 文件读取
  - [x] FIRST集合&FOLLOW集生成
  - [x] DFA生成
  - [x] ACTION和GOTO的生成
  - [x] 表结果生成
  - [x] 文件写入


## 开发日志

**2019.06.06**：完成ACTION集合GOTO集合的生成以及对应表结果，文件写入，LR(1)表产生器完成

**2019.06.05**：再次修改数据结构，完成FIRST集合&FOLLOW集的生成以及DFA的生成

**2019.06.04**：发现错误，修改数据结构，重写FIRST集合&FOLLOW集的生成

**2019.06.03**：完成FOLLOW集，继续看书

**2019.06.01-2019.06.02**：六一放假

**2019.05.31**：看书上课看书，完成FIRST集

**2019.05.30**：准备LR(1)表产生，看书上课看书

**2019.05.29**：吹逼摸鱼看书

**2019.05.27-2019.05.28**：完成词法分析器部分。