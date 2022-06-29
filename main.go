package main

import (
	"fmt"
	"os"
	"strconv"

	waio "github.com/PauGalopaDev/waioNEAT"
)

func main() {
	nbMatchs := 5
	nMatch := 0
	if len(os.Args) > 1 {
		nbMatchs, _ = strconv.Atoi(os.Args[1])
	}
	waio.Init()

	// Robot Params
	ins := map[string]float64{"look": 1, "energy": 1, "bias": 1}
	outs := map[string]float64{"move": 1, "rotleft": 1, "rotright": 1}

	// Make Genome
	firstGenome := waio.MakeGenome(ins, outs)

	// Make Match
	match := MakeMatch(20, 0.05, firstGenome, 5)

	// Start Matches
	for nMatch <= nbMatchs {
		round := 0
		for len(match.Robots) > 0 {
			match.Update()
			fmt.Printf("Match %d/%d:\tRound:%d\n", nMatch, nbMatchs, round)
			fmt.Printf("%v\n", match)
			round++
		}
		match = MakeMatch(20, 0.05, match.Genomes[len(match.Genomes)-1], 5)
		nMatch += 1
	}
}

/*
TODO:
Add argument handling
*/
