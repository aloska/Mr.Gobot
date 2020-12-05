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

/*
//константы регистров нужны ли? todo номер - просто индекс в массиве?
type Regg uint64
//чтобы узнать номер регистра, нужно по модулю разделить: ADD 12, 45, 3 - на самом деле суть ADD X[12%16], X[45%16], X[3%16] = ADD X[12], X[13], X[3]
const(
	X0	Regg = 0 			//X[0]
	X1	Regg = 1
	X2	Regg = 2
	X3	Regg = 3
	X4	Regg = 4
	X5	Regg = 5
	X6	Regg = 6
	X7	Regg = 7
	X8	Regg = 8
	X9	Regg = 9
	X10	Regg = 10
	X11	Regg = 11
	X12	Regg = 12
	X13	Regg = 13
	X14	Regg = 14
	X15	Regg = 15
	....
)
*/



