package main

import (
	"github.com/ajstarks/svgo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

var (
	r         = gin.Default()
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
	r.POST("/hello", helloPage)
	r.GET("/hello", helloPage)

	r.POST("/receptor-gen", receptorgen)
	r.POST("/neuron-gen", neurongen)
	r.POST("/preffector-gen", preffectorgen)

}

func helloPage (c* gin.Context){
	s := svg.New(c.Writer)
	s.Start(1000, 10000)
	s.Circle(250, 250, 250, "fill:none;stroke:red;stroke-width:4")
	s.Gstyle("fill:black;font-size:16pt;text-anchor:middle;font-family:monospace")
		s.Text(250,250, "Ну, здравствуй!")
	s.Gend()
	s.End()
}

func getfromstringslice(s []string, c *gin.Context, i int) int{
	if i>=len(s){
		svgError(c)
		return 0
	}
	ret, err:= strconv.Atoi(s[i])
	if err!=nil{
		svgError(c)
		return 0
	}
	return ret
}

func receptorgen (c* gin.Context){
	Ndataslice:=strings.Split(c.PostForm("Ndata"),",")
	Maxslice:=strings.Split(c.PostForm("Max"),",")
	Iterslice:=strings.Split(c.PostForm("Iter"),",")
	Ax1stslice:=strings.Split(c.PostForm("Ax1st"),",")
	AxShiftslice:=strings.Split(c.PostForm("AxShift"),",")
	AxNextShiftslice:=strings.Split(c.PostForm("AxNextShift"),",")

	
	Ndata:= getfromstringslice(Ndataslice, c, 0)
	NdataW:= getfromstringslice(Ndataslice, c, 1)
	NdataWB:= getfromstringslice(Ndataslice, c, 2)

	Maxi:= getfromstringslice(Maxslice, c, 0)
	Maxj:= getfromstringslice(Maxslice, c, 1)
	Maxk:= getfromstringslice(Maxslice, c, 2)

	Iteri:= getfromstringslice(Iterslice, c, 0)
	Iterj:= getfromstringslice(Iterslice, c, 1)
	Iterk:= getfromstringslice(Iterslice, c, 2)

	Ax1stX:= getfromstringslice(Ax1stslice, c, 0)
	Ax1stY:= getfromstringslice(Ax1stslice, c, 1)

	AxShiftX:= getfromstringslice(AxShiftslice, c, 0)
	AxShiftY:= getfromstringslice(AxShiftslice, c, 1)

	AxNextShiftX:= getfromstringslice(AxNextShiftslice, c, 0)
	AxNextShiftY:= getfromstringslice(AxNextShiftslice, c, 1)

	s := svg.New(c.Writer)
	s.Start(1000, 10000)

	cur:=0
	for i:=Ndata; i<= Maxi; i=i+Iteri{
		for j:=NdataW; j<= Maxj; j=j+Iterj{
			for k:=NdataWB; k<=Maxk; k=k+Iterk {

				drawIJK(s,i*10,j*10,k*2+5)
				cur++
			}
		}
	}

	x:=Ax1stX
	y:= Ax1stY

	for c:=0; c<cur; c++ {
		xs:=x
		ys:=y
		for i:=0;i<32;i++{
			drawXYRStyle(s,x*5,y*5,1,"stroke:green;stroke-width:0.1")
			if c==0{
				s.Line(Ndata*10,NdataW*10, x*5,y*5,"stroke-width:0.3;stroke:grey")
			}
			x=x+AxShiftX
			y=y+AxShiftY

		}
		x=xs+AxNextShiftX
		y=ys+AxNextShiftY
	}

	s.End()
}

func neurongen (c* gin.Context) {
	Layslice := strings.Split(c.PostForm("Lay"), ",")

	Soma1stslice := strings.Split(c.PostForm("Soma1st"), ",")
	SomaNextShiftslice := strings.Split(c.PostForm("SomaNextShift"), ",")
	SomaLayerShiftslice := strings.Split(c.PostForm("SomaLayerShift"), ",")

	Dend1stslice := strings.Split(c.PostForm("Dend1st"), ",")
	DendShiftslice := strings.Split(c.PostForm("DendShift"), ",")
	DendNextShiftslice := strings.Split(c.PostForm("DendNextShift"), ",")
	DendLayerShifttslice := strings.Split(c.PostForm("DendLayerShift"), ",")

	Ax1stslice := strings.Split(c.PostForm("Ax1st"), ",")
	AxShiftslice := strings.Split(c.PostForm("AxShift"), ",")
	AxNextShiftslice := strings.Split(c.PostForm("AxNextShift"), ",")
	AxLayerShifttslice := strings.Split(c.PostForm("AxLayerShift"), ",")

	Layers := getfromstringslice(Layslice, c, 0)
	Laywidth := getfromstringslice(Layslice, c, 1)

	Soma1stX := getfromstringslice(Soma1stslice, c, 0)
	Soma1stY := getfromstringslice(Soma1stslice, c, 1)
	SomaNextShiftX := getfromstringslice(SomaNextShiftslice, c, 0)
	SomaNextShiftY := getfromstringslice(SomaNextShiftslice, c, 1)
	SomaLayerShiftX := getfromstringslice(SomaLayerShiftslice, c, 0)
	SomaLayerShiftY := getfromstringslice(SomaLayerShiftslice, c, 1)

	Dend1stX := getfromstringslice(Dend1stslice, c, 0)
	Dend1stY := getfromstringslice(Dend1stslice, c, 1)
	DendShiftX := getfromstringslice(DendShiftslice, c, 0)
	DendShiftY := getfromstringslice(DendShiftslice, c, 1)
	DendNextShiftX := getfromstringslice(DendNextShiftslice, c, 0)
	DendNextShiftY := getfromstringslice(DendNextShiftslice, c, 1)
	DendLayerShiftX := getfromstringslice(DendLayerShifttslice, c, 0)
	DendLayerShiftY := getfromstringslice(DendLayerShifttslice, c, 1)

	Ax1stX := getfromstringslice(Ax1stslice, c, 0)
	Ax1stY := getfromstringslice(Ax1stslice, c, 1)
	AxShiftX := getfromstringslice(AxShiftslice, c, 0)
	AxShiftY := getfromstringslice(AxShiftslice, c, 1)
	AxNextShiftX := getfromstringslice(AxNextShiftslice, c, 0)
	AxNextShiftY := getfromstringslice(AxNextShiftslice, c, 1)
	AxLayerShiftX := getfromstringslice(AxLayerShifttslice, c, 0)
	AxLayerShiftY := getfromstringslice(AxLayerShifttslice, c, 1)

	s := svg.New(c.Writer)
	s.Start(1000, 10000)

	//генерируем сомы
	x:=Soma1stX
	y:=Soma1stY
	for lay:=0; lay<Layers;lay++{
		xs:=x
		ys:=y
		for i:=0;i<Laywidth;i++{
			drawSomaStyle(s,x*5,y*5,5,"stroke:red;stroke-width:0.1")
			x=x+SomaNextShiftX
			y=y+SomaNextShiftY
		}
		x=xs+SomaLayerShiftX
		y=ys+SomaLayerShiftY
	}

	//генерируем дендриты
	x=Dend1stX
	y=Dend1stY
	for lay:=0; lay<Layers;lay++{
		xl:=x
		yl:=y
		for i:=0;i<Laywidth;i++{
			xs:=x
			ys:=y
			for d:=0;d<16;d++ {
				drawXYRStyle(s, x*5, y*5, 1, "stroke:orange;stroke-width:0.1")
				if lay == 0 && i == 0 {
					s.Line(Soma1stX*5, Soma1stY*5, x*5, y*5, "stroke-width:0.3;stroke:grey")
				}
				x = x + DendShiftX
				y = y + DendShiftY
			}
			x = xs + DendNextShiftX
			y = ys + DendNextShiftY
		}
		x=xl+DendLayerShiftX
		y=yl+DendLayerShiftY
	}

	//генерируем аксоны
	x=Ax1stX
	y=Ax1stY
	for lay:=0; lay<Layers;lay++{
		xl:=x
		yl:=y
		for i:=0;i<Laywidth;i++{
			xs:=x
			ys:=y
			for d:=0;d<16;d++ {
				drawXYRStyle(s, x*5, y*5, 1, "fill:grey;stroke:orange;stroke-width:0.1")
				if lay == 0 && i == 0 {
					s.Line(Soma1stX*5, Soma1stY*5, x*5, y*5, "stroke-width:0.3;stroke:grey")
				}
				x = x + AxShiftX
				y = y + AxShiftY
			}
			x = xs + AxNextShiftX
			y = ys + AxNextShiftY
		}
		x=xl+AxLayerShiftX
		y=yl+AxLayerShiftY
	}

	s.End()
}

func preffectorgen (c* gin.Context) {

}

func svgError(c* gin.Context){
	s := svg.New(c.Writer)
	s.Start(1000, 10000)
	s.Circle(250, 250, 250, "fill:orange;stroke:red;stroke-width:4")
	s.Gstyle("fill:red;font-size:20pt;text-anchor:middle;font-family:monospace")
	s.Text(250,250, "Плохо форму заполнил!")
	s.Gend()
	s.End()
}

func drawIJK(s *svg.SVG, i int, j int, k int){
	s.Circle(i,j,k, "fill:yellow;stroke:blue;stroke-width:1")
}
func drawXYRStyle(s *svg.SVG, x int, y int, R int, style string){
	s.Circle(x,y,R, style)
}
func drawSomaStyle(s *svg.SVG, x int, y int, size int, style string){
	s.Rect(x-size/2,y-size/2,size,size,style)
}