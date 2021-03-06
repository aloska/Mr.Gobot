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
	//DENDRACHETA - ацх-метаб
	DENDRACHETA DendriteTypeEnum = 0x15
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
	//RECEPTORDATAUINT32BIGPOS - позитивный рецептор, нацеленный на беззнаковое целое 4-байтовое
	RECEPTORDATAUINT32BIGPOS ReceptorTypeEnum = 12
	//RECEPTORDATAUINT32BIGNEG - отрицательный рецептор, нацеленный на беззнаковое целое 4-байтовое
	RECEPTORDATAUINT32BIGNEG ReceptorTypeEnum = 13
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
