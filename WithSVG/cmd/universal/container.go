package universal

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/edsrzf/mmap-go"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unsafe"
)

type Memory struct {
	filename  string    //имя файла, где записана память
	bytearray mmap.MMap //мапа на этот файл
	V         []int64   //файл в виде слайса
}

type Chromosome struct {
	filename  string    //имя файла, где записан ген
	bytearray mmap.MMap //мапа на этот файл
	Codons    []Codon   //файл в виде слайса
}

type Codon struct {
	Code Comm
	Op1  uint64
	Op2  uint64
	Op3  int64
}

type IO struct {
	filename  string    //имя файла, где записан вход или выход
	bytearray mmap.MMap //мапа на этот файл
	V         []int64
}

type Solution struct {
	Path string
	Proc []Processor
	Chrom  []Chromosome
	Mem  []Memory
	In   []IO
	Out  []IO
	IsAsync bool
}

//для создания новых решателей из JSON
type Serialisator struct {
	Memories []int64 //слайс размеров памятей
	Ins      []int64 //слайс размеров входов
	Outs     []int64 //слайс размеров выходов
	Genes	 []string //файл с описанием генома (должен лежать в той же папке, что и json-файл)
	IsAsync	 bool	//асинхроный старт хромосом или последовательный

	/*файлы хромосом можут быть  с расширением .codons:
	codons: 22 0 0 2; 2 1 0 -14; ... (можно в несколько строк и без ';' )
	в этом случае все числа парсятся до возможности исполнения. Например, если комманды с кодом нет, то комманда генерится
	из остатка от деления на количество комманд

	или c расширением .sasm на ассемблере:

	asm:
		LDIN x0, 0, 2
		ADDI x1, x0, -14
		...

	Ассемблерный код должен быть валидным, поскольку пишется человеком или для человека

	Есть функция, переводящая любой ген (с правильными или неправильными кодонами) в валидный ассемблер.
	Если решатель решает задачу нормально - можно руками заоптимизировать его код и сделать новый решатель с отредактированным кодом,
	который будет работать быстрее.
	Для этого новый код нужно транслировать в файл genom и подставить вместо существующего
	*/
}

//парсинг кодонов со строки
func GetCodonsFromChromosomeString(gs *string) (*[]Codon, error){
	var cods []Codon

	fields:=strings.Fields(*gs)
	if len(fields)<5{
		return nil, errors.New("плохой формат данных генома: должно начинаться с 'codons:' или 'asm:' и далее не менее 4 полей")
	}else if fields[0]=="codons:"{
		i:=1
		state:=0
		for i<len(fields){
			str:=strings.TrimPrefix(strings.TrimSuffix(fields[i],";"),";")
			switch state{
			case 0://ожидаем комманду
				cods=append(cods,Codon{})
				val, err:=strconv.ParseInt(str, 10, 64)
				if err!=nil{
					return nil, err
				}
				cods[len(cods)-1].Code=Comm(val)
				state++
			case 1://первый операнд
				val, err:=strconv.ParseInt(str, 10, 64)
				if err!=nil{
					return nil, err
				}
				cods[len(cods)-1].Op1=uint64(val)
				state++
			case 2://2 операнд
				val, err:=strconv.ParseInt(str, 10, 64)
				if err!=nil{
					return nil, err
				}
				cods[len(cods)-1].Op2=uint64(val)
				state++
			case 3://3 операнд
				val, err:=strconv.ParseInt(str, 10, 64)
				if err!=nil{
					return nil, err
				}
				cods[len(cods)-1].Op3=val
				state=0
			}
			i++
		}

	}else {//может файл на асме?
		return GetCodonsFromAsmString(gs)
	}


	return &cods,nil
}

//парсинг кодонов из файла
func GetCodonsFromFile(filename string) (*[]Codon, error){

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	str := string(b)
	return GetCodonsFromChromosomeString(&str)
}

//NewSolution - создать нового решателя
//filejson - файл, с описанием решателя в формате json Serialisator
//в той же папке, где лежит этот файл, будет создана папка с именем как у файла (без расширения), с приставкой SOL- с новым решателем внутри
func NewSolution(filejson string) (*Solution, error){
	if !fileExists(filejson){
		log.Println("Нет файла описания "+filejson)
		return nil, errors.New("Нет файла описания ")
	}

	newfolder:=filepath.Dir(filejson)+"/SOL-"+strings.TrimSuffix(filepath.Base(filejson),filepath.Ext(filepath.Base(filejson)))
	if folderExists(newfolder){
		log.Println("Папка решателя уже существует! "+newfolder)
		return nil, errors.New("Папка решателя уже существует!")
	}

	file, err := ioutil.ReadFile(filejson)
	if err!=nil{
		return nil, err
	}

	ser:=Serialisator{}
	err = json.Unmarshal(file, &ser)
	if err!=nil{
		return nil,err
	}

	if len(ser.Memories)==0{
		return nil, errors.New("память должна быть хотя бы одна")
	}

	if len(ser.Ins)==0{
		return nil, errors.New("вход должен быть хотя бы один")
	}

	if len(ser.Outs)==0{
		return nil, errors.New("выход должен быть хотя бы один")
	}

	err=os.Mkdir(newfolder, os.ModePerm)
	if err!=nil{
		return nil, err
	}

	//парсим геном
	for i:=0;i<len(ser.Genes);i++ {
		chrom, err := GetCodonsFromFile(filepath.Dir(filejson) + "/" + ser.Genes[i])
		if err != nil {
			return nil, err
		}
		//если удачно отпарсили - создаем файл хромосомы
		err = StructsFileWrite(newfolder+"/"+strconv.Itoa(i)+".chromosome", chrom, binary.LittleEndian)
		if err != nil {
			return nil, err
		}
	}

	//создаем файлы процессоров
	for i:=0;i<len(ser.Genes);i++ {
		proc := Processor{}
		err = StructsFileWrite(newfolder+"/"+strconv.Itoa(i)+".processor", &proc, binary.LittleEndian)
		if err != nil {
			return nil, err
		}
	}

	//создаем файлы памяти
	for i:=0;i<len(ser.Memories);i++{
		if ser.Memories[i]<=0{
			log.Println("Неверный размер файла памяти в описании рещшателя: ", ser.Memories[i])
			return nil, errors.New("неверный размер файла памяти в описании рещшателя")
		}
		f, err:=os.Create(newfolder+"/"+strconv.Itoa(i)+".memory")
		if err!=nil{
			return nil, err
		}
		err=f.Truncate(ser.Memories[i]*8)
		if err!=nil{
			return nil, err
		}
		f.Close()
	}
	//создаем файлы входов
	for i:=0;i<len(ser.Ins);i++{
		if ser.Ins[i]<=0{
			log.Println("Неверный размер файла входа в описании рещшателя: ", ser.Ins[i])
			return nil, errors.New("неверный размер файла входа в описании рещшателя")
		}
		f, err:=os.Create(newfolder+"/"+strconv.Itoa(i)+".in")
		if err!=nil{
			return nil, err
		}
		err=f.Truncate(ser.Ins[i]*8)
		if err!=nil{
			return nil, err
		}
		f.Close()
	}
	//создаем файлы вЫходов
	for i:=0;i<len(ser.Outs);i++{
		if ser.Outs[i]<=0{
			log.Println("Неверный размер файла вЫхода в описании решателя: ", ser.Outs[i])
			return nil, errors.New("неверный размер файла вЫхода в описании рещшателя")
		}
		f, err:=os.Create(newfolder+"/"+strconv.Itoa(i)+".out")
		if err!=nil{
			return nil, err
		}
		err=f.Truncate(ser.Outs[i]*8)
		if err!=nil{
			return nil, err
		}
		f.Close()
	}

	//все было удачно - инициализируем решатель из созданных файлов и возвращаем указатель на него
	sol:=Solution{}
	if err=sol.Init(newfolder); err!=nil{
		return nil, err
	}

	sol.IsAsync=ser.IsAsync

	return &sol,nil
}



//инициализировать Solution из директории
func (so *Solution) Init(path string) error {
	//в директории должны быть обязательно файлы "0.chromosome", "0.processor", "0.memory","0.in","0.out"

	mems := []string{}
	ins := []string{}
	outs := []string{}
	procfiles := []string{}
	genfiles := []string{}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _, file := range files {
		if match, _ := regexp.MatchString("([0-9]+.processor)",
			file.Name()); match {
			procfiles = append (procfiles,file.Name())
		} else if match, _ := regexp.MatchString("([0-9]+.chromosome)",
			file.Name()); match {
			genfiles = append(genfiles, file.Name())
		} else if match, _ := regexp.MatchString("([0-9]+.memory)",
			file.Name()); match {
			mems = append(mems, file.Name())
		} else if match, _ := regexp.MatchString("([0-9]+.in)",
			file.Name()); match {
			ins = append(ins, file.Name())
		} else if match, _ := regexp.MatchString("([0-9]+.out)",
			file.Name()); match {
			outs = append(outs, file.Name())
		}
	}
	if len(ins) == 0 || len(outs) == 0 || len(mems) == 0 || len(procfiles) == 0 || len(genfiles) == 0 || len(procfiles)!=len(genfiles){
		return errors.New("не хватает файла или количество процессоров меньше количества хромосом. Обязательно должны быть \"0.chromosome\", \"0.processor\", \"0.memory\",\"0.in\",\"0.out\"")
	}

	if !sort.StringsAreSorted(genfiles) {
		sort.Strings(genfiles)
	}
	for i := 0; i < len(genfiles); i++ {
		so.Chrom = append(so.Chrom, Chromosome{})
		if err = so.Chrom[len(so.Chrom)-1].Init(path + "/" + genfiles[i]); err != nil {
			return err
		}
	}

	if !sort.StringsAreSorted(procfiles) {
		sort.Strings(procfiles)
	}
	for i := 0; i < len(procfiles); i++ {
		so.Proc = append(so.Proc, Processor{})
		if err = StructsFileRead(path+"/"+procfiles[i], &so.Proc[i], binary.LittleEndian); err != nil {
			return err
		}
	}

	if !sort.StringsAreSorted(mems) {
		sort.Strings(mems)
	}
	for i := 0; i < len(mems); i++ {
		so.Mem = append(so.Mem, Memory{})
		if err := so.Mem[len(so.Mem)-1].Init(path + "/" + mems[i]); err != nil {
			return err
		}
	}

	if !sort.StringsAreSorted(ins) {
		sort.Strings(ins)
	}
	for i := 0; i < len(ins); i++ {
		so.In = append(so.In, IO{})
		if err := so.In[len(so.In)-1].Init(path + "/" + ins[i]); err != nil {
			return err
		}
	}

	if !sort.StringsAreSorted(outs) {
		sort.Strings(outs)
	}
	for i := 0; i < len(outs); i++ {
		so.Out = append(so.Out, IO{})
		if err := so.Out[len(so.Out)-1].Init(path + "/" + outs[i]); err != nil {
			return err
		}
	}

	so.Path = path

	return nil
}

//из файла
func (io *IO) Init(fs string) error {
	var header reflect.SliceHeader

	f, err := openFile(os.O_RDWR, fs)
	if err != nil {
		return err
	}
	io.bytearray, err = mmap.Map(f, mmap.RDWR, 0)
	if err != nil {
		return err
	}

	io.filename = fs

	header.Data = (uintptr)(unsafe.Pointer(&io.bytearray[0]))
	header.Len = len(io.bytearray) / int(reflect.TypeOf(io.V).Elem().Size())
	header.Cap = header.Len
	io.V = *(*[]int64)(unsafe.Pointer(&header))

	io.filename = fs
	return nil
}

//из файла
func (m *Memory) Init(fs string) error {
	var header reflect.SliceHeader

	f, err := openFile(os.O_RDWR, fs)
	if err != nil {
		return err
	}
	m.bytearray, err = mmap.Map(f, mmap.RDWR, 0)
	if err != nil {
		return err
	}

	m.filename = fs

	header.Data = (uintptr)(unsafe.Pointer(&m.bytearray[0]))
	header.Len = len(m.bytearray) / int(reflect.TypeOf(m.V).Elem().Size())
	header.Cap = header.Len
	m.V = *(*[]int64)(unsafe.Pointer(&header))

	m.filename = fs
	return nil
}

//из файла
func (g *Chromosome) Init(fs string) error {
	var header reflect.SliceHeader

	f, err := openFile(os.O_RDWR, fs)
	if err != nil {
		return err
	}
	g.bytearray, err = mmap.Map(f, mmap.RDWR, 0)
	if err != nil {
		return err
	}

	g.filename = fs

	header.Data = (uintptr)(unsafe.Pointer(&g.bytearray[0]))
	header.Len = len(g.bytearray) / int(reflect.TypeOf(g.Codons).Elem().Size())
	header.Cap = header.Len
	g.Codons = *(*[]Codon)(unsafe.Pointer(&header))

	g.filename = fs
	return nil
}

func (so *Solution)Save(){
	for i:=0;i<len(so.In);i++{
		so.In[i].bytearray.Flush()
	}
	for i:=0;i<len(so.Mem);i++{
		so.Mem[i].bytearray.Flush()
	}
	for i:=0;i<len(so.Out);i++{
		so.Out[i].bytearray.Flush()
	}

	for i:=0;i<len(so.Proc);i++ {
		StructsFileOverwrite(so.Path+"/"+strconv.Itoa(i)+".processor", &so.Proc, binary.LittleEndian)
	}
}

func (so *Solution)Exit(){
	for i:=0;i<len(so.In);i++{
		so.In[i].bytearray.Unmap()
	}
	for i:=0;i<len(so.Mem);i++{
		so.Mem[i].bytearray.Unmap()
	}
	for i:=0;i<len(so.Out);i++{
		so.Out[i].bytearray.Unmap()
	}
}
