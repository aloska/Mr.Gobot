package main

import (
	"WithSVG/cmd/universal"
	"fmt"
	measure "github.com/hbollon/go-edlib"
	"math/rand"
)

func main() {
/*
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
	*/

	g1,_:=universal.MakeGenotypeFromStrings(";sdkfjbmalenks;klkdfgn⚤фыолывроцукекkjjkываылоцруузизшгуккнdfнгнjооозщсшмаитбдтывсиммячбапрдгозщⰡdfkjgrkjjkdfjsdk⚤сшмаитмуывсиммячбаⰡ",
		";sdkfjbnksfemalealkdfgn⚤фыолывроцукекkjjkdfнгнзщсшмаитбдтывсиммячбапрдгозщⰡdfkjgvhkkjkdfjsdk⚤сшмаитженабдтывсиммячбаⰡ",
	";sdk;djflfdklkdfgn⚤фыолыврfнгнзщтымячбапрдгозщⰡdfkjgvhkkjejjkdfjsdk⚤сшмамужбдтывсиммячбаⰡ",
		";sdkfjbdfgn⚤фыолывроцукнгнзщсшмаитбывсиммячбапрдгозщⰡdfkjgvhkkkjjkdfjsdk⚤сшмаитженабдтыммячбаⰡ"	,
		"паспаспаmalenks;klkdfgn⚤фыолывр гфыыоенгоцукекkjjkdfнгнзщс шмаитбдтывсиммячбапрдгозщⰡdfkjgrkjjkdfjsdk⚤сшмаитмуывсиммячбаⰡ",
		"прпрапрроksfemalealkdfgn⚤фыолывроцукекkjjkdfнгнзщсшмтбдтывсиммячбапрдгозщⰡdfkjgvhkkjkdfjsdk⚤сшоывагценывмаитженабдтывсиммячбаⰡ")
	g2,_:=universal.MakeGenotypeFromStrings("рпорроjflfddfgn⚤фыолыврfнгнз щтымячбапрдгозщⰡdfkjgvhkkjejjkdfjsdk⚤сшмамужбдтывсиммячбаⰡ",
		"апрпрпрпраgn⚤фыолывроцукнгнзщсшмаитбывсиммячбапрдгозщⰡdfkjgvhkkkdfjsdk⚤с шмаитженабдгвенулыврташгооырал тслтыммячбаⰡ"	,
	"паспаспаmaleыапиаыпиnиыпиифдуклоks;klkdfgn⚤фыолывр гфыыоенгоцудлфукр мзшпцгукр мжфш шмжфавбмифущ окоп098у84п34ьмдв. омяамфвамфвкифуккекkjjkdfнгнзщс шмаитбдтывсиммячбапрдгозщⰡdfkjgrkjjkdfjsdk⚤сшмаитмуывсиммячбаⰡ",
		"прпрапрроksfemalealkdfgn⚤фыолывроцукекkjjkdfнгнзщсшмтбдтывсиммячбапрдгозщⰡdfkjgvhkkjkdfjsdk⚤сшоывагценывмаитженабдтывсиммячбаⰡ",
	"рпорроjflfddfgn⚤фыолыврfнгнз щтымячбапрдгозщⰡdfkjgvhkkjejjkdfjsdk⚤сшмамужбдтывсиммячбаⰡ",
		"апрпрпрпраgn⚤фыолывроцукнгнзщсшмафцжущшкр зушйгкштп щйшеуо мэЖгргнве иФ лвоартм шдвгар мфвккам эщуфкзщлк09п304дмитбывсиммячбапрдгозщⰡdfkjgvhkkkdfjsdk⚤с шмаитженабдгвенулыврташгооырал тслтыммячбаⰡ"	)
	g3,_:=universal.MakeGenotypeFromStrings("паспаспаmaleыапиаыпиnиыпиифдуклоks;klkdfgn⚤фыолывр гфыыоенгоцудлфукр мзшпцгукр мжфш шмжфавбмифущ окоп098у84п34ьмдв. омяамфвамфвкифуккекkjjkdfнгнзщс шмаитбдтывсиммячбапрдгозщⰡdfkjgrkjjkdfjsdk⚤сшмаитмуывсиммячбаⰡ",
		"прпрапр⚤сioa;rwejmqwertyuiop[]asdfghjkl;'zxcvbnm,.йцукенгшщзхъфывапролджэячсмитьбюшоывагценывмаитженабдтывсиммячбаⰡроksfemalealkdfgn⚤фыолывроцукекkjjkdfнгнзщсшмтбдтывсиммячбапрдгозщⰡdfkjgvhkkjkdfjsdk⚤сшоывагценывмаитженабдтывсиммячбаⰡ")
	g4,_:=universal.MakeGenotypeFromStrings("рпорроjflfddfgn⚤фыолыврfнгнз щтымячбапрдгозщⰡdfkjgvhkkjejjkdfjsdk⚤сшмамужбдтывсиммячбаⰡ",
		"апрпр⚤сшоывагцюбьтимсчяэждлорпавыфъхзщшгнекуцйенывмаитженабдтывсиммячбаⰡпрпраgn⚤фыолывроцукнгнзщсшмафцжущшкр зушйгкштп щйшеуо мэЖгргнве иФ лвоартм шдвгар мфвккам эщуфкзщлк09п304дмитбывсиммячбапрдгозщⰡdfkjgvhkkkdfjsdk⚤с шмаитженабдгвенулыврташгооырал тслтыммячбаⰡ"	)

	var strtOrg []universal.Genotype

	strtOrg=append(strtOrg,g1)
	strtOrg=append(strtOrg,g2)
	strtOrg=append(strtOrg,g3)
	strtOrg=append(strtOrg,g4)
	evo:=universal.Evolution{Populations: strtOrg}
	evo.Functional=tional3

	evo2:=evo



	//evo.ForcePolyCross(20)
	//evo.ForcePolyCross(40)
	//evo.ForcePolyCross(50)
	i:=0
	for !evo.Step(0.985, 1000,true) && i<2000{
		i++
		if evo.Catastrofe==universal.ITERBETWEENCATASTROFE{
			fmt.Println("катастрофа: ", len(*evo.BestGenom), " maxpoly: ",maxpoly)
		}
		fmt.Println(i, ":\t",evo.BestFit(),"\t",len(evo.Populations)," sc:",universal.SpeciesConst," gc:",universal.GenusConst)
		fmt.Println(evo.BestGenom)
		fmt.Println((*evo.BestGenom)[0].M.Chromosome)
		//удалим пустые хромосомы, без генов?
		for a:=0;a<len(evo.Populations);a++{
			for b:=0;b<len(evo.Populations[a]);b++{
				if len(evo.Populations[a][b].M.Genes)==0 && len(evo.Populations[a][b].F.Genes)==0{
					evo.Populations[a]=append(evo.Populations[a][:b],evo.Populations[a][b+1:]...)
				}
			}
		}


		evo2.Step(0.98, 50,false)
		if i%27==0{
			evo.Populations=append(evo2.Populations[:15], evo.Populations...)

		}
		if i%23==0{
			evo2.Populations=append(evo.Populations[:5],evo2.Populations...)

		}

		//удалим пустые хромосомы, без генов?
		for a:=0;a<len(evo2.Populations);a++{
			for b:=0;b<len(evo2.Populations[a]);b++{
				if len(evo2.Populations[a][b].M.Genes)==0 && len(evo2.Populations[a][b].F.Genes)==0{
					evo2.Populations[a]=append(evo2.Populations[a][:b],evo2.Populations[a][b+1:]...)
				}
			}
		}

		maxpoly=0
	}
	fmt.Println(evo.BestGenom)

}
var maxpoly=0

func tional (g universal.Genotype) float32{
	var jw float32=0
	var jwn1,jwn2 float32

	for _,p:=range g {
		for _,gene:= range p.M.Genes {
			if rand.Intn(60)<50 {
				jwn1 = measure.JaroSimilarity(gene, "днк-полимераза")
				jwn2 = measure.JaroSimilarity(gene, "мозг запазухой")
			}else{
				jwn1 = measure.JaroWinklerSimilarity(gene, "днк-полимераза")
				jwn2 = measure.JaroWinklerSimilarity(gene, "мозг запазухой")
			}
			if jwn1>jwn2{
				if jwn1>jw {
					jw = jwn1
				}
			}else if jwn2>jw {
				jw = jwn2
			}
		}
		for _,gene:= range p.F.Genes {
			if rand.Intn(60)<50 {
				jwn1 = measure.JaroSimilarity(gene, "днк-полимераза")
				jwn2 = measure.JaroSimilarity(gene, "мозг запазухой")
			}else{
				jwn1 = measure.JaroWinklerSimilarity(gene, "днк-полимераза")
				jwn2 = measure.JaroWinklerSimilarity(gene, "мозг запазухой")
			}
			if jwn1>jwn2{
				if jwn1>jw {
					jw = jwn1
				}
			}else if jwn2>jw {
				jw = jwn2
			}
		}
	}
	la:=len(g)
	if la>maxpoly{
		maxpoly=la
	}

	return jw
}

func tional2 (g universal.Genotype) float32{
	var jw float32=0
	колгенов:=0

	for _,p:=range g {
		for _,gene:= range p.M.Genes {
			jw+= measure.JaroSimilarity(gene, "человекообразный генотип высшего уровня")
			колгенов++
		}
		for _,gene:= range p.F.Genes {
			jw+= measure.JaroSimilarity(gene, "человекообразный генотип высшего уровня")
			колгенов++
		}
	}
	la:=len(g)
	if la>maxpoly{
		maxpoly=la
	}
	if колгенов==0{
		return 0
	}
	return jw/float32(колгенов)
}

func tional3 (g universal.Genotype) float32{
	var jw float32=0


	for _,p:=range g {
		var s string
		for _,gene:= range p.M.Genes {
			s+=gene
		}
		jw1:= measure.JaroSimilarity(s, "стекло запотело")
		s=""
		for _,gene:= range p.F.Genes {
			s+=gene
		}
		jw2:= measure.JaroSimilarity(s, "стекло запотело")
		if jw1>jw2{
			jw=jw1
		}else{
			jw=jw2
		}
	}
	la:=len(g)
	if la>maxpoly{
		maxpoly=la
	}

	return jw
}
