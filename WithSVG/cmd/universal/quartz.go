package universal

func (s *Solution) Run(){
	for int(s.Proc.PC)<len(s.Gen.Codons){
		s.Step()
	}
}

func (s *Solution) Step(){
	//42 - кол-во комманд в нашем RISC-процессоре
	com:=s.Gen.Codons[s.Proc.PC]
	switch com.Code%42 {
	case NOP,NOP1,NOP2,NOP3,NOP4,NOP5,NOP6: //а что уж тут поделаешь?
		s.Proc.PC++	//хех!
	case ADD:
		s.Proc.ADD(com.Op1, com.Op2,uint64(com.Op3))
	case ADDI:
		s.Proc.ADDI(com.Op1, com.Op2,com.Op3)
	case SUB:
		s.Proc.SUB(com.Op1, com.Op2,uint64(com.Op3))
	case SUBI:
		s.Proc.SUBI(com.Op1, com.Op2,com.Op3)
	case MUL:
		s.Proc.MUL(com.Op1, com.Op2,uint64(com.Op3))
	case MULI:
		s.Proc.MULI(com.Op1, com.Op2,com.Op3)
	case DIV:
		s.Proc.DIV(com.Op1, com.Op2,uint64(com.Op3))
	case DIVI:
		s.Proc.DIVI(com.Op1, com.Op2,com.Op3)
	case REM:
		s.Proc.REM(com.Op1, com.Op2,uint64(com.Op3))
	case REMI:
		s.Proc.REMI(com.Op1, com.Op2,com.Op3)
	case AND:
		s.Proc.AND(com.Op1, com.Op2,uint64(com.Op3))
	case ANDI:
		s.Proc.ANDI(com.Op1, com.Op2,com.Op3)
	case OR:
		s.Proc.OR(com.Op1, com.Op2,uint64(com.Op3))
	case ORI:
		s.Proc.ORI(com.Op1, com.Op2,com.Op3)
	case XOR:
		s.Proc.XOR(com.Op1, com.Op2,uint64(com.Op3))
	case XORI:
		s.Proc.XORI(com.Op1, com.Op2,com.Op3)
	case SLL:
		s.Proc.SLL(com.Op1, com.Op2,uint64(com.Op3))
	case SLLI:
		s.Proc.SLLI(com.Op1, com.Op2,com.Op3)
	case SRL:
		s.Proc.SRL(com.Op1, com.Op2,uint64(com.Op3))
	case SRLI:
		s.Proc.SRLI(com.Op1, com.Op2,com.Op3)
	case LI:
		s.Proc.LI(com.Op1, com.Op3)
	case LDM:
		s.LDM(com.Op1, com.Op2,uint64(com.Op3))
	case LDIN:
		s.LDIN(com.Op1, com.Op2,uint64(com.Op3))
	case STM:
		s.STM(com.Op1, com.Op2,uint64(com.Op3))
	case STOUT:
		s.STOUT(com.Op1, com.Op2,uint64(com.Op3))
	case BEQ:
		s.BEQ(com.Op1, com.Op2,com.Op3)
	case BGE:
		s.BGE(com.Op1, com.Op2,com.Op3)
	case BLT:
		s.BLT(com.Op1, com.Op2,com.Op3)
	case BNE:
		s.BNE(com.Op1, com.Op2,com.Op3)
	case JMP:
		s.JMP(com.Op3)//Внимание! В JMP адресс в 3 операнде (как впрочем и везде - просто здесь только один операнд, но он в кодоне 3-ий!!)
	case SEQ:
		s.Proc.SEQ(com.Op1, com.Op2,uint64(com.Op3))
	case SGE:
		s.Proc.SGE(com.Op1, com.Op2,uint64(com.Op3))
	case SLT:
		s.Proc.SLT(com.Op1, com.Op2,uint64(com.Op3))
	case SNE:
		s.Proc.SNE(com.Op1, com.Op2,uint64(com.Op3))
	case PUSH:
		s.Proc.PUSH(com.Op1)
	case POP:
		s.Proc.POP(com.Op1)
	}
}
