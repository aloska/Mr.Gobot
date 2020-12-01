package agent

/*номер синаптического поля (ядра или входа или выхода), в синапсах которого находятся аксоны или дендриты

зарегестрированные константы не отражают реальный номер поля, а просто являются как бы ссылками, говорящими
программисту, что эти отростки относятся к неким стандартным полям

0xffff - отростки в поле синапсов вегетативной системы

0xfffe - отростки расположены в общем синаптическом файле всех входов мозга (для реализации связей между разными входами)
Поскольку нейроны могут быть расположены не только в мозге, но и в входных устройствах, у всех входных устройств может быть свое
синаптическое поле, и есть общее синаптическое поле (в простых агентах это может не использоваться, но возможность такая есть)

0
*/
type SynEnum uint16

const (
	//SYNVEGETATIC - синаптическое поле вегетативной системы
	SYNVEGETATIC SynEnum = 0xffff
	//SYNINPUTS - общее поле чувств
	SYNINPUTS SynEnum = 0xfffe
	//SYNOUTPUTS - общее поле действий
	SYNOUTPUTS SynEnum = 0xfffd
)

/*NeuronTypeEnum - константы типов нейронов
 */
type NeuronTypeEnum byte

const (
	//NEURONSTEM - стволовая клетка
	NEURONSTEM NeuronTypeEnum = 0x10
	//NEURONACETILHOLIN - ацетилхолиновый
	NEURONACETILHOLIN NeuronTypeEnum = 0x21
)

/*DendriteTypeEnum - константы типов дендрита
0x01-0x0f - резерв
0x1E-резерв
0x1F-резерв
*/
type DendriteTypeEnum byte

const (
	//DENDRONTO - никакой тип пока, онтогенез
	DENDRONTO DendriteTypeEnum = 0
	//DENDRGLUION - глютамат-ионный
	DENDRGLUION DendriteTypeEnum = 0x10
	//DENDRGLUMETA - глютамат-метаб
	DENDRGLUMETA DendriteTypeEnum = 0x11
	//DENDRGABAION - ГАМК-ионный
	DENDRGABAION DendriteTypeEnum = 0x12
	//DENDRGABAMETA - ГАМК-метаб
	DENDRGABAMETA DendriteTypeEnum = 0x13
	//DENDRACHION - ацх-ионный
	DENDRACHION DendriteTypeEnum = 0x14
	//DENDRACHMETA - ацх-метаб
	DENDRACHMETA DendriteTypeEnum = 0x15
	//DENDRAAION - АА-ионный
	DENDRAAION DendriteTypeEnum = 0x16
	//DENDRAAMETA - АА-метаб
	DENDRAAMETA DendriteTypeEnum = 0x17
	//DENDRNEION - NE-ионный
	DENDRNEION DendriteTypeEnum = 0x18
	//DENDRNEMETA - NE-метаб
	DENDRNEMETA DendriteTypeEnum = 0x19
	//DENDRDOPION - DOP-ионный
	DENDRDOPION DendriteTypeEnum = 0x1a
	//DENDRDOPMETA - DOP-метаб
	DENDRDOPMETA DendriteTypeEnum = 0x1b
	//DENDRSERION - SER-ионный
	DENDRSERION DendriteTypeEnum = 0x1c
	//DENDRSERMETA - SER-метаб
	DENDRSERMETA DendriteTypeEnum = 0x1d
)

/*DataTypeEnum - константы типов ячеек с данными входных и выходных файлов
 */
type DataTypeEnum uint16

const (
	/*DATAUINT32BIG - вот такой тип:
	type DataUInt32 struct {
		data [8]uint32
	}
	*/
	DATAUINT32BIG DataTypeEnum = 6

	//DATAUINT32 - uint32
	DATAUINT32 DataTypeEnum=7
)

/*ReceptorTypeEnum - константы типов рецепторов
 */
type ReceptorTypeEnum uint16

const (
	/*RECEPTORDATAUINT32BIGBIT - рецептор, нацеленный на беззнаковое целое 4-байтовое (1 из 8), и в этом числе просматривает каждый бит,
	плювая вещества в соответствующий аксон, если бит=1
	type DataUInt32 struct {
		Data [8]uint32  - нацелен на какое-то 1 из 8
	}
	Ndata - указывает на номер ячейки
	NdataW - указывает на номер в массиве Data[8]
	if D[Ndata].Data[NdataW] & 0x1 == 1 и т.д.
	*/
	RECEPTORDATAUINT32BIGBIT ReceptorTypeEnum = 12

	/*RECEPTORDATAUINT32BIGPOS - позитивный рецептор, нацеленный на беззнаковое целое 4-байтовое (1 из 8), просто плюет пропорционально числу в ячейке данных
	type DataUInt32 struct {
		Data [8]uint32  - нацелен на какое-то 1 из 8
	}
	Ndata - указывает на номер ячейки
	NdataW - указывает на номер в массиве Data[8]
	D[Ndata].Data[NdataW]
	 */
	RECEPTORDATAUINT32BIGPOS ReceptorTypeEnum = 14
	//RECEPTORDATAUINT32BIGNEG - отрицательный рецептор, нацеленный на беззнаковое целое 4-байтовое
	RECEPTORDATAUINT32BIGNEG ReceptorTypeEnum = 15


	/*RECEPTORDATABYTE11122234 - TODO рецептор, нацеленный на байт, специального типа, хорошо подходит для аудио и многих других типов
	работа которого выглядит так:
	Ndata - указывает на номер ячейки
	NdataW - указывает на байт в ячейке
	анализируемое число беззнаковое целое 32 байтовое:
	byte = D[Ndata].Data[NdataW]
	в byte рецептор анализирует каждый бит, и плюет в аксоны так:

	номер |
	бита  | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 |
	------+---+---+---+---+---+---+---+---+
	номер | 0 | 1 | 2 | 3 | 5 | 7 | 9 | 12|
	аксона|   |   |   | 4 | 6 | 8 | 10| 13|
	  на  |   |   |   |   |   |   | 11| 14|
	   0  |   |   |   |   |   |   |   | 15|
	------+---+---+---+---+---+---+---+---+
	номер | 16| 17| 18| 19| 21| 23| 25| 28|
	аксона|   |   |   | 20| 22| 24| 26| 29|
	  на  |   |   |   |   |   |   | 27| 30|
	   1  |   |   |   |   |   |   |   | 31|
	------+---+---+---+---+---+---+---+---+

	т.е. младшие биты влияют на меньшее количество аксонов, а старшие на большее, приблизительно экспоненциально
	Итого, 16 аксонов плюют вещества в случае 1 в своем бите,
	и 16 других аксонов плюют в случае 0 в своем ьите
	 */
)

/*PreffectorTypeEnum - константы типов преффекторов
*/
type PreffectorTypeEnum uint16

const (
	//PREFFECTORUINT32POS - позитивный преффектор, нацеленный на беззнаковое целое 4-байтовое
	PREFFECTORUINT32POS PreffectorTypeEnum=12
	//PREFFECTORUINT32NEG - отрицательный преффектор, нацеленный на беззнаковое целое 4-байтовое
	PREFFECTORUINT32NEG PreffectorTypeEnum=13
)


const (
	NACHANOPEN int = -65		//уровень заряда открытия Натриевого канала
	NACHANCLOSE int = 15		//уровень заряда закрытия Натриевого канала
	KCHANOPEN int = 20			//уровень заряда открытия Калиевого канала
	KCHANCLOSE int = -55		//уровень заряда закрытия Калиевого канала
	CHARGENORM int = -75		//нормальный (базовый) уровень заряда

	ADDNAKATPASA int8 = 10      //сколько раз за цикл делаем насосом

	NAVALCHANREOPEN byte = 20 	//уровень натрия в клетке, ниже которого возможна повторная реполяризация в том же цикле
	KVALCHANREOPEN byte = 30    //уровень калия в клетке, выше которого возможна повторная реполяризация в том же цикле

	NAORG byte =240   			//уровень натрия в межклеточной жидкости
	KORG byte =10				//уровень калия в межклеточной жидкости

	MINSOMACHARGE=-10           //минимальный заряд сомы, ниже которого кальциевые каналы сомы не открываются, и кальций не попадает из сомы в дендриты
	MAXCA byte = 30				//максимальное количество кальция, выше которого ставится под вопрос депрессия
	MINCA byte = 3				//минимальное количество кальция, ниже которого ставится под вопрос потенциация
	MAXCASOMA byte = 100 		//количество кальция, попадающего в дендрит из сомы, во время ПД сомы
	MAXMEDI byte = 30			//максимальное количество медиатора, выше которого ставится под вопрос депрессия или потенциация
	MINMEDI byte = 5				//минимальное количество медиатора, выше которого алгоритм обучения продолжает работу
	MINMEDIREACT byte = 10		//минимальное количество медиатора, способное открыть хотя бы один канал кальция

	MAXALLDENDRCHARGE int = 40		//максимальный суммарный заряд всех дендритов
	MAXDENDRCHARGE int = 120   //максимальный заряд дендрита

	MAXAXDIMEDIFORCE byte=30	//кличество предшественника медиатора в синапсе, при котором аксон запускаеп программу увеличивания силы выплевывания
)