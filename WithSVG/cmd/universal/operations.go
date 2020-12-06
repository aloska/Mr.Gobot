package universal


type Comm uint64

//команды состоят из 4 чисел 1-операция, 2 - куда положить результат, 2 другие - операнды
const (
	NOP 	Comm = 0 		//нет операции

	//арифметические и логические
	ADD 	Comm = 1		//сложение  ADD x1, x2, x3  -> x1=x2+x3
	ADDI	Comm = 2		//сложение с числом ADD x1, x2, с ->  x1=x2+с
	SUB		Comm = 3		//вычитание SUB x1,x2,x3 -> x1=x2+x3
	SUBI	Comm = 4		//вычитание с числом SUBI x1, x2, с ->  x1=x2-с
	MUL		Comm = 5		//умножение MUL x1,x2,x3 -> x1=x2*x3
	MULI	Comm = 6		//умножение на число MUL x1,x2,c -> x1=x2/c
	DIV		Comm = 7		//деление DIV x1,x2,x3 -> x1=x2*x3
	DIVI	Comm = 8		//деление на число DIV x1,x2,с -> x1=x2/с
	REM		Comm = 9		//остаток REM x1,x2,x3 -> x1=x2%x3
	REMI	Comm = 10		//остаток от деления на число REM x1,x2,с -> x1=x2%с
	AND		Comm = 11		//логическое и AND x1,x2,x3 -> x1=x2&x3
	ANDI	Comm = 12		//логическое и с числом ANDI x1,x2,c -> x1=x2&c
	OR		Comm = 13		//логическое или OR x1,x2,x3 -> x1=x2 | x3
	ORI		Comm = 14		//логическое или с числом ORI x1,x2,c -> x1=x2 | c
	XOR		Comm = 15		//логическое искл XOR x1,x2,x3 -> x1=x2^x3
	XORI	Comm = 16		//логическое искл с числом ORI x1,x2,c -> x1=x2^c
	SLL		Comm = 17		//сдвиг влево SLL x1,x2,x3 -> x1=x2<<x3
	SLLI	Comm = 18		//сдвиг влево с числом SLL x1,x2,с -> x1=x2<<с
	SRL		Comm = 19		//сдвиг вправо SLL x1,x2,x3 -> x1=x2<<x3
	SRLI	Comm = 20		//сдвиг вправо с числом SLL x1,x2,с -> x1=x2<<с

	//чтение и запись - все зациклено по делению с остатком. Внимание. Mnumber = Mnumber % len(M) и AddrMem = AddrMem % len(V)
	//невозможно адриссовать "вникуда" - если номер памяти или адрес в памяти превышает допустимое - автоматически будет взят остаток
	LDM		Comm = 21		//считать из памяти в регистр (LoaD from Memory) LDM x1, Mnumber, AddrMem -> x1=M[Mnumber].V[AddrMem]
	LDIN	Comm = 22 		//считать из входа (LoaD from In) LDIN x1, Snumber, AddrIn -> x1=In[Inumber].V[AddrIn]
	//считать значение выхода нельзя!
	STM		Comm = 23		//записать в память (STore to Memory) STM Mnumber, AddrMem, x1 ->  M[Mnumber].V[AddrMem] = x1
	STOUT	Comm = 24       //записать в выход (STore to Out) STOUT Onumber, AddrOut, x1 ->  Out[Onumber].V[AddrOut] = x1
	//записать значение во вход нельзя!

	//команды перехода - все зациклено, перейти в никуда нельзя.
	BEQ		Comm = 25		//перейти на столько то вперед или назад, если значения регистров равны. BEQ x1,x2,JumpAddr -> if x1==x2 jump to PC+JumpAddr
	BGE		Comm = 26 		//перейти если больше или равно BGE x1,x2,JumpAddr -> if x1>=x2 jump to PC+JumpAddr
	BLT		Comm = 27 		//перейти если меньше BLT x1,x2,JumpAddr -> if x1<x2 jump to PC+JumpAddr
	BNE		Comm = 28 		//перейти если не равно BNE x1,x2,JumpAddr -> if x1!=x2 jump to PC+JumpAddr
	JMP		Comm = 29		//перейти безусловно JMP JumpAddr

	//Комманды сравнения без перехода
	SEQ		Comm = 30		//Установить 1 если равно. SEQ x1,x2,x3 -> x1=(x2==x3)?1:0
	SGE		Comm = 31		//Установить 1 если больше или равно. SGE x1,x2,x3 -> x1=(x2>=x3)?1:0
	SLT		Comm = 32		//Установить 1 если меньше чем. SLT x1,x2,x3 -> x1=(x2<x3)?1:0
	SNE		Comm = 33		//Установить 1 если не ранво. SNE x1,x2,x3 -> x1=(x2!=x3)?1:0

	//операции со стеком
	PUSH 	Comm = 34		//засунуть в стек PUSH x1, any, any - два последних значения в гене этой коммады не имеют смысла, там могут быть любые числа
	POP		Comm = 35		//достать из стека POP x1, any,any -  два последних значения в гене этой коммады не имеют смысла, там могут быть любые числа

	//резервные комманды, ничего не делающие, но дающие возможность разнообразить геном
	NOP1	Comm = 36
	NOP2	Comm = 37
	NOP3	Comm = 38
	NOP4	Comm = 39
	NOP5	Comm = 40
	NOP6	Comm = 41
	NOP7	Comm = 42

	/*CISC мб лучше??
	ADDI	Comm = 		//сложение с входом ADDI x1, x2, addrInput  x1=x2 + (addrInput)
	ADDIC	Comm = 		//сложение с входом числа ADDI x1, addrInput, c   x1=(addrInput) + c
	ADDM	Comm = 		//сложение с числом из файла ADDM x1, x2, addrMem
	ADDMC	Comm = 		//сложение числа из файла с константой ADDMС x1, addrMem, c
	ADDMI	Comm = 		//сложение числа из файла с входом ADDMI x1, addrMem, addrInput
	....
	*/
)

//ADD - сложение
//x1, x2, x3 - номера регистров общего назначения (если больше 32 - по кругу)
func (p *Processor) ADD(x1, x2, x3 uint64){
	p.PC++
	p.X[x1%32]=p.X[x2%32]+p.X[x3%32]
}

//ADDI - сложение
func (p *Processor) ADDI(x1, x2 uint64, c int64){
	p.PC++
	p.X[x1%32]=p.X[x2%32]+c
}

func (p *Processor) SUB(x1, x2, x3 uint64){
	p.PC++
	p.X[x1%32]=p.X[x2%32]-p.X[x3%32]
}

func (p *Processor) SUBI(x1, x2 uint64, c int64){
	p.PC++
	p.X[x1%32]=p.X[x2%32]-c
}

func (p *Processor) MUL(x1, x2, x3 uint64){
	p.PC++
	p.X[x1%32]=p.X[x2%32] * p.X[x3%32]
}

func (p *Processor) MULI(x1, x2 uint64, c int64){
	p.PC++
	p.X[x1%32]=p.X[x2%32]*c
}

func (p *Processor) DIV(x1, x2, x3 uint64){
	p.PC++
	if p.X[x3%32]==0{
		p.X[x1%32]=0
	}else {
		p.X[x1%32] = p.X[x2%32] / p.X[x3%32]
	}
}

func (p *Processor) DIVI(x1, x2 uint64, c int64){
	p.PC++
	if c==0{
		p.X[x1%32]=0
	}else {
		p.X[x1%32] = p.X[x2%32] / c
	}
}

func (p *Processor) REM(x1, x2, x3 uint64){
	p.PC++
	if  p.X[x3%32]==0{
		p.X[x1%32]=0
	}else {
		p.X[x1%32] = p.X[x2%32] % p.X[x3%32]
	}
}

func (p *Processor) REMI(x1, x2 uint64, c int64){
	p.PC++
	if c==0{
		p.X[x1%32]=0
	}else {
		p.X[x1%32] = p.X[x2%32] % c
	}
}

func (p *Processor) AND(x1, x2, x3 uint64){
	p.PC++
	p.X[x1%32]=p.X[x2%32] & p.X[x3%32]
}

func (p *Processor) ANDI(x1, x2 uint64, c int64){
	p.PC++
	p.X[x1%32]=p.X[x2%32] & c
}

func (p *Processor) OR(x1, x2, x3 uint64){
	p.PC++
	p.X[x1%32]=p.X[x2%32] | p.X[x3%32]
}

func (p *Processor) ORI(x1, x2 uint64, c int64){
	p.PC++
	p.X[x1%32]=p.X[x2%32] | c
}

func (p *Processor) XOR(x1, x2, x3 uint64){
	p.PC++
	p.X[x1%32]=p.X[x2%32] ^ p.X[x3%32]
}

func (p *Processor) XORI(x1, x2 uint64, c int64){
	p.PC++
	p.X[x1%32]=p.X[x2%32] ^ c
}

func (p *Processor) SLL(x1, x2, x3 uint64){
	p.PC++
	p.X[x1%32]=p.X[x2%32] << uint64(p.X[x3%32])
}

func (p *Processor) SLLI(x1, x2 uint64, c int64){
	p.PC++
	p.X[x1%32]=p.X[x2%32] << uint64(c)
}

func (p *Processor) SRL(x1, x2, x3 uint64){
	p.PC++
	p.X[x1%32]=p.X[x2%32] >> uint64(p.X[x3%32])
}

func (p *Processor) SRLI(x1, x2 uint64, c int64){
	p.PC++
	p.X[x1%32]=p.X[x2%32] >> uint64(c)
}

func (p *Processor) SEQ(x1, x2, x3 uint64){
	p.PC++
	if p.X[x2%32]==p.X[x3%32]{
		p.X[x1%32]=1
	}else{
		p.X[x1%32]=0
	}
}

func (p *Processor) SGE(x1, x2, x3 uint64){
	p.PC++
	if p.X[x2%32]>=p.X[x3%32]{
		p.X[x1%32]=1
	}else{
		p.X[x1%32]=0
	}
}

func (p *Processor) SLT(x1, x2, x3 uint64){
	p.PC++
	if p.X[x2%32]<p.X[x3%32]{
		p.X[x1%32]=1
	}else{
		p.X[x1%32]=0
	}
}

func (p *Processor) SNE(x1, x2, x3 uint64){
	p.PC++
	if p.X[x2%32]!=p.X[x3%32]{
		p.X[x1%32]=1
	}else{
		p.X[x1%32]=0
	}
}

//круговой пуш - указатель стека переходит через xff
func (p *Processor) PUSH(x1 uint64){
	p.PC++
	p.SI++
	p.S[p.SI]=p.X[x1%32]
}

func (p *Processor) POP(x1 uint64) {
	p.PC++
	p.X[x1%32]=p.S[p.SI]
	p.SI--
}

//операции чтения/записи выполняет решатель
func (s *Solution) LDM(x1, mnumber, maddr uint64){
	s.Proc.PC++
	mn:=mnumber%uint64(len(s.Mem))
	ma:=maddr%uint64(len(s.Mem[mn].V))
	s.Proc.X[x1%32]=s.Mem[mn].V[ma]
}

func (s *Solution) LDIN(x1, innumber, inaddr uint64){
	s.Proc.PC++
	mn:=innumber%uint64(len(s.In))
	ma:=inaddr%uint64(len(s.In[mn].V))
	s.Proc.X[x1%32]=s.In[mn].V[ma]
}

func (s *Solution) STM(mnumber, maddr, x1 uint64){
	s.Proc.PC++
	mn:=mnumber%uint64(len(s.Mem))
	ma:=maddr%uint64(len(s.Mem[mn].V))
	s.Mem[mn].V[ma]=s.Proc.X[x1%32]
}

func (s *Solution) STOUT(onumber, oaddr, x1 uint64){
	s.Proc.PC++
	mn:=onumber%uint64(len(s.Out))
	ma:=oaddr%uint64(len(s.Out[mn].V))
	s.Out[mn].V[ma]=s.Proc.X[x1%32]
}

/*
по сути изменения указателя на следующую комманду
*/
func (s *Solution) BEQ(x1, x2 uint64, jumpAddr int64){
	if jumpAddr==0{//а то блокировка))
		s.Proc.PC++
		return
	}
	if s.Proc.X[x1%32]==s.Proc.X[x2%32]{
		//сначала выровняем jumpAddr по длине гена
		jumpAddr=jumpAddr%int64(len(s.Gen.Codons))
		if int64(s.Proc.PC)+jumpAddr<0{
			s.Proc.PC=uint64(int64(len(s.Gen.Codons))+int64(s.Proc.PC)+jumpAddr) //не знаю почему, но работает))
		}else{
			s.Proc.PC=uint64(int64(s.Proc.PC)+jumpAddr)%uint64(len(s.Gen.Codons))
		}
	}else{
		s.Proc.PC++
	}
}

func (s *Solution) BGE(x1, x2 uint64, jumpAddr int64){
	if jumpAddr==0{//а то блокировка))
		s.Proc.PC++
		return
	}
	if s.Proc.X[x1%32]>=s.Proc.X[x2%32]{
		//сначала выровняем jumpAddr по длине гена
		jumpAddr=jumpAddr%int64(len(s.Gen.Codons))
		if int64(s.Proc.PC)+jumpAddr<0{
			s.Proc.PC=uint64(int64(len(s.Gen.Codons))+int64(s.Proc.PC)+jumpAddr)
		}else{
			s.Proc.PC=uint64(int64(s.Proc.PC)+jumpAddr)%uint64(len(s.Gen.Codons))
		}
	}else{
		s.Proc.PC++
	}
}

func (s *Solution) BLT(x1, x2 uint64, jumpAddr int64){
	if jumpAddr==0{//а то блокировка))
		s.Proc.PC++
		return
	}
	if s.Proc.X[x1%32]<s.Proc.X[x2%32]{
		//сначала выровняем jumpAddr по длине гена
		jumpAddr=jumpAddr%int64(len(s.Gen.Codons))
		if int64(s.Proc.PC)+jumpAddr<0{
			s.Proc.PC=uint64(int64(len(s.Gen.Codons))+int64(s.Proc.PC)+jumpAddr)
		}else{
			s.Proc.PC=uint64(int64(s.Proc.PC)+jumpAddr)%uint64(len(s.Gen.Codons))
		}
	}else{
		s.Proc.PC++
	}
}

func (s *Solution) BNE(x1, x2 uint64, jumpAddr int64){
	if jumpAddr==0{//а то блокировка))
		s.Proc.PC++
		return
	}
	if s.Proc.X[x1%32]!=s.Proc.X[x2%32]{
		//сначала выровняем jumpAddr по длине гена
		jumpAddr=jumpAddr%int64(len(s.Gen.Codons))
		if int64(s.Proc.PC)+jumpAddr<0{
			s.Proc.PC=uint64(int64(len(s.Gen.Codons))+int64(s.Proc.PC)+jumpAddr)
		}else{
			s.Proc.PC=uint64(int64(s.Proc.PC)+jumpAddr)%uint64(len(s.Gen.Codons))
		}
	}else{
		s.Proc.PC++
	}
}

func (s *Solution) JMP(jumpAddr int64){
	if jumpAddr==0{//а то блокировка))
		s.Proc.PC++
		return
	}
	//сначала выровняем jumpAddr по длине гена
	jumpAddr=jumpAddr%int64(len(s.Gen.Codons))
	if int64(s.Proc.PC)+jumpAddr<0{
		s.Proc.PC=uint64(int64(len(s.Gen.Codons))+int64(s.Proc.PC)+jumpAddr)
	}else{
		s.Proc.PC=uint64(int64(s.Proc.PC)+jumpAddr)%uint64(len(s.Gen.Codons))
	}
}

