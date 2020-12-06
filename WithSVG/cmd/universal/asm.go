package universal

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func GetCodonsFromAsmString(s *string) (*[]Codon, error){
	var cods []Codon

	ss:=strings.Split(*s, "\n")
	i:=0
	for strings.TrimSpace(ss[i])!="asm:"{
		i++
		if i>=len(ss){
			return nil, errors.New("не обнаружена директива 'asm:'")
		}
	}
	i++
	//сейчас в i номер строки сразу за asm:
	for i<len(ss){
		sl:=strings.Fields(strings.TrimSpace(ss[i]))
		if len(sl)==0{ //у нас минимум комманда или метка
			i++
			continue
		}else if strings.HasPrefix(sl[0],"//"){//пропускаем комментарии
			i++
			continue
		}
		switch sl[0] {
		case  "NOP","nop":
			cods=append(cods,Codon{NOP,0,0,0})
		case "ADD","add":
			ms := regexp.MustCompile(`[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<4{//потому что x1,x2 , X3 вернет ["x1, x2,X3" "1" "2" "3"]
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v ADD ", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			op2,_:=strconv.ParseUint(ms[2],10,64)
			op3,_:=strconv.ParseInt(ms[3],10,64)
			cods=append(cods,Codon{ADD,op1,op2,op3})
		case "ADDI","addi":
			ms := regexp.MustCompile(`[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)\s*,\s*(-?[0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<4{
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v ADDI ", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			op2,_:=strconv.ParseUint(ms[2],10,64)
			op3,_:=strconv.ParseInt(ms[3],10,64)
			cods=append(cods,Codon{ADDI,op1,op2,op3})
		case "SUB","sub":
			ms := regexp.MustCompile(`[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<4{//потому что x1,x2 , X3 вернет [x1, x2,X3 1 2 3]
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v SUB ", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			op2,_:=strconv.ParseUint(ms[2],10,64)
			op3,_:=strconv.ParseInt(ms[3],10,64)
			cods=append(cods,Codon{SUB,op1,op2,op3})
		case "SUBI","subi":
			ms := regexp.MustCompile(`[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)\s*,\s*(-?[0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<4{
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v SUBI ", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			op2,_:=strconv.ParseUint(ms[2],10,64)
			op3,_:=strconv.ParseInt(ms[3],10,64)
			cods=append(cods,Codon{SUBI,op1,op2,op3})
		case "MUL","mul":
			ms := regexp.MustCompile(`[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<4{//потому что x1,x2 , X3 вернет [x1, x2,X3 1 2 3]
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v MUL ", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			op2,_:=strconv.ParseUint(ms[2],10,64)
			op3,_:=strconv.ParseInt(ms[3],10,64)
			cods=append(cods,Codon{MUL,op1,op2,op3})
		case "MULI","muli":
			ms := regexp.MustCompile(`[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)\s*,\s*(-?[0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<4{
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v MULI ", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			op2,_:=strconv.ParseUint(ms[2],10,64)
			op3,_:=strconv.ParseInt(ms[3],10,64)
			cods=append(cods,Codon{MULI,op1,op2,op3})
		case "DIV","div":
			ms := regexp.MustCompile(`[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<4{//потому что x1,x2 , X3 вернет [x1, x2,X3 1 2 3]
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v DIV ", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			op2,_:=strconv.ParseUint(ms[2],10,64)
			op3,_:=strconv.ParseInt(ms[3],10,64)
			cods=append(cods,Codon{DIV,op1,op2,op3})
		case "DIVI","divi":
			ms := regexp.MustCompile(`[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)\s*,\s*(-?[0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<4{
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v DIVI ", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			op2,_:=strconv.ParseUint(ms[2],10,64)
			op3,_:=strconv.ParseInt(ms[3],10,64)
			cods=append(cods,Codon{DIVI,op1,op2,op3})
		case "REM","rem":
			ms := regexp.MustCompile(`[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<4{//потому что x1,x2 , X3 вернет [x1, x2,X3 1 2 3]
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v REM ", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			op2,_:=strconv.ParseUint(ms[2],10,64)
			op3,_:=strconv.ParseInt(ms[3],10,64)
			cods=append(cods,Codon{REM,op1,op2,op3})
		case "REMI","remi":
			ms := regexp.MustCompile(`[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)\s*,\s*(-?[0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<4{
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v REMI ", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			op2,_:=strconv.ParseUint(ms[2],10,64)
			op3,_:=strconv.ParseInt(ms[3],10,64)
			cods=append(cods,Codon{REMI,op1,op2,op3})
		case "AND","and":
			ms := regexp.MustCompile(`[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<4{//потому что x1,x2 , X3 вернет [x1, x2,X3 1 2 3]
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v AND ", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			op2,_:=strconv.ParseUint(ms[2],10,64)
			op3,_:=strconv.ParseInt(ms[3],10,64)
			cods=append(cods,Codon{AND,op1,op2,op3})
		case "ANDI","andi":
			ms := regexp.MustCompile(`[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)\s*,\s*(-?[0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<4{
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v ANDI ", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			op2,_:=strconv.ParseUint(ms[2],10,64)
			op3,_:=strconv.ParseInt(ms[3],10,64)
			cods=append(cods,Codon{ANDI,op1,op2,op3})
		case "OR","or":
			ms := regexp.MustCompile(`[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<4{//потому что x1,x2 , X3 вернет [x1, x2,X3 1 2 3]
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v OR ", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			op2,_:=strconv.ParseUint(ms[2],10,64)
			op3,_:=strconv.ParseInt(ms[3],10,64)
			cods=append(cods,Codon{OR,op1,op2,op3})
		case "ORI","ori":
			ms := regexp.MustCompile(`[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)\s*,\s*(-?[0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<4{
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v ORI ", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			op2,_:=strconv.ParseUint(ms[2],10,64)
			op3,_:=strconv.ParseInt(ms[3],10,64)
			cods=append(cods,Codon{ORI,op1,op2,op3})
		case "XOR","xor":
			ms := regexp.MustCompile(`[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<4{//потому что x1,x2 , X3 вернет [x1, x2,X3 1 2 3]
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v XOR ", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			op2,_:=strconv.ParseUint(ms[2],10,64)
			op3,_:=strconv.ParseInt(ms[3],10,64)
			cods=append(cods,Codon{XOR,op1,op2,op3})
		case "XORI","xori":
			ms := regexp.MustCompile(`[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)\s*,\s*(-?[0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<4{
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v XORI ", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			op2,_:=strconv.ParseUint(ms[2],10,64)
			op3,_:=strconv.ParseInt(ms[3],10,64)
			cods=append(cods,Codon{XORI,op1,op2,op3})
		case "SLL","sll":
			ms := regexp.MustCompile(`[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<4{//потому что x1,x2 , X3 вернет [x1, x2,X3 1 2 3]
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v SLL ", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			op2,_:=strconv.ParseUint(ms[2],10,64)
			op3,_:=strconv.ParseInt(ms[3],10,64)
			cods=append(cods,Codon{SLL,op1,op2,op3})
		case "SLLI","slli":
			ms := regexp.MustCompile(`[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)\s*,\s*([0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<4{
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v SLLI ", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			op2,_:=strconv.ParseUint(ms[2],10,64)
			op3,_:=strconv.ParseInt(ms[3],10,64)
			cods=append(cods,Codon{SLLI,op1,op2,op3})
		case "SRL","srl":
			ms := regexp.MustCompile(`[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<4{//потому что x1,x2 , X3 вернет [x1, x2,X3 1 2 3]
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v SRL ", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			op2,_:=strconv.ParseUint(ms[2],10,64)
			op3,_:=strconv.ParseInt(ms[3],10,64)
			cods=append(cods,Codon{SRL,op1,op2,op3})
		case "SRLI","srli":
			ms := regexp.MustCompile(`[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)\s*,\s*([0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<4{
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v SRLI ", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			op2,_:=strconv.ParseUint(ms[2],10,64)
			op3,_:=strconv.ParseInt(ms[3],10,64)
			cods=append(cods,Codon{SRLI,op1,op2,op3})
		case "LDM","ldm":
			ms := regexp.MustCompile(`[x|X]([0-9]+)\s*,\s*([0-9]+)\s*,\s*([0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<4{
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v Должно быть: LDM xn, m, addr", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			op2,_:=strconv.ParseUint(ms[2],10,64)
			op3,_:=strconv.ParseInt(ms[3],10,64)
			cods=append(cods,Codon{LDM,op1,op2,op3})
		case "LDIN","ldin":
			ms := regexp.MustCompile(`[x|X]([0-9]+)\s*,\s*([0-9]+)\s*,\s*([0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<4{
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v Должно быть: LDIN xN, i, addr", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			op2,_:=strconv.ParseUint(ms[2],10,64)
			op3,_:=strconv.ParseInt(ms[3],10,64)
			cods=append(cods,Codon{LDIN,op1,op2,op3})
		case "STM","stm":
			ms := regexp.MustCompile(`([0-9]+)\s*,\s*([0-9]+)\s*,\s*[x|X]([0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<4{
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v STM", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			op2,_:=strconv.ParseUint(ms[2],10,64)
			op3,_:=strconv.ParseInt(ms[3],10,64)
			cods=append(cods,Codon{STM,op1,op2,op3})
		case "STOUT","stout":
			ms := regexp.MustCompile(`([0-9]+)\s*,\s*([0-9]+)\s*,\s*[x|X]([0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<4{
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v STOUT", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			op2,_:=strconv.ParseUint(ms[2],10,64)
			op3,_:=strconv.ParseInt(ms[3],10,64)
			cods=append(cods,Codon{STOUT,op1,op2,op3})
		case "BEQ","beq"://у нас ветвление без меток, как в дебаггере, потому что геномика может столько ветвлений наставить - заколупаешься из них метки генерить
		//так что только отностельные адреса
			ms := regexp.MustCompile(`[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)\s*,\s*(-?[0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<4{
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v BEQ ", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			op2,_:=strconv.ParseUint(ms[2],10,64)
			op3,_:=strconv.ParseInt(ms[3],10,64)
			cods=append(cods,Codon{BEQ,op1,op2,op3})
		case "BGE","bge":
			ms := regexp.MustCompile(`[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)\s*,\s*(-?[0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<4{
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v BGE ", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			op2,_:=strconv.ParseUint(ms[2],10,64)
			op3,_:=strconv.ParseInt(ms[3],10,64)
			cods=append(cods,Codon{BGE,op1,op2,op3})
		case "BLT","blt":
			ms := regexp.MustCompile(`[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)\s*,\s*(-?[0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<4{
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v BLT ", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			op2,_:=strconv.ParseUint(ms[2],10,64)
			op3,_:=strconv.ParseInt(ms[3],10,64)
			cods=append(cods,Codon{BLT,op1,op2,op3})
		case "BNE","bne":
			ms := regexp.MustCompile(`[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)\s*,\s*(-?[0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<4{
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v BNE ", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			op2,_:=strconv.ParseUint(ms[2],10,64)
			op3,_:=strconv.ParseInt(ms[3],10,64)
			cods=append(cods,Codon{BNE,op1,op2,op3})
		case "JMP","jmp":
			ms := regexp.MustCompile(`(-?[0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<2{
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v JMP ", i))
			}
			op3,_:=strconv.ParseInt(ms[1],10,64)
			cods=append(cods,Codon{JMP,0,0,op3})//Внимание! адресс в третьем операнде в кодоне всегда!
		case "SEQ","seq":
			ms := regexp.MustCompile(`[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<4{//потому что x1,x2 , X3 вернет [x1, x2,X3 1 2 3]
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v SEQ ", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			op2,_:=strconv.ParseUint(ms[2],10,64)
			op3,_:=strconv.ParseInt(ms[3],10,64)
			cods=append(cods,Codon{SEQ,op1,op2,op3})
		case "SGE","sge":
			ms := regexp.MustCompile(`[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<4{//потому что x1,x2 , X3 вернет [x1, x2,X3 1 2 3]
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v SGE ", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			op2,_:=strconv.ParseUint(ms[2],10,64)
			op3,_:=strconv.ParseInt(ms[3],10,64)
			cods=append(cods,Codon{SGE,op1,op2,op3})
		case "SLT","slt":
			ms := regexp.MustCompile(`[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<4{//потому что x1,x2 , X3 вернет [x1, x2,X3 1 2 3]
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v SLT ", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			op2,_:=strconv.ParseUint(ms[2],10,64)
			op3,_:=strconv.ParseInt(ms[3],10,64)
			cods=append(cods,Codon{SLT,op1,op2,op3})
		case "SNE","sne":
			ms := regexp.MustCompile(`[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)\s*,\s*[x|X]([0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<4{//потому что x1,x2 , X3 вернет [x1, x2,X3 1 2 3]
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v SNE ", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			op2,_:=strconv.ParseUint(ms[2],10,64)
			op3,_:=strconv.ParseInt(ms[3],10,64)
			cods=append(cods,Codon{SNE,op1,op2,op3})
		case "PUSH","push":
			ms := regexp.MustCompile(`[x|X]([0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<2{
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v PUSH ", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			cods=append(cods,Codon{PUSH,op1,0,0})
		case "POP","pop":
			ms := regexp.MustCompile(`[x|X]([0-9]+)`).FindStringSubmatch(ss[i])
			if len(ms)<2{
				return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v POP ", i))
			}
			op1,_:=strconv.ParseUint(ms[1],10,64)
			cods=append(cods,Codon{POP,op1,0,0})
		default:
			return nil, errors.New( fmt.Sprintf( "Ошибка: стр. %v неизвестная команда ", i))
		}
		i++
	}
	return &cods,nil
}

func GetReadableFromCodons(cods []Codon) *[2]string{
	ret:=[2]string{"codons: ",
				   "/*This is assembler only for Solution RISC architecture*/\nasm:\n"}
	for i:=0;i<len(cods);i++{
		switch cods[i].Code%42{
		case NOP,NOP1,NOP2,NOP3,NOP4,NOP5,NOP6,NOP7:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code,cods[i].Op1,cods[i].Op2,cods[i].Op3)
			ret[1]+="\tNOP\n"
		case ADD:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code, cods[i].Op1%32, cods[i].Op2%32,uint64(cods[i].Op3%32))
			ret[1]+=fmt.Sprintf("\tADD\tx%v, x%v, x%v\n", cods[i].Op1%32, cods[i].Op2%32,uint64(cods[i].Op3%32))
		case ADDI:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code, cods[i].Op1%32, cods[i].Op2%32,cods[i].Op3)
			ret[1]+=fmt.Sprintf("\tADDI\tx%v, x%v, %v\n", cods[i].Op1%32, cods[i].Op2%32,cods[i].Op3)
		case SUB:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code, cods[i].Op1%32, cods[i].Op2%32,uint64(cods[i].Op3%32))
			ret[1]+=fmt.Sprintf("\tSUB\tx%v, x%v, x%v\n", cods[i].Op1%32, cods[i].Op2%32,uint64(cods[i].Op3%32))
		case SUBI:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code, cods[i].Op1%32, cods[i].Op2%32,cods[i].Op3)
			ret[1]+=fmt.Sprintf("\tSUBI\tx%v, x%v, %v\n", cods[i].Op1%32, cods[i].Op2%32,cods[i].Op3)
		case MUL:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code, cods[i].Op1%32, cods[i].Op2%32,uint64(cods[i].Op3%32))
			ret[1]+=fmt.Sprintf("\tMUL\tx%v, x%v, x%v\n", cods[i].Op1%32, cods[i].Op2%32,uint64(cods[i].Op3%32))
		case MULI:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code, cods[i].Op1%32, cods[i].Op2%32,cods[i].Op3)
			ret[1]+=fmt.Sprintf("\tMULI\tx%v, x%v, %v\n", cods[i].Op1%32, cods[i].Op2%32,cods[i].Op3)
		case DIV:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code, cods[i].Op1%32, cods[i].Op2%32,uint64(cods[i].Op3%32))
			ret[1]+=fmt.Sprintf("\tDIV\tx%v, x%v, x%v\n", cods[i].Op1%32, cods[i].Op2%32,uint64(cods[i].Op3%32))
		case DIVI:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code, cods[i].Op1%32, cods[i].Op2%32,cods[i].Op3)
			ret[1]+=fmt.Sprintf("\tDIVI\tx%v, x%v, %v\n", cods[i].Op1%32, cods[i].Op2%32,cods[i].Op3)
		case REM:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code, cods[i].Op1%32, cods[i].Op2%32,uint64(cods[i].Op3%32))
			ret[1]+=fmt.Sprintf("\tREM\tx%v, x%v, x%v\n", cods[i].Op1%32, cods[i].Op2%32,uint64(cods[i].Op3%32))
		case REMI:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code, cods[i].Op1%32, cods[i].Op2%32,cods[i].Op3)
			ret[1]+=fmt.Sprintf("\tREMI\tx%v, x%v, %v\n", cods[i].Op1%32, cods[i].Op2%32,cods[i].Op3)
		case AND:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code, cods[i].Op1%32, cods[i].Op2%32,uint64(cods[i].Op3%32))
			ret[1]+=fmt.Sprintf("\tAND\tx%v, x%v, x%v\n", cods[i].Op1%32, cods[i].Op2%32,uint64(cods[i].Op3%32))
		case ANDI:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code, cods[i].Op1%32, cods[i].Op2%32,cods[i].Op3)
			ret[1]+=fmt.Sprintf("\tANDI\tx%v, x%v, %v\n", cods[i].Op1%32, cods[i].Op2%32,cods[i].Op3)
		case OR:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code, cods[i].Op1%32, cods[i].Op2%32,uint64(cods[i].Op3%32))
			ret[1]+=fmt.Sprintf("\tOR\tx%v, x%v, x%v\n", cods[i].Op1%32, cods[i].Op2%32,uint64(cods[i].Op3%32))
		case ORI:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code, cods[i].Op1%32, cods[i].Op2%32,cods[i].Op3)
			ret[1]+=fmt.Sprintf("\tORI\tx%v, x%v, %v\n", cods[i].Op1%32, cods[i].Op2%32,cods[i].Op3)
		case XOR:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code, cods[i].Op1%32, cods[i].Op2%32,uint64(cods[i].Op3%32))
			ret[1]+=fmt.Sprintf("\tXOR\tx%v, x%v, x%v\n", cods[i].Op1%32, cods[i].Op2%32,uint64(cods[i].Op3%32))
		case XORI:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code, cods[i].Op1%32, cods[i].Op2%32,cods[i].Op3)
			ret[1]+=fmt.Sprintf("\tXORI\tx%v, x%v, %v\n", cods[i].Op1%32, cods[i].Op2%32,cods[i].Op3)
		case SLL:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code, cods[i].Op1%32, cods[i].Op2%32,uint64(cods[i].Op3%32))
			ret[1]+=fmt.Sprintf("\tSLL\tx%v, x%v, x%v\n", cods[i].Op1%32, cods[i].Op2%32,uint64(cods[i].Op3%32))
		case SLLI:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code, cods[i].Op1%32, cods[i].Op2%32,uint64(cods[i].Op3))
			ret[1]+=fmt.Sprintf("\tSLLI\tx%v, x%v, %v\n", cods[i].Op1%32, cods[i].Op2%32,uint64(cods[i].Op3))
		case SRL:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code, cods[i].Op1%32, cods[i].Op2%32,uint64(cods[i].Op3%32))
			ret[1]+=fmt.Sprintf("\tSRL\tx%v, x%v, x%v\n", cods[i].Op1%32, cods[i].Op2%32,uint64(cods[i].Op3%32))
		case SRLI:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code, cods[i].Op1%32, cods[i].Op2%32,uint64(cods[i].Op3))
			ret[1]+=fmt.Sprintf("\tSRLI\tx%v, x%v, %v\n", cods[i].Op1%32, cods[i].Op2%32,uint64(cods[i].Op3))
		case LDM:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code, cods[i].Op1%32, cods[i].Op2,uint64(cods[i].Op3))
			ret[1]+=fmt.Sprintf("\tLDM\tx%v, %v, %v\n", cods[i].Op1%32, cods[i].Op2,uint64(cods[i].Op3))
		case LDIN:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code, cods[i].Op1%32, cods[i].Op2,uint64(cods[i].Op3))
			ret[1]+=fmt.Sprintf("\tLDIN\tx%v, %v, %v\n", cods[i].Op1%32, cods[i].Op2,uint64(cods[i].Op3))
		case STM:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code, cods[i].Op1, cods[i].Op2,uint64(cods[i].Op3%32))
			ret[1]+=fmt.Sprintf("\tSTM\t%v, %v, x%v\n", cods[i].Op1, cods[i].Op2,uint64(cods[i].Op3%32))
		case STOUT:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code, cods[i].Op1, cods[i].Op2,uint64(cods[i].Op3%32))
			ret[1]+=fmt.Sprintf("\tSTOUT\t%v, %v, x%v\n", cods[i].Op1, cods[i].Op2,uint64(cods[i].Op3%32))
		case BEQ:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code, cods[i].Op1%32, cods[i].Op2%32,cods[i].Op3)
			ret[1]+=fmt.Sprintf("\tBEQ\tx%v, x%v, %v\n", cods[i].Op1%32, cods[i].Op2%32,cods[i].Op3)
		case BGE:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code, cods[i].Op1%32, cods[i].Op2%32,cods[i].Op3)
			ret[1]+=fmt.Sprintf("\tBGE\tx%v, x%v, %v\n", cods[i].Op1%32, cods[i].Op2%32,cods[i].Op3)
		case BLT:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code, cods[i].Op1%32, cods[i].Op2%32,cods[i].Op3)
			ret[1]+=fmt.Sprintf("\tBLT\tx%v, x%v, %v\n", cods[i].Op1%32, cods[i].Op2%32,cods[i].Op3)
		case BNE:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code, cods[i].Op1%32, cods[i].Op2%32,cods[i].Op3)
			ret[1]+=fmt.Sprintf("\tBNE\tx%v, x%v, %v\n", cods[i].Op1%32, cods[i].Op2%32,cods[i].Op3)
		case JMP:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code, 0, 0,cods[i].Op3)
			ret[1]+=fmt.Sprintf("\tJMP\t%v\n",cods[i].Op3)
		case SEQ:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code, cods[i].Op1%32, cods[i].Op2%32,uint64(cods[i].Op3%32))
			ret[1]+=fmt.Sprintf("\tSEQ\tx%v, x%v, x%v\n", cods[i].Op1%32, cods[i].Op2%32,uint64(cods[i].Op3%32))
		case SGE:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code, cods[i].Op1%32, cods[i].Op2%32,uint64(cods[i].Op3%32))
			ret[1]+=fmt.Sprintf("\tSGE\tx%v, x%v, x%v\n", cods[i].Op1%32, cods[i].Op2%32,uint64(cods[i].Op3%32))
		case SLT:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code, cods[i].Op1%32, cods[i].Op2%32,uint64(cods[i].Op3%32))
			ret[1]+=fmt.Sprintf("\tSLT\tx%v, x%v, x%v\n", cods[i].Op1%32, cods[i].Op2%32,uint64(cods[i].Op3%32))
		case SNE:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code, cods[i].Op1%32, cods[i].Op2%32,uint64(cods[i].Op3%32))
			ret[1]+=fmt.Sprintf("\tSNE\tx%v, x%v, x%v\n", cods[i].Op1%32, cods[i].Op2%32,uint64(cods[i].Op3%32))
		case PUSH:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code,cods[i].Op1,0,0)
			ret[1]+=fmt.Sprintf("\tPUSH\tx%v\n",cods[i].Op1)
		case POP:
			ret[0]+=fmt.Sprintf("%v %v %v %v; ",cods[i].Code,cods[i].Op1,0,0)
			ret[1]+=fmt.Sprintf("\tPOP\tx%v\n",cods[i].Op1)
		}

	}
	return &ret
}