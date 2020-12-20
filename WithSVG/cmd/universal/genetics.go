package universal

import (
	"errors"
	measure "github.com/hbollon/go-edlib"
	"math"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

type Chromosome string

/* Внимание! Названия не обязательно соответствуют настоящей генетике - они выбраны для своего удобства
*/

type Pairoid struct{
	M Monoid
	F Monoid
}

type Monoid struct {
	Chromosome  //сырая хромосома
	Genes []string	//гены, готовые к работе
}

//наши генотипы имеют сколько угодно парных наборов
type Genotype []Pairoid

//их гаметы состоят из наборов непарных
type Gameta []Monoid

//NewMonoid - создает моноид из хромосомы (парсит хромосому и достает из нее гены)
func NewMonoid(chr Chromosome) (Monoid, error){
	g:=Monoid{Chromosome: chr}
	//сначала транскрипция
	//парсим ДНК от старт до стоп-кодонов
	reg:=regexp.MustCompile("[\u26a2-\u26a4](.+?)[\u2c00-\u2c5e]")
	matches:=reg.FindAllString(string(chr), -1)//matches содержит все сырые гены, вместе со старт и стоп-кодонами
	for _,v:=range matches{
		out := strings.Map(func(r rune) rune {
			if utf8.RuneLen(r) < 3 { //сплайсинг - вырезаем все интроны (руны, длиной больше 2)
				return r
			}
			return -1
		}, v)

		g.Genes=append(g.Genes, out)
	}
	if len(g.Genes)==0{
		return g,errors.New("there is no any gene in string")
	}
	return g,nil
}

//смешивает два набора генов в шахматном порядке 
func (d* Pairoid) concatMF() []string{
	var res []string
	if len(d.M.Genes)>len(d.F.Genes){
		for i := 0; i < len(d.F.Genes); i++ {
			if i%2==0{
				res=append(res,d.M.Genes[i])
			}else{
				res=append(res,d.F.Genes[i])
			}
		}
		res=append(res, d.M.Genes[len(d.F.Genes):]...)
	}else{
		for i := 0; i < len(d.M.Genes); i++ {
			if i%2==0{
				res=append(res,d.M.Genes[i])
			}else{
				res=append(res,d.F.Genes[i])
			}
		}
		res=append(res, d.F.Genes[len(d.M.Genes):]...)
	}
	return res
}


//за основу берется a-хромосома
func Crossingover(a Chromosome, b Chromosome) Chromosome{
	rand.Seed(time.Now().UnixNano())
	//достанем руны из хромосом
	ar:=[]rune(a)
	br:=[]rune(b)
	//определим минимальную из двух длин
	lena:=0
	if len(ar)<len(br){
		lena=len(ar)
	}else{
		lena=len(br)
	}

	//от 1 до 3 максимально кроссинговеров
	mc:=rand.Intn(3)

	for i:=0;i<=mc; i++{
		//чем дальше от середины - тем больше вероятность кроссинговера
		ra:=int(rand.ExpFloat64()*float64(lena)*(math.Pow(-1,float64(i))))
		ia:=rand.Intn(lena/CROSSLENDIV)+2
		if ra< -lena/2 {
			//просто c 0 отрежем
			ra=0
		}else if ra>lena/2 {
			//просто с конца обрежем
			ra=lena-ia-1
		}else {
			ra=lena/2+ra
		}
		//собственно вырезка и замена
		for j:=ra;j<=ra+ia && j<lena;j++{
			c:=ar[j]
			ar[j]=br[j]
			br[j]=c
		}
	}
	rand.Seed(time.Now().UnixNano())//ну еще больше случайности, ну!
	//ну и мутация при кроссинговере случается
	if rand.Intn(1000001)>1000000-MUTAFACTORCROSS{ //случается ли?
		//количество рун, подверженых мутации
		m:=rand.Intn(MUTAPOINTRUNEMAX)+1 //при кроссинговере у нас только точечная мутация
		ind:=rand.Intn(lena-m)//случайно выбираем индекс начала
		mutype:=rand.Intn(10)//выбираем способ мутации
		switch mutype{
		case 0://делеция (крайне редкое событие)
			ar=append(ar[:ind], ar[ind+m:]...)//вырезали из слайса руны от ind до ind+m
		case 1,2,3,4://создание тандемного повтора
			ar=append(ar[:ind+m], append(ar[ind: ind+m], ar[ind+m:]...)...)
		case 5,6,7,8://инверсия
			inv:=make([]rune,0)
			m++//поскольку 1 символ сам с собой не поменяется, 1 означает поменять местами 2 символа...
			for i := ind+m-1; i >= ind; i-- {
				inv=append(inv, ar[i])
			}
			ar=append(ar[:ind], append(inv, ar[ind+m:]...)...)
		default://случайная замена руны-нуклеотида (тоже редко, как и делеция))
			for i := ind; i < ind+m; i++ {
				ar[i]=rune(rand.Intn(0xffff))
			}
		}
	}
	return Chromosome(ar)
}

//нормальный (без катастроф) мейоз, возвращает случайную гамету со случайной мутацией
func (g Genotype) Meyosis() Gameta{
	var gameta Gameta
	rand.Seed(time.Now().UnixNano())
	for i:=0; i<len(g);i++ {
		//для каждой пары F и M создаем свой моноид
		//подбросим монетку
		if rand.Intn(2)==1{//отца победила
			//и еще монетка - есть ли кросинговер с материнской
			if rand.Intn(2)==1{
				//есть кроссинговер - смешиваем
				chr:=Crossingover(g[i].M.Chromosome, g[i].F.Chromosome)
				//и создаем гаплоид
				gap, err:=NewMonoid(chr)
				if err!=nil{//пробуем еще раз создать, в случае неудачи
					chr=Crossingover(g[i].M.Chromosome, g[i].F.Chromosome)
					gap, _=NewMonoid(chr) // и плевать на ошибку теперь, будет как будет, без рабочих генов
				}
				gameta=append(gameta,gap)
			}else{//нет кроссинговера - просто добавим мужскую
				gameta=append(gameta, g[i].M)
			}
		}else{//матя победила
			//и еще монетка - есть ли кросинговер с отца
			if rand.Intn(2)==1{
				//есть кроссинговер - смешиваем
				chr:=Crossingover(g[i].M.Chromosome, g[i].F.Chromosome)
				//и создаем гаплоид
				gap, err:=NewMonoid(chr)
				if err!=nil{//пробуем еще раз создать, в случае неудачи
					chr=Crossingover(g[i].F.Chromosome, g[i].M.Chromosome)
					gap, _=NewMonoid(chr) // и плевать на ошибку теперь, будет как будет, без рабочих генов
				}
				gameta=append(gameta,gap)
			}else{//нет кроссинговера - просто добавим мужскую
				gameta=append(gameta, g[i].F)
			}
		}

	}
	return gameta
}


//нормальное скрещивание, при котором количество хромосом одинаковое у обоих родителей
func Pairing(P1, P2 Genotype) (Genotype, error){
	rand.Seed(time.Now().UnixNano())
	if len(P1)!=len(P2){//если число хромосом не одинаково - то нормально скрещиваться эти организмы не могут
		return nil, errors.New("cannot normal pairing non-same chromosome numbers, use revolution-pairing for this")
	}
	//получим случайные гаметы от обоих родителей
	gP1:=P1.Meyosis()
	gP2:=P2.Meyosis()
	F:=Genotype{}//будущий новый генотип
	for i:=0;i<len(gP1);i++{
		/*сравниваем хромосомы как строки мерой Джаро-Винклера:
		https://ru.wikipedia.org/wiki/%D0%A1%D1%85%D0%BE%D0%B4%D1%81%D1%82%D0%B2%D0%BE_%D0%94%D0%B6%D0%B0%D1%80%D0%BE_%E2%80%94_%D0%92%D0%B8%D0%BD%D0%BA%D0%BB%D0%B5%D1%80%D0%B0
		 */
		jw:=measure.JaroWinklerSimilarity(string(gP1[i].Chromosome), string(gP2[i].Chromosome))
		if jw >= SpeciesConst{//хромосомы достаточно похожи, просто выбираем, кто будет M, а кто F в новом генотипе
			if rand.Intn(100)>50 {//todo нужны ли мутации здесь или достаточно мейозных?
				F=append(F, Pairoid{M:P1[i].M, F:P2[i].F})
			}else{
				F=append(F, Pairoid{M:P1[i].F, F:P2[i].M})
			}

		}else if jw>=GenusConst{//хромосомы не достаточно похожи, но поскольку одного рода, применим специальное скрещивание
			//делим хромосомы на 3 равные части, и молимся, чтобы гены не порвались (хотя может и нужно, чтоб порвались - хз жеж)
			runeP1:=[]rune(string(gP1[i].Chromosome))
			runeP2:=[]rune(string(gP1[i].Chromosome))
			lena3P1:=len(runeP1)/3
			lena3P2:=len(runeP2)/3

			//первая треть
			p11, p21, err:= LCSCrossing(string(runeP1[:lena3P1]),string(runeP2[:lena3P2]))
			if err!=nil{
				return nil,err
			}
			//вторая треть
			p12, p22, err:= LCSCrossing(string(runeP1[lena3P1:lena3P1*2]),string(runeP2[lena3P2:lena3P2*2]))
			if err!=nil{
				return nil,err
			}
			//третья треть
			p13, p23, err:= LCSCrossing(string(runeP1[lena3P1*2:]),string(runeP2[lena3P2*2:]))
			if err!=nil{
				return nil,err
			}
			//создаем моноиды
			gap1, _:=NewMonoid(Chromosome(p11+p12+p13)) //нам неособо важно, есть ли рабочие гены в этой хромосоме при родовом скрещивании
			gap2, _:=NewMonoid(Chromosome(p21+p22+p23)) //нам неособо важно, есть ли рабочие гены в этой хромосоме при родовом скрещивании
			//выбираем папу-маму
			if rand.Intn(100)>50 {
				F=append(F, Pairoid{M:gap1, F:gap2})
			}else{
				F=append(F, Pairoid{M:gap2, F:gap1})
			}

		}else{//хромосомы не находятся в одном роду согласно представлениям о текущей эпохе, они не могут скрещиваться совсем
			return nil, errors.New("cannot pairing chromosome: different genus")
		}

	}
	return F, nil
}

//при родовом скрещивании, из двух участков от разных гомологичных хромосом создает 2 более близких участка, взвешивая наибольшие последовательности
//lcs при этом еще подвергается значительной мутации - ну как значительной?
//https://play.golang.org/p/pQ7x0otHK9X - песочница для теста
func LCSCrossing(p1n,p2n string)(string, string, error){
	rand.Seed(time.Now().UnixNano())
	//находим наибольшую общую последовательность замечательной функцией из пакета edlib
	if lcs,err:=measure.LCSBacktrack(p1n,p2n);err==nil{
		if len(lcs)/2<len(p1n)/3{//если наибольшая последовательность на целую треть меньше данного участка, то схожесть участков слабая
			//поэтому выберем, что делать - обе заменим на LCS или добавить к каждой LCS
			if rand.Intn(100)>50 {
				p1n = Mutation(lcs, MUTAFACTORCHROM, MUTACHROMRUNEMAX)
				p2n = Mutation(lcs, MUTAFACTORCHROM, MUTACHROMRUNEMAX)
			}else{
				p1n = p1n+Mutation(lcs, MUTAFACTORCHROM, MUTACHROMRUNEMAX)
				p2n = p2n+Mutation(lcs, MUTAFACTORCHROM, MUTACHROMRUNEMAX)
			}
		}else{//схожесть годная, взвешиваем чей участок более близок к LCS
			le1:=measure.LCSEditDistance(p1n,lcs)
			le2:=measure.LCSEditDistance(p2n,lcs)
			if le1<le2{//p1n ближе к общей последовательности, оставляем ее, а p2n меняем на LCS
				p2n=Mutation(lcs, MUTAFACTORCHROM, MUTACHROMRUNEMAX)
			}else if le1==le2{//если одинаково близки - выбираем случайно, кто поменяется
				if rand.Intn(100)>50{
					p1n=Mutation(lcs, MUTAFACTORCHROM, MUTACHROMRUNEMAX)
				}else{
					p2n=Mutation(lcs, MUTAFACTORCHROM, MUTACHROMRUNEMAX)
				}
			}else{
				p1n=Mutation(lcs, MUTAFACTORCHROM, MUTACHROMRUNEMAX)
			}
		}
	}else{
		return "", "", err
	}
	return p1n,p2n,nil
}

//функция мутации - возвращает мутированную строку или ту же самую (случайно)
//Внимание - у мейоза своя функция мутации!
//mutafactor - сколько раз на миллион случается мутация
//maxrune - максимално возможное кол-во рун, затронутых мутацией
func Mutation(chr string, mutafactor int, maxrune int) string{
	rand.Seed(time.Now().UnixNano())
	ar:=[]rune(chr)
	lena:=len(ar)
	if rand.Intn(1000001)>1000000-mutafactor{ //случается ли мутация?
		//количество рун, подверженых мутации
		m:=rand.Intn(maxrune)+1 //сколько рун затрагивается мутацией?
		if m>lena/2{
			m=lena/2	//но не больше половины длины исходной строки
		}
		ind:=rand.Intn(lena-m)//случайно выбираем индекс начала
		mutype:=rand.Intn(20)//выбираем способ мутации
		switch mutype{
		case 0,1,2://делеция (редкое событие)
			ar=append(ar[:ind], ar[ind+m:]...)//вырезали из слайса руны от ind до ind+m
		case 3,4,5,6,7,8,9,10://создание тандемного повтора
			ar=append(ar[:ind+m], append(ar[ind: ind+m], ar[ind+m:]...)...)
		case 11,12://инверсия
			inv:=make([]rune,0)
			m++//поскольку 1 символ сам с собой не поменяется, 1 означает поменять местами 2 символа...
			for i := ind+m-1; i >= ind; i-- {
				inv=append(inv, ar[i])
			}
			ar=append(ar[:ind], append(inv, ar[ind+m:]...)...)
		default://случайная замена руны-нуклеотида (довольно часто: 13,14,15,16,17,18,19)
			for i := ind; i < ind+m; i++ {
				ar[i]=rune(rand.Intn(0xffff))
			}
		}
	}
	return string(ar)
}

//создает генотип из набора строк
//строки попарно представляют из себя хромосомы M и F, соответственно количество строк должно быть кратно 2
func MakeGenotypeFromStrings(strs ...string) (Genotype, []error){
	var (
		G Genotype
		errs []error
	)
	lena:=len(strs)
	len2:=len(strs)/2
	if lena<2 || 2*len2!=lena{
		return nil, []error{errors.New("number of strings must be 2n")}
	}
	for i:=0;i<lena;i+=2{
		monM, err:=NewMonoid(Chromosome(strs[i]))
		if err!=nil{
			errs=append(errs, err)
		}
		monF, err:=NewMonoid(Chromosome(strs[i+1]))
		if err!=nil{
			errs=append(errs, err)
		}
		G=append(G,Pairoid{M:monM,F:monF})
	}
	return G,errs
}