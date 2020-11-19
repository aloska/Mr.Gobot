package agent

import (
	mmap "github.com/edsrzf/mmap-go"
	"os"
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
