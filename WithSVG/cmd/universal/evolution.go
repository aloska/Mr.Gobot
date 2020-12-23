package universal

import (
	"math"
	"math/rand"
	"sort"
)

const(
	SPECIESCONSTINIT float32=0.60
	GENUSCONSTINIT float32=0.47
)
var(
	SpeciesConst float32=SPECIESCONSTINIT	/*видовая константа
	При нормальном скрещивании, хромосомы сравниваются как строки мерой Джаро-Винклера.
	Если эта мера больше равна SpeciesConst, то хромосомы наследуются без изменений потомком
	Если эта мера меньше SpeciesConst, то смотрится на дургую константу -
	*/
	GenusConst float32=GENUSCONSTINIT /*родовая константа
	При "почти" нормальном скрещивании, если мера равенства хромосом ниже GenusConst, то скрещивание не происходит.

	Если мера лежит между GenusConst и SpeciesConst - запускается особый алгоритм скрещивания, называемый родовым (разные виды
	скрещивают между собой в пределах одного рода)

	Межродовое скрещивание возможно только в эпоху революций (катастроф) и не считается нормальным
	*/

)


/*
*Эволюция понятия не имеет, что мы там оптимизируем - при ее создании нужно передать ей набор популяционный стартовый и
функцию, возвращающую тем меньше значения, чем хуже данный генотип справляется с задачей известной тому кот запустил эволюцию
Evolution также реализует интерфейс Sorter
*/
type Evolution struct{
	Populations []Genotype	/*текущий набор всех генетических организмов

	*/
	initial []Genotype //сохраняемые время от времени копии лучших, которые добавляются во время революций к слайсу

	Functional func(g Genotype) float32	/*функционал, который участвует в сравнении и сортировке организмов
	должен возвращать число тем меньше, чем хуже генотип соответствует поставленной задаче
	Можно создать набор тестов для генотипа, взвесить его на весах суда божьего и вернуть число крутости
	При этом так может статься, что лучшие организмы не смогут скрещиваться нормально, поскольку в разных родах
	*/

	bestFit float32
	BestGenom *Genotype
	Catastrofe int
	fluent int			//текущая итерация

	functionalValues map[int]float32	//результаты вычисления функционалов для каждого из генотипа
}
func (e *Evolution) BestFit() float32 {return e.bestFit}

func (e *Evolution) Len() int{
	return len(e.Populations)
}

func (e *Evolution) Swap(i, j int){
	e.Populations[i],e.Populations[j]=e.Populations[j],e.Populations[i]
}

func (e *Evolution) Less(i, j int) bool{
	var fitI,fitJ float32
	if f,ok:=e.functionalValues[i]; ok{
		fitI=f
	}else {
		fitI = e.Functional(e.Populations[i])
		e.functionalValues[i]=fitI
	}
	if f,ok:=e.functionalValues[j]; ok{
		fitJ=f
	}else {
		fitJ = e.Functional(e.Populations[j])
		e.functionalValues[j]=fitJ
	}
	if fitI>fitJ { //Внимание! здесь должно быть больше, чтобы отсортировало по убыванию
		if fitI>e.bestFit{
			e.bestFit=fitI
			e.BestGenom=&e.Populations[j]
		}
		return true
	}
	if fitJ>e.bestFit{
		e.bestFit=fitJ
		e.BestGenom=&e.Populations[i]
	}
	return false
}


//один шаг эволюции
//toFit - величина значения функционала, больше равно которой функция вернет true
//max - максимально допустимое количество животных в популяции
func (e *Evolution) Step(toFit float32, max int, withGlobalChanging bool) bool{

	if e.fluent==0 || e.fluent==2 || e.fluent==5 || e.fluent==9  || e.fluent%50==0{
		e.initial=append(e.initial, e.Populations[0])
	}

	max=int(float64(max)/math.Log(float64(e.fluent+2)))

	//создаем новые популяции обычным скрещиванием (все со всеми)
	lena:=len(e.Populations)//сколько животных до скрещивания

	k:=0//сколько дополнительно появилось
	kmax:=0

	for i:=0;i<lena-1;i++{
		for j:=i+1;j<lena-i;j++{
			if rand.Intn(100)<50 {
				G, err := Pairing(e.Populations[i], e.Populations[j])
				if err == nil {
					e.Populations = append(e.Populations, G)
					k++
				}
				kmax++
			}
		}
	}
	//если потомков больше предков - устраживаем константы скрещивания, так что на следующем шаге потомков будет меньше генерится
	if withGlobalChanging {
		if k > lena && SpeciesConst<1{
			SpeciesConst += 0.01
			GenusConst += 0.01
		} else if GenusConst>0.2{
			SpeciesConst -= 0.01
			GenusConst -= 0.01
		}
	}
	if e.fluent>4 && (k<kmax/20 || k*7>=kmax*8 || SpeciesConst>=0.99 || GenusConst<=0.11){//если прям сильно меньше или сильно больше - катастрофа
		if e.Catastrofe>0{//прошлый раз была жеж
			e.Catastrofe--
		}else {
			e.Populations = append(e.Populations,e.initial...)//добавим лучших из истории
			e.Populations = append(e.Populations,*e.BestGenom)
			for i := 0; i < lena-1; i++ {
				for j := i; j < lena; j++ {
					if rand.Intn(150)<50 {
						G, err := PolyPairing(e.Populations[i], e.Populations[j]) //полиплоидизация
						if err == nil {
							e.Populations = append(e.Populations, G)

						}
					}
				}
			}
			//после катастрофы устанавливаем константы эволюции в нормальное значени
			if withGlobalChanging {
				SpeciesConst = SPECIESCONSTINIT
				GenusConst = GENUSCONSTINIT
			}
			e.Catastrofe = ITERBETWEENCATASTROFE
		}
	}else{
		e.Catastrofe--
	}

	e.functionalValues=make(map[int]float32) //для того, чтобы вычислять функционал только один раз
	sort.Sort(e)//теперь популяции отсортированы так, что лучшие лежат первыми

	//обрежем популяцию снизу, если она превысила размер
	if len(e.Populations)>max{

			e.Populations = e.Populations[:max]

	}

	e.fluent++

	//после сортировки нам известна лучшая величина и лучший геном, он первый
	if e.bestFit>=toFit{
		return true
	}
	return false
}


func (e *Evolution) ForcePolyCross(max int){
	lena:=len(e.Populations)
	for i := 0; i < lena-1; i++ {
		for j := i + 1; j < lena; j++ {
			G, err := PolyPairing(e.Populations[i], e.Populations[j]) //полиплоидизация
			if err == nil {
				e.Populations = append(e.Populations, G)

			}
		}
	}
	sort.Sort(e)
	if len(e.Populations)>max {
		e.Populations = e.Populations[:max]
	}
}