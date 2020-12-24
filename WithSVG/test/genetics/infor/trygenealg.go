package main

import (
	u "WithSVG/cmd/universal"
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
)

var (
	sol u.Solution
)

func main() {

	in:=u.IO{V:[]int64{0,0}}
	out:=u.IO{V:[]int64{0}}
	m:=u.Memory{V:[]int64{0}}
	sol=u.Solution{IsAsync: false,
		Proc: []u.Processor{},
		In: []u.IO{in},
		Out: []u.IO{out},
		Mem: []u.Memory{m},
		Algs: []u.Algorithm{},
	}
/*
  g0,_:=u.MakeGenotypeFromStrings("⚤̫ӹʍѠ̇xⰡ ⚤xghnⰡ ⚤xdgfhnA001Ⱑ ⚤A101Ⱑ ⚤C000Ⱑ ⚤C001Ⱑ ⚤с шмаитженабдгвенулыврташгооырал тслтыммячбаⰡ",
   "⚤̫ӹʍѠ̇xⰡ ⚤xghnⰡ ⚤xdgfhnA001Ⱑ ⚤A101Ⱑ ⚤C000Ⱑ ⚤C001Ⱑ ⚤с шмаитженабдгвенулыврташгооырал тслтыммячбаⰡ",
   " ⚤.001Ⱑ ⚤.010Ⱑ ⚤.101Ⱑ ⚤.110Ⱑ ⚤/001Ⱑ ⚤/011Ⱑ ⚤/101Ⱑ ⚤0001Ⱑ",
   " ⚤.010Ⱑ ⚤.101Ⱑ ⚤.110Ⱑ ⚤/001Ⱑ ⚤/011Ⱑ ⚤/101Ⱑ")
*/
//*
	g1,_:=u.MakeGenotypeFromStrings("fjhglkdafhgadfg⚤+++Ⱑkj2i937yr78fuhndskmc,nfw2y4fn98wuopj⚤,000Ⱑfpmoi34uf98wynirfoipsef4rngfofi1ou3mpfv",
		"ofi1ou3mpfvpos⚤,000Ⱑ2wep4fno2w9iu4nf92fmd2i3unrf7t4yugfokr⚤,001Ⱑ5342iq1o39umf8wvyn7nowmpea.",
		"mof2i4fnyu98ywbgoiwfe⚤,111Ⱑ2wlenfio2iufnow23fnoqw8ne3yf92893f⚤-001Ⱑ nowien98wvmwpoirgo834⚤-011Ⱑc",
		";sdkf⚤A100Ⱑx xghn xdgfhn⚤A001Ⱑxgh ghdgh s⚤A101Ⱑfgh sh⚤C000Ⱑsgh sfghsq34⚤C001ⰡxfgsfgтженабдтыммячбаⰡ"	,
		"паспаспаmalenks;klkdfgn⚤фыолывр гфыыоенгоцукекkjjkdfнгнзщс шмаитбдтывсиммячбапрдгозщⰡdfkjgrkjjkdfjsdk⚤сшмаитмуывсиммячбаⰡ",
		"прпрапрроksfemalealkdfgn⚤фыолывроцукекkjjkdfнгнзщсшмтбдтывсиммячбапрдгозщⰡdfkjgvhkkjkdfjsdk⚤сшоывагценывмаитженабдтывсиммячбаⰡ")
	g2,_:=u.MakeGenotypeFromStrings("рпорроjflfddfgn⚤фыолыврfнгнз щтымячбапрдгозщⰡdfkjgvhkkjejjkdfjsdk⚤сшмамужбдтывсиммячбаⰡ",
		"ne3yf92893f⚤-001Ⱑ nowien98wvmwpoirgo834⚤-011Ⱑcnowi4uno234vmwepir3om49u⚤-101Ⱑ24f24gvsуеифуеревклн⚤-110Ⱑфупкфуки"	,
		"паспаспаmaleыапиаыпиnиыпиифдуклоks;klkdfgn⚤фыолывр гфыыоенгоцудлфукр мзшпцгукр мжфш шмжфавбмифущ окоп098у84п34ьмдв. омяамфвамфвкифуккекkjjkdfнгнзщс шмаитбдтывсиммячбапрдгозщⰡdfkjgrkjjkdfjsdk⚤сшмаитмуывсиммячбаⰡ",
		"прпрапрроksfemalealkdfgn⚤фыолывроцукекkjjkdfнгнзщсшмтбдтывсиммячбапрдгозщⰡdfkjgvhkkjkdfjsdk⚤сшоывагценывмаитженабдтывсиммячбаⰡ",
		"ук⚤.010Ⱑ⚤.101Ⱑфваифвпифеит⚤.110Ⱑфвптфвтпфптьпогьвраимв⚤/001Ⱑdsfygjufg⚤/011Ⱑjmk,gjoiugyftd⚤/101Ⱑadbnmhiu245⚤0001Ⱑ6567876543",
		"⚤0010Ⱑgrytumnb567⚤0110Ⱑ456765432⚤0011Ⱑfghjmnfbd234⚤0111Ⱑgdfhnbv3456⚤0000Ⱑfhgmffd3456⚤1002Ⱑbnfghnfgb6543⚤1012Ⱑvfbcnhjnh6523123⚤1102Ⱑ123hngfgd⚤1112Ⱑ12345⚤A000Ⱑncgh cgh ")
	g3,_:=u.MakeGenotypeFromStrings("пасп⚤A000Ⱑ+⚤A101Ⱑ+⚤C000Ⱑ+⚤C001Ⱑаспаmaleыапиаыпиnиыпиифд⚤C000Ⱑ+⚤C001Ⱑуклоks;klkdfgn⚤фыолывр гфыыоенгоцудлфукр мзшпцгукр мжфш шмжфавбмифущ окоп098у84п34ьмдв. омяамфвамфвкифуккекkjjkdfнгнзщс шмаитбдтывсиммячбапрдгозщⰡdfkjgrkjjkdfjsdk⚤сшмаитмуывсиммячбаⰡ",
		"прпр⚤A000Ⱑ+⚤A101Ⱑ+а⚤0010Ⱑgrytumnb567⚤0110Ⱑ456765432⚤0011Ⱑfghj⚤C000Ⱑ+⚤C001Ⱑmnfbd234⚤0111Ⱑgdfhnbv3456⚤0000Ⱑfhgmffd3456⚤1002Ⱑbnfghnfgb6543ячбапрдгозщⰡdfkjgvhkkjkdfjsdk⚤сшоывагценывмаитженабдтывсиммячбаⰡ",
		"апрпрпрпраgn⚤фыолывроцукнгнзщсшмафцжущшкр зушйгкштп щйшеуо мэЖгргнве иФ лвоартм шдвгар мфвккам эщуфкзщлк09п304дмитбывсиммячбапрдгозщⰡdfkjgvhkkkdfjsdk⚤с шмаитженабдгвенулыврташгооырал тслтыммячбаⰡ",
		"kdfjhglkdafhgadfg⚤+++Ⱑkj2i937yr78fuhndskmc,nfw2y4fn98wuopj⚤,000Ⱑfpmoi34uf98wynirfoipsef4rngfofi1ou3mpfvpos⚤,000Ⱑ2wep4fno2w9iu4nf92fmd2i3unrf7t4yugfokr⚤,001Ⱑ5342iq1o39umf8wvyn7nowmpea.")
	g4,_:=u.MakeGenotypeFromStrings("рпоррmnfbd234⚤0111Ⱑgdfhnbv3456⚤0000Ⱑfhgmffd3456⚤1002Ⱑbnfghnfgb6543⚤1012Ⱑvfbcnhjnh6523123⚤1102Ⱑ123hngfgd⚤1112Ⱑ12345⚤A000Ⱑncgh cgh ммячбаⰡ",
		"апрп⚤A100Ⱑx xghn xdgfhn⚤A001Ⱑxgh ghdgh s⚤A101Ⱑfgh sh⚤C000Ⱑsgh sfghsq34⚤C001Ⱑxfgsfgнзщсшмафцжущшкр зушйгкштп щйшеуо мэЖгргнве иФ лвоартм шдвгар мфвккам эщуфкзщлк09п304дмитбывсиммячбапрдгозщⰡdfkjgvhkkkdfjsdk⚤с шмаитженабдгвенулыврташгооырал тслтыммячбаⰡ",
		"прфукифкнгшщгшнгкфуеркпщпюош⚤.001Ⱑфукпуфкпфук⚤.010Ⱑ⚤.101Ⱑфваифвпифеит⚤.110Ⱑфвптфвтпфптьпогьвраимв⚤/001Ⱑdsfygjufg⚤/011Ⱑjmk,gjoiugyftd⚤/101Ⱑadbnmhiu245⚤0001Ⱑ6567876543",
		"⚤0010Ⱑgrytumnb567⚤0110Ⱑ456765432⚤0011Ⱑfghjmnfbd234⚤0111Ⱑgdfhnbv3456⚤0000Ⱑfhgmffd3456⚤1002Ⱑbnfghnfgb6543⚤1012Ⱑvfbcnhjnh6523123⚤1102Ⱑ123hngfgd⚤1112Ⱑ12345⚤A000Ⱑncgh cgh абдгвенулыврташгооырал тслтыммячбаⰡ")
//*/
	var strtOrg []u.Genotype

	//strtOrg=append(strtOrg,g0)
	//*
	strtOrg=append(strtOrg,g1)
	strtOrg=append(strtOrg,g2)
	strtOrg=append(strtOrg,g3)
	strtOrg=append(strtOrg,g4)
//*/
/*
	strtOrg,err:=u.ReadGenotypesFromFile("uni.txt")
	if err!=nil{
		fmt.Println(err)
		os.Exit(5)
	}
*/
	evo:=u.Evolution{Populations: strtOrg}
	evo.Functional=tionalAlg

	evo2:=evo



	//evo.ForcePolyCross(20)
	//evo.ForcePolyCross(40)
	//evo.ForcePolyCross(50)
	i:=0
	for !evo.Step(1, 200,true) && i<2000{
		i++
		if evo.Catastrofe==u.ITERBETWEENCATASTROFE{
			fmt.Println("катастрофа: ", len(evo.Populations[0]), " maxpoly: ",maxpol)
		}
		fmt.Println(i, ":\t",evo.BestFit(),"\t",len(evo.Populations)," sc:",u.SpeciesConst," gc:",u.GenusConst)
		fmt.Println(evo.Populations[0])
		//удалим пустые хромосомы, без генов?
		for a:=0;a<len(evo.Populations);a++{
			for b:=0;b<len(evo.Populations[a]);b++{
				if len(evo.Populations[a][b].M.Genes)==0 && len(evo.Populations[a][b].F.Genes)==0{
					evo.Populations[a]=append(evo.Populations[a][:b],evo.Populations[a][b+1:]...)
				}
			}
		}

		if i%5==0{
			evo.ForcePolyCross(200)
		}

		evo2.Step(1, 50,false)
		if i%13==0{
			evo.Populations=append(evo.Populations,evo2.Populations[:3]...)

		}
		if i%17==0{
			evo2.Populations=append(evo2.Populations,evo.Populations[:5]...)

		}

		//удалим пустые хромосомы, без генов?
		for a:=0;a<len(evo2.Populations);a++{
			for b:=0;b<len(evo2.Populations[a]);b++{
				if len(evo2.Populations[a][b].M.Genes)==0 && len(evo2.Populations[a][b].F.Genes)==0{
					evo2.Populations[a]=append(evo2.Populations[a][:b],evo2.Populations[a][b+1:]...)
				}
			}
		}

		maxpol=0
	}

	fmt.Println("\nbest index: ",evo.BestIndex())

	f, _ := os.Create("try.txt")


	w := bufio.NewWriter(f)

	for ii:=0;ii<len(evo.Populations[0]);ii++{
		fmt.Println(evo.Populations[0][ii].M)
		fmt.Println(evo.Populations[0][ii].F)
		fmt.Fprintln(w,evo.Populations[0][ii].M)
		fmt.Fprintln(w,evo.Populations[0][ii].F)
	}
	fmt.Println(evo.Populations[0])
	fmt.Fprintln(w,evo.Populations[0])
	comm,_:=evo.Populations[0].MakeAlgorithms(true,false)
	for ii:=0;ii<len(comm);ii++{
		str:=u.GetReadableFromCommands(comm[ii])
		fmt.Println(str)
		fmt.Fprintln(w,str)
	}
	w.Flush()
	f.Close()
}
var (
	maxpol=0
	prime=[]int{2,3,5,7,11,13,17,19,23,29,31,37,43,49,53,59,61,67,71,79,89}
)

func tionalAlg (g u.Genotype) float32{
	var jw float32=0



	comm,_:=g.MakeAlgorithms(true, false)
	if len(comm)==0{
		return 0
	}

	//5 раз
	for i:=0;i<5;i++ {
		solu:=sol
		for _,a:=range comm{
			solu.Algs=append(solu.Algs, u.Algorithm{Commands: a})
			solu.Proc=append(solu.Proc, u.Processor{})
		}

		end := make(chan int)
		quit := make(chan int,2)

		a := prime[rand.Intn(5) + 1]
		rand.Seed(time.Now().UnixNano())
		b := prime[rand.Intn(9) + 6]
		c := a * b
		solu.In[0].V[0] = int64(a)
		solu.In[0].V[1] = int64(b)
		solu.Out[0].V[0] = 999999999

		go solu.RunSync(end, quit)

		старт := time.Now()
		//не больше 100 миллисекунд
		normalEnd := false
		for time.Since(старт).Milliseconds() < 20 {
			select {
			default:
			case <-end:
				normalEnd = true
				break
			}
		}
		if!normalEnd{
			quit <- 1
			<-end
			close(end)
			close(quit)
			return 0
		}
		close(end)
		close(quit)

		res:=math.Abs(float64(solu.Out[0].V[0]-int64(c)))
		if res==0{
			jw+=1
		}else{
			jw+=float32(1.0/res)
		}
	}
	la:=len(g)
	if la>maxpol{
		maxpol=la
	}
	return jw/5.0
}
