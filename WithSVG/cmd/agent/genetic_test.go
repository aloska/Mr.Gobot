package agent

import (
	"testing"
)

func TestGenDataGenerate(t *testing.T) {
	filename := "c:/ALOSKA/my/tmp/GenData.genes"
	data := GenData{DATAUINT32BIG, 0, 0, 100, 1, 1, 0, 0}

	var datanew GenData
	if err := StructsFileWrite(filename, &data); err != nil {
		t.Error("can't generate GenData ", err)
	} else if err = StructsFileRead(filename, &datanew); err == nil {
		//fmt.Print(datanew)
		if datanew.Datatype != DATAUINT32BIG || datanew.Fps != 100 {
			t.Error("bad readed GenData, need fps==100, has ", datanew.Fps)
		}
	} else {
		t.Error("can't read GenData ", err)
	}

}
