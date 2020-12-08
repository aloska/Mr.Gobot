package universal

import "sync"

func (s *Solution) Run(){
	var wg sync.WaitGroup
	for i:=0;i<len(s.Chrom);i++ {
		if s.IsAsync{
			wg.Add(1)
			go func(chr int){
				defer wg.Done()
				for int(s.Proc[chr].PC) < len(s.Chrom[chr].Codons) {
					s.Step(chr)
				}
			}(i)
		}else {
			for int(s.Proc[i].PC) < len(s.Chrom[i].Codons) {
				s.Step(i)
			}
		}
	}
	if s.IsAsync{
		wg.Wait()
	}
}

func (s *Solution) Step(chr int){
	//42 - кол-во комманд в нашем RISC-процессоре
	com:=s.Chrom[chr].Codons[s.Proc[chr].PC]
	switch com.Code%42 {
	case NOP: //а что уж тут поделаешь?
		s.Proc[chr].PC++	//хех!
	case ADD:
		s.Proc[chr].ADD(com.Op1, com.Op2,uint64(com.Op3))
	case ADDI:
		s.Proc[chr].ADDI(com.Op1, com.Op2,com.Op3)
	case SUB:
		s.Proc[chr].SUB(com.Op1, com.Op2,uint64(com.Op3))
	case SUBI:
		s.Proc[chr].SUBI(com.Op1, com.Op2,com.Op3)
	case MUL:
		s.Proc[chr].MUL(com.Op1, com.Op2,uint64(com.Op3))
	case MULI:
		s.Proc[chr].MULI(com.Op1, com.Op2,com.Op3)
	case DIV:
		s.Proc[chr].DIV(com.Op1, com.Op2,uint64(com.Op3))
	case DIVI:
		s.Proc[chr].DIVI(com.Op1, com.Op2,com.Op3)
	case REM:
		s.Proc[chr].REM(com.Op1, com.Op2,uint64(com.Op3))
	case REMI:
		s.Proc[chr].REMI(com.Op1, com.Op2,com.Op3)
	case AND:
		s.Proc[chr].AND(com.Op1, com.Op2,uint64(com.Op3))
	case ANDI:
		s.Proc[chr].ANDI(com.Op1, com.Op2,com.Op3)
	case OR:
		s.Proc[chr].OR(com.Op1, com.Op2,uint64(com.Op3))
	case ORI:
		s.Proc[chr].ORI(com.Op1, com.Op2,com.Op3)
	case XOR:
		s.Proc[chr].XOR(com.Op1, com.Op2,uint64(com.Op3))
	case XORI:
		s.Proc[chr].XORI(com.Op1, com.Op2,com.Op3)
	case SLL:
		s.Proc[chr].SLL(com.Op1, com.Op2,uint64(com.Op3))
	case SLLI:
		s.Proc[chr].SLLI(com.Op1, com.Op2,com.Op3)
	case SRL:
		s.Proc[chr].SRL(com.Op1, com.Op2,uint64(com.Op3))
	case SRLI:
		s.Proc[chr].SRLI(com.Op1, com.Op2,com.Op3)
	case LI:
		s.Proc[chr].LI(com.Op1, com.Op3)
	case LDM:
		s.LDM(chr, com.Op1, com.Op2,uint64(com.Op3))
	case LDIN:
		s.LDIN(chr, com.Op1, com.Op2,uint64(com.Op3))
	case STM:
		s.STM(chr, com.Op1, com.Op2,uint64(com.Op3))
	case STOUT:
		s.STOUT(chr, com.Op1, com.Op2,uint64(com.Op3))
	case LDMX:
		s.LDMX(chr, com.Op1, com.Op2,uint64(com.Op3))
	case LDINX:
		s.LDINX(chr, com.Op1, com.Op2,uint64(com.Op3))
	case STMX:
		s.STMX(chr, com.Op1, com.Op2,uint64(com.Op3))
	case STOUTX:
		s.STOUTX(chr, com.Op1, com.Op2,uint64(com.Op3))
	case BEQ:
		s.BEQ(chr, com.Op1, com.Op2,com.Op3)
	case BGE:
		s.BGE(chr, com.Op1, com.Op2,com.Op3)
	case BLT:
		s.BLT(chr, com.Op1, com.Op2,com.Op3)
	case BNE:
		s.BNE(chr, com.Op1, com.Op2,com.Op3)
	case BLE:
		s.BLE(chr, com.Op1, com.Op2,com.Op3)
	case BGT:
		s.BGT(chr, com.Op1, com.Op2,com.Op3)
	case JMP:
		s.JMP(chr, com.Op3)//Внимание! В JMP адресс в 3 операнде (как впрочем и везде - просто здесь только один операнд, но он в кодоне 3-ий!!)
	case SEQ:
		s.Proc[chr].SEQ(com.Op1, com.Op2,uint64(com.Op3))
	case SGE:
		s.Proc[chr].SGE(com.Op1, com.Op2,uint64(com.Op3))
	case SLT:
		s.Proc[chr].SLT(com.Op1, com.Op2,uint64(com.Op3))
	case SNE:
		s.Proc[chr].SNE(com.Op1, com.Op2,uint64(com.Op3))
	case PUSH:
		s.Proc[chr].PUSH(com.Op1)
	case POP:
		s.Proc[chr].POP(com.Op1)
	}
}
