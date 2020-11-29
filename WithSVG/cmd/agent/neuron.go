package agent

import "math/rand"

//NumberToXY - переводит index в одномерном массиве в двумерные координаты x, y
func NumberToXY(N uint32, maxX uint32) (uint32, uint32) {
	return N % maxX, N / maxX
}

/*

	Y
	^
	|
	|         (13, 17) -> 13 + 17 * 20 = 353
	|                                    353 % 20 = 13
	|                                    353 / 20 = 17
	0,0 -----------> X
	       maxX=20
*/

//XYToNumber - переводит двумерные координаты  в index одномерного массива
func XYToNumber(X uint32, Y uint32, maxX uint32) uint32 {
	return X + Y*maxX
}

//Dendrite - структура описания дендритов
//8 байт
type Dendrite struct {
	/*
		НЕЛЬЗЯ МЕНЯТЬ ПОРЯДОК СЛЕДОВАНИЯ ПОЛЕЙ!!! И ДОБАВЛЯТЬ, КАК-ТО ИЗМЕНЯТЬ ЭТО!!! ИСПОЛЬЗУЕТСЯ unsafe!!! для маппинга структуры на часть файла
	*/
	//typed - тип дендрита
	/*
		0-никакой тип пока, онтогенез
		0x01-0x0f - резерв

		0x10-глютамат-ионный
		0x11-глютамат-метаб
		0x12-ГАМК-ионный
		0x13-ГАМК-метаб
		0x14-ацх-ионный
		0x15-ацх-метаб
		0x16-АА-ионный
		0x17-АА-метаб
		0x18-NE-ионный
		0x19-NE-метаб
		0x1A-DOP-ионный
		0x1B-DOP-метаб
		0x1C-SER-ионный
		0x1D-SER-метаб
		0x1E-резерв
		0x1F-резерв

		Старшие пол-байта говорят о количестве рецепторов на дендриите. Напр, 0x30 - глютаматный ионный в 3 раза больше рецепторов, чем 0x10
	*/
	Typed DendriteTypeEnum
	Ca byte   //количество кальция в этом шипике
	State byte	//состояние дендрита - используется ВМ
	Charge int8	//заряд на дендрите -  заряды всех дендритов средневзвешиваются, и в соме ПД или нет
	N  uint32 //номер синапса в файле синапсов. Координаты его вычисляются по номеру:  y = N / maxX   x = N % maxX

}

//Axon - структура описания аксонов
//8 байт
type Axon struct {
	/*
		НЕЛЬЗЯ МЕНЯТЬ ПОРЯДОК СЛЕДОВАНИЯ ПОЛЕЙ!!! И ДОБАВЛЯТЬ, КАК-ТО ИЗМЕНЯТЬ ЭТО!!! ИСПОЛЬЗУЕТСЯ unsafe!!! для маппинга структуры на часть файла
	*/
	/*
		vesiculs - заполняется из сомы клетки и обратным захватом
	*/
	Vesiculs byte   //состояние визикул с нейромедиатором
	Ca       byte   //количество кальция в этом шипике аксона -
	Na		 byte //количество натрия в шипике
	K    	 byte //калий
	N        uint32 //номер синапса в файле синапсов. Координаты его вычисляются по номеру:  y = N / maxX   x = N % maxX
}

//DoAChMediate - может вызвать только клетка ацх-эргическая, выброс ацетилхолина в синапс
/*
a *Axon - аксон, от куда выброс
to *Chemical - синапс, куда плювать
byte maxVesicul - максимально возможное количество, которое клетка хочет выплюнуть. Может не равняться реально выплюнотому

return byte - возвращает количество реально выброшенного вещества
*/
func (ax *Axon) DoAChMediate(to *Chemical, maxVesicul byte) byte {
	if to.WASTE >= 0xfff0 || to.ACh > 0xf0 || ax.Vesiculs == 0 {
		return 0
	}
	//определяем реально возможное, ближайшее к maxVesicul количество выбрасываемого вещества
	a := 0xff - to.ACh
	if maxVesicul > ax.Vesiculs {
		maxVesicul = ax.Vesiculs
	}
	if maxVesicul > a {
		maxVesicul = a
	}
	//заходит кальций в аксон, КАЛЬЦИЕВЫЙ ТОК!
	if int16(ax.Ca)+int16(maxVesicul) > 0xff {
		ax.Ca = 0xff
	} else {
		ax.Ca = ax.Ca + maxVesicul
	}
	//плюем!
	ax.Vesiculs = ax.Vesiculs - maxVesicul
	to.ACh = to.ACh + maxVesicul
	return maxVesicul
}

//AChSynt - может вызвать только аксон АЦХ-эргический, синтезирует внутри аксона ацетилхолин, беря холин из синапса
//вот здесь и кальций понадобится!!! При плювании, Ca входит в аксон, а при выходе Ca - заходит холин
//(ну а мы не имитируем вход холина на аксоне, а сразу синтезируем)
//еще сома клетки может захватывать холин и раздавать по аксонам своим - см. (c *Chemical) AChSynt()
func (ax *Axon) AChSynt(from *Chemical) bool {
	if from.CHOL < 1 || ax.Vesiculs > 0xf2 || ax.Ca < 1 {
		return false
	}
	from.CHOL = from.CHOL - 1
	ax.Vesiculs = ax.Vesiculs + 1
	ax.Ca = ax.Ca - 1
	return true
}

//Neuron - нейрон. Все клетки, кроме рецепторов и эффекторов подходят под описание.
/*У стволовых неспециализированных клеток нет пока ни дендритов, ни аксонов. Она живет себе в цикле кребса, как-то обменивается со средой, и потихоньку ее
дендриты начинают рости. Когда дендрит встречает вещество, означающее, что  рядом находится чей-то аксон. Например, глютамат, который не поставляется
организмом, а только генерируется внутри другого аксона, то дендрит этой клетки становится глютаматным.
Пока клетка не подключится минимум пятью дендритами - она существует как глиальная клетка, обслуживая некоторые аксоны нейронов.
Когда клетка подключится таким образом минимум пятью дендритами, то выясняется ее специализация. Чем больше дендритов какого-то одного типа в ней есть,
такой специализацией она и начинает обладать.
Меняется ее typen на какой-то конкретный, более 0x20
С этого момента дендриты перестают рости и начинают рости аксоны, начиная с точки максимального удаления самого далекого дендрита. Они уже точно плюют то
вещество, типом которого стала клетка. Соответсвенно, клетка начинает обслуживать аксоны и генерировать внутри себя этот нейромедиатор.

Поскольку от тела всех работающих неройнов распрстраняется NO и CO, то аксоны могут расти по градиенту этих газов (там где больше газа - больше
надежды на встречу с клеткой-мишенью)

Дендриты и аксоны растут от тела клетки, и за один цикл делают не больше одного шага

Аксональная пластичность
------------------------
Аксон может установить связь с дендритом и поддерживать ее, пока в ответ на его нейромедиатор получает приемлемый ответ в виде
вещества-предшественника или анандамида. Если в щели не будет вещества-предшественника (например, для глютамата и ГАМК - это глютамин), или
хотя бы АА, то аксон смотрит, есть ли там сам нейромедиатор. Если есть, то это ему сигнал о прекращении синапса.
Т.е. прошлый раз он туда плевал, и там оно и осталось, и нет попыток его плевки разрушить

Дендритная пластичность
-----------------------
Долгосрочная потенциация и депрессия.
Существуют определенные уровни кальция в дендрите, который модулирует работу синапса:
1. Ca > max, NE > max
Большое количество кальция внутри дендрита говорит о том, что клетка недавно разряжалась. Если в присутствии большого количества кальция
дендрит получает мощные сигналы нейромедиатора, то дендрит может выключать ионные и метаботропные каналы. Тем самым формируется краткосрочная
депрессия и дендрит проводит меньше  тока натрия и калия, равно как и кальция.
2. Если краткосрочная депрессия будет продолжаться, то в определенные момент дендрит может полностью исключить свои рецептеры и отключиться
от связи с аксоном. Это будет долгосрочная депрессия.
3. Ca < min, Ne > max Краткосрочная потенциация происходит, когда дендрит видит большое количество нейромедиатора в синапсе и малое количество кальция внутри себя.
Тогда дендрит может увеличить силу синапса, добавляя рецепт-каналов, реагирующих на нейромедиатор.

Кальций, таким образом, является некими дендритными часами, сообщающими каждому дендриту о времени последнего разряда. И если кальция много -
значит разряд был недавно, и значит мы не сможем разрядится сейчас, и поэтому тот нейромедиатор, который нам подсовывает аксон бесполезен или
даже губителен.

Низкий уровень кальция сигнализирует о том, что клетка в общем готова к работе и что она довольно давно разряжалась.

4. max > Ca > min, max > NE > min
Поэтому, обученный нейрон в основном работает вне режима потенциации и депресии. Когда у дендритов средний уровень медиатора на синапсе
совпадает с низким уровнем кальция. В этом случае синапс не усиливается и не ослабляется.

5. Чем сильнее синапс (чем больше у него каналов), тем быстрее кальций из дендрита выводится наружу, и тем быстрее дендрит будет способен
участвовать в генерации ПД клетки.

Если клетка делает ПД, а в дендрите много кальция, то его каналы не открываются, и не усиливают ПД.
6. Ca > max, NE < min
Высокий уровень кальция и низкий уровень нейромедиатора не влияет на пластичность, но существенно влияет на результат работы клетки в целом:
в то время, когда клетка разряжается, большое количество кальция в дендрите не дает возможности натрию войти и вытеснить кальций. Чтобы натрий
вошел, должен быть нейромедиатор, открывающий натриевый канал.
Слабый дендрит будет медленно выпускать кальций и, как следствие будет реже готовым к приему сигнала и реже участвовать в жизни всей клетки.
Сильный дендрит будет быстрее выпускать кальций и будет чаще учавстовать в жизни клетки, и принимать сигналы с высокой частотой.

У нас используется совершенно простая модель кальциевых каналов. В отличие от натрия и калия, кальций не будет участвовать в обмене со средой.
Он просто берется из неоткуда (а там много кальция))) и уходит в никуда. Нам важен кальций больше как маркер времени и характеристика пластичности.
*/
//поля сгруппированы по 4 байта, всего 296 байт
type Neuron struct {
	/*
		НЕЛЬЗЯ МЕНЯТЬ ПОРЯДОК СЛЕДОВАНИЯ ПОЛЕЙ!!! И ДОБАВЛЯТЬ, КАК-ТО ИЗМЕНЯТЬ ЭТО!!! ИСПОЛЬЗУЕТСЯ unsafe!!! для маппинга структуры на часть файла
	*/
	Typen NeuronTypeEnum //byte [0]Тип ячейки, для стволовой клетки - 0x10
	State byte           /*состояние клетки - используется ВМ для хранения инфы о текущем состоянии

	*/
	SynNumber SynEnum    //номер синаптического поля Ширину и высоту поля хранят сами структуры поля, и еще гены нейронов
	Gen uint16/* Номер гена (индекс гена в слайсе genes []GenNeuron), к которому относится нейрон
	Это важный параметр, по которому нейрон может узнать гено-специфические способы работы
	Устанавлвается во время генерации клетки из гена
	*/
	SynNumberAxons SynEnum //Номер поля, где аксоны

	Chemic    Chemical     //[]Структура для химии
	N         uint32       //номер нейрона в файле синапсов. Координаты его вычисляются по номеру:  y = N / maxX   x = N % maxX
	D         uint16       //[]Состояния дендритов. каждый бит отвечает за дендрит. Если он 1 - дендрит есть и его надо обрабатывать
	A         uint16       //[]Состояния аксонов. каждый бит отвечает за 1 аксон. Если он 1 - аксон есть и его надо обрабатывать
	Dendrites [16]Dendrite //[]Координаты дендритов в общем файле синапсов, их типы и кальций в них
	Axons     [16]Axon     //[]Координаты аксонов в общем файле синапсов, количество везикул и кальций
	//[]
}
//CalcCharge - функция вычисления заряда клетки в средне-стабильной среде
func (n *Neuron) CalcCharge() int{
	ret:=int((float32(n.Chemic.Na) - float32(NAORG) + float32(n.Chemic.K) - float32(KORG)-130)/1.7)
	return ret
}

func (n *Neuron) Gradient(){
	if n.Chemic.Na>NAORG{
		n.Chemic.Na-=1
	}
	if n.Chemic.Na<NAORG{
		n.Chemic.Na+=1
	}
	if n.Chemic.K>KORG{
		n.Chemic.K-=1
	}
	if n.Chemic.K<KORG{
		n.Chemic.K+=1
	}
}

//DendrCharge - Заряд на всех дендритах
func (n *Neuron) DendrCharge() int{
	//пробегаемся по дендритам и вычисляем их сумарное зарядное влияние на сому (среднее арифм или геом)
	k:=1
	chdend:=0
	for i:=0;i<16;i++{
		if n.D&(1<<i) != 0 { //проверяем, что данный дендрит включен
			k++
			chdend+=int(n.Dendrites[i].Charge)
		}
	}
	chdend=chdend/k //k будет больше на 1 от реально существующих включеных дендритов норм? todo
	if chdend>MAXALLDENDRCHARGE {//ограничение по максимуму влияния дендритов
		return MAXALLDENDRCHARGE
	}
	if chdend < -MAXALLDENDRCHARGE{
		return -MAXALLDENDRCHARGE
	}
	return chdend
}

//DoDendrites -
func (n *Neuron) DoDendrites(gene *GenNeuron){
	somacharge:=n.CalcCharge()
	for i:=0;i<16;i++{
		if n.D&(1<<i) != 0 { //проверяем, что данный дендрит включен
			if somacharge>MINSOMACHARGE{//если заряд сомы больше MINSOMACHARGE, значит она делает ПД, и ее кальциевые каналы открыты
				n.Dendrites[i].Ca=MAXCASOMA //и значит кальций входит в дендрит через сому
			}else if n.Dendrites[i].Ca>1 {//кальций выходит из каналов
				n.Dendrites[i].Ca-=1
				if n.Dendrites[i].Ca>MAXCA*2{
					n.Dendrites[i].Ca=byte((int16(n.Dendrites[i].Ca)*7 + int16(MINCA))/8)
				}
			}
			//сила рецепторов
			Power:=byte(0xf0 & n.Dendrites[i].Typed)

			//чем сильней рецептор, тем быстрее выходит кальций
			if n.Dendrites[i].Ca> Power/0x10{
				n.Dendrites[i].Ca-=Power/0x10
			}else {
				n.Dendrites[i].Ca=1
			}
			var medi byte =0 //величина медиатора
			switch n.Dendrites[i].Typed {
			case DENDRGABAION, DENDRGABAMETA:
				//todo
			case DENDRACHION, DENDRACHMETA:
				medi=org.synapsesMap[n.SynNumber].syn[n.Dendrites[i].N].ACh
			}
			//если мы в долгосрочной депрессии, мы пропускаем обработку синапса совсем, аксону это не понравится и он отключится (или нет - это нне наши проблемы)
			if n.Dendrites[i].State<200 && medi>0{ //мы не в долгосрочной депресии и есть ли медиатор в синапсе?
				isBreakMedi:=false //тормозный ли медиатор
				//эстераза разрушает ацх, а если мы в депрессии - то не разрушаем
				switch n.Dendrites[i].Typed {
				case DENDRGABAION, DENDRGABAMETA:
					//todo
					isBreakMedi=true
				case DENDRACHION, DENDRACHMETA:
					org.synapsesMap[n.SynNumber].syn[n.Dendrites[i].N].CHOL=medi //+n.c.CHOL //здесь в реале надо прибавить существующий холин
					org.synapsesMap[n.SynNumber].syn[n.Dendrites[i].N].ACh=0
				}

				//достаточно ли медиатора, чтобы открыть ионные каналы (зависит также от силы дендрита)
				if medi*(Power/0x10) >= MINMEDIREACT  {
					//количество заряда (относительное). поскольку ацх - возбужд, заряд от нуля до +120
					ch := int(float32(Power)/10* float32(medi) / (float32(n.Dendrites[i].Ca) + 4))
					//ограничение заряда 
					if ch > MAXDENDRCHARGE {
						n.Dendrites[i].Charge = int8(MAXDENDRCHARGE)
					} else {
						n.Dendrites[i].Charge = int8(ch)
					}
					if isBreakMedi{//если тормозный, то заряд жеж отрицательный
						n.Dendrites[i].Charge = -n.Dendrites[i].Charge
					}

				}else{
					n.Dendrites[i].Charge = 0
				}

			}else{//если медиатора нет или мы в ДД - наши каналы закрыты, заряд КП 0
				n.Dendrites[i].Charge = 0
			}

			//обучение проходит только если медиатор присутствовал в щели в минимальном количестве
			//иначе получим такую штуку, что мы приспосабливаемся не к работе пресинаптического нейрона
			if medi>=MINMEDI{
				//если в ДД, уменьшаем ожидание
				if n.Dendrites[i].State>200 {
					n.Dendrites[i].State--
				}else if n.Dendrites[i].State>102 && n.Dendrites[i].State<150{//если в ДП, тоже уменьшаем
					n.Dendrites[i].State--
				} else{
					//состояние обучения
					switch n.Dendrites[i].State{
					case 1: //обычная работа
						if n.Dendrites[i].Ca>MAXCA && medi>MAXMEDI{
							//здесь краткосрочная депрессия ставится на проверку
							n.Dendrites[i].State=10
						}else if n.Dendrites[i].Ca<MINCA && medi>MAXMEDI{
							//краткосрочная потенциация ставится на проверку
							n.Dendrites[i].State=100
						}
					case 9://конец краткосрочной депрессии
						if n.Dendrites[i].Ca>MAXCA && medi>MAXMEDI{
							//здесь краткосрочная депрессия продолжается
							n.Dendrites[i].State+=1
						}else {
							n.Dendrites[i].State=1//нормальный режим
						}
					case 10, 11, 12:
						if n.Dendrites[i].Ca>MAXCA && medi>MAXMEDI{
							//здесь краткосрочная депрессия продолжается
							n.Dendrites[i].State+=1
						}else {
							n.Dendrites[i].State-=1
						}
					case 13:
						if n.Dendrites[i].Ca>MAXCA && medi>MAXMEDI{
							//таки депрессия
							if Power == 0x10 {//уже некуда депрессировать, переходим в долговременную депрессию
								n.Dendrites[i].State=255
								break
							}else {
								Power=Power-0x10
								n.Dendrites[i].State-=1 //дадим шанс не выключить еще рецепторов на следующем шаге
							}
							n.Dendrites[i].Typed = DendriteTypeEnum(Power + (byte(n.Dendrites[i].Typed)&0x0f))
						}else {
							n.Dendrites[i].State-=1
						}
					case 100:
						if n.Dendrites[i].Ca<MINCA && medi>MAXMEDI{
							//краткосрочная потенциация продолжает проверку
							n.Dendrites[i].State=101
						} else{
							//показалось
							n.Dendrites[i].State=1
						}
					case 101:
						if n.Dendrites[i].Ca<MINCA && medi>MAXMEDI{
							//таки да, нейрон пресинаптический плюет когда надо, увеличиваем кол-во рецепторов
							if Power<0xf0{
								Power+=0x10
							} else {
								//итак максимально сильный дендрит... TODO!
							}
							n.Dendrites[i].Typed = DendriteTypeEnum(Power + (byte(n.Dendrites[i].Typed)&0x0f))
							n.Dendrites[i].State=110 //это же долговременная потенциация? доверяем пресинаптическому нейрону некоторое время
						} else{
							//показалось
							n.Dendrites[i].State=1
						}
					case 102:
						n.Dendrites[i].State=1 //ДП (долговременная потенциация) окончена, обычная работа
					case 200: //долговременная депрессия (ДД) конец, но мы идем в краткосрочную
						n.Dendrites[i].State=9

					default:
						if n.Dendrites[i].Ca>MAXCA && medi>MAXMEDI{
							//здесь краткосрочная депрессия ставится на проверку
							n.Dendrites[i].State=10
						}else if n.Dendrites[i].Ca<MINCA && medi>MAXMEDI{
							//краткосрочная потенциация ставится на проверку
							n.Dendrites[i].State=100
						}
						n.Dendrites[i].State=1
					}
				}
				//поскольку в синапсе есть медиатор, он открыл каналы кальция, и кальций входит в дендрит из синапса
				//и в принципе, может немного выйти, если нейромедиатора мало, а кальция много
				n.Dendrites[i].Ca=byte((int(n.Dendrites[i].Ca)*2+int(medi))/3)
			}
			
		}
	}
}

func (n *Neuron) NaOpened(){
	if n.Chemic.Na < 150{
		n.Chemic.Na+=90
	}else if n.Chemic.Na < 210{
		n.Chemic.Na+=40
	}else if n.Chemic.Na < 240{
		n.Chemic.Na+=20
	}else if n.Chemic.Na < 250{
		n.Chemic.Na+=3
	}
}
func (n *Neuron) KOpened(){
	if n.Chemic.K>150{
		n.Chemic.K-=40
	}else if n.Chemic.K>75{
		n.Chemic.K-=25
	}else if n.Chemic.K>20{
		n.Chemic.K-=8
	}else if n.Chemic.K>5{
		n.Chemic.K-=3
	}
}

//DoLiveCicle - главная функция жизни нейрона
func (n *Neuron) DoLiveCicle(gene *GenNeuron){
	charge:=n.CalcCharge()
	//качаем насосом Na-K
	if n.CalcCharge()>CHARGENORM-20{//если уже почти норма, то нефиг качать, оно через каналы выйдет
		//в геноме написано, сколько раз за цикл делаем дополнительно или уменьшительно
		n.Chemic.DoNaKATPasa()//один раз по-любому
		for i:=gene.AddNaKATPasa+ADDNAKATPASA;i>0;i++ {
			n.Chemic.DoNaKATPasa()
		}
	}
	//общий заряд
	charge+=n.DendrCharge()

	switch n.State {
	case 1://нормальная работа
		if charge<CHARGENORM{ //глубокая реполяризация, каналы открыты для выравнивания
			n.Gradient()
		}else if charge>NACHANOPEN{
			n.State=10
			n.NaOpened()
			charge=n.CalcCharge()
			if charge>=NACHANCLOSE{
				n.State=20
			}
		}
	case 10://начало деполяризации
		if charge>=NACHANOPEN && charge<=NACHANCLOSE{
			n.NaOpened()
			charge=n.CalcCharge()
		}
		if charge>=KCHANOPEN {
			n.KOpened()
			charge=n.CalcCharge()

			if charge> NACHANOPEN && n.Chemic.Na<NAVALCHANREOPEN && n.Chemic.K>KVALCHANREOPEN{
				n.State=1
			}
		}
		if charge>=NACHANCLOSE{
			n.State=20
		}

	case 20://только калиевый ток
		n.KOpened()
		charge=n.CalcCharge()
		if charge<KCHANCLOSE {
			n.State=1
		}
		if charge> NACHANOPEN && n.Chemic.Na<NAVALCHANREOPEN && n.Chemic.K>KVALCHANREOPEN{
			n.State=10
		}

	default:
		n.State=1

	}
	if rand.Intn(100)>10{
		n.Gradient()
	}
}

func (n *Neuron) DoAxons(gene *GenNeuron){

}