package main

import (
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
	ins := map[string]float64{"look": 1, "energy": 1}
	outs := map[string]float64{"move": 1, "rotleft": 1, "rotright": 1}

	// Make Genome
	firstGenome := waio.MakeGenome(ins, outs)

	// Make Match
	match := MakeMatch(50, 0.05, firstGenome, 5)

	// Start Matches
	for nMatch <= nbMatchs {
		for len(match.Robots) > 0 {
			match.Update()
		}
		match = MakeMatch(50, 0.05, match.Genomes[len(match.Genomes)-1], 5)
		nMatch += 1
	}
}
