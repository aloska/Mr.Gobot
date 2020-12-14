package universal

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Chromosome string

type Dyploid struct{
	M Gaploid
	F Gaploid
}

type Gaploid struct {
	Chromosome  //сырая хромосома
	Genes []string	//гены, готовые к работе
}

type Genotype []Dyploid

//NewGaploid - создает гаплоид из хромосомы (парсит хромосому и достает из нее гены)
func NewGaploid(chr string) (*Gaploid, error){
	g:=Gaploid{Chromosome: Chromosome(chr)}
	//сначала транскрипция
	//парсим ДНК от старт до стоп-кодонов
	reg:=regexp.MustCompile("[\u26a2-\u26a4](.+?)[\u2c00-\u2c5e]")
	matches:=reg.FindAllString(chr, -1)//matches содержит все сырые гены, вместе со старт и стоп-кодонами
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
		return &g,errors.New("there is no any gene in string")
	}
	return &g,nil
}

//смешивает два набора генов в шахматном порядке 
func (d* Dyploid) concatMF() []string{
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

//создание всех алгоритмов из генотипа
func (g Genotype) MakeAlgorithms() ([][]Command, []error){
	var res [][]Command //здесь будет len(g.Genom) алгоритмов
	var ers []error
	for i:=0; i<len(g);i++{
		//для каждого диплоидного набора свой алгоритм
		var a []Command //сюда сложим комманды алгоритма

		//смешиваем отцовские и материнские гены в одном алгоритме
		//но может оказаться, что количество генов не одинаковое, а в какой-то хромосоме совсем отсутствуют, поэтому спец-функция
		genes:=g[i].concatMF() //набор генов, определяющий алгоритм, но это не все)) может какой ген поломанный - выяснится на этапе трансляции
		for _, gen:= range genes{
			//трансляция
			if com,err:=GeneTranslation(gen); err==nil{
				a=append(a, com ...)//если ошибок трансляции нет - добавим к будущему алгоритму
			}else{
				ers=append(ers, err)
			}

		}
		//добавим в общую корзину
		if len(a)>0{
			res=append(res,a)
		}else{
			//ни один ген в хромосоме (мужской и женской) не справился Внимание todo - удалить диплоид из генотипа?
		}

	}
	return res, ers
}

//формат комманды (означает сколько байт должно быть в каждом операнде)
var COMMFORMAT = map[Comm]string{
	NOP: "0 0 0",

	ADD: "1 1 1", SUB: "1 1 1",MUL: "1 1 1",DIV: "1 1 1",REM: "1 1 1",AND: "1 1 1",OR: "1 1 1",XOR: "1 1 1",
	SLL: "1 1 1",SRL: "1 1 1",SEQ: "1 1 1",SGE: "1 1 1",SLT: "1 1 1",SNE: "1 1 1",LDMX: "1 1 1",LDINX: "1 1 1",
	STMX: "1 1 1",STOUTX: "1 1 1",

	ADDI: "1 1 8", SUBI: "1 1 8",MULI: "1 1 8",DIVI: "1 1 8",REMI: "1 1 8",ANDI: "1 1 8",ORI: "1 1 8",XORI: "1 1 8",
	SLLI: "1 1 8",SRLI: "1 1 8",BEQ: "1 1 8",BGE: "1 1 8",BLT: "1 1 8",BNE: "1 1 8",BLE: "1 1 8",BGT: "1 1 8",

	LDM: "1 8 8", LDIN:"1 8 8",

	STM: "8 8 1", STOUT:"8 8 1",

	JMP:"8 0 0", //только засунуть в третий операнд эти 8 байт!!

	PUSH:"1 0 0", POP:"1 0 0",

	LI:"1 8 0", //только засунуть в третий операнд эти 8 байт!!
}

//трансляция гена (вызывать только после транскрипции и сплайсинга!)
//во входной строке ожидаются только 1 и 2-байтовые руны
func GeneTranslation(gene string) (alg []Command, er error){
	r:=[]rune(gene)
	i:=0//индекс слайса рун из гена
	state:=0
	index:=-1 //индекс текущей (транслируемой с гена) команды процессора
	for i<len(r){
		switch state{
		case 0://ожидается код комманды
			alg=append(alg, Command{0,0,0,0})
			index++
			alg[index].Code=Comm(r[i])
			switch COMMFORMAT[alg[index].Code%COUNTCOMMAND] {//какой формат у команды?
			case "0 0 0":
				//операндов нет - ничего не делаем, просто переход к следующей комманде
			case "1 1 1", "1 1 8", "1 8 8", "1 0 0", "1 8 0":
				state=1
			case "8 8 1":
				state=10
			case "8 0 0":
				state=30	//потому что в третий операнд надо засунуть
			}
			i++
		case 1://ожидается 1 операнд byte (регистровый)
			alg[index].Op1=uint64(r[i])
			switch COMMFORMAT[alg[index].Code%COUNTCOMMAND] {
			case "1 0 0":
				//на этом все
			case "1 1 1", "1 1 8":
				state=2
			case "1 8 8":
				state=20
			case "1 8 0":
				state=30 //потому что в третий операнд надо засунуть
			}
			i++
		case 2://ожидается 2 операнд byte (регистровый)
			alg[index].Op2=uint64(r[i])
			switch COMMFORMAT[alg[index].Code%COUNTCOMMAND] {
			case "1 1 1":
				state=3
			case "1 1 8":
				state=30
			}
			i++
		case 3://ожидается 3 операнд byte (регистровый)
			alg[index].Op3=int64(r[i])//здесь не будет отрицательных чисел, просто в 3 операнде они бывают, см. case 30
			state=0
			i++
		case 10://ожидается 1 операнд uint64, нужно 8 рун для парсинга
			if i+7>=len(r){//нам надо 8 рун, а осталось меньше 8
				return nil, errors.New(fmt.Sprintf( "gene broken, state: %v, runen: %v", state,i))
			}
			//у нас все LittleEndian
			alg[index].Op1=uint64(r[i]&0xff) | uint64(r[i+1]&0xff)<<8 | uint64(r[i+2]&0xff)<<16 | uint64(r[i+3]&0xff)<<24 |
				uint64(r[i+4]&0xff)<<32 | uint64(r[i+5]&0xff)<<40 | uint64(r[i+6]&0xff)<<48 | uint64(r[i+7]&0xff)<<56
			state=20
			i+=8
		case 20://ожидается 2 операнд uint64
			if i+7>=len(r){
				return nil, errors.New(fmt.Sprintf( "gene broken, state: %v, runen: %v", state,i))
			}
			alg[index].Op2=uint64(r[i]&0xff) | uint64(r[i+1]&0xff)<<8 | uint64(r[i+2]&0xff)<<16 | uint64(r[i+3]&0xff)<<24 |
				uint64(r[i+4]&0xff)<<32 | uint64(r[i+5]&0xff)<<40 | uint64(r[i+6]&0xff)<<48 | uint64(r[i+7]&0xff)<<56
			switch COMMFORMAT[alg[index].Code%COUNTCOMMAND] {//какой формат у команды?
			case "8 8 1":
				state=3
			case "1 8 8":
				state=30
			}
			i+=8
		case 30://ожидается 3 операнд int64
			if i+7>=len(r){
				return nil, errors.New(fmt.Sprintf( "gene broken, state: %v, runen: %v", state,i))
			}
			alg[index].Op3=int64(r[i]&0xff) | int64(r[i+1]&0xff)<<8 | int64(r[i+2]&0xff)<<16 | int64(r[i+3]&0xff)<<24 |
				int64(r[i+4]&0xff)<<32 | int64(r[i+5]&0xff)<<40 | int64(r[i+6]&0xff)<<48 | int64(r[i+7]&0xff)<<56
			state=0
			i+=8
		}
	}
	if state!=0{ //если мы вышли из цикла не в ожидании кода команды, значит ген поломаный
		return nil, errors.New(fmt.Sprintf( "gene broken because end unexpectedly, runen: %v, byten: %v", state,i))
	}
	return
}