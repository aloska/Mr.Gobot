package universal

import (
	"errors"
	measure "github.com/hbollon/go-edlib"
	"golang.org/x/tools/container/intsets"
	"math/rand"
	"sort"
	"time"
)

/*
Эпоха революционных изменений, или эпоха генетической катастрофы.
Должны быть константы видовая и родовая достаточно низкие, чтобы все, что здесь описано могло происходить
*/

//катастрофический мейоз - гамета может быть образована как увеличением обычного набора, так и делецией некоторых хромосом
func (g Genotype) PolyMeyosis() Gameta{
	rand.Seed(time.Now().UnixNano())
	var ga Gameta
	//используем 4 гаметы стандартные, как при обычном мейозе, и они могут быть случайно все одинаковые
	//теоритически может получится гамета со всеми 4 внутри (сразу смешанный квадраплоид)

	for i:=0;i<4;i++ {
		if rand.Intn(100)>=POLYADDGAMETPERCENT {//есть ли добавление этой гаметы в набор полиплоидный
			gat := g.Meyosis()//получим стандартную случайную гамету
			lena := len(gat)
			if rand.Intn(100) > POLYDELCHROMGAMETPERCENT { //есть ли удаление некоторых хромосом?
				m:=rand.Intn(lena)+1//количество удаляемых хромосом
				if m>=lena{
					break//полностью пустая
				}
				ind:=rand.Intn(lena-m)
				gat=append(gat[:ind], gat[ind+m:]...)
			}
			ga = append(ga, gat...)
		}
	}
	if len(ga)==0{//а такое возможно
		 return g.PolyMeyosis()//ну попробуем еще раз
	}
	//теперь мутации, много мутации
	for i:=0;i<len(ga);i++{
		ga[i].Chromosome=Chromosome(Mutation(string(ga[i].Chromosome), POLYMUTAFACTORMEYOS, POLYMUTAMEYOSRUNEMAX))
		//ну и пробуем моноиды создать, шоуж
		ga[i],_=NewMonoid(ga[i].Chromosome) //плевать, если нет рабочих генов - это как бы не мудрено)))
	}

	return ga
}

//НЕнормальное скрещивание, при котором уж, еж и черепаха может скрестится (и не только попарно, а сразу все ыыы)
//это не шутка - можно передать функции любое число родителей
//лучше использовать в период катастроф или когда стандартное скрещивание невозможно
func PolyPairing(P ...Genotype) (Genotype, error){
	rand.Seed(time.Now().UnixNano())
	if len(P)<2{
		return nil, errors.New("cannot pairing less than 2 parents")
	}
	var F Genotype	//общий потомок
	var ga []Gameta //все гаметы
	//катастрофический мейоз
	lena:=intsets.MaxInt //найдем минимальное количество хромосом
	lenMax:=0  //и максимальное
	for _,p:=range P{
		g:=p.PolyMeyosis()
		l:=len(g)
		if l<lena{
			lena=l
		}
		if l>lenMax{
			lenMax=l
		}
		ga=append(ga,g)

	}
	//считаем, что хромосомы количеством lena парные или почти парные, и идут подряд, скрещиваем их
	for i:=0;i<lena;i++{
		var di Pairoid
		//отсортируем хромосомы по наиболее похожим
		jw:=JWSorter{}
		for _,g:=range ga{
			jw=append(jw, g[i].Chromosome)
		}
		sort.Sort(jw)
		//теперь в jw хромосомы расположены по порядку от самой похожей на всех до самой не похожей на всех
		//самая похожая будет M
		di.M,_=NewMonoid(jw[0])
		//а женская на то и женская, чтобы соединить в себе все остальные
		di.F,_=NewMonoid(Chromosome(PolyAggregation(jw[1:],0,2)))
		F=append(F,di)
	}
	//обработка непарных хромосом
	for i:=lena;i<lenMax;i++ {
		var di Pairoid
		//отсортируем хромосомы по наиболее похожим
		jw:=JWSorter{}
		for _, g := range ga {
			if len(g) <= i {
				break
			} //в этой гамете нет лишних хромосом
			jw=append(jw, g[i].Chromosome)
		}
		if len(jw)<2 {break}//с одной хромосомой каши не сваришь
		sort.Sort(jw)
		di.M,_=NewMonoid(jw[0])
		di.F,_=NewMonoid(Chromosome(PolyAggregation(jw[1:],0,2)))
		F=append(F,di)
	}

	return F,nil
}

//для реализации сортировки хромосом по мере Джаро-Винклера
//https://play.golang.org/p/5pfSGyx2IaO
type JWSorter []Chromosome

func (jws JWSorter) Len() int {return len(jws)}
func (jws JWSorter) Swap(i, j int) {jws[i], jws[j] = jws[j], jws[i]}
func (jws JWSorter) Less(i, j int) bool {
	lcs,_:=measure.LCSBacktrack(string(jws[i]), string(jws[j]))
	jwi:=measure.JaroWinklerSimilarity(string(jws[i]), lcs)
	jwj:=measure.JaroWinklerSimilarity(string(jws[j]), lcs)
	if jwi<jwj{
		return true
	}
	return false
}

//соединяет хромосомы по принципу: от 1-ой num-ная часть, от 2-ой num+1-ая часть...
func PolyAggregation(ch []Chromosome, ind, num int) string{
	ru:=[]rune(string(ch[0]))
	if len(ru)<=ind{
		return ""
	}
	if len(ch)==1 || ind+len(ru)/num>=len(ru){
		return string(ru[ind:])
	}
	nump:=num+1
	return string(ru[ind:ind+len(ru)/num]) + PolyAggregation(ch[1:],ind+len(ru)/num,nump)
}