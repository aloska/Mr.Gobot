package agent

/*CoreEnum - константы номеров ядра
 */
type CoreEnum uint16

const (
	//COREMY - аксоны направлены в то же поле, в котором находится клетка
	COREMY CoreEnum = 0xffff
	//COREINPUTS - аксоны направлены в общее поле чувств
	COREINPUTS CoreEnum = 0xfffe
	//COREOUTPUTS - аксоны направлены в общее поле действий
	COREOUTPUTS CoreEnum = 0xfffd
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
