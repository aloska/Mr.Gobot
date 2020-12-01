package main

import (
	"WithSVG/cmd/agent"
	"strconv"

	"net/http"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

var (
	a agent.Agent
	sm map[agent.SynEnum] *agent.Synapses
	c agent.Chemical
	g agent.GenNeuron
	s agent.Synapses
	o agent.Organism
	n agent.Neuron
	che []agent.Chemical
)

func main() {
	sm=make(map[agent.SynEnum] *agent.Synapses)

	a=agent.Agent{}
	c=agent.Chemical{ACh: 0, Na: 3, K:250, GLUC: 0x100, CHOL: 0,O2: 0x100, OMEGA: 0x100	}
	g=agent.GenNeuron{}
	s=agent.Synapses{}
	o=agent.Organism{}
	n=agent.Neuron{
		Typen: agent.NEURONACETILHOLIN,
		State:0,
		SynNumber: 0,
		Gen:0,
		SynNumberAxons: 0,
		Chemic: c,
		N:0,
		D:0xffff,
		A:0xffff}

	for i:=0;i<16;i++{
		n.Dendrites[i].N=uint32(i+1)
		n.Axons[i].N=uint32(i+15) //сами себе аксонами на вход 2 дендритов, для проверки
		n.Dendrites[i].Typed=agent.DENDRACHION
		n.Axons[i].Vesiculs=5
	}
	che=make([]agent.Chemical,40)
	for i:=0;i<40;i++{
		che[i]=c
	}
	s.SetDebugSyn(che)
	sm[0]=&s
	o.SetDebugSynMap(sm)
	a.SetDebugOrganism(&o)


	http.HandleFunc("/", basa)
	http.ListenAndServe(":8082", nil)

}

func basa(w http.ResponseWriter, _ *http.Request) {

	axState:=make(map[int] []opts.LineData)

	dendCharge:=make(map[int] []opts.LineData)
	dendCa:=make(map[int] []opts.LineData)
	dendACh:=make(map[int] []opts.LineData)
	dendCHOL:=make(map[int] []opts.LineData)
	for i:=0;i<16;i++{
		axState[i]=make([]opts.LineData,0)
		dendCharge[i]=make([]opts.LineData,0)
		dendCa[i]=make([]opts.LineData,0)
		dendACh[i]=make([]opts.LineData,0)
		dendCHOL[i]=make([]opts.LineData,0)
	}

	X:=make([]opts.LineData,0)
	cellCharge:=make([]opts.LineData,0)
	cellNa:=make([]opts.LineData,0)
	cellK:=make([]opts.LineData,0)

	for i:=0;i<200;i++{
		if (i>68&& i<70) ||(i>72&& i<74) || (i>76&& i<78) ||
			(i>80&& i<82) || (i>84&& i<86) || (i>88 && i<90) || (i>92 && i<94)||
			(i==95) ||(i==97) || (i==99) ||(i==101) ||(i==103) || (i==105){
			/*
				(i>68&& i<70) ||(i>72&& i<74) || (i>76&& i<78) ||
						(i>80&& i<82) || (i>84&& i<86) || (i>88 && i<90) || (i>92 && i<94)||
						(i==95) ||(i==97) || (i==99) ||(i==101) ||(i==103) || (i==105){

					(i>68&& i<70) ||(i>72&& i<74) || (i>76&& i<78) ||
					(i>80&& i<82) || (i>84&& i<86) || (i>88 && i<90) || (i>92 && i<94){

					i==5 || i==75  || i==150 || i==250 || i==350 {
			*/
			for j:=5;j<11;j++ {
				che[j].ACh = 50
			}

			/*} else if(i>150&& i<153) {
			*/
		} else{


		}

		n.DoDendrites(&g)
		n.DoLiveCicle(&g)
		n.DoAxons(&g)

		X=append(X,opts.LineData{Value:i})
		cellCharge=append(cellCharge,opts.LineData{Value:n.CalcCharge()})
		cellNa=append(cellNa,opts.LineData{Value:n.Chemic.Na})
		cellK=append(cellK,opts.LineData{Value:n.Chemic.K})

		for k:=0;k<16;k++{
			dendCharge[k]=append(dendCharge[k],opts.LineData{Value:n.Dendrites[k].Charge})
			dendCa[k]=append(dendCa[k],opts.LineData{Value:n.Dendrites[k].Ca})
			dendACh[k]=append(dendACh[k],opts.LineData{Value:che[k+1].ACh})
			dendCHOL[k]=append(dendCHOL[k],opts.LineData{Value:che[k+1].CHOL})

			axState[k]=append(axState[k],opts.LineData{Value:n.Axons[k].State})
		}
	}

	page:=components.NewPage()

	line := charts.NewLine()
	line.SetXAxis(X).
		SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeInfographic}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Заряд клетки",
		}))
	line.SetXAxis(X).
		AddSeries("Заряд клетки", cellCharge).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: false}))

	lineNaK := charts.NewLine()
	lineNaK.SetXAxis(X).
		SetGlobalOptions(
			charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeInfographic}),
			charts.WithTitleOpts(opts.Title{
				Title:    "Na-K",
			}))
	lineNaK.SetXAxis(X).
		AddSeries("Na", cellNa).
		AddSeries("K", cellK).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: false}))

	lineD := charts.NewLine()
	lineD.SetXAxis(X).
		SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeInfographic}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Заряд Дендритов",
		}))
	for i:=0;i<16;i++{
		lineD.AddSeries(strconv.Itoa(i), dendCharge[i])
	}

	lineDCa := charts.NewLine()
	lineDCa.SetXAxis(X).
		SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeInfographic}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Кальций Дендритов",
		}))
	for i:=0;i<16;i++{
		lineDCa.AddSeries(strconv.Itoa(i), dendCa[i])
	}

	lineDACh := charts.NewLine()
	lineDACh.SetXAxis(X).
		SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeInfographic}),
		charts.WithTitleOpts(opts.Title{
			Title:    "АЦХ Дендритов",
		}))
	for i:=0;i<16;i++{
		lineDACh.AddSeries(strconv.Itoa(i), dendACh[i])
	}

	lineDCHOL := charts.NewLine()
	lineDCHOL.SetXAxis(X).
		SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeInfographic}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Холин Дендритов",
		}))
	for i:=0;i<16;i++{
		lineDCHOL.AddSeries(strconv.Itoa(i), dendCHOL[i])
	}

	lineAstate := charts.NewLine()
	lineAstate.SetXAxis(X).
		SetGlobalOptions(
			charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeInfographic}),
			charts.WithTitleOpts(opts.Title{
				Title:    "Состояние аксонов",
			}))
	for i:=0;i<16;i++{
		lineAstate.AddSeries(strconv.Itoa(i), axState[i])
	}


	page.AddCharts(line,lineNaK,lineD,lineDCa,lineDACh,lineDCHOL,lineAstate).
		SetLayout(components.PageFlexLayout).
		PageTitle="Телеметрия нейрона"
	page.Render(w)
}
