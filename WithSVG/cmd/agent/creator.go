package agent

import "encoding/binary"

//генерация файла рецепторов
func (r *Receptors) createReceptorsFile() error {
	var recs          []Receptor

	for genindex, rrr:=range r.genes{
		cur:=0
		for i:=rrr.Ndata; i<= rrr.Maxi; i=i+rrr.Iteri{
			for j:=rrr.NdataW; j<= rrr.Maxj; j=j+rrr.Iterj{
				for k:=rrr.NdataWB; k<=rrr.Maxk; k=k+rrr.Iterk {
						recs = append(recs, Receptor{
						Threshold: rrr.Recep.Threshold,
						Typer: rrr.Typer,
						Coren: rrr.Coren,
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
