package main

import (
	"WithSVG/cmd/universal"
	"fmt"
)

func main() {
	sol, err:=universal.NewSolution("c:/ALOSKA/work/solutions/ololo.json")
	if err!=nil{
		fmt.Println(err)
	}
	fmt.Println(sol.Gen.Codons)
}
