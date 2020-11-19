package agent

import "encoding/binary"

func createGenes() {

}

func createSensesGenes() {

}

func createBrainGenes() {

}

func createEffectorsGenes() {

}

func easyOrganism() {
	const (
		path = "c:/ALOSKA/my/tmp/testOrganism/"
	)
	var (
		genData     GenData
		genReceptor GenReceptor
	)
	genData.Datatype = DATAUINT32BIG
	genData.Dataneed = 0
	genData.Runifchange = 0
	genData.Fps = 10
	genData.Httpchan = 0
	genData.Len = 1
	genData.Serv1, genData.Serv2 = 0, 0
	StructsFileWrite(path+"Senses/Input-0/GenData.genes", &genData, binary.LittleEndian)

	genReceptor.Typer = RECEPTORDATAUINT32BIGPOS
	genReceptor.Typemedi = NEURONACETILHOLIN
	genReceptor.Coren = 0

}
