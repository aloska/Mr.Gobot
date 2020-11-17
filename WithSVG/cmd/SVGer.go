package main

import (
	"WithSVG/cmd/agent"
	"github.com/ajstarks/svgo"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func  setsettings(c* gin.Context)  {
	maxxyslice:=strings.Split(c.PostForm("MaxXY"),",")
	if sc, err:=strconv.Atoi(c.PostForm("Scale"));err==nil{
		scale = sc
	}	else{
	  svgError(c)
	  return
	}
	ViewX=getfromstringslice(maxxyslice, c, 0)
	ViewY=getfromstringslice(maxxyslice, c, 0)
	drawall(c)
}

func  delentities(c* gin.Context)  {
	if dr, err:=strconv.Atoi(c.PostForm("DelRecept"));err==nil{
		if dr>=0 && dr<len(greceptors) {
			greceptors = append(greceptors[:dr], greceptors[dr+1:]...)
		}
	}
	if dn, err:=strconv.Atoi(c.PostForm("DelNeuron"));err==nil{
		if dn>=0 && dn<len(gneurons) {
			gneurons = append(gneurons[:dn], gneurons[dn+1:]...)
		}
	}
	if dp, err:=strconv.Atoi(c.PostForm("DelPreffector"));err==nil{
		if dp>=0 && dp<len(gpreffectors) {
			gpreffectors = append(gpreffectors[:dp], gpreffectors[dp+1:]...)
		}
	}

	drawall(c)
}

func helloPage (c* gin.Context){
	s := svg.New(c.Writer)
	s.StartviewUnit (100,100,"%",0, 0,ViewX*scale,ViewY*scale)
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

	if globalErr{
		globalErr=false
		return
	}

	greceptors=append(greceptors, agent.GenReceptor{})
	
	greceptors[len(greceptors)-1].Ndata=uint32(Ndata)
	greceptors[len(greceptors)-1].NdataW=byte(NdataW)
	greceptors[len(greceptors)-1].NdataWB=byte(NdataWB)
	greceptors[len(greceptors)-1].Maxi=uint32(Maxi)
	greceptors[len(greceptors)-1].Maxj=byte(Maxj)
	greceptors[len(greceptors)-1].Maxk=byte(Maxk)
	greceptors[len(greceptors)-1].Iteri=uint32(Iteri)
	greceptors[len(greceptors)-1].Iterj=byte(Iterj)
	greceptors[len(greceptors)-1].Iterk=byte(Iterk)

	greceptors[len(greceptors)-1].MaxX=uint32(ViewX)
	greceptors[len(greceptors)-1].Ax1stX=uint32(Ax1stX)
	greceptors[len(greceptors)-1].Ax1stY=uint32(Ax1stY)
	greceptors[len(greceptors)-1].AxShiftX=uint16(AxShiftX)
	greceptors[len(greceptors)-1].AxShiftY=uint16(AxShiftY)
	greceptors[len(greceptors)-1].AxNextShiftX=uint16(AxNextShiftX)
	greceptors[len(greceptors)-1].AxNextShiftY=uint16(AxNextShiftY)

	drawall(c)

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

	if globalErr{
		globalErr=false
		return
	}

	gneurons=append(gneurons, agent.GenNeuron{})
	
	gneurons[len(gneurons)-1].MaxX=uint32(ViewX)

	gneurons[len(gneurons)-1].Layers=uint16(Layers)
	gneurons[len(gneurons)-1].Laywidth=uint32(Laywidth)

	gneurons[len(gneurons)-1].Soma1stX=uint32(Soma1stX)
	gneurons[len(gneurons)-1].Soma1stY=uint32(Soma1stY)
	gneurons[len(gneurons)-1].SomaNextShiftX=uint16(SomaNextShiftX)
	gneurons[len(gneurons)-1].SomaNextShiftY=uint16(SomaNextShiftY)
	gneurons[len(gneurons)-1].SomaLayerShiftX=uint16(SomaLayerShiftX)
	gneurons[len(gneurons)-1].SomaLayerShiftY=uint16(SomaLayerShiftY)


	gneurons[len(gneurons)-1].Dend1stX=uint32(Dend1stX)
	gneurons[len(gneurons)-1].Dend1stY=uint32(Dend1stY)
	gneurons[len(gneurons)-1].DendShiftX=uint16(DendShiftX)
	gneurons[len(gneurons)-1].DendShiftY=uint16(DendShiftY)
	gneurons[len(gneurons)-1].DendNextShiftX=uint16(DendNextShiftX)
	gneurons[len(gneurons)-1].DendNextShiftY=uint16(DendNextShiftY)
	gneurons[len(gneurons)-1].DendLayerShiftX=uint16(DendLayerShiftX)
	gneurons[len(gneurons)-1].DendLayerShiftY=uint16(DendLayerShiftY)

	gneurons[len(gneurons)-1].Ax1stX=uint32(Ax1stX)
	gneurons[len(gneurons)-1].Ax1stY=uint32(Ax1stY)
	gneurons[len(gneurons)-1].AxShiftX=uint16(AxShiftX)
	gneurons[len(gneurons)-1].AxShiftY=uint16(AxShiftY)
	gneurons[len(gneurons)-1].AxNextShiftX=uint16(AxNextShiftX)
	gneurons[len(gneurons)-1].AxNextShiftY=uint16(AxNextShiftY)
	gneurons[len(gneurons)-1].AxLayerShiftX=uint16(AxLayerShiftX)
	gneurons[len(gneurons)-1].AxLayerShiftY=uint16(AxLayerShiftY)

	drawall(c)

}

func preffectorgen (c* gin.Context) {
	Ndataslice:=strings.Split(c.PostForm("Ndata"),",")
	Maxslice:=strings.Split(c.PostForm("Max"),",")
	Iterslice:=strings.Split(c.PostForm("Iter"),",")
	Dend1stslice:=strings.Split(c.PostForm("Dend1st"),",")
	DendShiftslice:=strings.Split(c.PostForm("DendShift"),",")
	DendNextShiftslice:=strings.Split(c.PostForm("DendNextShift"),",")


	Ndata:= getfromstringslice(Ndataslice, c, 0)
	NdataW:= getfromstringslice(Ndataslice, c, 1)
	NdataWB:= getfromstringslice(Ndataslice, c, 2)

	Maxi:= getfromstringslice(Maxslice, c, 0)
	Maxj:= getfromstringslice(Maxslice, c, 1)
	Maxk:= getfromstringslice(Maxslice, c, 2)

	Iteri:= getfromstringslice(Iterslice, c, 0)
	Iterj:= getfromstringslice(Iterslice, c, 1)
	Iterk:= getfromstringslice(Iterslice, c, 2)

	Dend1stX:= getfromstringslice(Dend1stslice, c, 0)
	Dend1stY:= getfromstringslice(Dend1stslice, c, 1)

	DendShiftX:= getfromstringslice(DendShiftslice, c, 0)
	DendShiftY:= getfromstringslice(DendShiftslice, c, 1)

	DendNextShiftX:= getfromstringslice(DendNextShiftslice, c, 0)
	DendNextShiftY:= getfromstringslice(DendNextShiftslice, c, 1)

	if globalErr{
		globalErr=false
		return
	}

	gpreffectors=append(gpreffectors, agent.GenPreffector{})

	gpreffectors[len(gpreffectors)-1].Ndata=uint32(Ndata)
	gpreffectors[len(gpreffectors)-1].NdataW=byte(NdataW)
	gpreffectors[len(gpreffectors)-1].NdataWB=byte(NdataWB)
	gpreffectors[len(gpreffectors)-1].Maxi=uint32(Maxi)
	gpreffectors[len(gpreffectors)-1].Maxj=byte(Maxj)
	gpreffectors[len(gpreffectors)-1].Maxk=byte(Maxk)
	gpreffectors[len(gpreffectors)-1].Iteri=uint32(Iteri)
	gpreffectors[len(gpreffectors)-1].Iterj=byte(Iterj)
	gpreffectors[len(gpreffectors)-1].Iterk=byte(Iterk)

	gpreffectors[len(gpreffectors)-1].MaxX=uint32(ViewX)
	gpreffectors[len(gpreffectors)-1].Dend1stX=uint32(Dend1stX)
	gpreffectors[len(gpreffectors)-1].Dend1stY=uint32(Dend1stY)
	gpreffectors[len(gpreffectors)-1].DendShiftX=uint16(DendShiftX)
	gpreffectors[len(gpreffectors)-1].DendShiftY=uint16(DendShiftY)
	gpreffectors[len(gpreffectors)-1].DendNextShiftX=uint16(DendNextShiftX)
	gpreffectors[len(gpreffectors)-1].DendNextShiftY=uint16(DendNextShiftY)

	drawall(c)

}

func drawall(c*gin.Context){
	s := svg.New(c.Writer)
	s.StartviewUnit(100,100, "%", 0,0,ViewX*scale,ViewY*scale)

	var sdrawX, sdrawY int
	//рецепторы
	for _, rrr:=range greceptors{
		cur:=0
		for i:=rrr.Ndata; i<= rrr.Maxi; i=i+rrr.Iteri{
			for j:=rrr.NdataW; j<= rrr.Maxj; j=j+rrr.Iterj{
				for k:=rrr.NdataWB; k<=rrr.Maxk; k=k+rrr.Iterk {
					if cur==RdrawConnector{
						sdrawX=int(i)
						sdrawY=int(j)
					}
					drawReceptIJK(s,int(i),int(j),int(k))
					cur++
				}
			}
		}

		x:=rrr.Ax1stX
		y:= rrr.Ax1stY

		for c:=0; c<cur; c++ {
			xs:=x
			ys:=y
			for i:=0;i<32;i++{
				drawDendrAxStyle(s,int(x),int(y),1,"fill:none;stroke:green;stroke-width:0.2")
				if c==RdrawConnector{
					drawConnectReceptor(s, sdrawX,sdrawY, int(x),int(y),"stroke-width:0.2;stroke:yellow")

				}
				x=x+uint32(rrr.AxShiftX)
				y=y+uint32(rrr.AxShiftY)

			}
			x=xs+uint32(rrr.AxNextShiftX)
			y=ys+uint32(rrr.AxNextShiftY)
		}
	}  

	//нейроны
	for _, nnn:=range gneurons{
		//генерируем сомы
		x:=nnn.Soma1stX
		y:=nnn.Soma1stY
		for lay:=0; lay<int(nnn.Layers);lay++{
			xs:=x
			ys:=y
			for i:=0;i<int(nnn.Laywidth);i++{
				if lay==NdrawConnector[0] && i==NdrawConnector[1]{
					sdrawX=int(x)
					sdrawY=int(y)
				}
				drawSomaStyle(s,int(x),int(y),1,"fill:none;stroke:red;stroke-width:0.2")
				x=x+uint32(nnn.SomaNextShiftX)
				y=y+uint32(nnn.SomaNextShiftY)
			}
			x=xs+uint32(nnn.SomaLayerShiftX)
			y=ys+uint32(nnn.SomaLayerShiftY)
		}

		//генерируем дендриты
		x=nnn.Dend1stX
		y=nnn.Dend1stY
		for lay:=0; lay<int(nnn.Layers);lay++{
			xl:=x
			yl:=y
			for i:=0;i<int(nnn.Laywidth);i++{
				xs:=x
				ys:=y
				for d:=0;d<16;d++ {
					drawDendrAxStyle(s, int(x), int(y), 1, "fill:none;stroke:blue;stroke-width:0.2")
					if lay == NdrawConnector[0] && i == NdrawConnector[1] {
						drawConnectSoma(s,sdrawX, sdrawY, int(x), int(y),
							"stroke-width:0.2;stroke:blue;stroke-dasharray:2")
					}
					x = x + uint32(nnn.DendShiftX)
					y = y + uint32(nnn.DendShiftY)
				}
				x = xs + uint32(nnn.DendNextShiftX)
				y = ys + uint32(nnn.DendNextShiftY)
			}
			x=xl+uint32(nnn.DendLayerShiftX)
			y=yl+uint32(nnn.DendLayerShiftY)
		}

		//генерируем аксоны
		x=nnn.Ax1stX
		y=nnn.Ax1stY
		for lay:=0; lay<int(nnn.Layers);lay++{
			xl:=x
			yl:=y
			for i:=0;i<int(nnn.Laywidth);i++{
				xs:=x
				ys:=y
				for d:=0;d<16;d++ {
					drawDendrAxStyle(s, int(x), int(y), 1, "fill:none;stroke:green;stroke-width:0.2")
					if lay == NdrawConnector[0] && i == NdrawConnector[1] {
						drawConnectSoma(s,sdrawX, sdrawY, int(x), int(y),
							"stroke-width:0.2;stroke:green;stroke-dasharray:2")

					}
					x = x + uint32(nnn.AxShiftX)
					y = y + uint32(nnn.AxShiftY)
				}
				x = xs + uint32(nnn.AxNextShiftX)
				y = ys + uint32(nnn.AxNextShiftY)
			}
			x=xl+uint32(nnn.AxLayerShiftX)
			y=yl+uint32(nnn.AxLayerShiftY)
		}
	}

	//преффекторы
	for _, ppp:=range gpreffectors{
		cur:=0
		for i:=ppp.Ndata; i<= ppp.Maxi; i=i+ppp.Iteri{
			for j:=ppp.NdataW; j<= ppp.Maxj; j=j+ppp.Iterj{
				for k:=ppp.NdataWB; k<=ppp.Maxk; k=k+ppp.Iterk {
					if cur==PdrawConnector{
						sdrawX=int(i)
						sdrawY=int(j)
					}
					drawPreffectIJK(s,int(i),int(j),int(k))
					cur++
				}
			}
		}

		x:=ppp.Dend1stX
		y:= ppp.Dend1stY

		for c:=0; c<cur; c++ {
			xs:=x
			ys:=y
			for i:=0;i<32;i++{
				drawDendrAxStyle(s,int(x),int(y),1,"stroke:blue;stroke-width:0.2")
				if c==PdrawConnector{
					drawConnectPreffector(s, sdrawX,sdrawY, int(x),int(y),"stroke-width:0.2;stroke:red")

				}
				x=x+uint32(ppp.DendShiftX)
				y=y+uint32(ppp.DendShiftY)

			}
			x=xs+uint32(ppp.DendNextShiftX)
			y=ys+uint32(ppp.DendNextShiftY)
		}
	}

	s.End()
}


func svgError(c* gin.Context){
	globalErr=true
	s := svg.New(c.Writer)
	s.StartviewUnit (100,100,"%",0,0,ViewX*scale, ViewY*scale)
	s.Circle(250, 250, 250, "fill:orange;stroke:red;stroke-width:4")
	s.Gstyle("fill:red;font-size:20pt;text-anchor:middle;font-family:monospace")
	s.Text(250,250, "Плохо форму заполнил!")
	s.Gend()
	s.End()
}

func drawReceptIJK(s *svg.SVG, i int, j int, k int){
	s.Circle(i*scale,j*scale,k*scale/2+scale, "fill:yellow;stroke:blue;stroke-width:1")
}
func drawPreffectIJK(s *svg.SVG, i int, j int, k int){
	s.Circle(ViewX*scale-i*scale,j*scale,k*scale/2+scale, "fill:red;stroke:blue;stroke-width:1")
}
func drawDendrAxStyle(s *svg.SVG, x int, y int, R int, style string){
	s.Circle(x*scale,y*scale,R, style)
}
func drawSomaStyle(s *svg.SVG, x int, y int, size int, style string){
	s.Rect(x*scale-size/2,y*scale-size/2,size,size,style)
}

func drawConnectReceptor(s *svg.SVG, xReceptor int, yReceptor int, xAxon int, yAxon int, style string){
	s.Line(xReceptor*scale,yReceptor*scale, xAxon*scale,yAxon*scale,style)
}

func drawConnectPreffector(s *svg.SVG, xPreffector int, yPreffector int, xDendrite int, yDendrite int, style string){
	s.Line(ViewX*scale-xPreffector*scale,yPreffector*scale, xDendrite*scale,yDendrite*scale,style)
}

func drawConnectSoma(s *svg.SVG, x1 int, y1 int, x2 int, y2 int, style string){
	s.Line(x1*scale,y1*scale, x2*scale,y2*scale,style)
}

