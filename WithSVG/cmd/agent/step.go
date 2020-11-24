package agent

func (r *Receptor) Step(o *Organism, gene *GenReceptor){
	o.wgo.Done()
}

func (p *Preffector) Step(o *Organism, gene *GenPreffector){
	o.wgo.Done()
}

func (n *Neuron) Step(o *Organism, gene *GenNeuron){
	o.wgo.Done()
}

func (res *Receptors) Step(o *Organism){
	//по всем рецепторам
	for i:=0;i<len(res.recs);i++{
		o.wgo.Add(1)
		go res.recs[i].Step(o, &res.genes[res.recs[i].Gen])
	}
	o.wgo.Done()
}

func (pres *Preffectors) Step(o *Organism){
	//по всем преффекторам
	for i:=0;i<len(pres.prefs);i++{
		o.wgo.Add(1)
		go pres.prefs[i].Step(o, &pres.genes[pres.prefs[i].Gen])
	}
	o.wgo.Done()
}

func (ce *Cells) Step(o *Organism){
	//по всем нейронам
	for i:=0;i<len(ce.neurons);i++{
		o.wgo.Add(1)
		go ce.neurons[i].Step(o, &ce.genes[ce.neurons[i].Gen])
	}
	o.wgo.Done()
}

func (in *Input) Step(o *Organism){
	//по всем рецепторам
	for i:=0;i<len(in.receptors);i++{
		o.wgo.Add(1)
		go in.receptors[i].Step(o)
	}
	//и нейронам
	for i:=0;i<len(in.cells);i++{
		o.wgo.Add(1)
		go in.cells[i].Step(o)
	}
	o.wgo.Done()
}

func (ef *Effector) Step(o *Organism){
	//по всем преффекторам
	for i:=0;i<len(ef.preffectors);i++{
		o.wgo.Add(1)
		go ef.preffectors[i].Step(o)
	}
	//и нейронам
	for i:=0;i<len(ef.cells);i++{
		o.wgo.Add(1)
		go ef.cells[i].Step(o)
	}
	o.wgo.Done()
}

//Step - один шаг жизни чувств
func (s *Senses) Step(){
	//бежим по всем входам
	for i:=0;i<len(s.inputs);i++{
		s.organism.wgo.Add(1)
		go s.inputs[i].Step(s.organism)
	}
	//и по нейронам общим, если есть
	for i:=0;i<len(s.cells);i++{
		s.organism.wgo.Add(1)
		go s.cells[i].Step(s.organism)
	}
	//это в конце
	s.organism.wgo.Done()
}



func (co *Core) Step (o *Organism){
	//по всем Cells
	for i:=0;i<len(co.cells);i++{
		o.wgo.Add(1)
		go co.cells[i].Step(o)
	}
	o.wgo.Done()
}

//Step - один шаг жизни мозга
func (b *Brain) Step(){
    //по всем ядрам
	for i:=0;i<len(b.cores);i++{
		b.organism.wgo.Add(1)
		go b.cores[i].Step(b.organism)
	}
	//это в конце
	b.organism.wgo.Done()
}


//Step - один шаг жизни действий
func (ac *Actions) Step(){
	//бежим по всем выходам
	for i:=0;i<len(ac.effectors);i++{
		ac.organism.wgo.Add(1)
		go ac.effectors[i].Step(ac.organism)
	}
	//и по нейронам общим, если есть
	for i:=0;i<len(ac.cells);i++{
		ac.organism.wgo.Add(1)
		go ac.cells[i].Step(ac.organism)
	}
	//это в конце
	ac.organism.wgo.Done()
}

//Step - один шаг жизни вегетатики
func (v *Vegetatic) Step(){
	//бежим по всем выходам
	for i:=0;i<len(v.effectors);i++{
		v.organism.wgo.Add(1)
		go v.effectors[i].Step(v.organism)
	}
	//и по нейронам общим, если есть
	for i:=0;i<len(v.cells);i++{
		v.organism.wgo.Add(1)
		go v.cells[i].Step(v.organism)
	}

	//это в конце
	v.organism.wgo.Done()
}
