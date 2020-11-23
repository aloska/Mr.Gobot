package agent

import (
	"fmt"
	mmap "github.com/edsrzf/mmap-go"
	"io/ioutil"
	"regexp"
	"sort"
	"strconv"
)

/*Synapses - структура с файлом синапсов, хранящая все  Chemistry, самая интенсивно используемая

 */
type Synapses struct {
	filedesc string  //файл описания синапсов в формате syn-[0-9]+x[0-9]+.[0-9]+
	bytearrayTypicalChe mmap.MMap  //замапленный файл синапсов
	TypicalChe *Chemical //в файле описания содержится типичный состав химии, и он в этом поле отражен

	number    SynEnum     //уникальный номер синаптического поля (ядра или входа или выхода)
	filename  string     //имя файла, где записаны синапсы
	bytearray mmap.MMap  //замапленный файл синапсов
	syn       []Chemical //тот же файл в удобном виде, как слайс Chemical с веществами
	maxX      uint32     //ширина синаптического поля
	maxY      uint32     //высота
}

func (sy *Synapses) Init(o *Organism) bool{
	//сюда входим с известными путями к файлам

	//читаем сожержимое файла описания в typicalChe
	if err:=sy.mmapTypicalChe(); err!=nil{
		o.agent.errorr("Synapses не может прочитать TypicalChe из файла описания: "+err.Error())
		o.agent.log.Error("Synapses не может прочитать TypicalChe из файла описания: "+err.Error())
		return false
	}

	//создаем mmap синапсов (если его нет - функция сама создаст)
	if err:=sy.mmapSynapse(); err!=nil{
		o.agent.errorr("Synapses не может сделать ммап: "+err.Error())
		o.agent.log.Error("Synapses не может сделать ммап: "+err.Error())
		return false
	}

	//добавляем это поле синапсов в быстрый доступ организма
	o.synapsesMap[sy.number]=sy


	return true
}

/*Cells - клетки, ммап на файл клеток, и описывающий их геном
 */
type Cells struct {
	filenameGens  string      //имя файла, где записаны гены
	bytearrayGens mmap.MMap   //замапленный файл генов
	genes         []GenNeuron //тот же файл в удобном виде, как слайс Gen с отдельными генами для каждого вида

	filenameCells  string    //имя файла, где записаны клетки
	bytearrayCells mmap.MMap //замапленный файл клеток
	neurons        []Neuron  //тот же файл в удобном виде, как слайс Neuron с отдельными генами для каждого вида

}

func (ce *Cells) Init(o *Organism) bool{
//сюда входим с известными путями к файлам
	//создаем mmap на гены нейронов
	if err:=ce.mmapGenNeuron(); err!=nil{
		o.agent.errorr("Cells не может создать mmap на гены: "+err.Error())
		o.agent.log.Error("Cells не может создать mmap на гены: "+err.Error())
		return false
	}

	//создаем файл нейронов, если нет
	if !fileExists(ce.filenameCells){
		o.agent.info("Файла нейронов пока нет. Создаем. "+ce.filenameCells)
		o.agent.log.Info("Файла нейронов пока нет. Создаем. "+ce.filenameCells)

		if err:=ce.createNeuronsFile(); err!=nil{
			o.agent.errorr("Cells не может создать файл нейронов: "+err.Error())
			o.agent.log.Error("Cells не может создать файл нейронов: "+err.Error())
			return false
		}
	}
	//создаем ммап на нейроны
	if err:=ce.mmapNeuron(); err!=nil{
		o.agent.errorr("Cells не может создать mmap нейронов: "+err.Error())
		o.agent.log.Error("Cells не может создать mmap нейронов: "+err.Error())
		return false
	}

	//добавляем себя в слайс быстрого доступа
	o.cellsSlice=append(o.cellsSlice, ce)
	//увеличиваем общее количество клеток организма
	o.countall+=len(ce.neurons)

	return true
}

/*Core - ядро организма. Может быть много ядер
Может хранить только один файл синапсов и сколько угодно файлов клеток и генов

*/
type Core struct {
	number uint16		//номер ядра
	path     string    //папка, где расположены все файлы ядры (за пределами этой папки другие ядра)
	synapses Synapses  //обслуживает файл с синапсами этого ядра (у каждого ядра только один файл синапсов)
	cells    []Cells   //слайс обслуживает файлы с клетками и геномами этих клеток
}

func (co *Core) Init(o *Organism) bool{
//сюда входим с известными путем и номером ядра
	//ищем гены нейронов и файл-описание синапса
	gneuronfiles:=[]string{}
	syndescfile:=""
	files, _ := ioutil.ReadDir(co.path)
	for _, file := range files {
		if match, _ := regexp.MatchString("(GenNeuron-[0-9]+.genes)",
			file.Name()); match{
			gneuronfiles = append(gneuronfiles,file.Name())
		}else if match, _ := regexp.MatchString("(syn-[0-9]+x[0-9]+.[0-9]+)",
			file.Name()); match{
			syndescfile = file.Name()
		}
	}
	//остортируем
	if !sort.StringsAreSorted(gneuronfiles){
		sort.Strings(gneuronfiles)
	}
	//добавляем
	re := regexp.MustCompile("[0-9]+")
	for _, neus:= range gneuronfiles{

		co.cells=append(co.cells,
			Cells{
				filenameGens: co.path+"/"+neus,
				filenameCells: co.path+"/Neuron-"+re.FindString(neus)+".neurons"	})
	}
	//инициализируем
	for i:=0;i<len(co.cells);i++{
		if !co.cells[i].Init(o){
			o.agent.errorr(co.cells[i].filenameGens+" не может инициализироваться")
			o.agent.log.Error(co.cells[i].filenameGens+" не может инициализироваться")
			return false
		}
	}

	//отпарсим номер и размер
	re = regexp.MustCompile(`[0-9]+`)
	ss:=re.FindAllString(syndescfile, -1)
	mx,_:=strconv.Atoi(ss[0])
	my,_:=strconv.Atoi(ss[1])
	num,_:=strconv.Atoi(ss[2])
	co.synapses = Synapses{
		number: SynEnum(num),
		maxX: uint32(mx),
		maxY: uint32(my),
		filename: co.path+"/Synapse-"+ss[2]+".chemical",
		filedesc: co.path+"/"+syndescfile}
	//инициализируем синапсы
	if !co.synapses.Init(o){
		o.agent.errorr(co.synapses.filename+" не может инициализироваться")
		o.agent.log.Error(co.synapses.filename+" не может инициализироваться")
		return false
	}
	return  true
}

/*Brain - большая структура с целым мозгом из всех ядер, входящих в его состав
 */
type Brain struct {
	path     string    //папка, где расположены все ядра, а входные и выходные устройства входят в состав организма
	cores    []Core    //мозг состоит из ядер
	organism *Organism //ссылка на весь родительский организм
}

//Init - инициализация мозга системы
func (b *Brain) Init(o *Organism) bool{
	b.organism=o
	//Найдем все ядра
	cores:=[]string{}
	files, _ := ioutil.ReadDir(o.path+"/Brain")
	for _, file := range files {
		if file.IsDir(){
			cores = append(cores,file.Name())
		} else{
			//здесь не должно быть файлов. Мы сюда не попадаем, потому что сверху проверили
		}
	}
	//остортируем ядра
	if !sort.StringsAreSorted(cores){
		sort.Strings(cores)
	}
	//добавляем ядра
	re := regexp.MustCompile("[0-9]+")
	for _,cor:=range cores{
		num, _:=strconv.Atoi(re.FindString(cor))
		b.cores=append(b.cores,
			Core{
				number: uint16(num),
				path:o.path+"/Brain/"+cor})

	}
	//инициализируем ядра
	for i:=0;i<len(b.cores);i++{
		if !b.cores[i].Init(o){
			o.agent.errorr(strconv.Itoa(int(b.cores[i].number))+" не может инициализироваться")
			o.agent.log.Error(strconv.Itoa(int(b.cores[i].number))+" не может инициализироваться")
			return false
		}
	}
	o.agent.info("/Brain готов к работе")
	o.agent.log.Info("/Brain готовы к работе")
	return true
}

/*DataInput - описание входа
Поскольку в Го нет женериков, прийдется хранить поля всех видов данных
*/
type DataInput struct {
	typeData DataTypeEnum //тип входа 1-DataByte...

	filenameGen  string    //имя файла, где записан ген входа
	bytearrayGen mmap.MMap //замапленный файл гена
	genData      *GenData  //тот же файл ввиде гена (у ячеек данных он один)

	filenameData  string    //имя файла, где записаны ячейки входа
	bytearrayData mmap.MMap //замапленный файл ячеек

	/*Данные с входа:
	Прочти это, перед тем, как думать, можно ли что-то здесь заоптимизировать!
	--------------------------------------------------------------------------

	слайс рефлектирован на файл данных
	на самом деле будет использован только один слайс в этой структуре, в заыисимости от typeData
	Почему не через интерфейсы? Потому что нам нужен слайс ссылающийся на mmap файла, а не хранящий ЗНАЧЕНИЯ
	Ссылки на интерфесы тоже не подходят, потому что ссылка на интерфейс ссылается имеенно на интерфейс, а он хранит ЗНАЧЕНИЕ

	По мере реализации разных типов входных данных, сюда будут добавляться слайсы на эти типы.
	Тот слайс, который используется обозначен в typeData

	Немного коряво, потому что пустые слайсы хранят дополнительно 24 байта. Поэтому, если у нас есть 10 разных типов входов, то каждый
	из них хранит 10*24=240 байт минимум ненужных данных, а все вместе 2к байт ненужных.
	Если у нас 100 входов и 10 типов входов, то: 10*24=240 Б на каждом, итого 240*100=24Кб
	Вобщем-то, не так много, если учесть, что мы платим за удобство работы со слайсом данных вместо работы с сырыми байтами.

	Хотя, можно было бы отдать все на "воспитание" самим рецепторам, храня только слайс байтов данных. Но это жуть как неудобно.
	Опять же, у конечного универсального решателя какие входы и сколько?
	Зрение - это один вход RGBA
	Слух - один вход байтов и возможно фурье
	ну датчики положения тела, штук 1000 - а это один тип
	итого "лишних" байтов на 1000-2000 входов наберется не более 100*24*2000 ~ 2-3Мб - по сравнению с терабайтами самого организма,
	не смешите меня больше своим неудобством в угоду экономии, ааааааааа)))

	Опять же, с помощью такого подходв можно входом сделать прямо го-шную структуру! Описать ее тип данных и сделать специальный рецептор,
	реагирующий на структуру целиком. Вот где удобство.

	А там, может в следующих версиях Го подгонят дженерики - и все переделаем))
	Вопрос закрыт.
	*/
	//должен быть список всех ТИПОВ входов, используемых в агенте
	dataByte   []DataByte
	dataUInt32 []DataUInt32
	datauint32 []Datauint32
	dataBit    []DataBit
}

//Init -
func (d *DataInput) Init(o* Organism) bool{
//сюда входим с известными путями к файлам

	//создаем mmap на ген данных
	if err:=d.mmapGenData(); err!=nil{
		o.agent.errorr("DataInput не может создать mmap: "+err.Error())
		o.agent.log.Error("DataInput не может создать mmap: "+err.Error())
		return false
	}
	//создаем mmap на файл данных, если файла нет, функция сама создаст его
	if err:=d.mmapData(); err!=nil{
		o.agent.errorr("DataInput не может создать mmap: "+err.Error())
		o.agent.log.Error("DataInput не может создать mmap: "+err.Error())
		return false
	}

	d.typeData=d.genData.Datatype

	return true
}

/*Receptors - рецепторы, ммап на файл рецепторов, и описывающий их геном
 */
type Receptors struct {
	filenameGens  string        //имя файла, где записаны гены
	bytearrayGens mmap.MMap     //замапленный файл генов
	genes         []GenReceptor //тот же файл в удобном виде, как слайс Gen с отдельными генами для каждого вида рецепторов

	filenameRecs  string     //имя файла, где записаны рецепторы
	bytearrayRecs mmap.MMap  //замапленный файл клеток
	recs          []Receptor //тот же файл в удобном виде, как слайс Receptor с отдельными генами для каждого вида
}

func (re *Receptors) Init(o *Organism) bool{
//сюда входим с известными путями к файлам
	//создаем mmap на ген рецепторов
	if err:=re.mmapGenReceptor(); err!=nil{
		o.agent.errorr("Receptors не может создать mmap: "+err.Error())
		o.agent.log.Error("Receptors не может создать mmap: "+err.Error())
		return false
	}

	//создаем файл рецепторов, если нет
	if !fileExists(re.filenameRecs){
		o.agent.info("Файла рецепторов пока нет. Создаем. "+re.filenameRecs)
		o.agent.log.Info("Файла рецепторов пока нет. Создаем. "+re.filenameRecs)

		if err:=re.createReceptorsFile(); err!=nil{
			o.agent.errorr("Receptors не может создать файл рецепторов: "+err.Error())
			o.agent.log.Error("Receptors не может создать файл рецепторов: "+err.Error())
			return false
		}
	}

	//создаем ммап на рецепторы
	if err:=re.mmapReceptors(); err!=nil{
		o.agent.errorr("Receptors не может создать mmap на рецепторы: "+err.Error())
		o.agent.log.Error("Receptors не может создать mmap на рецепторы: "+err.Error())
		return false
	}

	//увеличиваем общее количество клеток организма
	o.countall+=len(re.recs)

	return true
}

/*Input - вход состоит из входных файлов, рецепторов и (если нужно, нейронов)

 */
type Input struct {
	number    uint16      //номер входа
	path      string      //папка, где расположены все файлы этого входа
	dataInput DataInput   //вход
	receptors []Receptors //все рецепторы этого входа
	synapses  Synapses    //обслуживает файл с синапсами этого входа (если они есть!!)
	cells     []Cells     //слайс обслуживает файлы с нейронами и геномами этих нейронов, если они заданы
}

func (in* Input) Init(o *Organism) bool{
//сюда входим с известными путем и номером входа
	//скажем DataInput где его файлы
	in.dataInput.filenameGen=in.path+"/GenData.genes"
	in.dataInput.filenameData=in.path+"/Data.data"
	//инициализация
	if !in.dataInput.Init(o){
		o.agent.log.Error("InputData не может инициализироваться")
		o.agent.errorr("InputData не может инициализироваться")
		return false
	}

	//ищем гены рецепторов
	receptorfiles:=[]string{}
	isSyn:=false
	syndescfile:="" //и заодно файл-описание синапсов, если есть
	files, _ := ioutil.ReadDir(in.path)
	for _, file := range files {
		if match, _ := regexp.MatchString("(GenReceptor-[0-9]+.genes)",
			file.Name()); match{
			receptorfiles = append(receptorfiles,file.Name())
		}else if match, _ := regexp.MatchString("(syn-[0-9]+x[0-9]+.[0-9]+)",
			file.Name()); match{
			syndescfile = file.Name()
			isSyn=true
		}
	}
	//остортируем
	if !sort.StringsAreSorted(receptorfiles){
		sort.Strings(receptorfiles)
	}
	//добавляем
	re := regexp.MustCompile("[0-9]+")
	for _, recs:= range receptorfiles{

		in.receptors=append(in.receptors,
							Receptors{
								filenameGens: in.path+"/"+recs,
								filenameRecs: in.path+"/Receptor-"+re.FindString(recs)+".receptors"	})
	}
	//инициализируем
	for i:=0;i<len(in.receptors);i++{
		if !in.receptors[i].Init(o){
			o.agent.errorr(in.receptors[i].filenameGens+" не может инициализироваться")
			o.agent.log.Error(in.receptors[i].filenameGens+" не может инициализироваться")
			return false
		}
	}


	//есть ли свое синаптическое поле?
	if isSyn{
		//ищем гены нейронов
		gneuronfiles:=[]string{}

		files, _ := ioutil.ReadDir(in.path)
		for _, file := range files {
			if match, _ := regexp.MatchString("(GenNeuron-[0-9]+.genes)",
				file.Name()); match{
				gneuronfiles = append(gneuronfiles,file.Name())
			}
		}
		//остортируем
		if !sort.StringsAreSorted(gneuronfiles){
			sort.Strings(gneuronfiles)
		}
		//добавляем
		re := regexp.MustCompile("[0-9]+")
		for _, neus:= range gneuronfiles{

			in.cells=append(in.cells,
				Cells{
					filenameGens: in.path+"/"+neus,
					filenameCells: in.path+"/Neuron-"+re.FindString(neus)+".neurons"	})
		}
		//инициализируем
		for i:=0;i<len(in.cells);i++{
			if !in.cells[i].Init(o){
				o.agent.errorr(in.cells[i].filenameGens+" не может инициализироваться")
				o.agent.log.Error(in.cells[i].filenameGens+" не может инициализироваться")
				return false
			}
		}

		//отпарсим номер и размер
		re = regexp.MustCompile(`[0-9]+`)
		ss:=re.FindAllString(syndescfile, -1)
		mx,_:=strconv.Atoi(ss[0])
		my,_:=strconv.Atoi(ss[1])
		num,_:=strconv.Atoi(ss[2])
		in.synapses = Synapses{
			number: SynEnum(num),
			maxX: uint32(mx),
			maxY: uint32(my),
			filename: in.path+"/Synapse-"+ss[2]+".chemical",
			filedesc: in.path+"/"+syndescfile}
		//инициализируем синапсы
		if !in.synapses.Init(o){
			o.agent.errorr(in.synapses.filename+" не может инициализироваться")
			o.agent.log.Error(in.synapses.filename+" не может инициализироваться")
			return false
		}
	}
	return true
}

/*Senses - все входы (ощущения)
 */
type Senses struct {
	inputs []Input //все входы

	synapses Synapses //синаптическое поле всех входов
	cells    []Cells  /*слайс обслуживает файлы с нейронами и геномами этих нейронов, если они заданы, для общего синаптичесского поля всех входов
	этих нейронов может не быть, и тогда это значит, что общего синаптического поля входов нет
	такое поведение может использоваться для очень простых агентов
	*/
	organism *Organism //ссылка на весь родительский организм
}

//Init - инициализация Чувств системы
func (s *Senses) Init(o *Organism) bool{
	s.organism=o
	//Найдем все входы
	inputs:=[]string{}
	syndescfile:="" //и заодно файл-описание синапсов, если есть
	isSyn:=false
	files, _ := ioutil.ReadDir(o.path+"/Senses")
	for _, file := range files {
		if file.IsDir(){
			inputs = append(inputs,file.Name())
		} else if match, _ := regexp.MatchString("(syn-[0-9]+x[0-9]+.[0-9]+)",
			file.Name()); match{
			syndescfile = file.Name()
			isSyn=true
		}
	}

	//остортируем входы
	if !sort.StringsAreSorted(inputs){
		sort.Strings(inputs)
	}
	//добавляем входы
	re := regexp.MustCompile("[0-9]+")
	for _,inp:=range inputs{
		num, _:=strconv.Atoi(re.FindString(inp))
		s.inputs=append(s.inputs,
			Input{
				number: uint16(num),
				path:o.path+"/Senses/"+inp})

	}
	//инициализируем входы
	for i:=0;i<len(s.inputs);i++{
		if !s.inputs[i].Init(o){
			o.agent.errorr(strconv.Itoa(int(s.inputs[i].number))+" не может инициализироваться")
			o.agent.log.Error(strconv.Itoa(int(s.inputs[i].number))+" не может инициализироваться")
			return false
		}
	}

	//есть ли свое синаптическое поле?
	if isSyn{
		//ищем гены нейронов
		gneuronfiles:=[]string{}

		files, _ := ioutil.ReadDir(o.path+"/Senses")
		for _, file := range files {
			if match, _ := regexp.MatchString("(GenNeuron-[0-9]+.genes)",
				file.Name()); match{
				gneuronfiles = append(gneuronfiles,file.Name())
			}
		}
		//остортируем
		if !sort.StringsAreSorted(gneuronfiles){
			sort.Strings(gneuronfiles)
		}
		//добавляем
		re := regexp.MustCompile("[0-9]+")
		for _, neus:= range gneuronfiles{

			s.cells=append(s.cells,
				Cells{
					filenameGens: o.path+"/Senses/"+neus,
					filenameCells: o.path+"/Senses/Neuron-"+re.FindString(neus)+".neurons"	})
		}
		//инициализируем
		for i:=0;i<len(s.cells);i++{
			if !s.cells[i].Init(o){
				o.agent.errorr(s.cells[i].filenameGens+" не может инициализироваться")
				o.agent.log.Error(s.cells[i].filenameGens+" не может инициализироваться")
				return false
			}
		}

		//отпарсим номер и размер
		re = regexp.MustCompile(`[0-9]+`)
		ss:=re.FindAllString(syndescfile, -1)
		mx,_:=strconv.Atoi(ss[0])
		my,_:=strconv.Atoi(ss[1])
		num,_:=strconv.Atoi(ss[2])
		s.synapses = Synapses{
			number: SynEnum(num),
			maxX: uint32(mx),
			maxY: uint32(my),
			filename: o.path+"/Senses/Synapse-"+ss[2]+".chemical",
			filedesc: o.path+"/Senses/"+syndescfile}
		//инициализируем синапсы
		if !s.synapses.Init(o){
			o.agent.errorr(s.synapses.filename+" не может инициализироваться")
			o.agent.log.Error(s.synapses.filename+" не может инициализироваться")
			return false
		}
		//добавим также в мапу синапсов для быстрого доступа,
		//в итоге это поле синапсов будет лежать в мапе и под своим номером, и под номером SYNINPUTS=0xfffe
		o.synapsesMap[SYNINPUTS]=&s.synapses
	}


	o.agent.info("/Senses готовы к работе")
	o.agent.log.Info("/Senses готовы к работе")

	return true
}


/*DataOutput - описание выхода
 */
type DataOutput struct {
	typeData DataTypeEnum //тип выхода

	filenameGen  string      //имя файла, где записан ген выхода
	bytearrayGen mmap.MMap   //замапленный файл гена
	genData      *GenDataOut //тот же файл ввиде гена (у ячеек данных он один)

	filenameData  string    //имя файла, где записаны ячейки выхода
	bytearrayData mmap.MMap //замапленный файл ячеек

	//должен быть список всех ТИПОВ выходов, используемых в агенте
	datauint32 []Datauint32
	dataRune   []DataRune
	dataUInt32 []DataUInt32
}
func (d *DataOutput) Init(o* Organism) bool{
	//сюда входим с известными путями к файлам

	//создаем mmap на ген данных
	if err:=d.mmapGenData(); err!=nil{
		o.agent.errorr("DataInput не может создать mmap: "+err.Error())
		o.agent.log.Error("DataInput не может создать mmap: "+err.Error())
		return false
	}
	//создаем mmap на файл данных, если файла нет, функция сама создаст его
	if err:=d.mmapData(); err!=nil{
		o.agent.errorr("DataInput не может создать mmap: "+err.Error())
		o.agent.log.Error("DataInput не может создать mmap: "+err.Error())
		return false
	}

	d.typeData=d.genData.Datatype

	return true
}

/*Preffectors - преффекторы, ммап на файл преффекторов, и описывающий их геном
файлов преффекторов для одного выхода может быть много
*/
type Preffectors struct {
	filenameGens  string          //имя файла, где записаны гены
	bytearrayGens mmap.MMap       //замапленный файл генов
	genes         []GenPreffector //тот же файл в удобном виде, как слайс Gen с отдельными генами для каждого вида преффекторов

	filenamePres  string       //имя файла, где записаны преффекторы
	bytearrayPres mmap.MMap    //замапленный файл клеток
	prefs         []Preffector //тот же файл в удобном виде, как слайс Preffector с отдельными генами для каждого вида
}

func (pre *Preffectors) Init(o *Organism) bool{
	//сюда входим с известными путями к файлам
	//создаем mmap на ген преффекторов
	if err:=pre.mmapGenPreffector(); err!=nil{
		o.agent.errorr("Preffectors не может создать mmap: "+err.Error())
		o.agent.log.Error("Preffectors не может создать mmap: "+err.Error())
		return false
	}

	//создаем файл рецепторов, если нет
	if !fileExists(pre.filenamePres){
		o.agent.info("Файла преффекторов пока нет. Создаем. "+pre.filenamePres)
		o.agent.log.Info("Файла преффекторов пока нет. Создаем. "+pre.filenamePres)

		if err:=pre.createPreffectorsFile(); err!=nil{
			o.agent.errorr("Preffectors не может создать файл преффекторов: "+err.Error())
			o.agent.log.Error("Preffectors не может создать файл преффекторов: "+err.Error())
			return false
		}
	}

	//создаем ммап на преффекторы
	if err:=pre.mmapPreffectors(); err!=nil{
		o.agent.errorr("Preffectors не может создать mmap на преффекторов: "+err.Error())
		o.agent.log.Error("Preffectors не может создать mmap на преффекторов: "+err.Error())
		return false
	}
	//увеличиваем общее количество клеток организма
	o.countall+=len(pre.prefs)

	return true
}

/*Effector - считывает данные со своих префекторов и складывает значения в выходной файл
Эффектор только один на выход, а префекторов много

Конкретная реализация эффекторов описана в effector.go
------------------------------------------------------
*/
type Effector struct {
	number     uint16     //номер выхода
	path       string     //папка, где расположены все файлы этого выхода
	dataOutput DataOutput //выход

	preffectors []Preffectors //преффекторы этого выхода (может быть много не только генов, но и геномов)

	synapses Synapses //обслуживает файл с синапсами этого выхода (если они есть!!)
	cells    []Cells  //слайс обслуживает файлы с нейронами и геномами этих нейронов, если они заданы
}

func (ef *Effector) Init(o *Organism) bool{
	//сюда входим с известными путем и номером выхода
	//скажем DataInput где его файлы
	ef.dataOutput.filenameGen=ef.path+"/GenDataOut.genes"
	ef.dataOutput.filenameData=ef.path+"/DataOut.data"
	//инициализация
	if !ef.dataOutput.Init(o){
		o.agent.log.Error("dataOutput не может инициализироваться")
		o.agent.errorr("dataOutput не может инициализироваться")
		return false
	}
	//ищем гены преффекторов
	preffectorfiles:=[]string{}
	isSyn:=false
	syndescfile:="" //и заодно файл-описание синапсов, если есть
	files, _ := ioutil.ReadDir(ef.path)
	for _, file := range files {
		if match, _ := regexp.MatchString("(GenPreffector-[0-9]+.genes)",
			file.Name()); match{
			preffectorfiles = append(preffectorfiles,file.Name())
		}else if match, _ := regexp.MatchString("(syn-[0-9]+x[0-9]+.[0-9]+)",
			file.Name()); match{
			syndescfile = file.Name()
			isSyn=true
		}
	}
	//остортируем
	if !sort.StringsAreSorted(preffectorfiles){
		sort.Strings(preffectorfiles)
	}
	//добавляем
	re := regexp.MustCompile("[0-9]+")
	for _, recs:= range preffectorfiles{

		ef.preffectors=append(ef.preffectors,
			Preffectors{
				filenameGens: ef.path+"/"+recs,
				filenamePres: ef.path+"/Preffector-"+re.FindString(recs)+".preffectors"	})
	}
	//инициализируем
	for i:=0;i<len(ef.preffectors);i++{
		if !ef.preffectors[i].Init(o){
			o.agent.errorr(ef.preffectors[i].filenameGens+" не может инициализироваться")
			o.agent.log.Error(ef.preffectors[i].filenameGens+" не может инициализироваться")
			return false
		}
	}

	//есть ли свое синаптическое поле?
	if isSyn{
		//ищем гены нейронов
		gneuronfiles:=[]string{}

		files, _ := ioutil.ReadDir(ef.path)
		for _, file := range files {
			if match, _ := regexp.MatchString("(GenNeuron-[0-9]+.genes)",
				file.Name()); match{
				gneuronfiles = append(gneuronfiles,file.Name())
			}
		}
		//остортируем
		if !sort.StringsAreSorted(gneuronfiles){
			sort.Strings(gneuronfiles)
		}
		//добавляем
		re := regexp.MustCompile("[0-9]+")
		for _, neus:= range gneuronfiles{

			ef.cells=append(ef.cells,
				Cells{
					filenameGens: ef.path+"/"+neus,
					filenameCells: ef.path+"/Neuron-"+re.FindString(neus)+".neurons"	})
		}
		//инициализируем
		for i:=0;i<len(ef.cells);i++{
			if !ef.cells[i].Init(o){
				o.agent.errorr(ef.cells[i].filenameGens+" не может инициализироваться")
				o.agent.log.Error(ef.cells[i].filenameGens+" не может инициализироваться")
				return false
			}
		}

		//отпарсим номер и размер
		re = regexp.MustCompile(`[0-9]+`)
		ss:=re.FindAllString(syndescfile, -1)
		mx,_:=strconv.Atoi(ss[0])
		my,_:=strconv.Atoi(ss[1])
		num,_:=strconv.Atoi(ss[2])
		ef.synapses = Synapses{
			number: SynEnum(num),
			maxX: uint32(mx),
			maxY: uint32(my),
			filename: ef.path+"/Synapse-"+ss[2]+".chemical",
			filedesc: ef.path+"/"+syndescfile}
		//инициализируем синапсы
		if !ef.synapses.Init(o){
			o.agent.errorr(ef.synapses.filename+" не может инициализироваться")
			o.agent.log.Error(ef.synapses.filename+" не может инициализироваться")
			return false
		}
	}

	return true
}

/*Actions - все выходы (действия организма, кроме внутренних)
 */
type Actions struct {
	effectors []Effector //все выходы
	synapses Synapses //синаптическое поле всех выходов (может не быть)
	cells    []Cells  /*слайс обслуживает файлы с нейронами и геномами этих нейронов, если они заданы,
	для общего синаптичесского поля всех выходов
	этих нейронов может не быть, и тогда это значит, что общего синаптического поля вЫходов нет
	такое поведение может использоваться для большинства агентов.
	Но для сложных агентов оно нужно - что-то типа мозжечка, корректирующего сложные синхронные слаженные поведения многих выходов
	*/

	organism *Organism //ссылка на весь родительский организм
}

//Init - инициализация Действий системы
func (ac *Actions) Init(o *Organism) bool{
	ac.organism=o

	//Найдем все выходы
	effecs:=[]string{}
	syndescfile:=""
	isSyn:=false
	files, _ := ioutil.ReadDir(o.path+"/Actions")
	for _, file := range files {
		if file.IsDir(){
			effecs = append(effecs,file.Name())
		} else if match, _ := regexp.MatchString("(syn-[0-9]+x[0-9]+.[0-9]+)",
			file.Name()); match{
			syndescfile = file.Name()
			isSyn=true
		}
	}
	//остортируем выходы
	if !sort.StringsAreSorted(effecs){
		sort.Strings(effecs)
	}
	//добавляем выходы
	re := regexp.MustCompile("[0-9]+")
	for _,eff:=range effecs{
		num, _:=strconv.Atoi(re.FindString(eff))
		ac.effectors=append(ac.effectors,
			Effector{
				number: uint16(num),
				path: o.path+"/Actions/"+eff})

	}
	//инициализируем выходы
	for i:=0;i<len(ac.effectors);i++{
		if !ac.effectors[i].Init(o){
			o.agent.errorr(strconv.Itoa(int(ac.effectors[i].number))+" не может инициализироваться")
			o.agent.log.Error(strconv.Itoa(int(ac.effectors[i].number))+" не может инициализироваться")
			return false
		}
	}

	//есть ли свое синаптическое поле?
	if isSyn{
		//ищем гены нейронов
		gneuronfiles:=[]string{}

		files, _ := ioutil.ReadDir(o.path+"/Actions")
		for _, file := range files {
			if match, _ := regexp.MatchString("(GenNeuron-[0-9]+.genes)",
				file.Name()); match{
				gneuronfiles = append(gneuronfiles,file.Name())
			}
		}
		//остортируем
		if !sort.StringsAreSorted(gneuronfiles){
			sort.Strings(gneuronfiles)
		}
		//добавляем
		re := regexp.MustCompile("[0-9]+")
		for _, neus:= range gneuronfiles{

			ac.cells=append(ac.cells,
				Cells{
					filenameGens: o.path+"/Actions/"+neus,
					filenameCells: o.path+"/Actions/Neuron-"+re.FindString(neus)+".neurons"	})
		}
		//инициализируем
		for i:=0;i<len(ac.cells);i++{
			if !ac.cells[i].Init(o){
				o.agent.errorr(ac.cells[i].filenameGens+" не может инициализироваться")
				o.agent.log.Error(ac.cells[i].filenameGens+" не может инициализироваться")
				return false
			}
		}

		//отпарсим номер и размер
		re = regexp.MustCompile(`[0-9]+`)
		ss:=re.FindAllString(syndescfile, -1)
		mx,_:=strconv.Atoi(ss[0])
		my,_:=strconv.Atoi(ss[1])
		num,_:=strconv.Atoi(ss[2])
		ac.synapses = Synapses{
			number: SynEnum(num),
			maxX: uint32(mx),
			maxY: uint32(my),
			filename: o.path+"/Actions/Synapse-"+ss[2]+".chemical",
			filedesc: o.path+"/Actions/"+syndescfile}
		//инициализируем синапсы
		if !ac.synapses.Init(o){
			o.agent.errorr(ac.synapses.filename+" не может инициализироваться")
			o.agent.log.Error(ac.synapses.filename+" не может инициализироваться")
			return false
		}

		//добавим также в мапу синапсов для быстрого доступа,
		//в итоге это поле синапсов будет лежать в мапе и под своим номером, и под номером SYNOUTPUTS=0xfffd
		o.synapsesMap[SYNOUTPUTS]=&ac.synapses
	}


	o.agent.info("/Actions готовы к работе")
	o.agent.log.Info("/Actions готовы к работе")

	return true
}

/*Vegetatic - вегетативная нервная система

 */
type Vegetatic struct {
	effectors []Effector //эффекторы вегетативной системы (сердце, дыхание, питание, очищение)

	synapses Synapses //обслуживает файл с синапсами этого выхода (если они есть!!)
	cells    []Cells  //слайс обслуживает файлы с нейронами и геномами этих нейронов, если они заданы

	organism *Organism //ссылка на весь родительский организм
}

//Init - инициализация вегетативной системы
func (v *Vegetatic) Init(o *Organism) bool{
//инициализация очень похожа на Actions (один в один))
	v.organism = v.organism

	//Найдем все выходы
	effecs:=[]string{}
	syndescfile:=""
	isSyn:=false
	files, _ := ioutil.ReadDir(o.path+"/Vegetatic")
	for _, file := range files {
		if file.IsDir(){
			effecs = append(effecs,file.Name())
		} else if match, _ := regexp.MatchString("(syn-[0-9]+x[0-9]+.[0-9]+)",
			file.Name()); match{
			syndescfile = file.Name()
			isSyn=true
		}
	}
	//остортируем выходы
	if !sort.StringsAreSorted(effecs){
		sort.Strings(effecs)
	}
	//добавляем выходы
	re := regexp.MustCompile("[0-9]+")
	for _,eff:=range effecs{
		num, _:=strconv.Atoi(re.FindString(eff))
		v.effectors=append(v.effectors,
			Effector{
				number: uint16(num),
				path: o.path+"/Vegetatic/"+eff})

	}
	//инициализируем выходы
	for i:=0;i<len(v.effectors);i++{
		if !v.effectors[i].Init(o){
			o.agent.errorr(strconv.Itoa(int(v.effectors[i].number))+" не может инициализироваться")
			o.agent.log.Error(strconv.Itoa(int(v.effectors[i].number))+" не может инициализироваться")
			return false
		}
	}

	//есть ли свое синаптическое поле?
	if isSyn{
		//ищем гены нейронов
		gneuronfiles:=[]string{}

		files, _ := ioutil.ReadDir(o.path+"/Vegetatic")
		for _, file := range files {
			if match, _ := regexp.MatchString("(GenNeuron-[0-9]+.genes)",
				file.Name()); match{
				gneuronfiles = append(gneuronfiles,file.Name())
			}
		}
		//остортируем
		if !sort.StringsAreSorted(gneuronfiles){
			sort.Strings(gneuronfiles)
		}
		//добавляем
		re := regexp.MustCompile("[0-9]+")
		for _, neus:= range gneuronfiles{

			v.cells=append(v.cells,
				Cells{
					filenameGens: o.path+"/Vegetatic/"+neus,
					filenameCells: o.path+"/Vegetatic/Neuron-"+re.FindString(neus)+".neurons"	})
		}
		//инициализируем
		for i:=0;i<len(v.cells);i++{
			if !v.cells[i].Init(o){
				o.agent.errorr(v.cells[i].filenameGens+" не может инициализироваться")
				o.agent.log.Error(v.cells[i].filenameGens+" не может инициализироваться")
				return false
			}
		}

		//отпарсим номер и размер
		re = regexp.MustCompile(`[0-9]+`)
		ss:=re.FindAllString(syndescfile, -1)
		mx,_:=strconv.Atoi(ss[0])
		my,_:=strconv.Atoi(ss[1])
		num,_:=strconv.Atoi(ss[2])
		v.synapses = Synapses{
			number: SynEnum(num),
			maxX: uint32(mx),
			maxY: uint32(my),
			filename: o.path+"/Vegetatic/Synapse-"+ss[2]+".chemical",
			filedesc: o.path+"/Vegetatic/"+syndescfile}
		//инициализируем синапсы
		if !v.synapses.Init(o){
			o.agent.errorr(v.synapses.filename+" не может инициализироваться")
			o.agent.log.Error(v.synapses.filename+" не может инициализироваться")
			return false
		}

		//добавим также в мапу синапсов для быстрого доступа,
		//в итоге это поле синапсов будет лежать в мапе и под своим номером, и под номером SYNVEGETATIC=0xffff
		o.synapsesMap[SYNVEGETATIC]=&v.synapses
	}

	o.agent.info("/Vegetatic готовы к работе")
	o.agent.log.Info("/Vegetatic готовы к работе")
	return true
}
/*Organism - самая полная структура, состоящая из мозга, входных и выходных устройств

 */
type Organism struct {
	path      string     //папка, где расположен весь организм
	brain     Brain     //мозг из ядер
	senses    Senses    //ощущения (входы, рецепторы...)
	actions   Actions   //действия
	vegetatic Vegetatic //вегетативная система

	synapsesMap map[SynEnum] *Synapses	/*мапа со всеми реально существующими
	синаптическими полями для быстрого доступа
	(добавляет сюда тот, в ком есть поле)
	а сами синапсы лежат в тех структурах, в папках которых они есть
	Если у входа есть синаптическое поле - он его и создает и отвечает за него
	*/
	cellsSlice []*Cells /*слайс указателей всех реально существующих
	Cells для быстрого доступа
	Добавляет в слайс сама Cells во время Init*/

	countall int //общее количество клеток организма (нейроны, рецепторы, преффекторы)

	agent *Agent
}

//Init - проверка и маппинг всех файлов
func (o *Organism) Init(a *Agent) bool{
	o.path=a.path
	o.agent=a
	o.synapsesMap=make(map[SynEnum]*Synapses)

	if !o.senses.Init(o) {
		return false
	}
	if !o.brain.Init(o) {
		return false
	}
	if !o.actions.Init(o) {
		return false
	}
	if !o.vegetatic.Init(o){
		return false
	}

	return true
}

//Check проверка работоспособности
func (o *Organism) Check(a *Agent, c chan error) {
	//проверим нейроны
	for i:=0;i<len(o.cellsSlice);i++ {
		//проверим, что все дендриты и аксоны не выходят за границы синаптических полей
		for j:=0; j<len(o.cellsSlice[i].neurons);j++{
			//количество синапсов в синаптическом поле, в котором находится нейрон
			lenn:=uint32(len(o.synapsesMap[o.cellsSlice[i].neurons[j].SynNumber].syn))
			if o.cellsSlice[i].neurons[j].N>= lenn{
				//номер клетки больше количества синапсов в поле!
				o.agent.log.Error(fmt.Sprintf("геном: %v, номер гена: %v", o.cellsSlice[i].filenameGens, o.cellsSlice[i].neurons[j].Gen))
				c<-fmt.Errorf("номер нейрона %v больше количества синапсов в поле %v!", o.cellsSlice[i].neurons[j].N, o.cellsSlice[i].neurons[j].SynNumber)
			}
			//бежим по дендритам
			for _, d:=range o.cellsSlice[i].neurons[j].Dendrites{
				if d.N>=lenn{
					//номер дендрита больше количества синапсов
					o.agent.log.Error(fmt.Sprintf("геном: %v, номер гена: %v", o.cellsSlice[i].filenameGens, o.cellsSlice[i].neurons[j].Gen))
					c<-fmt.Errorf("номер дендрита %v нейрона больше количества синапсов в поле %v!", d.N, o.cellsSlice[i].neurons[j].SynNumber)
				}
			}
			//количество синапсов в синаптическом поле, в котором находится аксоны нейрона
			lenna:=uint32(len(o.synapsesMap[o.cellsSlice[i].neurons[j].SynNumberAxons].syn))
			//бежим по аксонам
			for _, d:=range o.cellsSlice[i].neurons[j].Axons{
				if d.N>=lenna{
					//номер дендрита больше количества синапсов
					o.agent.log.Error(fmt.Sprintf("геном: %v, номер гена: %v", o.cellsSlice[i].filenameGens, o.cellsSlice[i].neurons[j].Gen))
					c<-fmt.Errorf("номер аксона %v нейрона больше количества синапсов в поле %v!", d.N, o.cellsSlice[i].neurons[j].SynNumber)
				}
			}
			//все в порядке
			c<-nil
		}
	}
	//проверим рецепторы
	for _, inp:=range o.senses.inputs{
		for _, recs:= range inp.receptors{
			for k:=0;k<len(recs.recs);k++ {
				lenn := uint32(len(o.synapsesMap[recs.recs[k].SynNumber].syn))
				//бежим по аксонам
				for _, d := range recs.recs[k].Axons {
					if d.N >= lenn {
						//номер дендрита больше количества синапсов
						o.agent.log.Error(fmt.Sprintf("геном: %v, номер гена: %v", recs.filenameGens, recs.recs[k].Gen))
						c <- fmt.Errorf("номер аксона %v нейрона больше количества синапсов в поле %v!", d.N, recs.recs[k].SynNumber )
					}
				}
				//все в порядке
				c<-nil
			}
		}
	}
	//проверим преффекторы
	for _, ef:=range o.actions.effectors{
		for _, pre:= range ef.preffectors{
			for k:=0;k<len(pre.prefs);k++ {
				lenn := uint32(len(o.synapsesMap[pre.prefs[k].SynNumber].syn))
				//бежим по дендритам
				for _, d := range pre.prefs[k].Dendrites {
					if d.N >= lenn {
						//номер дендрита больше количества синапсов
						o.agent.log.Error(fmt.Sprintf("геном: %v, номер гена: %v", pre.filenameGens, pre.prefs[k].Gen))
						c <- fmt.Errorf("номер дендрита %v нейрона больше количества синапсов в поле %v!", d.N, pre.prefs[k].SynNumber )
					}
				}
				//все в порядке
				c<-nil
			}
		}
	}
	close(c)
}

//Live Организм начинает жить (комманда сверху)
//Вызывается через go
func (o *Organism) Live() {
//сюда попападаем только после инициализации организма
	for {
		select {
		case <- o.agent.sleep :
			//TODO сделать flush всему и остановиться
			//при попадании сюда, блокируемся и ждем комманды live или quit
			select {
			case <-o.agent.live:
				//комманда жить!
				//можно здесь ничего не делать, мы покинем select и попадем в default, где основная жизнь
			case <-o.agent.quit:
				//сюда лучше попадать после сна, иначе данные не сохранятся
				o.agent.wga.Done()
				return
			}
		case <-o.agent.quit:
			//сюда лучше попадать после сна, иначе данные не сохранятся
			o.agent.wga.Done()
			return
		default:
			//TODO главная работа начинается здесь
		}
	}
}

//Sleep - орагнизм идет спать (команда сверху)
func (o *Organism) Sleep(){

}

/*
ПРИМЕР ОРГАНИЗАЦИИ ФАЙЛОВ И ПАПОК ОРГАНИЗМА
\---Organism-1
    +---Actions
    |   +---Effector-0
    |   |       Data.data
    |   |       GenData.genes
    |   |       GenNeuron-0.genes
    |   |       GenPreffector-0.genes
    |   |       GenPreffector-1.genes
    |   |       Neuron-0.neurons
    |   |       Preffector-0.preffectors
    |   |       Preffector-1.preffectors
    |   |       Synapse-10.chemical
    |   |
    |   \---Effector-1
    |           Data.data
    |           GenData.genes
    |           GenPreffector-0.genes
    |           Preffector-0.preffectors
    |
    +---Brain
    |   +---Core-0
    |   |       GenNeuron-0.genes
    |   |       GenNeuron-1.genes
    |   |       GenNeuron-2.genes
    |   |       Neuron-0.neurons
    |   |       Neuron-1.neurons
    |   |       Neuron-2.neurons
    |   |       Synapse-3.chemical
    |   |
    |   +---Core-1
    |   |       GenNeuron-0.genes
    |   |       GenNeuron-1.genes
    |   |       GenNeuron-2.genes
    |   |       Neuron-0.neurons
    |   |       Neuron-1.neurons
    |   |       Neuron-2.neurons
    |   |       Synapse-4.chemical
    |   |
    |   +---Core-2
    |   |       GenNeuron-0.genes
    |   |       Neuron-0.neurons
    |   |       Synapse-5.chemical
    |   |
    |   +---Core-3
    |   |       GenNeuron-0.genes
    |   |       GenNeuron-1.genes
    |   |       GenNeuron-2.genes
    |   |       GenNeuron-3.genes
    |   |       Neuron-0.neurons
    |   |       Neuron-1.neurons
    |   |       Neuron-2.neurons
    |   |       Neuron-3.neurons
    |   |       Synapse-6.chemical
    |   |
    |   +---Core-4
    |   |       GenNeuron-0.genes
    |   |       GenNeuron-1.genes
    |   |       GenNeuron-2.genes
    |   |       GenNeuron-3.genes
    |   |       Neuron-0.neurons
    |   |       Neuron-1.neurons
    |   |       Neuron-2.neurons
    |   |       Neuron-3.neurons
    |   |       Synapse-7.chemical
    |   |
    |   +---Core-5
    |   |       GenNeuron-0.genes
    |   |       GenNeuron-1.genes
    |   |       GenNeuron-2.genes
    |   |       GenNeuron-3.genes
    |   |       Neuron-0.neurons
    |   |       Neuron-1.neurons
    |   |       Neuron-2.neurons
    |   |       Neuron-3.neurons
    |   |       Synapse-8.chemical
    |   |
    |   \---Core-6
    |           GenNeuron-0.genes
    |           Neuron-0.neurons
    |           Synapse-9.chemical
    |
    +---Senses
    |   |   GenNeuron-0.genes
    |   |   GenNeuron-1.genes
    |   |   GenNeuron-2.genes
    |   |   Neuron-0.neurons
    |   |   Neuron-1.neurons
    |   |   Neuron-2.neurons
    |   |   Synapse-2.chemical
    |   |
    |   +---Input-0
    |   |       Data.data
    |   |       GenData.genes
    |   |       GenNeuron-0.genes
    |   |       GenReceptor-0.genes
    |   |       GenReceptor-1.genes
    |   |       Neuron-0.neurons
    |   |       Receptor-0.receptors
    |   |       Receptor-1.receptors
    |   |       Synapse-0.chemical
    |   |
    |   +---Input-1
    |   |       Data.data
    |   |       GenData.genes
    |   |       GenReceptor-0.genes
    |   |       Receptor-0.receptors
    |   |
    |   +---Input-2
    |   |       Data.data
    |   |       GenData.genes
    |   |       GenNeuron-0.genes
    |   |       GenReceptor-0.genes
    |   |       Neuron-0.neurons
    |   |       Receptor-0.receptors
    |   |       Synapse-1.chemical
    |   |
    |   \---Input-3
    |           Data.data
    |           GenData.genes
    |           GenReceptor-0.genes
    |           GenReceptor-1.genes
    |           Receptor-0.receptors
    |           Receptor-1.receptors
    |
    \---Vegetatic
        |   GenNeuron-0.genes
        |   GenNeuron-1.genes
        |   Neuron-0.neurons
        |   Neuron-1.neurons
        |   Synapse-12.chemical
        |
        +---Effector-0
        |       Data.data
        |       GenData.genes
        |       GenNeuron-0.genes
        |       GenPreffector-0.genes
        |       GenPreffector-1.genes
        |       Neuron-0.neurons
        |       Preffector-0.preffectors
        |       Preffector-1.preffectors
        |       Synapse-11.chemical
        |
        \---Effector-1
                Data.data
                GenData.genes
                GenPreffector-0.genes
                Preffector-0.preffectors


*/
