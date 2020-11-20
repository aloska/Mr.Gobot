package main

import (
	"WithSVG/cmd/agent"
	"encoding/binary"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	r         = gin.Default()
	greceptors []agent.GenReceptor
	gneurons []agent.GenNeuron
	gpreffectors []agent.GenPreffector
	globalErr =false
	scale=2
	ViewX=64
	ViewY=64
	RdrawConnector = 0
	NdrawConnector = [2]int{0,0}
	PdrawConnector=0
	gridshow=true
	gridX=0
	gridY=0
	gridW=64
	gridH=64
	gridN=8

)
func main() {
	getRoutes() //издеся маршруты роутим
	r.Run()
}

func getRoutes() {

	r.LoadHTMLFiles("view/index.html", "view/index2.html", "view/index3.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", map[string]interface{}{})
	})
	r.GET("/i2", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index2.html", map[string]interface{}{})
	})
	r.GET("/i3", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index3.html", map[string]interface{}{})
	})
	r.StaticFile("styles.css","view/styles.css")
	r.StaticFile("plain-draggable.min.js","view/plain-draggable.min.js")
	r.POST("/hello", helloPage)
	r.GET("/hello", helloPage)

	r.POST("/receptor-gen", receptorgen)
	r.POST("/neuron-gen", neurongen)
	r.POST("/preffector-gen", preffectorgen)
	r.POST("/genfiles-generate",genfilesgenerate)
	r.POST("/gendatain-generate",gendataingenerate)
	r.POST("/gendataout-generate",gendataoutgenerate)

	r.POST("/set-settings", setsettings)
	r.POST("/del-entities", delentities)

	r.GET("/hellopng",hellopng)
	r.GET("/ajax/:command", getajax)

	r.POST("/set-draw-connectors",setdrawconnectors)
	r.POST("/set-draw-grid",setdrawgrid)

	r.POST("/read-genes", readgenes)
}

func readgenes(c *gin.Context){
	filename:=strings.TrimSpace(c.PostForm("GenePath"))

	if _, err := os.Stat(filename); err == nil {
		if file, err := os.Open(filename); err == nil {
			defer file.Close()

			switch strings.TrimSpace(c.PostForm("TypeGene")) {
			case "R":
				err=nil
				for i:=0; err==nil;i++   {
					rec:=agent.GenReceptor{}
					err= agent.StructsFileReadEOF(file,&rec,binary.LittleEndian)
					if err==nil{
						greceptors = append(greceptors, rec)
					}
				}
				break
			case "N":
				err=nil
				for i:=0; err==nil;i++   {
					rec:=agent.GenNeuron{}
					err= agent.StructsFileReadEOF(file,&rec,binary.LittleEndian)
					if err==nil{
						gneurons = append(gneurons, rec)
					}
				}
				break
			case "P":
				err=nil
				for i:=0; err==nil;i++   {
					rec:=agent.GenPreffector{}
					err= agent.StructsFileReadEOF(file,&rec,binary.LittleEndian)
					if err==nil{
						gpreffectors = append(gpreffectors, rec)
					}
				}
				break
			}
			drawall(c)
		}else{
			svgError(c,"Не могу открыть файл")
			return
		}
	} else if os.IsNotExist(err) {
		svgError(c,"нет такого файла")
	}
}

func getajax(c* gin.Context){
	command := c.Param("command")
	switch command {
	case "countgenes":
		c.String(http.StatusOK,"["+strconv.Itoa(len(greceptors))+", "+strconv.Itoa(len(gneurons))+", "+strconv.Itoa(len(gpreffectors))+"]")
		break
	}
}

func setdrawconnectors(c* gin.Context){
	rslice:=strings.Split(c.PostForm("RDrawConn"),",")
	nslice:=strings.Split(c.PostForm("NDrawConn"),",")
	pslice:=strings.Split(c.PostForm("PDrawConn"),",")
	RdrawConnector=int(getfromstringslice(rslice,c,0))
	NdrawConnector[0]= int(getfromstringslice(nslice,c,0))
	NdrawConnector[1]= int(getfromstringslice(nslice,c,1))
	PdrawConnector=int(getfromstringslice(pslice,c,0))
	if globalErr{
		globalErr=false
		return
	}

	drawall(c)
}

func setdrawgrid(c *gin.Context){
	gslice:=strings.Split(c.PostForm("DrawGrid"),",")
	gridX=int(getfromstringslice(gslice,c,0))
	gridY=int(getfromstringslice(gslice,c,1))
	gridW=int(getfromstringslice(gslice,c,2))
	gridH=int(getfromstringslice(gslice,c,3))
	gridN=int(getfromstringslice(gslice,c,4))
	if globalErr{
		globalErr=false
		return
	}

	gridshow=!gridshow
	drawall(c)
}

func gendataingenerate(c *gin.Context)  {
	dslice:=strings.Split(c.PostForm("Dataf"),",")
	data:=agent.GenData{
		Datatype: agent.DataTypeEnum(getfromstringslice(dslice, c, 0)),
		Dataneed: byte(getfromstringslice(dslice, c, 1)),
		Runifchange: byte(getfromstringslice(dslice, c, 2)),
		Fps: uint16(getfromstringslice(dslice, c, 3)),
		Httpchan: uint16(getfromstringslice(dslice, c, 4)),
		Len: uint32(getfromstringslice(dslice, c, 5))}
	if globalErr{
		globalErr=false
		return
	}

	if err:=agent.StructsFileWrite("./tmp/GenData.genes",data,binary.LittleEndian); err!=nil{
		svgError(c,err.Error())
		return
	}
	svgInfo(c, "Файл гена в папке ./tmp")
}

func gendataoutgenerate(c *gin.Context)  {
	dslice:=strings.Split(c.PostForm("Dataf"),",")
	data:=agent.GenDataOut{
		Datatype: agent.DataTypeEnum(getfromstringslice(dslice, c, 0)),
		Fps: uint16(getfromstringslice(dslice, c, 1)),
		Httpchan: uint16(getfromstringslice(dslice, c, 2)),
		Len: uint32(getfromstringslice(dslice, c, 3))}
	if globalErr{
		globalErr=false
		return
	}
	if err:=agent.StructsFileWrite("./tmp/GenDataOut.genes",data,binary.LittleEndian); err!=nil{
		svgError(c,err.Error())
		return
	}
	svgInfo(c, "Файл гена в папке ./tmp")
}

func genfilesgenerate(c* gin.Context)  {
	rslice:=strings.Split(c.PostForm("Receptor"),",")
	nslice:=strings.Split(c.PostForm("Neuron"),",")
	pslice:=strings.Split(c.PostForm("Preffector"),",")
	cslice:=strings.Split(c.PostForm("Chemical"),",")

	if len(greceptors)>0 {
		receptor := agent.Receptor{
			Typer: agent.ReceptorTypeEnum(getfromstringslice(rslice, c, 0)),
			Coren: agent.CoreEnum(getfromstringslice(rslice, c, 1)),
			Typemedi: agent.NeuronTypeEnum(getfromstringslice(rslice, c, 2)),
			Force: uint16(getfromstringslice(rslice, c, 3)),
			Divforce: uint16(getfromstringslice(rslice, c, 4)),
			Threshold: uint64(getfromstringslice(rslice, c, 5)),
			A: uint32(getfromstringslice(rslice, c, 6)),
			Axons: [32]agent.Axon{}}
		for i:=0; i< len(greceptors);i++ {
			greceptors[i].Typer = agent.ReceptorTypeEnum(getfromstringslice(rslice, c, 0))
			greceptors[i].Coren = agent.CoreEnum(getfromstringslice(rslice, c, 1))
			greceptors[i].Typemedi = agent.NeuronTypeEnum(getfromstringslice(rslice, c, 2))
			greceptors[i].MaxX = uint32(ViewX)
			greceptors[i].Recep = receptor
		}
		if globalErr{
			globalErr=false
			return
		}

		if err:=agent.StructsFileWrite("./tmp/receptors.genes",greceptors,binary.LittleEndian); err!=nil{
			svgError(c,err.Error())
			return
		}
	}

	if len(gneurons)>0 {
		chemic := agent.Chemical{
			uint16(getfromstringslice(cslice, c, 0)),
			uint16(getfromstringslice(cslice, c, 1)),
			uint16(getfromstringslice(cslice, c, 2)),
			uint16(getfromstringslice(cslice, c, 3)),
			byte(getfromstringslice(cslice, c, 4)),
			byte(getfromstringslice(cslice, c, 5)),
			byte(getfromstringslice(cslice, c, 6)),
			byte(getfromstringslice(cslice, c, 7)),
			byte(getfromstringslice(cslice, c, 8)),
			byte(getfromstringslice(cslice, c, 9)),
			byte(getfromstringslice(cslice, c, 10)),
			byte(getfromstringslice(cslice, c, 11)),
			byte(getfromstringslice(cslice, c, 12)),
			byte(getfromstringslice(cslice, c, 13)),
			byte(getfromstringslice(cslice, c, 14)),
			byte(getfromstringslice(cslice, c, 15)),
			byte(getfromstringslice(cslice, c, 16)),
			byte(getfromstringslice(cslice, c, 17)),
			byte(getfromstringslice(cslice, c, 18)),
			byte(getfromstringslice(cslice, c, 19)),
		}

		neuron := agent.Neuron{
			Typen: agent.NeuronTypeEnum(getfromstringslice(nslice, c, 0)),

			Coren: agent.CoreEnum(getfromstringslice(nslice, c, 1)),
			Chemic: chemic,

			D: uint16(getfromstringslice(nslice, c, 3)),
			Dendrites: [16]agent.Dendrite{},
			A: uint16(getfromstringslice(nslice, c, 4)),
			Axons: [16]agent.Axon{}	}
		for i:=0; i< len(gneurons);i++  {
			gneurons[i].Typen = agent.NeuronTypeEnum(getfromstringslice(nslice, c, 0))
			gneurons[i].Coren = agent.CoreEnum(getfromstringslice(nslice, c, 1))
			gneurons[i].MaxX = uint32(ViewX)
			gneurons[i].MaxXOtherCore = uint32(getfromstringslice(nslice, c, 2))
			gneurons[i].Neur = neuron
		}
		if globalErr{
			globalErr=false
			return
		}

		if err:=agent.StructsFileWrite("./tmp/neurons.genes",gneurons,binary.LittleEndian); err!=nil{
			svgError(c,err.Error())
			return
		}
	}

	if len(gpreffectors)>0 {
		preffector := agent.Preffector{
			Typep: agent.PreffectorTypeEnum(getfromstringslice(pslice, c, 0)),
			Coren: agent.CoreEnum(getfromstringslice(pslice, c, 1)),
			D:uint32(getfromstringslice(pslice, c, 2)),
			Dendrites: [32]agent.Dendrite{}	}
		for i:=0;i<len(gpreffectors);i++ {
			gpreffectors[i].Typep = agent.PreffectorTypeEnum(getfromstringslice(pslice, c, 0))
			gpreffectors[i].Coren = agent.CoreEnum(getfromstringslice(pslice, c, 1))
			gpreffectors[i].Prefec = preffector
		}
		if globalErr{
			globalErr=false
			return
		}

		if err := agent.StructsFileWrite("./tmp/preffectors.genes", gpreffectors,binary.LittleEndian); err != nil {
			svgError(c,err.Error())
			return
		}
	}

	svgInfo(c, "Файлы в папке ./tmp")

}





