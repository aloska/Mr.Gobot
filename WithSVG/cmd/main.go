package main

import (
	"WithSVG/cmd/agent"
	"github.com/ajstarks/svgo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

var (
	r         = gin.Default()
	greceptors []agent.GenReceptor
	gneurons []agent.GenNeuron
	gpreffectors []agent.GenPreffector
	globalErr =false
	scale=5
	ViewX=1000
	ViewY=1000
	RdrawConnector = 0
	NdrawConnector = [2]int{0,0}
	PdrawConnector=0
	gridshow=false
	gridX=0
	gridY=0
	gridW=10
	gridH=10
	gridN=2

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
	r.POST("genfiles-generate",genfilesgenerate)

	r.POST("/set-settings", setsettings)
	r.POST("/del-entities", delentities)

	r.GET("/hellopng",hellopng)
	r.GET("/ajax/:command", getajax)

	r.POST("/set-draw-connectors",setdrawconnectors)
	r.POST("/set-draw-grid",setdrawgrid)
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
	drawall(c)
}

func setdrawgrid(c *gin.Context){
	gslice:=strings.Split(c.PostForm("DrawGrid"),",")
	gridX=int(getfromstringslice(gslice,c,0))
	gridY=int(getfromstringslice(gslice,c,1))
	gridW=int(getfromstringslice(gslice,c,2))
	gridH=int(getfromstringslice(gslice,c,3))
	gridN=int(getfromstringslice(gslice,c,4))
	gridshow=!gridshow
	drawall(c)
}

func genfilesgenerate(c* gin.Context)  {
	rslice:=strings.Split(c.PostForm("Receptor"),",")
	nslice:=strings.Split(c.PostForm("Neuron"),",")
	pslice:=strings.Split(c.PostForm("Preffector"),",")
	cslice:=strings.Split(c.PostForm("Chemical"),",")

	if len(greceptors)>0 {
		receptor := agent.Receptor{
			agent.ReceptorTypeEnum(getfromstringslice(rslice, c, 0)),
			agent.CoreEnum(getfromstringslice(rslice, c, 1)),
			0,
			0,
			0,
			0,
			agent.NeuronTypeEnum(getfromstringslice(rslice, c, 2)),
			uint16(getfromstringslice(rslice, c, 3)),
			uint16(getfromstringslice(rslice, c, 4)),
			uint64(getfromstringslice(rslice, c, 5)),
			uint32(getfromstringslice(rslice, c, 6)),
			[32]agent.Axon{},
		}
		for i:=0; i< len(greceptors);i++ {
			greceptors[i].Typer = agent.ReceptorTypeEnum(getfromstringslice(rslice, c, 0))
			greceptors[i].Coren = agent.CoreEnum(getfromstringslice(rslice, c, 1))
			greceptors[i].Typemedi = agent.NeuronTypeEnum(getfromstringslice(rslice, c, 2))
			greceptors[i].MaxX = uint32(ViewX)
			greceptors[i].Recep = receptor
		}
		if err:=agent.StructsFileWrite("./tmp/receptors.genes",greceptors); err!=nil{
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
			agent.NeuronTypeEnum(getfromstringslice(nslice, c, 0)),
			0,
			agent.CoreEnum(getfromstringslice(nslice, c, 1)),
			chemic,
			0,
			uint16(getfromstringslice(nslice, c, 3)),
			[16]agent.Dendrite{},
			uint16(getfromstringslice(nslice, c, 4)),
			[16]agent.Axon{},
		}
		for i:=0; i< len(gneurons);i++  {
			gneurons[i].Typen = agent.NeuronTypeEnum(getfromstringslice(nslice, c, 0))
			gneurons[i].Coren = agent.CoreEnum(getfromstringslice(nslice, c, 1))
			gneurons[i].MaxX = uint32(ViewX)
			gneurons[i].MaxXOtherCore = uint32(getfromstringslice(nslice, c, 2))
			gneurons[i].Neur = neuron
		}
		if err:=agent.StructsFileWrite("./tmp/neurons.genes",gneurons); err!=nil{
			svgError(c,err.Error())
			return
		}
	}

	if len(gpreffectors)>0 {
		preffector := agent.Preffector{
			agent.PreffectorTypeEnum(getfromstringslice(pslice, c, 0)),
			agent.CoreEnum(getfromstringslice(pslice, c, 1)),
			0,
			0,
			0,
			0,
			uint64(getfromstringslice(pslice, c, 2)),
			[64]agent.Dendrite{},
		}
		for i:=0;i<len(gpreffectors);i++ {
			gpreffectors[i].Typep = agent.PreffectorTypeEnum(getfromstringslice(pslice, c, 0))
			gpreffectors[i].Coren = agent.CoreEnum(getfromstringslice(pslice, c, 1))
			gpreffectors[i].Prefec = preffector
		}

		if err := agent.StructsFileWrite("./tmp/preffectors.genes", gpreffectors); err != nil {
			svgError(c,err.Error())
			return
		}
	}

	s := svg.New(c.Writer)
	s.StartviewUnit (100,100,"%",0,0,ViewX, ViewY)
	s.Circle(ViewX/2, ViewY/2, ViewY/3, "fill:yellow;stroke:green;stroke-width:4")
	s.Gstyle("fill:green;font-size:"+strconv.Itoa(ViewX/16*scale)+"pt;text-anchor:middle;font-family:monospace")
	s.Text(ViewX/2,ViewY/2, "Файлы готовы! См. папку /tmp")
	s.Gend()
	s.End()

}





