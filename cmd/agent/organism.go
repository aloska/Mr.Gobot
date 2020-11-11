package agent

import (
	mmap "github.com/edsrzf/mmap-go"
)

/*Synapses - структура с файлом синапсов, хранящая все  Chemistry, самая интенсивно используемая

 */
type Synapses struct {
	number    uint16     //уникальный номер синаптического поля (ядра или входа или выхода)
	filename  string     //имя файла, где записаны синапсы
	bytearray mmap.MMap  //замапленный файл синапсов
	syn       []Chemical //тот же файл в удобном виде, как слайс Chemical с веществами
	maxX      uint32     //ширина синаптического поля
	maxY      uint32     //высота
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

/*Core - ядро организма. Может быть много ядер
Может хранить только один файл синапсов и сколько угодно файлов клеток и генов

*/
type Core struct {
	path     string    //папка, где расположены все файлы ядры (за пределами этой папки другие ядра)
	synapses Synapses  //обслуживает файл с синапсами этого ядра (у каждого ядра только один файл синапсов)
	cells    []Cells   //слайс обслуживает файлы с клетками и геномами этих клеток
	organism *Organism //ссылка на весь родительский организм
}

/*Brain - большая структура с целым мозгом из всех ядер, входящих в его состав
 */
type Brain struct {
	path     string    //папка, где расположены все ядра, а входные и выходные устройства входят в состав организма
	cores    []Core    //мозг состоит из ядер
	organism *Organism //ссылка на весь родительский организм
}

/*DataInput - описание входа
Поскольку в Го нет женериков, прийдется хранить поля всех видов данных
*/
type DataInput struct {
	typeData uint16 //тип входа 1-DataByte...

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
	dataBit    []DataBit
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

/*Senses - все входы (ощущения)
 */
type Senses struct {
	path   string  //папка, где расположены все входы
	inputs []Input //все входы

	synapsesInputs Synapses //синаптическое поле всех входов
	cellsInputs    []Cells  /*слайс обслуживает файлы с нейронами и геномами этих нейронов, если они заданы, для общего синаптичесского поля всех входов
	этих нейронов может не быть, и тогда это значит, что общего синаптического поля входов нет
	такое поведение может использоваться для очень простых агентов
	*/
	organism *Organism //ссылка на весь родительский организм
}

/*DataOutput - описание выхода
 */
type DataOutput struct {
	typeData uint16 //тип выхода

	filenameGen  string      //имя файла, где записан ген выхода
	bytearrayGen mmap.MMap   //замапленный файл гена
	genData      *GenDataOut //тот же файл ввиде гена (у ячеек данных он один)

	filenameData  string    //имя файла, где записаны ячейки выхода
	bytearrayData mmap.MMap //замапленный файл ячеек

	//должен быть список всех ТИПОВ выходов, используемых в агенте
	datauint32 []Datauint32
	dataRune   []DataRune
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

/*Actions - все выходы (действия организма, кроме внутренних)
 */
type Actions struct {
	path      string     //папка, где расположены все входы
	effectors []Effector //все выходы

	organism *Organism //ссылка на весь родительский организм
}

/*Organism - самая полная структура, состоящая из мозга, входных и выходных устройств

 */
type Organism struct {
	brain   Brain   //мозг из ядер
	senses  Senses  //ощущения (входы, рецепторы...)
	actions Actions //действия
}

//Live ...
func (o *Organism) Live() {

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
    \---Senses
        |   GenNeuron-0.genes
        |   GenNeuron-1.genes
        |   GenNeuron-2.genes
        |   Neuron-0.neurons
        |   Neuron-1.neurons
        |   Neuron-2.neurons
        |   Synapse-2.chemical
        |
        +---Input-0
        |       Data.data
        |       GenData.genes
        |       GenNeuron-0.genes
        |       GenReceptor-0.genes
        |       GenReceptor-1.genes
        |       Neuron-0.neurons
        |       Receptor-0.receptors
        |       Receptor-1.receptors
        |       Synapse-0.chemical
        |
        +---Input-1
        |       Data.data
        |       GenData.genes
        |       GenReceptor-0.genes
        |       Receptor-0.receptors
        |
        +---Input-2
        |       Data.data
        |       GenData.genes
        |       GenNeuron-0.genes
        |       GenReceptor-0.genes
        |       Neuron-0.neurons
        |       Receptor-0.receptors
        |       Synapse-1.chemical
        |
        \---Input-3
                Data.data
                GenData.genes
                GenReceptor-0.genes
                GenReceptor-1.genes
                Receptor-0.receptors
                Receptor-1.receptors


*/
