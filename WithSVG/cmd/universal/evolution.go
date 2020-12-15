package universal

import (
	"math"
	"math/rand"
	"time"
)


//за основу берется a-хромосома
func Crossingover(a Chromosome, b Chromosome) Chromosome{
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
	rand.Seed(time.Now().UnixNano())
	//от 1 до 3 максимально кроссинговеров
	mc:=rand.Intn(3)

	for i:=0;i<=mc; i++{
		//чем дальше от середины - тем больше вероятность кроссинговера
		ra:=int(rand.ExpFloat64()*float64(lena)*(math.Pow(-1,float64(i))))
		ia:=rand.Intn(lena/20)+2
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
	return Chromosome(ar)
}

//возвращает случайную гамету со случайной мутацией
func (g Genotype) Meyosis() Gameta{
	var gameta Gameta
	for i:=0; i<len(g);i++ {
		//для каждой пары F и M создаем свой гаплоид
		//подбросим монетку
		if rand.Intn(2)==1{//отца победила
			//и еще монетка - есть ли кросинговер с материнской
			if rand.Intn(2)==1{
				//есть кроссинговер - смешиваем
				chr:=Crossingover(g[i].M.Chromosome, g[i].F.Chromosome)
				//и создаем гаплоид
				gap, err:=NewGaploid(chr)
				if err!=nil{//пробуем еще раз создать, в случае неудачи
					chr=Crossingover(g[i].M.Chromosome, g[i].F.Chromosome)
					gap, _=NewGaploid(chr) // и плевать на ошибку теперь, будет как будет, без рабочих генов
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
				gap, err:=NewGaploid(chr)
				if err!=nil{//пробуем еще раз создать, в случае неудачи
					chr=Crossingover(g[i].F.Chromosome, g[i].M.Chromosome)
					gap, _=NewGaploid(chr) // и плевать на ошибку теперь, будет как будет, без рабочих генов
				}
				gameta=append(gameta,gap)
			}else{//нет кроссинговера - просто добавим мужскую
				gameta=append(gameta, g[i].F)
			}
		}

	}
	return gameta
}