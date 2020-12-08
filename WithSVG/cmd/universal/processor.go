package universal

type Registers struct{
	X 		[32]int64	//регистры общего назначения
	PC		uint64		//указатель текущей позиции в геноме (выполняемой программе)
	SI		byte		//указатель стека
	Serv1	byte
	Serv2	uint16
	Serv3	uint32
}

type Processor struct{
	Registers
	S	[256]int64
}




