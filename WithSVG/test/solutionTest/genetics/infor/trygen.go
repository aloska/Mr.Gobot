package main

import (
	"WithSVG/cmd/universal"
	"fmt"
)

func main() {

	var (
		algs [][]universal.Command
		G universal.Genotype
	)


	gen1:=`

Траливали какая-то

asm:
	
	nop
	add x1, x2, x26
	addi x15, x21, -14587
	addi x0, x22, 2458
	sub x3,x4,x1
	subi x1,x2,-45
	subi x1,x2,45
`
	gen2:=`
asm:
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
`
	gen3:=`
asm:
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
`
	gen4:=`
asm:
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

	alg1,_:=universal.GetCommandsFromAlgorithmString(&gen1)
	alg2,_:=universal.GetCommandsFromAlgorithmString(&gen2)
	alg3,_:=universal.GetCommandsFromAlgorithmString(&gen3)
	alg4,_:=universal.GetCommandsFromAlgorithmString(&gen4)
	algs=append(algs,*alg1)
	algs=append(algs,*alg2)
	algs=append(algs,*alg3)
	algs=append(algs,*alg4)
	introns:=universal.MakeEmptyIntrons(algs)
	introns[0][1]="Ｔｒｙｉｎｇ２．１"
	mscribes:=universal.MakeEmptyScribes(algs)
	mscribes[0]="Ну попробуем шо уж там?"
	fscribes:=universal.MakeEmptyScribes(algs)

	pa, _:=universal.MakePairoidFromAlgs(algs,introns,mscribes,fscribes)
	G=append(G,pa)
	algnew,err:=G.MakeAlgorithms()
	if err!=nil{
		fmt.Println(err)
	}
	fmt.Println(pa.M)
	fmt.Println(pa.F)

	for _,v:=range algnew{
		readable:=universal.GetReadableFromCommands(v)
		fmt.Println(readable)
	}

}

