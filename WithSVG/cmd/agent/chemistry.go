package agent

import (
	"math/rand"
)

/*
Межклеточная жидкость мышцы лягушки содержит 120 ммол/л натрия и 2,5 ммол/ли калия, а внутри клеток
9,2 ммол/л и 140 ммол/л соответственно
*/

//Chemical - структура веществ всякого синапса и клетки
type Chemical struct {
	/*
		НЕЛЬЗЯ МЕНЯТЬ ПОРЯДОК СЛЕДОВАНИЯ ПОЛЕЙ!!! И ДОБАВЛЯТЬ, КАК-ТО ИЗМЕНЯТЬ ЭТО!!! ИСПОЛЬЗУЕТСЯ unsafe!!! для маппинга структуры на часть файла
	*/
	GLUC  uint16 //[0]Глюкоза (энергетически важная молекула, участвует во многих процессах клетки)
	O2    uint16 //[2]окислитель глюкозы, без него процессы генерации веществ из глюкозы не идут
	WASTE uint16 //[4]всякая кака, которую организм выводит
	OMEGA uint16 //[6]жирная кислота, использующаяся для синтеза мембран всех клеток. Из нее также в некоторых клетках синтезируется анандамид, который очень быстро разлагается на этаноламин и арахидоновую кислоту. Из этаноламина получается холин.
	K     byte   //[8]ион калия, выводит каку, создает ток разряда, перекачивается Na-K-ATPasa
	Na    byte   //[9]ион натрия, вводит нямку, создает ток разряда, перекачивается Na-K-ATPasa
	CO    byte   //[10](рассеивается при высвобождении по всем близлешащим клеткам), связан с окислением глюкозы, выделяется на дендритах в ответ на открытие на них ионных каналов, в основном в дофаминовых и серотониновых
	NO    byte   //[11](рассеивается при высвобождении по всем близлешащим клеткам), связан с глутаматом
	ACh   byte   //[12]Ацетилхолин (в основном в перефирических файлах входных и выходных устройств), синтезируется из глюкозы и холина
	CHOL  byte   //[13]Холин. Синтезируется из анандамида или OMEGA. Нужен для работы ацетилзолиновых синапсов
	AA    byte   //[14]анандамид, медиатор, действующий на аксон ГАМК возбуждающе, а на аксон глютамата торможающе)
	GLU   byte   //[15]Глютамат (главный медиатор, в информационной модели он возбуждающий, но мы не будем болше говорить о возбуждении и торможении, а только о различных состояниях жизни клетки), синтезируется из глутамина снаружи клетки или из аспартата внутри
	GABA  byte   //[16]ГАМК, синтезируется из глютамата, быстро распадается, переходя в глюкозу
	GLN   byte   //[17]ГЛУТАМИН - синтезируется из аспартата снаружи, выводится в некотором количестве, участвует в цикле Глутамат-захвата. Глутаматные нейроны синтезируют глутамат только из глутамина. Глутамат в синапсе превращается в глутамин.
	NE    byte   //[18]Норадреналин, синтезируется из дофамина. Всасывается обратно и частично выводится.
	DOP   byte   //[19]Дофамин, синтезируется из тирозина. Всасывается обратно в аксонах и частично выводится. Из него синтезируется норадреналин.
	SER   byte   //[20]Серотонин. Синтезируется из триптофана.  Всасывается обратно и частично выводится.
	TYR   byte   //[21]Тирозин - для синтеза дофамина
	TRP   byte   //[22]Триптофан - для синтеза серотонина.
	ASP   byte   //[23]Аспартат - для синтеза глутамата
}

//Unwaste - функция, которую может вызвать только организм, очищает межклеточную среду
func (c *Chemical) Unwaste() {
	c.WASTE = c.WASTE / 2
}

//UnCONO - функция, которую может вызвать только организм, очищает межклеточную среду от газов
func (c *Chemical) UnCONO() {
	c.NO = c.NO / 5 * 4
	c.CO = c.CO / 5 * 4
}

//AddGluckose - может вызвать только организм, добавляет глюкозы, в зависимости от загрязненности и количества глюкозы
func (c *Chemical) AddGluckose() {
	gg := c.GLUC + 100 - c.WASTE/50 - c.GLUC/70 //обратная функция от загрязненности и количества уже имеющейся глюкозы
	if gg > c.GLUC {
		c.GLUC = gg
	}
}

//AddOxygen - может вызвать только организм, добавляет ксилорода, в зависимости от загрязненности и количества кислорода
func (c *Chemical) AddOxygen() {
	gg := c.O2 + 100 - c.WASTE/70 - c.O2/90 //обратная функция от загрязненности и количества уже имеющейся кислорода
	if gg > c.O2 {
		c.O2 = gg
	}
}

//AddOmega - может вызвать только организм, добавляет жироной кислоты, в зависимости от загрязненности и количества ее
func (c *Chemical) AddOmega() {
	gg := c.OMEGA + 100 - c.WASTE/70 - c.OMEGA/90 //обратная функция от загрязненности и количества ее имеющейся
	if gg > c.OMEGA {
		c.OMEGA = gg
	}
}

//DoKrebs - может вызвать только клетка, пожирание внутренней глюкозы и кислорода для нужд клетки. Если удалось завершить цикл Кребса, возвращает true
func (c *Chemical) DoKrebs() bool {
	/*
		Это основной цикл клетки, от которого зависит работа Na-K-ATP-asa
		Если он не может правильно завершиться, клетка не сможет качать Na и K
	*/
	if c.GLUC == 0 {
		//Если нет глюкозы, попробуем получить ее из OMEGA. Но в этом случае клетка пропускает цикл насоса все равно
		c.GetGlucFromOmega()
		return false
	} else if c.O2 == 0 {
		//совсем плохо дело, ну или нет - если клетка уже заряжена
		return false
	} else if c.WASTE >= 0xfff0 {
		//ох и грязный! Грязь не увеличиваем, но выполняем цикл с 30% вероятностью
		if rand.Intn(100) > 30 {
			return false
		}
		c.WASTE = c.WASTE - 1 //здесь уменьшили, чтобы она не изменилась при вычислении
	}

	c.GLUC = c.GLUC - 1
	c.O2 = c.O2 - 1
	c.WASTE = c.WASTE + 1
	return true
}

//GetGlucFromOmega - может вызвать только клетка (любая), в случае нехватки глюкозы, производит глюкозу из жира
func (c *Chemical) GetGlucFromOmega() bool {
	if c.GLUC > 30 { //в клетке уже есть глюкоза, чего надо еще?
		return false
	}
	if c.OMEGA == 0 || c.O2 < 11 || c.WASTE >= 0xfff0 {
		return false
	}
	c.OMEGA = c.OMEGA - 1
	c.O2 = c.O2 - 11
	c.WASTE = c.WASTE + 1
	c.GLUC = c.GLUC + 3
	return true

}

//GetGlucFromGABA - может вызвать только клетка GABA-эргическая, в случае нехватки глюкозы или скопления GABA
func (c *Chemical) GetGlucFromGABA() bool {
	if c.GLUC > 0xa000 { //в клетке слишком много глюкозы
		return false
	}
	if c.GABA < 2 || c.O2 < 4 || c.WASTE >= 0xfff0 || c.NO > 0xfb {
		return false
	}
	c.GABA = c.GABA - 2
	c.O2 = c.O2 - 4
	c.WASTE = c.WASTE + 1
	c.NO = c.NO + 2
	c.GLUC = c.GLUC + 1
	return true

}

//AChSynt - может вызвать только клетка АЦХ-эргическая, синтезирует внутри ацетилхолин, если цикл Кребса рабочий (там из этого цикла вещества берутся)
func (c *Chemical) AChSynt() bool {
	if !c.DoKrebs() {
		return false
	}
	if c.CHOL < 1 || c.ACh > 0xfe {
		return false
	}
	c.CHOL = c.CHOL - 1
	c.ACh = c.ACh + 1
	return true
}

//DoAChEstarasa - может вызвать только дендрит клетки АЦХ-эргической, разрушает ацетилхолин в синапсе
func (c *Chemical) DoAChEstarasa() bool {
	if c.ACh < 1 || c.CHOL > 0xfe || c.WASTE >= 0xfff0 {
		return false
	}

	c.CHOL = c.CHOL + 1
	c.ACh = c.ACh - 1
	c.WASTE = c.WASTE + 1
	return true
}

//DoAAtoCHOLEstarasa - может вызвать только аксон клетки АЦХ-эргической, разрушает Анандамид в синапсе, под действием кислорода
func (c *Chemical) DoAAtoCHOLEstarasa() bool {
	if c.AA < 1 || c.CHOL > 0xf0 || c.O2 < 15 || c.CO > 230 {
		return false
	}

	c.CHOL = c.CHOL + 1
	c.AA = c.AA - 1
	c.O2 = c.O2 - 15
	c.CO = c.CO + 17
	return true
}

//DoAAtoGLNEstarasa - может вызвать только аксон клетки ГАМК- или Глутамат-эргической, разрушает Анандамид в синапсе, под действием кислорода
func (c *Chemical) DoAAtoGLNEstarasa() bool {
	if c.AA < 1 || c.GLN > 0xf0 || c.O2 < 18 || c.CO > 230 {
		return false
	}

	c.GLN = c.GLN + 1
	c.AA = c.AA - 1
	c.O2 = c.O2 - 18
	c.CO = c.CO + 17
	return true
}

//DoNaKATPasa - может вызвать только клетка, перекачка натрия и калия
func (c *Chemical) DoNaKATPasa(to *Chemical) bool {
	/*
		Очень важный цикл перекачки Натрия наружу и калия внутрь
	*/
	if !c.DoKrebs() { //этот насос - самый чуствительный к наличию АТФ, образующейся в митохондриях
		return false
	}
	if to.Na > 0xfc || c.Na < 3 || to.K < 2 || c.K > 0xfe || c.WASTE >= 0xfff0 || to.WASTE >= 0xfff0 {
		return false
	}
	c.Na = c.Na - 3
	to.Na = to.Na + 3
	to.K = to.K - 2
	c.K = c.K + 2
	return true
}

//MakeMemrane - может вызвать только клетка, сколько то там раз за цикл. Можно случайно. Типа функция починки мембраны
func (c *Chemical) MakeMemrane() bool {
	if c.OMEGA < 1 || c.O2 < 3 || c.WASTE >= 0xfff0 {
		return false
	}

	c.OMEGA = c.OMEGA - 1
	c.O2 = c.O2 - 3
	c.WASTE = c.WASTE + 1
	return true
}

//AASynt - могут вызвать дендриты клеток в ответ на переизбыток активности клетки-передатчика, синтезирует внутри клетки анандамид, и сразу
//выбрасывается дендритом
func (c *Chemical) AASynt() bool {
	if c.OMEGA < 1 || c.NO < 1 || c.CO < 2 || c.WASTE >= 0xfffe {
		return false
	}
	c.OMEGA = c.OMEGA - 1
	c.NO = c.NO - 1
	c.CO = c.CO - 2
	c.WASTE = c.WASTE + 1
	return true
}

//GLNSyntASP - глютамат и ГАМК эргические клетки могут синтезировать глутамин из аспартата
func (c *Chemical) GLNSyntASP() bool {
	if c.ASP < 1 || c.GLUC < 1 || c.NO < 1 || c.WASTE >= 0xfffd || c.GLN > 0xf0 {
		return false
	}
	c.ASP = c.ASP - 1
	c.GLUC = c.GLUC - 1
	c.NO = c.NO - 21
	c.WASTE = c.WASTE + 1
	c.GLN = c.GLN + 1
	return true
}
