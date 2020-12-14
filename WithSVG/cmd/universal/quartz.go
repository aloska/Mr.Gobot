package universal

import "sync"

func (s *Solution) Run(){
	var wg sync.WaitGroup
	for i:=0;i<len(s.Algs);i++ {
		if s.IsAsync{
			wg.Add(1)
			go func(alg int){
				defer wg.Done()
				for int(s.Proc[alg].PC) < len(s.Algs[alg].Commands) {
					s.Step(alg)
				}
			}(i)
		}else {
			for int(s.Proc[i].PC) < len(s.Algs[i].Commands) {
				s.Step(i)
			}
		}
	}
	if s.IsAsync{
		wg.Wait()
	}
}

func (s *Solution) Step(alg int){
	com:=s.Algs[alg].Commands[s.Proc[alg].PC]
	switch com.Code%COUNTCOMMAND { //43 - кол-во комманд в нашем RISC-процессоре
	case NOP: //а что уж тут поделаешь?
		s.Proc[alg].PC++	//хех!
	case ADD:
		s.Proc[alg].ADD(com.Op1, com.Op2,uint64(com.Op3))
	case ADDI:
		s.Proc[alg].ADDI(com.Op1, com.Op2,com.Op3)
	case SUB:
		s.Proc[alg].SUB(com.Op1, com.Op2,uint64(com.Op3))
	case SUBI:
		s.Proc[alg].SUBI(com.Op1, com.Op2,com.Op3)
	case MUL:
		s.Proc[alg].MUL(com.Op1, com.Op2,uint64(com.Op3))
	case MULI:
		s.Proc[alg].MULI(com.Op1, com.Op2,com.Op3)
	case DIV:
		s.Proc[alg].DIV(com.Op1, com.Op2,uint64(com.Op3))
	case DIVI:
		s.Proc[alg].DIVI(com.Op1, com.Op2,com.Op3)
	case REM:
		s.Proc[alg].REM(com.Op1, com.Op2,uint64(com.Op3))
	case REMI:
		s.Proc[alg].REMI(com.Op1, com.Op2,com.Op3)
	case AND:
		s.Proc[alg].AND(com.Op1, com.Op2,uint64(com.Op3))
	case ANDI:
		s.Proc[alg].ANDI(com.Op1, com.Op2,com.Op3)
	case OR:
		s.Proc[alg].OR(com.Op1, com.Op2,uint64(com.Op3))
	case ORI:
		s.Proc[alg].ORI(com.Op1, com.Op2,com.Op3)
	case XOR:
		s.Proc[alg].XOR(com.Op1, com.Op2,uint64(com.Op3))
	case XORI:
		s.Proc[alg].XORI(com.Op1, com.Op2,com.Op3)
	case SLL:
		s.Proc[alg].SLL(com.Op1, com.Op2,uint64(com.Op3))
	case SLLI:
		s.Proc[alg].SLLI(com.Op1, com.Op2,com.Op3)
	case SRL:
		s.Proc[alg].SRL(com.Op1, com.Op2,uint64(com.Op3))
	case SRLI:
		s.Proc[alg].SRLI(com.Op1, com.Op2,com.Op3)
	case LI:
		s.Proc[alg].LI(com.Op1, com.Op3) //внимание второй операнд не используется
	case LDM:
		s.LDM(alg, com.Op1, com.Op2,uint64(com.Op3))
	case LDIN:
		s.LDIN(alg, com.Op1, com.Op2,uint64(com.Op3))
	case STM:
		s.STM(alg, com.Op1, com.Op2,uint64(com.Op3))
	case STOUT:
		s.STOUT(alg, com.Op1, com.Op2,uint64(com.Op3))
	case LDMX:
		s.LDMX(alg, com.Op1, com.Op2,uint64(com.Op3))
	case LDINX:
		s.LDINX(alg, com.Op1, com.Op2,uint64(com.Op3))
	case STMX:
		s.STMX(alg, com.Op1, com.Op2,uint64(com.Op3))
	case STOUTX:
		s.STOUTX(alg, com.Op1, com.Op2,uint64(com.Op3))
	case BEQ:
		s.BEQ(alg, com.Op1, com.Op2,com.Op3)
	case BGE:
		s.BGE(alg, com.Op1, com.Op2,com.Op3)
	case BLT:
		s.BLT(alg, com.Op1, com.Op2,com.Op3)
	case BNE:
		s.BNE(alg, com.Op1, com.Op2,com.Op3)
	case BLE:
		s.BLE(alg, com.Op1, com.Op2,com.Op3)
	case BGT:
		s.BGT(alg, com.Op1, com.Op2,com.Op3)
	case JMP:
		s.JMP(alg, com.Op3)//Внимание! В JMP адресс в 3 операнде (как впрочем и везде - просто здесь только один операнд, но он в кодоне 3-ий!!)
	case SEQ:
		s.Proc[alg].SEQ(com.Op1, com.Op2,uint64(com.Op3))
	case SGE:
		s.Proc[alg].SGE(com.Op1, com.Op2,uint64(com.Op3))
	case SLT:
		s.Proc[alg].SLT(com.Op1, com.Op2,uint64(com.Op3))
	case SNE:
		s.Proc[alg].SNE(com.Op1, com.Op2,uint64(com.Op3))
	case PUSH:
		s.Proc[alg].PUSH(com.Op1)
	case POP:
		s.Proc[alg].POP(com.Op1)
	}
}
