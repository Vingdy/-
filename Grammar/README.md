# 文法

v:关键字var

p:关键字print

x:变量

c:常量

w:关键字while

S->vx

S->x=E

S->vx=E

S->pE

~~S->wEdoE~~

~~S->ifEE~~

~~S->ifEEelseE~~

E->E+T

E->E-T

T->F

T->T*F

T->T/F

T->F

F->(E)

F->x

F->c