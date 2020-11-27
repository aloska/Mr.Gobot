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
	NACHANCLOSE int = 20
	KCHANOPEN int = 35
	KCHANCLOSE int = -60
	CHARGENORM int = -75
	
	NAVALCHANREOPEN byte = 20
	KVALCHANREOPEN byte = 30
	
	NAORG byte =240
	KORG byte =10

)

func NAKATFasa(n *N, c *C){
	
		if n.Na>5 { 
			n.Na-=3
		}
		if n.K<250{
			n.K+=2		
		}
		
		if c.Na<200{
			c.Na+=3
		}
		if c.K>20{
			c.K-=2
		}
	BalanseC(c)
}

func CalcCharge(n *N, c *C) int{
	//return int((float32(n.Na)*1.9 - float32(c.Na)*1.6 + float32(n.K)*1.9 - float32(c.K)*1.6)/3.8)	
	//return int((float32(n.Na)*2 - float32(c.Na)*1.5 + float32(n.K)*2 - float32(c.K)*1.5)/4)	
	return int((float32(n.Na) - float32(NAORG) + float32(n.K) - float32(KORG)-
				float32(c.Na) - float32(c.K))/2.5+15)	
}

func Gradient(n *N, c *C){
	if int(n.Na)>(int(NAORG)+int(c.Na))/2{ 
		n.Na-=1
		if c.Na<200{
			c.Na+=1
		}
	}
	if int(n.Na)<(int(NAORG)+int(c.Na))/2{
		n.Na+=1
		if c.Na>20{
			c.Na-=1
		}
	}
	if int(n.K)>(int(KORG)+int(c.K))/2{ 
		n.K-=1
		if c.K<200{
			c.K+=1
		}		
	}
	if int(n.K)<(int(KORG)+int(c.K))/2{
		n.K+=1
		if c.K>20{
			c.K-=1
		}
	}
	
}
func NaOpened(n *N, c *C){
	if n.Na < 150{
		n.Na+=60
		if c.Na>60{
			c.Na-=30
		}
	}else if n.Na < 210{
		n.Na+=38
		if c.Na>60{
			c.Na-=20
		}
	}else if n.Na < 240{
		n.Na+=15
		if c.Na>60{
			c.Na-=10
		}
	}else if n.Na < 250{
		n.Na+=3		
	}
	BalanseC(c)
}
func KOpened(n *N, c *C){
	if n.K>150{		
		n.K-=100
		if c.K<150{
			c.K+=40
		}
	}else if n.K>75{
		n.K-=65
		if c.K<200{
			c.K+=20
		}
	}else if n.K>20{		
		n.K-=15
		if c.K<200{
			c.K+=10
		}
	}else if n.K>5{		
		n.K-=3
	}
	BalanseC(c)
}
func BalanseC(c *C){
	c.Na=byte((float32(NAORG)+float32(c.Na)*12)/13.)
	c.K=byte((float32(KORG)+float32(c.K)*5)/6.)
}
func Step(n *N, c *C, addCharge int ){

    NAKATFasa(n, c)
	NAKATFasa(n, c)
	NAKATFasa(n, c)
	
	charge:=CalcCharge(n, c)+addCharge
	
	switch n.State{
		case 1://норм
			BalanseC(c)
			if charge<=CHARGENORM{ //глубокая реполяризация, каналы открыты для выравнивания 
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
				/*
				if charge> NACHANOPEN && n.Na<NAVALCHANREOPEN && n.K>KVALCHANREOPEN{
					n.State=1
				}
				*/
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
			/*
			if charge> NACHANOPEN && n.Na<NAVALCHANREOPEN && n.K>KVALCHANREOPEN{
				n.State=10
			}*/
			
			
				
		default:
				n.State=1		
	}
	if rand.Intn(100)>50{
		Gradient(n,c)
	}
	
}

func main() {
	n:=N{250,5, 0}
	n2:=N{250,5, 0}
	c:=C{KORG,NAORG}
	dataCH := []float64{}
	dataCH2 := []float64{}
	dataNa := []float64{}
	dataK := []float64{}
	dataPIC:=[]float64{}
	dataNac := []float64{}
	dataKc := []float64{}
	for i:=0;i<200;i++{
		if (i>56&& i<58) ||(i>60&& i<62) || (i>64&& i<66) || 
		(i>68&& i<70) ||(i>72&& i<74) || (i>76&& i<78) || /*
			(i>80&& i<82) || (i>84&& i<86)  ||
			(i>96&& i<98) || (i>100&& i<102) || (i>104 && i<106) || (i>108 && i<110)||
			(i>112&& i<114) || (i>116&& i<118) || (i>120 && i<122) || */(i>124 && i<126){
			Step(&n,&c,20)
			Step(&n2,&c,20)
			dataPIC=append(dataPIC, 30)
		} else if (i>88 && i<90) || (i>92 && i<94){
			Step(&n,&c,20)
			Step(&n2,&c,0)
			dataPIC=append(dataPIC, -30)
		} else if(i>190&& i<193) {
			Step(&n,&c,20)
			Step(&n2,&c,0)
			dataPIC=append(dataPIC, 20)
		} else{
			Step(&n,&c,0)
			Step(&n2,&c,0)
			dataPIC=append(dataPIC, 0)
		}
		dataCH=append(dataCH, float64(CalcCharge(&n,&c)))
		dataCH2=append(dataCH2, float64(CalcCharge(&n2,&c)))
		dataNa=append(dataNa, float64(n.Na))
		dataK=append(dataK, float64(n.K))
		dataNac=append(dataNac, float64(c.Na))
		dataKc=append(dataKc, float64(c.K))
	}
	graphCH2 := asciigraph.Plot(dataCH2, asciigraph.Height(15), asciigraph.Width(200), asciigraph.Caption("Заряд клетки 2-ой"))
	graphCH := asciigraph.Plot(dataCH, asciigraph.Height(15), asciigraph.Width(200), asciigraph.Caption("Заряд клетки"))
	graphNa := asciigraph.Plot(dataNa, asciigraph.Height(5), asciigraph.Width(200), asciigraph.Caption("Натрий"))
	graphK := asciigraph.Plot(dataK, asciigraph.Height(5), asciigraph.Width(200), asciigraph.Caption("Калий"))
	graphPIC := asciigraph.Plot(dataPIC, asciigraph.Height(5), asciigraph.Width(200), asciigraph.Caption("Раздражение"))
	graphNac := asciigraph.Plot(dataNac, asciigraph.Height(5), asciigraph.Width(200), asciigraph.Caption("Натрий снаружи"))
	graphKc := asciigraph.Plot(dataKc, asciigraph.Height(5), asciigraph.Width(200), asciigraph.Caption("Калий снаружи"))

	fmt.Println("--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------")
    fmt.Println(graphCH2)
	fmt.Println("--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------")
    fmt.Println(graphCH)
	fmt.Println("--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println(graphNa)
	fmt.Println("--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println(graphK)
	fmt.Println("--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println(graphPIC)
	fmt.Println("--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println(graphNac)
	fmt.Println("--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println(graphKc)
	fmt.Println("--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------")
}
