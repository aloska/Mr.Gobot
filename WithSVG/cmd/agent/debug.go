package agent

func (sy *Synapses)SetDebugSyn(che []Chemical){
	sy.syn=che
}

//SetDebugOrganism - чтобы тестировать в каком-то левом организме разное всякое
func (a *Agent) SetDebugOrganism(o *Organism){
	org=o
}

func (o *Organism) SetDebugSynMap(sm map[SynEnum] *Synapses){
	o.synapsesMap=sm
}