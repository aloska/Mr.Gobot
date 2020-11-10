package agent

/*CoreEnum - константы номеров ядра
 */
type CoreEnum uint16

const (
	//COREMY - аксоны направлены в то же поле, в котором находится клетка
	COREMY CoreEnum = 0xffff
	//COREINPUTS - аксоны направлены в поле чувств
	COREINPUTS CoreEnum = 0xfffe
	//COREOUTPUTS - аксоны направлены в поле действий
	COREOUTPUTS CoreEnum = 0xfffd
)
