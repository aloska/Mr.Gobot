package agent

import (
	mmap "github.com/edsrzf/mmap-go"
	"os"
	"reflect"
	"unsafe"
)

func openFile(flags int, filepath string) (*os.File, error) {
	f, err := os.OpenFile(filepath, flags, 0644)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (d *DataInput) mmapGenData() error{
	//открываем файл
	f, err := openFile(os.O_RDWR, d.filenameGen)
	if err!=nil{
		return err
	}
	//создаем mmap
	d.bytearrayGen, err = mmap.Map(f, mmap.RDWR, 0)
	if err!=nil{
		return err
	}
	//делаем unsafe на структуру
	d.genData=(*GenData)(unsafe.Pointer(&d.bytearrayGen[0]))
	return nil
}

func (d *DataOutput) mmapGenData() error{
	//открываем файл
	f, err := openFile(os.O_RDWR, d.filenameGen)
	if err!=nil{
		return err
	}
	//создаем mmap
	d.bytearrayGen, err = mmap.Map(f, mmap.RDWR, 0)
	if err!=nil{
		return err
	}
	//делаем unsafe на структуру
	d.genData=(*GenDataOut)(unsafe.Pointer(&d.bytearrayGen[0]))
	return nil
}

func (d *DataInput) mmapData() error{
	//от типа ячеек зависит размер файла с данными
	var (
		size int64
		dataUInt32 DataUInt32
		datauint32 Datauint32
		header reflect.SliceHeader //TODO - возможно его надо в саму структуру положить? Если мусорщик удалит, потеряем контроль над стурктурой
	)
	//ВНИМАНИЕ!! такой же свитч есть и дальше. Если сюда добавляешь, то и туда добавь!
	switch d.genData.Datatype {
	case DATAUINT32BIG:
		size=int64(d.genData.Len)*int64(unsafe.Sizeof(dataUInt32))
		break
	case DATAUINT32:
		size=int64(d.genData.Len)*int64(unsafe.Sizeof(datauint32))
		break
	}
	//проверяем, есть ли уже файл с данными
	if !fileExists(d.filenameData){
		//создаем файл данных
		f, err:=os.Create(d.filenameData)
		if err!=nil{
			return err
		}
		err=f.Truncate(size)
		if err!=nil{
			return err
		}
		f.Close()
	}
	//открываем файл
	f, err := openFile(os.O_RDWR, d.filenameData)
	if err!=nil{
		return err
	}
	//создаем mmap
	d.bytearrayData, err = mmap.Map(f, mmap.RDWR, 0)
	if err!=nil{
		return err
	}

	//делаем unsafe на структуру (магия и только!)
	header.Data =(uintptr)(unsafe.Pointer(&d.bytearrayData[0]))
	switch d.genData.Datatype {
	case DATAUINT32BIG:
		header.Len = len(d.bytearrayData)/int(reflect.TypeOf(d.dataUInt32).Elem().Size())
		header.Cap = header.Len
		d.dataUInt32=*(*[]DataUInt32)(unsafe.Pointer(&header))
		break
	case DATAUINT32:
		header.Len = len(d.bytearrayData)/int(reflect.TypeOf(d.datauint32).Elem().Size())
		header.Cap = header.Len
		d.datauint32=*(*[]Datauint32)(unsafe.Pointer(&header))
		break
	}



	return nil
}

//близнецы с func (d *DataInput) mmapData()
func (d *DataOutput) mmapData() error{
	//от типа ячеек зависит размер файла с данными
	var (
		size int64
		dataUInt32 DataUInt32
		datauint32 Datauint32
		header reflect.SliceHeader //TODO - возможно его надо в саму структуру положить? Если мусорщик удалит, потеряем контроль над стурктурой
	)
	//ВНИМАНИЕ!! такой же свитч есть и дальше. Если сюда добавляешь, то и туда добавь!
	switch d.genData.Datatype {
	case DATAUINT32BIG:
		size=int64(d.genData.Len)*int64(unsafe.Sizeof(dataUInt32))
		break
	case DATAUINT32:
		size=int64(d.genData.Len)*int64(unsafe.Sizeof(datauint32))
		break
	}
	//проверяем, есть ли уже файл с данными
	if !fileExists(d.filenameData){
		//создаем файл данных
		f, err:=os.Create(d.filenameData)
		if err!=nil{
			return err
		}
		err=f.Truncate(size)
		if err!=nil{
			return err
		}
		f.Close()
	}
	//открываем файл
	f, err := openFile(os.O_RDWR, d.filenameData)
	if err!=nil{
		return err
	}
	//создаем mmap
	d.bytearrayData, err = mmap.Map(f, mmap.RDWR, 0)
	if err!=nil{
		return err
	}

	//делаем unsafe на структуру (магия и только!)
	header.Data =(uintptr)(unsafe.Pointer(&d.bytearrayData[0]))
	switch d.genData.Datatype {
	case DATAUINT32BIG:
		header.Len = len(d.bytearrayData)/int(reflect.TypeOf(d.dataUInt32).Elem().Size())
		header.Cap = header.Len
		d.dataUInt32=*(*[]DataUInt32)(unsafe.Pointer(&header))
		break
	case DATAUINT32:
		header.Len = len(d.bytearrayData)/int(reflect.TypeOf(d.datauint32).Elem().Size())
		header.Cap = header.Len
		d.datauint32=*(*[]Datauint32)(unsafe.Pointer(&header))
		break
	}
	return nil
}

func (re *Receptors) mmapGenReceptor() error{
	//открываем файл
	f, err := openFile(os.O_RDWR, re.filenameGens)
	if err!=nil{
		return err
	}
	//создаем mmap
	re.bytearrayGens, err = mmap.Map(f, mmap.RDWR, 0)
	if err!=nil{
		return err
	}
	//делаем unsafe на структуру
	var header reflect.SliceHeader //TODO - возможно его надо в саму структуру положить? Если мусорщик удалит, потеряем контроль над стурктурой
	header.Data =(uintptr)(unsafe.Pointer(&re.bytearrayGens[0]))
	header.Len = len(re.bytearrayGens)/int(reflect.TypeOf(re.genes).Elem().Size())
	header.Cap = header.Len
	re.genes=*(*[]GenReceptor)(unsafe.Pointer(&header))

	return nil
}

func (pre *Preffectors) mmapGenPreffector() error{
	//открываем файл
	f, err := openFile(os.O_RDWR, pre.filenameGens)
	if err!=nil{
		return err
	}
	//создаем mmap
	pre.bytearrayGens, err = mmap.Map(f, mmap.RDWR, 0)
	if err!=nil{
		return err
	}
	//делаем unsafe на структуру
	var header reflect.SliceHeader //TODO - возможно его надо в саму структуру положить? Если мусорщик удалит, потеряем контроль над стурктурой
	header.Data =(uintptr)(unsafe.Pointer(&pre.bytearrayGens[0]))
	header.Len = len(pre.bytearrayGens)/int(reflect.TypeOf(pre.genes).Elem().Size())
	header.Cap = header.Len
	pre.genes=*(*[]GenPreffector)(unsafe.Pointer(&header))

	return nil
}

func (re *Receptors) mmapReceptors() error{
	//открываем файл
	f, err := openFile(os.O_RDWR, re.filenameRecs)
	if err!=nil{
		return err
	}
	//создаем mmap
	re.bytearrayRecs, err = mmap.Map(f, mmap.RDWR, 0)
	if err!=nil{
		return err
	}
	//делаем unsafe на структуру
	var header reflect.SliceHeader //TODO - возможно его надо в саму структуру положить? Если мусорщик удалит, потеряем контроль над стурктурой
	header.Data =(uintptr)(unsafe.Pointer(&re.bytearrayRecs[0]))
	header.Len = len(re.bytearrayRecs)/int(reflect.TypeOf(re.recs).Elem().Size())
	header.Cap = header.Len
	re.recs=*(*[]Receptor)(unsafe.Pointer(&header))

	return nil
}

func (pre *Preffectors) mmapPreffectors() error{
	//открываем файл
	f, err := openFile(os.O_RDWR, pre.filenamePres)
	if err!=nil{
		return err
	}
	//создаем mmap
	pre.bytearrayPres, err = mmap.Map(f, mmap.RDWR, 0)
	if err!=nil{
		return err
	}
	//делаем unsafe на структуру
	var header reflect.SliceHeader //TODO - возможно его надо в саму структуру положить? Если мусорщик удалит, потеряем контроль над стурктурой
	header.Data =(uintptr)(unsafe.Pointer(&pre.bytearrayPres[0]))
	header.Len = len(pre.bytearrayPres)/int(reflect.TypeOf(pre.prefs).Elem().Size())
	header.Cap = header.Len
	pre.prefs=*(*[]Preffector)(unsafe.Pointer(&header))

	return nil
}

func (ce *Cells) mmapGenNeuron() error{
	//открываем файл
	f, err := openFile(os.O_RDWR, ce.filenameGens)
	if err!=nil{
		return err
	}
	//создаем mmap
	ce.bytearrayGens, err = mmap.Map(f, mmap.RDWR, 0)
	if err!=nil{
		return err
	}
	//делаем unsafe на структуру
	var header reflect.SliceHeader //TODO - возможно его надо в саму структуру положить? Если мусорщик удалит, потеряем контроль над стурктурой
	header.Data =(uintptr)(unsafe.Pointer(&ce.bytearrayGens[0]))
	header.Len = len(ce.bytearrayGens)/int(reflect.TypeOf(ce.genes).Elem().Size())
	header.Cap = header.Len
	ce.genes=*(*[]GenNeuron)(unsafe.Pointer(&header))

	return nil
}

func (ce *Cells) mmapNeuron() error{
	//открываем файл
	f, err := openFile(os.O_RDWR, ce.filenameCells)
	if err!=nil{
		return err
	}
	//создаем mmap
	ce.bytearrayCells, err = mmap.Map(f, mmap.RDWR, 0)
	if err!=nil{
		return err
	}
	//делаем unsafe на структуру
	var header reflect.SliceHeader //TODO - возможно его надо в саму структуру положить? Если мусорщик удалит, потеряем контроль над стурктурой
	header.Data =(uintptr)(unsafe.Pointer(&ce.bytearrayCells[0]))
	header.Len = len(ce.bytearrayCells)/int(reflect.TypeOf(ce.neurons).Elem().Size())
	header.Cap = header.Len
	ce.neurons=*(*[]Neuron)(unsafe.Pointer(&header))

	return nil
}

func (sy *Synapses) mmapTypicalChe() error{
	//открываем файл
	f, err := openFile(os.O_RDWR, sy.filedesc)
	if err!=nil{
		return err
	}
	//создаем mmap
	sy.bytearrayTypicalChe, err = mmap.Map(f, mmap.RDWR, 0)
	if err!=nil{
		return err
	}
	//делаем unsafe на структуру
	sy.TypicalChe=(*Chemical)(unsafe.Pointer(&sy.bytearrayTypicalChe[0]))
	return nil
}

func (sy *Synapses) mmapSynapse() error{
	//проверяем, есть ли уже файл с данными
	isneedfill:=false
	if !fileExists(sy.filename){
		//создаем файл данных
		f, err:=os.Create(sy.filename)
		if err!=nil{
			return err
		}
		err=f.Truncate(int64(sy.maxX)*int64(sy.maxY)*int64(unsafe.Sizeof(*sy.TypicalChe)))
		if err!=nil{
			return err
		}
		f.Close()
		isneedfill=true //надо заполнить после после маппинга, потому что оно нулячее
	}

	//открываем файл
	f, err := openFile(os.O_RDWR, sy.filename)
	if err!=nil{
		return err
	}
	//создаем mmap
	sy.bytearray, err = mmap.Map(f, mmap.RDWR, 0)
	if err!=nil{
		return err
	}
	//делаем unsafe на структуру
	var header reflect.SliceHeader //TODO - возможно его надо в саму структуру положить? Если мусорщик удалит, потеряем контроль над стурктурой
	header.Data =(uintptr)(unsafe.Pointer(&sy.bytearray[0]))
	header.Len = len(sy.bytearray)/int(reflect.TypeOf(sy.syn).Elem().Size())
	header.Cap = header.Len
	sy.syn=*(*[]Chemical)(unsafe.Pointer(&header))

	//теперь синапсы заливаем содержимым из типичной ячейки, потому что они только что созданы
	if isneedfill {
		for i:=0;i<len(sy.syn);i++{
			sy.syn[i]=*sy.TypicalChe
		}
		sy.bytearray.Flush()
	}
	return nil
}