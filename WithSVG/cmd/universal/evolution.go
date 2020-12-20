package universal

import "sort"

const(
	SPECIESCONSTINIT float32=0.84
	GENUSCONSTINIT float32=0.66
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
	Catastrofe =false
)


/*
*Эволюция понятия не имеет, что мы там оптимизируем - при ее создании нужно передать ей набор популяционный стартовый и
функцию, возвращающую тем меньше значения, чем хуже данный генотип справляется с задачей известной тому кот запустил эволюцию
Evolution также реализует интерфейс Sorter
*/
type Evolution struct{
	Populations []Genotype	/*текущий набор всех генетических организмов
	*/

	Functional func(g Genotype) float32	/*функционал, который участвует в сравнении и сортировке организмов
	должен возвращать число тем меньше, чем хуже генотип соответствует поставленной задаче
	Можно создать набор тестов для генотипа, взвесить его на весах суда божьего и вернуть число крутости
	При этом так может статься, что лучшие организмы не смогут скрещиваться нормально, поскольку в разных родах
	*/

	bestFit float32
	bestGenom *Genotype
}
func (e *Evolution) BestFit() float32 {return e.bestFit}

func (e *Evolution) Len() int{
	return len(e.Populations)
}

func (e *Evolution) Swap(i, j int){
	e.Populations[i],e.Populations[j]=e.Populations[j],e.Populations[i]
}

func (e *Evolution) Less(i, j int) bool{
	fitI:=e.Functional(e.Populations[i])
	fitJ:=e.Functional(e.Populations[j])
	if fitI>fitJ { //Внимание! здесь должно быть больше, чтобы отсортировало по убыванию
		if fitI>e.bestFit{
			e.bestFit=fitI
			e.bestGenom=&e.Populations[j]
		}
		return true
	}
	if fitJ>e.bestFit{
		e.bestFit=fitJ
		e.bestGenom=&e.Populations[i]
	}
	return false
}

//один шаг эволюции
//toFit - величина значения функционала, больше равно которой функция вернет true
//max - максимально допустимое количество животных в популяции
func (e *Evolution) Step(toFit float32, max int) bool{

	//создаем новые популяции обычным скрещиванием (все со всеми)
	lena:=len(e.Populations)//сколько животных до скрещивания
	k:=0//сколько дополнительно появилось
	kmax:=0
	for i:=0;i<lena-1;i++{
		for j:=i+1;j<lena-i;j++{
			G,err:=Pairing(e.Populations[i],e.Populations[j])
			if err==nil{
				e.Populations=append(e.Populations,G)
				k++
			}
			kmax++
		}
	}
	//если потомков больше предков - устраживаем константы скрещивания, так что на следующем шаге потомков будет меньше генерится
	if k>lena{
		SpeciesConst+=0.01
		GenusConst+=0.01
	}else { //если потомков меньше
		SpeciesConst-=0.01
		GenusConst-=0.01
	}
	if k<max/20 || k*5>=kmax*4 || SpeciesConst>=1 || GenusConst<=0{//если прям сильно меньше или сильно больше - катастрофа
		if Catastrofe{//прошлый раз была жеж
			Catastrofe=false
		}else {
			for i := 0; i < lena-1; i++ {
				for j := i + 1; j < lena-i; j++ {
					G, err := PolyPairing(e.Populations[i], e.Populations[j]) //полиплоидизация
					if err == nil {
						e.Populations = append(e.Populations, G)

					}
				}
			}
			//после катастрофы устанавливаем константы эволюции в нормальное значени
			SpeciesConst = SPECIESCONSTINIT
			GenusConst = GENUSCONSTINIT
			Catastrofe = true
		}
	}

	sort.Sort(e)//теперь популяции отсортированы так, что лучшие лежат первыми
	//обрежем популяцию снизу, если она превысила размер
	if len(e.Populations)>max{
		e.Populations=e.Populations[:max]
	}

	//после сортировки нам известна лучшая величина и лучший геном, он первый
	if e.bestFit>=toFit{
		return true
	}
	return false
}