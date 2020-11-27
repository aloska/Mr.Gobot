package main

import (
	"fmt"
	"github.com/guptarohit/asciigraph"
	"math/rand"
)

type C struct{
	K byte
	Na byte
}
type N struct{
	K byte
	Na byte
	State byte

}

const (
	NACHANOPEN int = -65
	NACHANCLOSE int = 10
	KCHANOPEN int = 5
	KCHANCLOSE int = -80
	CHARGENORM int = -75
	
	NAVALCHANREOPEN byte = 100
	KVALCHANREOPEN byte = 40

)

func NAKATFasa(n *N, c *C){
	if n.Na>5 && n.K<250 && c.Na<250 && c.K>5{
		n.Na-=3
		c.Na+=3
		n.K+=2
		c.K-=2
	}
}

func CalcCharge(n *N, c *C) int{
	//return int((float32(n.Na)*1.9 - float32(c.Na)*1.6 + float32(n.K)*1.9 - float32(c.K)*1.6)/3.8)	
	return int((float32(n.Na)*2 - float32(c.Na)*1.5 + float32(n.K)*2 - float32(c.K)*1.5)/4)	
}

func Gradient(n *N, c *C){
	if n.Na>c.Na && n.Na>5{ 
		n.Na-=1
		c.Na+=1
	}
	if n.Na<c.Na && c.Na>5{
		n.Na+=1
		c.Na-=1
	}
	if n.K>c.K && n.K>5{ 
		n.K-=1
		c.K+=1
	}
	if n.K<c.K && c.K>5{
		n.K+=1
		c.K-=1
	}
}
func NaOpened(n *N, c *C){
	if n.Na < 190 && c.Na>65{
		n.Na+=60
		c.Na-=60
	}else if n.Na < 210 && c.Na>40{
		n.Na+=38
		c.Na-=38
	}else if n.Na < 240 && c.Na>15{
		n.Na+=15
		c.Na-=15
	}else if n.Na < 250 && c.Na>5{
		n.Na+=3
		c.Na-=3
	}

}
func KOpened(n *N, c *C){
	if c.K < 200 && n.K>55{
		c.K+=50
		n.K-=50
	}else if c.K < 240 && n.K>20{
		c.K+=15
		n.K-=15
	}else if c.K < 250 && n.K>5{
		c.K+=3
		n.K-=3
	}

}
func Step(n *N, c *C, addCharge int ){
	NAKATFasa(n, c)
	
	charge:=CalcCharge(n, c)+addCharge
	
	switch n.State{
		case 1://норм
			if charge<=CHARGENORM{ //глубокая реполяризация, каналы открыты для выравнивания 
				Gradient(n,c)
				Gradient(n,c)
				Gradient(n,c)
				//Gradient(n,c)
				//Gradient(n,c)
				charge=CalcCharge(n, c)+addCharge
			} else if charge>NACHANOPEN{
				n.State=10
			}
		case 10://начало деполяризации
			if charge>=NACHANOPEN && charge<=NACHANCLOSE{
				NaOpened(n,c)
				charge=CalcCharge(n, c)+addCharge
			}
			if charge>=KCHANOPEN {
				KOpened(n, c)
				charge=CalcCharge(n, c)+addCharge
				if charge> NACHANOPEN && n.Na<NAVALCHANREOPEN && n.K>KVALCHANREOPEN{
					n.State=1
				}
			}
			if charge>=NACHANCLOSE{
				n.State=20
			}
		case 20://только калиевый ток
			KOpened(n, c)
			charge:=CalcCharge(n, c)+addCharge
			if charge<KCHANCLOSE {
				n.State=1
			}
			if charge> NACHANOPEN && n.Na<NAVALCHANREOPEN && n.K>KVALCHANREOPEN{
				n.State=1
			}
			
				
		default:
				n.State=1		
	}
	if rand.Intn(100)>10{
		Gradient(n,c)
	}
	NAKATFasa(n, c)
	NAKATFasa(n, c)
	NAKATFasa(n, c)
}

func main() {
	n:=N{150,100, 0}
	c:=C{100,150}
	dataCH := []float64{}
	dataNa := []float64{}
	dataK := []float64{}
	dataPIC:=[]float64{}
	for i:=0;i<200;i++{
		if (i>80&& i<85) || (i>90&& i<95) || (i>100 && i<105){
			Step(&n,&c,30)
			dataPIC=append(dataPIC, 30)
		} else if(i>150&& i<155) {
			Step(&n,&c,30)
			dataPIC=append(dataPIC, 30)
		} else{
			Step(&n,&c,0)
			dataPIC=append(dataPIC, 0)
		}
		dataCH=append(dataCH, float64(CalcCharge(&n,&c)))
		dataNa=append(dataNa, float64(n.Na))
		dataK=append(dataK, float64(n.K))
	}
	graphCH := asciigraph.Plot(dataCH, asciigraph.Height(15), asciigraph.Caption("Заряд клетки"))
	graphNa := asciigraph.Plot(dataNa, asciigraph.Height(5), asciigraph.Caption("Натрий"))
	graphK := asciigraph.Plot(dataK, asciigraph.Height(5), asciigraph.Caption("Калий"))
	graphPIC := asciigraph.Plot(dataPIC, asciigraph.Height(5), asciigraph.Caption("Раздражение"))

	fmt.Println("--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------")
    fmt.Println(graphCH)
	fmt.Println("--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println(graphNa)
	fmt.Println("--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println(graphK)
	fmt.Println("--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println(graphPIC)
	fmt.Println("--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------")
}
