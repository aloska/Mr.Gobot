package agent

import (
	"encoding/binary"
	"errors"
	"strconv"
)

//генерация файла рецепторов
func (r *Receptors) createReceptorsFile() error {
	var recs          []Receptor

	for genindex, rrr:=range r.genes{
		//сначала проверим, что ген в порядке
		switch {
		case rrr.MaxX==0:
			return errors.New("Ген рецептора №"+strconv.Itoa(genindex)+", "+
				r.filenameGens+" MaxX==0")
		case rrr.Ndata>rrr.Maxi:
			return errors.New("Ген рецептора №"+strconv.Itoa(genindex)+", "+
				r.filenameGens+" Ndata>Maxi")
		case rrr.NdataW>rrr.Maxj:
			return errors.New("Ген рецептора №"+strconv.Itoa(genindex)+", "+
				r.filenameGens+" NdataW>Maxj")
		case rrr.NdataWB>rrr.Maxk:
			return errors.New("Ген рецептора №"+strconv.Itoa(genindex)+", "+
				r.filenameGens+" NdataW>Maxj")
		case rrr.Iteri==0 || rrr.Iterj==0 || rrr.Iterk==0:
			return errors.New("Ген рецептора №"+strconv.Itoa(genindex)+", "+
				r.filenameGens+" Iter[ijk] должны быть > 0")

		}

		cur:=0
		for i:=rrr.Ndata; i<= rrr.Maxi; i=i+rrr.Iteri{
			for j:=rrr.NdataW; j<= rrr.Maxj; j=j+rrr.Iterj{
				for k:=rrr.NdataWB; k<=rrr.Maxk; k=k+rrr.Iterk {
						recs = append(recs, Receptor{
						Threshold: rrr.Recep.Threshold,
						Typer: rrr.Typer,
						SynNumber: rrr.SynNumber,
						Ndata: i,
						NdataW: j,
						NdataWb: k,
						Typemedi: rrr.Typemedi,
						Force: rrr.Recep.Force,
						Divforce: rrr.Recep.Divforce,
						Gen: uint16(genindex),
						A: rrr.Recep.A,
						Axons: [32]Axon{} })
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
				//уровни кальция, везикулы... берем из типиного рецептора в генах
				recs[c].Axons[i].Ca=rrr.Recep.Axons[i].Ca
				recs[c].Axons[i].Vesiculs=rrr.Recep.Axons[i].Vesiculs
				//Номер синапса в файле синапсов генерируется, исходя из координат
				recs[c].Axons[i].N=XYToNumber(x,y,rrr.MaxX)

				x=x+uint32(rrr.AxShiftX)
				y=y+uint32(rrr.AxShiftY)

			}
			x=xs+uint32(rrr.AxNextShiftX)
			y=ys+uint32(rrr.AxNextShiftY)
		}
	}

	if err:=StructsFileWrite(r.filenameRecs,recs,binary.LittleEndian); err!=nil{

		return err
	}
	return nil
}

//генерация файла нейронов
func (ce *Cells) createNeuronsFile() error{
	var neus  []Neuron

	for genindex, nnn:=range ce.genes{
		//сначала проверим, что ген в порядке
		switch {
		case nnn.MaxX==0:
			return errors.New("Ген нейрона №"+strconv.Itoa(genindex)+", "+
				ce.filenameGens+" MaxX==0")
		case nnn.Layers==0:
			return errors.New("Ген нейрона №"+strconv.Itoa(genindex)+", "+
				ce.filenameGens+" Layers==0")
		case nnn.Laywidth==0:
			return errors.New("Ген нейрона №"+strconv.Itoa(genindex)+", "+
				ce.filenameGens+" Laywidth==0")
		case nnn.SynNumber!=nnn.SynNumberAxons:
			//аксоны лежат в более другом поле, чем сома и дендриты,
			//а значит
			if nnn.MaxXOtherSyn==0{
				//такого не должно быть
				return errors.New("Ген нейрона №"+strconv.Itoa(genindex)+", "+
					ce.filenameGens+" SynNumber!=SynNumberAxons но MaxXOtherSyn==0")
			}
		}

		//генерируем сомы
		x:=nnn.Soma1stX
		y:=nnn.Soma1stY
		for lay:=0; lay<int(nnn.Layers);lay++{
			xs:=x
			ys:=y
			for i:=0;i<int(nnn.Laywidth);i++{
				neus=append(neus, Neuron{
					Typen: nnn.Typen,
					SynNumber: nnn.SynNumber,
					Gen: uint16(genindex),
					SynNumberAxons: nnn.SynNumberAxons,
					Chemic: nnn.Neur.Chemic,
					N: XYToNumber(x,y,nnn.MaxX),
					D: nnn.Neur.D,
					A: nnn.Neur.A,
					Dendrites: [16]Dendrite{},
					Axons: [16]Axon{} })

				x=x+uint32(nnn.SomaNextShiftX)
				y=y+uint32(nnn.SomaNextShiftY)
			}
			x=xs+uint32(nnn.SomaLayerShiftX)
			y=ys+uint32(nnn.SomaLayerShiftY)
		}

		//генерируем дендриты
		cur:=0 //для отслеживания индекса в генерируемом слайсе
		x=nnn.Dend1stX
		y=nnn.Dend1stY
		for lay:=0; lay<int(nnn.Layers);lay++{
			xl:=x
			yl:=y
			for i:=0;i<int(nnn.Laywidth);i++{
				xs:=x
				ys:=y
				for d:=0;d<16;d++ {
					neus[cur].Dendrites[d].Ca=nnn.Neur.Dendrites[d].Ca
					neus[cur].Dendrites[d].Typed=nnn.Neur.Dendrites[d].Typed
					neus[cur].Dendrites[d].N=XYToNumber(x,y,nnn.MaxX)

					x = x + uint32(nnn.DendShiftX)
					y = y + uint32(nnn.DendShiftY)
				}
				cur++
				x = xs + uint32(nnn.DendNextShiftX)
				y = ys + uint32(nnn.DendNextShiftY)
			}
			x=xl+uint32(nnn.DendLayerShiftX)
			y=yl+uint32(nnn.DendLayerShiftY)
		}

		//генерируем аксоны
		x=nnn.Ax1stX
		y=nnn.Ax1stY
		cur=0 //для отслеживания индекса в генерируемом слайсе
		for lay:=0; lay<int(nnn.Layers);lay++{
			xl:=x
			yl:=y
			for i:=0;i<int(nnn.Laywidth);i++{
				xs:=x
				ys:=y
				for d:=0;d<16;d++ {
					neus[cur].Axons[d].Ca=nnn.Neur.Axons[d].Ca
					neus[cur].Axons[d].Vesiculs=nnn.Neur.Axons[d].Vesiculs
					if neus[cur].SynNumber!= neus[cur].SynNumberAxons{
						//аксоны в другом файле синапсов, чем сома и дендриты
						neus[cur].Axons[d].N = XYToNumber(x, y, nnn.MaxXOtherSyn)
					}else {
						neus[cur].Axons[d].N = XYToNumber(x, y, nnn.MaxX)
					}

					x = x + uint32(nnn.AxShiftX)
					y = y + uint32(nnn.AxShiftY)
				}
				cur++
				x = xs + uint32(nnn.AxNextShiftX)
				y = ys + uint32(nnn.AxNextShiftY)
			}
			x=xl+uint32(nnn.AxLayerShiftX)
			y=yl+uint32(nnn.AxLayerShiftY)
		}
	}

	if err:=StructsFileWrite(ce.filenameCells,neus,binary.LittleEndian); err!=nil{
		return err
	}


	return nil
}
