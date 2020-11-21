package agent

import (
	"encoding/binary"
	"fmt"
	"os"
)



type structsFileReaderWriter interface{}

/*StructsFileWrite - пишет любые данные в файл пачкой. Файл не должен существовать!
Использование:
StructsFileWrite("/путь/к/файлу", &any) - передача по указателю
и то же самое, но значением:
StructsFileWrite("/путь/к/файлу", any)
Даже слайс структур, которые содержат слайсы можно зафигачить!
Протестировано на:
GenData
*/
func StructsFileWrite(filename string, fw structsFileReaderWriter, order binary.ByteOrder) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		f, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("can't create file: %v", filename)
		}
		defer f.Close()

		err = binary.Write(f, order, fw)
		if err != nil {
			return fmt.Errorf("can't write to file: %v", filename)
		}
	} else {
		return fmt.Errorf("file already exists: %v", filename)
	}
	return nil
}

/*StructsFileRead - читает любые данные из файла пачкой.
Размер данных должен быть определен!!! Нельзя записать в нерожденный слайс! нужен make([],len,cap) для слайса предварительно
Протестировано на:
GenData
*/
func StructsFileRead(filename string, fw structsFileReaderWriter, order binary.ByteOrder) error {
	if _, err := os.Stat(filename); err == nil {
		if file, err := os.Open(filename); err == nil {
			defer file.Close()
			err = binary.Read(file, order, fw)

			if err == nil {
				return nil
			}
			return err

		}
		return fmt.Errorf("can't open file %v", err)

	} else if os.IsNotExist(err) {
		return fmt.Errorf("file doesn't exist: %v", filename)

	}
	return fmt.Errorf("unknown error %v", filename)
}

//StructsFileReadEOF - читает любые данные до io.EOF - можно в цикле читать, пока ошибка не будет EOF,
//File должен быть открыт
func StructsFileReadEOF(file *os.File, fw structsFileReaderWriter, order binary.ByteOrder) error {
	err:= binary.Read(file, order, fw)
	if err == nil {
			return nil
		}
	return err
}


func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}