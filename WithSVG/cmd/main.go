package main

import (
	"WithSVG/cmd/agent"
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

	r.POST("/set-settings", setsettings)
	r.POST("/del-entities", delentities)

	r.GET("/hellopng",hellopng)
	r.GET("/ajax/:command", getajax)
	r.POST("/set-draw-connectors",setdrawconnectors)
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
	RdrawConnector=getfromstringslice(rslice,c,0)
	NdrawConnector[0]= getfromstringslice(nslice,c,0)
	NdrawConnector[1]= getfromstringslice(nslice,c,1)
	PdrawConnector=getfromstringslice(pslice,c,0)
	drawall(c)
}







