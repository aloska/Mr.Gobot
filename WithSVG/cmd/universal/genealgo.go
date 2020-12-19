package universal

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

//создание всех алгоритмов из генотипа
//каждый алгоритм состоит из всех генов Пароида сразу
func (g Genotype) MakeAlgorithms() ([][]Command, []error){
	var res [][]Command //здесь будет len(g.Genom) алгоритмов
	var ers []error
	for i:=0; i<len(g);i++{
		//для каждого пароидного набора свой алгоритм
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
//https://play.golang.org/p/FLIZoxSxzSb  - тест здесь
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
			alg[index].Code=Comm(r[i])%COUNTCOMMAND
			switch COMMFORMAT[alg[index].Code] {//какой формат у команды?
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
			if unicode.IsDigit(r[i]){ //цифры означают цифры
				alg[index].Op1=uint64(r[i]-'0')
			}else{
				alg[index].Op1=uint64(r[i]%32)
			}
			switch COMMFORMAT[alg[index].Code] {
			case "1 0 0":
				//на этом все
				state=0
			case "1 1 1", "1 1 8":
				state=2
			case "1 8 8":
				state=20
			case "1 8 0":
				state=30 //потому что в третий операнд надо засунуть
			}
			i++
		case 2://ожидается 2 операнд byte (регистровый)
			if unicode.IsDigit(r[i]){ //цифры означают цифры
				alg[index].Op2=uint64(r[i]-'0')
			}else{
				alg[index].Op2=uint64(r[i]%32)
			}
			switch COMMFORMAT[alg[index].Code] {
			case "1 1 1":
				state=3
			case "1 1 8":
				state=30
			}
			i++
		case 3://ожидается 3 операнд byte (регистровый)
			if unicode.IsDigit(r[i]){ //цифры означают цифры
				alg[index].Op3=int64(r[i]-'0')
			}else {
				alg[index].Op3=int64(r[i]%32)
			}
			//здесь не будет отрицательных чисел, просто в 3 операнде они бывают, см. case 30
			state=0
			i++
		case 10://ожидается 1 операнд uint64, нужно 8 рун для парсинга, из каждой возьмутся только младшие байты
			if unicode.IsDigit(r[i]){ //цифры означают цифры, их только 1 руна нужна
				alg[index].Op1=uint64(r[i]-'0')
				i++
			}else{
				if i+7>=len(r){//нам надо 8 рун, а осталось меньше 8
					return nil, errors.New(fmt.Sprintf( "gene broken, state: %v, runen: %v", state,i))
				}
				//у нас все LittleEndian
				alg[index].Op1=uint64(r[i]&0xff) | uint64(r[i+1]&0xff)<<8 | uint64(r[i+2]&0xff)<<16 | uint64(r[i+3]&0xff)<<24 |
					uint64(r[i+4]&0xff)<<32 | uint64(r[i+5]&0xff)<<40 | uint64(r[i+6]&0xff)<<48 | uint64(r[i+7]&0xff)<<56
				i+=8
			}
			state=20

		case 20://ожидается 2 операнд uint64
			if unicode.IsDigit(r[i]){ //цифры означают цифры
				alg[index].Op2=uint64(r[i]-'0')
				i++
			}else{
				if i+7>=len(r){
					return nil, errors.New(fmt.Sprintf( "gene broken, state: %v, runen: %v", state,i))
				}
				alg[index].Op2=uint64(r[i]&0xff) | uint64(r[i+1]&0xff)<<8 | uint64(r[i+2]&0xff)<<16 | uint64(r[i+3]&0xff)<<24 |
					uint64(r[i+4]&0xff)<<32 | uint64(r[i+5]&0xff)<<40 | uint64(r[i+6]&0xff)<<48 | uint64(r[i+7]&0xff)<<56
				i+=8
			}
			switch COMMFORMAT[alg[index].Code%COUNTCOMMAND] {//какой формат у команды?
			case "8 8 1":
				state=3
			case "1 8 8":
				state=30
			}

		case 30://ожидается 3 операнд int64
			if unicode.IsDigit(r[i]){ //цифры означают цифры
				alg[index].Op3=int64(r[i]-'0')
				i++
			}else{
				if i+7>=len(r){
					return nil, errors.New(fmt.Sprintf( "gene broken, state: %v, runen: %v [-%v]", state,i,i+8-len(r)))
				}
				alg[index].Op3=int64(r[i]&0xff) | int64(r[i+1]&0xff)<<8 | int64(r[i+2]&0xff)<<16 | int64(r[i+3]&0xff)<<24 |
					int64(r[i+4]&0xff)<<32 | int64(r[i+5]&0xff)<<40 | int64(r[i+6]&0xff)<<48 | int64(r[i+7]&0xff)<<56
				i+=8
			}
			state=0

		}
	}
	if state!=0{ //если мы вышли из цикла не в ожидании кода команды, значит ген поломаный
		return nil, errors.New(fmt.Sprintf( "gene broken because end unexpectedly, state: %v, runen: %v", state,i))
	}
	return
}

/*создает Pairoid из набора алгоритмов (сразу с генами и т.п.), таким образом, что
- один алгоритм - один ген
- между коммандами помещаются интроны (только 3 и 4 байтовые UTF-8, одинаковые для M и F)
- между генами помещаются подписи (любые символы UTF-8, могут быть разными для M и F)
Материнская половинка хромосомы отличается от отцовской лишь незначительными интронами и некодирующими участками
В scribes должны быть переданы все промежутки между генами (типа подписи, комментарии...), и они могут быть пустыми "", но не nil!
В introns должны быть переданы все желаемые интроны
При необходимости получения чистых хромосом (только гены без интронов и промежутков между генами), можно генерить
scribes и introns функциями MakeEmptyScribes и MakeEmptyIntrons
Внимание - интроны должны быть 3 или 4 байтовыми, иначе они будут восприниматься как экзоны
Чтобы писать комментарии внутри генов интронами - можно пользоваться ASCII-зеркалом, начинающимся в UTF-8 с (ef bc 90)
*/
func MakePairoidFromAlgs(algs [][]Command, introns [][]string, Mscribes, Fscribes []string) (Pairoid, error){
	var P Pairoid
	lena:=len(algs)
	if lena<1 || lena!=len(introns) || lena!=len(Mscribes) || lena!=len(Fscribes){
		return P, errors.New("MakePairoidFromAlgs: bad lenght parameter(s)")
	}
	var m,f string
	for i:=0;i<lena;i++{
		//создаем ген в хромосоме с интронами
		var s string
		for j, cods:=range algs[i]{
			comanda:=cods.Code%COUNTCOMMAND
			s+=string(comanda+43)
			switch comanda{
			case NOP:
				//команду добавили, а больше нечего
			case ADD,SUB,MUL,DIV,REM,AND,OR,XOR,SLL,SRL,SEQ,SGE,SLT,SNE,LDMX,LDINX,STMX,STOUTX:
				if cods.Op1<10{//если операнд можно изобразить только одной цифрой
					s+=strconv.Itoa(int(cods.Op1))
				}else{
					s+=string(cods.Op1+32*2)
				}

				if cods.Op2<10{
					s+=strconv.Itoa(int(cods.Op2))
				}else{
					s+=string(cods.Op2+32*2)
				}

				if cods.Op3<10{
					s+=strconv.Itoa(int(cods.Op3))
				}else{
					s+=string(cods.Op3+32*2)
				}
			case ADDI,SUBI,MULI,DIVI,REMI,ANDI,ORI,XORI,SLLI,SRLI,BEQ,BGE,BLT,BNE,BLE,BGT:
				if cods.Op1<10{
					s+=strconv.Itoa(int(cods.Op1))
				}else{
					s+=string(cods.Op1+32*2)
				}

				if cods.Op2<10{
					s+=strconv.Itoa(int(cods.Op2))
				}else{
					s+=string(cods.Op2+32*2)
				}

				if cods.Op3<10 && cods.Op3>=0{
					s+=strconv.Itoa(int(cods.Op3))
				}else{//придется использовать 8 рун из кирилицы
					s+=string(0x400+(cods.Op3&0xff))+string(0x400+(cods.Op3>>8 & 0xff))+string(0x400+(cods.Op3>>16 & 0xff))+
						string(0x400+(cods.Op3>>24 & 0xff))+string(0x400+(cods.Op3>>32 & 0xff))+string(0x400+(cods.Op3>>40 & 0xff))+
						string(0x400+(cods.Op3>>48 & 0xff))+string(0x400+(cods.Op3>>56 & 0xff))
				}
			case LDM,LDIN:
				if cods.Op1<10{
					s+=strconv.Itoa(int(cods.Op1))
				}else{
					s+=string(cods.Op1+32*2)
				}

				if cods.Op2<10{//uint всегда больше 0
					s+=strconv.Itoa(int(cods.Op2))
				}else{
					s+=string(0x400+(cods.Op2&0xff))+string(0x400+(cods.Op2>>8 & 0xff))+string(0x400+(cods.Op2>>16 & 0xff))+
						string(0x400+(cods.Op2>>24 & 0xff))+string(0x400+(cods.Op2>>32 & 0xff))+string(0x400+(cods.Op2>>40 & 0xff))+
						string(0x400+(cods.Op2>>48 & 0xff))+string(0x400+(cods.Op2>>56 & 0xff))
				}

				if cods.Op3<10 && cods.Op3>=0{
					s+=strconv.Itoa(int(cods.Op3))
				}else{
					s+=string(0x400+(cods.Op3&0xff))+string(0x400+(cods.Op3>>8 & 0xff))+string(0x400+(cods.Op3>>16 & 0xff))+
						string(0x400+(cods.Op3>>24 & 0xff))+string(0x400+(cods.Op3>>32 & 0xff))+string(0x400+(cods.Op3>>40 & 0xff))+
						string(0x400+(cods.Op3>>48 & 0xff))+string(0x400+(cods.Op3>>56 & 0xff))
				}
			case STM,STOUT:
				if cods.Op1<10{
					s+=strconv.Itoa(int(cods.Op1))
				}else{
					s+=string(0x400+(cods.Op1&0xff))+string(0x400+(cods.Op1>>8 & 0xff))+string(0x400+(cods.Op1>>16 & 0xff))+
						string(0x400+(cods.Op1>>24 & 0xff))+string(0x400+(cods.Op1>>32 & 0xff))+string(0x400+(cods.Op1>>40 & 0xff))+
						string(0x400+(cods.Op1>>48 & 0xff))+string(0x400+(cods.Op1>>56 & 0xff))
				}

				if cods.Op2<10{
					s+=strconv.Itoa(int(cods.Op2))
				}else{
					s+=string(0x400+(cods.Op2&0xff))+string(0x400+(cods.Op2>>8 & 0xff))+string(0x400+(cods.Op2>>16 & 0xff))+
						string(0x400+(cods.Op2>>24 & 0xff))+string(0x400+(cods.Op2>>32 & 0xff))+string(0x400+(cods.Op2>>40 & 0xff))+
						string(0x400+(cods.Op2>>48 & 0xff))+string(0x400+(cods.Op2>>56 & 0xff))
				}

				if cods.Op3<10 && cods.Op3>=0{
					s+=strconv.Itoa(int(cods.Op3))
				}else{
					s+=string(cods.Op3+32*2)
				}
			case JMP:
				if cods.Op3<10 && cods.Op3>=0{
					s+=strconv.Itoa(int(cods.Op3))
				}else{
					s+=string(0x400+(cods.Op3&0xff))+string(0x400+(cods.Op3>>8 & 0xff))+string(0x400+(cods.Op3>>16 & 0xff))+
						string(0x400+(cods.Op3>>24 & 0xff))+string(0x400+(cods.Op3>>32 & 0xff))+string(0x400+(cods.Op3>>40 & 0xff))+
						string(0x400+(cods.Op3>>48 & 0xff))+string(0x400+(cods.Op3>>56 & 0xff))
				}
			case PUSH,POP:
				if cods.Op1<10{
					s+=strconv.Itoa(int(cods.Op1))
				}else{
					s+=string(cods.Op1+32*2)
				}
			case LI:
				if cods.Op1<10{
					s+=strconv.Itoa(int(cods.Op1))
				}else{
					s+=string(cods.Op1+32*2)
				}

				if cods.Op3<10 && cods.Op3>=0{
					s+=strconv.Itoa(int(cods.Op3))
				}else{
					s+=string(0x400+(cods.Op3&0xff))+string(0x400+(cods.Op3>>8 & 0xff))+string(0x400+(cods.Op3>>16 & 0xff))+
						string(0x400+(cods.Op3>>24 & 0xff))+string(0x400+(cods.Op3>>32 & 0xff))+string(0x400+(cods.Op3>>40 & 0xff))+
						string(0x400+(cods.Op3>>48 & 0xff))+string(0x400+(cods.Op3>>56 & 0xff))
				}
			}
			s+=introns[i][j]
		}
		m+="⚤"+s+"Ⱑ"+Mscribes[i]
		f+="⚤"+s+"Ⱑ"+Fscribes[i]
	}
	P.M,_=NewMonoid(Chromosome(m))//ошибок не может быть - ведь мы из алгоритмов сделали гены
	P.F,_=NewMonoid(Chromosome(f))

	return P,nil
}

//вспомагательная функция для MakePairoidFromAlgs (если желается пустых интронов)
//возвращает матрицу пустых строк в соответствии со слайсом слайсов Command
func MakeEmptyIntrons(algs [][]Command) [][]string{
	a := make([][]string, len(algs))
	for i := range a {
		a[i] = make([]string, len(algs[i]))
	}
	return a
}

//вспомагательная функция для MakePairoidFromAlgs (если желается пустых подписей)
//возвращает слай пустых m строк
func MakeEmptyScribes(algs [][]Command) []string{
	return make([]string,len(algs))
}

