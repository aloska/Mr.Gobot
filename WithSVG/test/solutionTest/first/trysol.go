package main

import (
	"WithSVG/cmd/universal"
	"fmt"
)

func main() {

	sol, err:=universal.NewSolution("c:/ALOSKA/work/solutions/minimum.json")
	//sol:=universal.Solution{}
	//err:=sol.Init("c:/ALOSKA/work/solutions/SOL-minimum")
	if err!=nil{
		fmt.Println(err)
	}
	defer sol.Exit()

	sol.In[0].V[0]=5
	sol.In[0].V[1]=-5
	sol.In[0].V[2]=56516
	sol.In[0].V[3]=5545
	sol.In[0].V[4]=-42
	sol.In[0].V[5]=2135
	sol.In[0].V[6]=154

	sol.Proc.PC=0	//иначе не будет ничего делать, если уже выключился - он свое состояние помнит)))
	sol.Run()
	sol.Save()

	fmt.Println(sol.Out[0].V[0])


/*
	asm:=`

Траливали какая-то

asm:
	
	nop
	add x1, x2, x26
	addi x15, x21, -14587
	addi x0, x22, 2458
	sub x3,x4,x1
	subi x1,x2,-45
	subi x1,x2,45
	mul x3,x4,x1
	muli x1,x2,-45
	muli x1,x2,45
	div x3,x4,x1
	divi x1,x2,-45
	divi x1,x2,45
	rem x3,x4,x1
	remi x1,x2,-45
	remi x1,x2,45
	and x3,x4,x1
	andi x1,x2,-45
	andi x1,x2,45
	or x3,x4,x1
	ori x1,x2,-45
	ori x1,x2,45
	xor x3,x4,x1
	xori x1,x2,-45
	xori x1,x2,45
	sll x3,x4,x1
	slli x1,x2,45
	srl x3,x4,x1
	srli x1,x2,45
	ldm x15,0,1
	ldin x15,5,10
	stm 4,4,x2
	stout 5,8, x5
	beq x1,x2, 12
	bge x1,x2, 12
	blt x1,x2, 12
	bne x1,x2, 12
	jmp -12
	seq x1, x2, x3
	sge x1, x2, x3
	slt x1, x2, x3
	sne x1, x2, x3
	push x31
	pop	x24
	li x23, -4578
	

`
	if codons, err:=universal.GetCodonsFromAsmString(&asm); err==nil{
		s:=universal.GetReadableFromCodons(*codons)
		fmt.Println(s[0])
		fmt.Println(s[1])
	}else{
		fmt.Println(err)
	}
*/
}
