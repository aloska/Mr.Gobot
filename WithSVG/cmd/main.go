package main

import (
	"WithSVG/cmd/agent"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	r         = gin.Default()
	greceptors []agent.GenReceptor
	gneurons []agent.GenNeuron
	gpreffectors []agent.GenPreffector
	globalErr =false
	scale=5
	MaxX=1000
)
func main() {
	getRoutes() //издеся маршруты роутим
	r.Run()
}

func getRoutes() {

	r.LoadHTMLFiles("view/index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", map[string]interface{}{})
	})
	r.StaticFile("styles.css","view/styles.css")
	r.POST("/hello", helloPage)
	r.GET("/hello", helloPage)

	r.POST("/receptor-gen", receptorgen)
	r.POST("/neuron-gen", neurongen)
	r.POST("/preffector-gen", preffectorgen)

	r.POST("/set-settings", setsettings)
	r.POST("/del-entities", delentities)

	r.GET("/hellopng",hellopng)
}







