package agent

func (r *Receptor) Step(gene *GenReceptor, datai *DataInput){
	//сначала все аксоны выполняют синтез вещества на окончании
	for i:=0;i<32;i++{
		if r.A&(1<<i) != 0 { //проверяем, что данный аксон включен
			switch r.Typemedi {
			case NEURONACETILHOLIN: //ацетилхолин
				//вычисляем ячейку в файле синапсов, от куда синтез идет o.synapsesMap[r.SynNumber].syn[r.Axons[i].N]
				//один раз - по любому
				r.Axons[i].AChSynt(&org.synapsesMap[r.SynNumber].syn[r.Axons[i].N])
				//и еще столько раз, сколько ген захотел
				for j:=byte(0);j<gene.SyntPerCicleAx;j++{
					//ну и не фиг стараться, если больше не получается
					if !r.Axons[i].AChSynt(&org.synapsesMap[r.SynNumber].syn[r.Axons[i].N]){
						break
					}
				}
			}
		}
	}
	//анализ данных и плювание
	switch r.Typer {
	case RECEPTORDATAUINT32BIGPOS:
		r.DoReceptorUInt32(&datai.dataUInt32[r.Ndata],org.synapsesMap[r.SynNumber])
	}
	org.wgo.Done()
}

func (p *Preffector) Step( gene *GenPreffector){
	org.wgo.Done()
}

func (n *Neuron) Step(gene *GenNeuron){
	org.wgo.Done()
}

func (res *Receptors) Step(datai *DataInput){
	//по всем рецепторам
	for i:=0;i<len(res.recs);i++{
		org.wgo.Add(1)
		go res.recs[i].Step( &res.genes[res.recs[i].Gen], datai)
	}
	org.wgo.Done()
}

func (pres *Preffectors) Step(){
	//по всем преффекторам
	for i:=0;i<len(pres.prefs);i++{
		org.wgo.Add(1)
		go pres.prefs[i].Step(&pres.genes[pres.prefs[i].Gen])
	}
	org.wgo.Done()
}

func (ce *Cells) Step(){
	//по всем нейронам
	for i:=0;i<len(ce.neurons);i++{
		org.wgo.Add(1)
		go ce.neurons[i].Step(&ce.genes[ce.neurons[i].Gen])
	}
	org.wgo.Done()
}

func (in *Input) Step(){
	//по всем рецепторам
	for i:=0;i<len(in.receptors);i++{
		org.wgo.Add(1)
		go in.receptors[i].Step(&in.dataInput)
	}
	//и нейронам
	for i:=0;i<len(in.cells);i++{
		org.wgo.Add(1)
		go in.cells[i].Step()
	}
	org.wgo.Done()
}

func (ef *Effector) Step(){
	//по всем преффекторам
	for i:=0;i<len(ef.preffectors);i++{
		org.wgo.Add(1)
		go ef.preffectors[i].Step()
	}
	//и нейронам
	for i:=0;i<len(ef.cells);i++{
		org.wgo.Add(1)
		go ef.cells[i].Step()
	}
	org.wgo.Done()
}

//Step - один шаг жизни чувств
func (s *Senses) Step(){
	//бежим по всем входам
	for i:=0;i<len(s.inputs);i++{
		org.wgo.Add(1)
		go s.inputs[i].Step()
	}
	//и по нейронам общим, если есть
	for i:=0;i<len(s.cells);i++{
		org.wgo.Add(1)
		go s.cells[i].Step()
	}
	//это в конце
	org.wgo.Done()
}



func (co *Core) Step (){
	//по всем Cells
	for i:=0;i<len(co.cells);i++{
		org.wgo.Add(1)
		go co.cells[i].Step()
	}
	org.wgo.Done()
}

//Step - один шаг жизни мозга
func (b *Brain) Step(){
    //по всем ядрам
	for i:=0;i<len(b.cores);i++{
		org.wgo.Add(1)
		go b.cores[i].Step()
	}
	//это в конце
	org.wgo.Done()
}


//Step - один шаг жизни действий
func (ac *Actions) Step(){
	//бежим по всем выходам
	for i:=0;i<len(ac.effectors);i++{
		org.wgo.Add(1)
		go ac.effectors[i].Step()
	}
	//и по нейронам общим, если есть
	for i:=0;i<len(ac.cells);i++{
		org.wgo.Add(1)
		go ac.cells[i].Step()
	}
	//это в конце
	org.wgo.Done()
}

//Step - один шаг жизни вегетатики
func (v *Vegetatic) Step(){
	//бежим по всем выходам
	for i:=0;i<len(v.effectors);i++{
		org.wgo.Add(1)
		go v.effectors[i].Step()
	}
	//и по нейронам общим, если есть
	for i:=0;i<len(v.cells);i++{
		org.wgo.Add(1)
		go v.cells[i].Step()
	}

	//это в конце
	org.wgo.Done()
}
