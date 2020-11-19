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
