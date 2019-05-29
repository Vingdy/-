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
  - [ ] 语法分析部分
  - [ ] 语义分析部分

## 开发日志

**2019.05.29**：吹逼摸鱼看书

**2019.05.27-2019.05.28**：完成词法分析器部分。