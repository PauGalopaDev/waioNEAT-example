package main

import (
	waio "github.com/PauGalopaDev/waioNEAT"
)

const (
	Up = iota
	Right
	Down
	Left
)

type Robot struct {
	Pos     Pair
	Dir     int
	Match   *Match
	Genome  *waio.Genome
	Brain   waio.Network
	Inputs  map[string]*float64
	Outputs map[string]*float64
	Energy  int
}

func MakeRobot(g *waio.Genome) *Robot {
	result := &Robot{}

	result.Genome = g
	result.Brain = *waio.MakeNetwork(g)

	result.Inputs = make(map[string]*float64, len(result.Brain.Input))
	result.Outputs = make(map[string]*float64, len(result.Brain.Output))

	for param, neuron := range result.Brain.Input {
		result.Inputs[param] = &neuron.Value
	}

	for param, neuron := range result.Brain.Output {
		result.Outputs[param] = &neuron.Value
	}

	result.Energy = 100
	return result
}

func (r *Robot) Look() {
	dir := Pair{0, 0}
	switch r.Dir {
	case Up:
		dir.i -= 1
	case Down:
		dir.i += 1
	case Left:
		dir.j -= 1
	case Right:
		dir.j += 1
	}

	tp := r.Pos
	tp.Add(dir)
	for r.Match.PosOk(tp) {
		if r.Match.Grid[tp.i][tp.j].Energy {
			if _, b := r.Inputs["look"]; b {
				*r.Inputs["look"] = 1.0
			} else {
				return
			}
		}
		tp.Add(dir)
	}

	if _, b := r.Inputs["look"]; b {
		*r.Inputs["look"] = 0.0
	}
}

func (r *Robot) EvalEnergy() {
	if _, b := r.Inputs["energy"]; b {
		*r.Inputs["energy"] = float64(r.Energy) / 100
	}
}

func (r *Robot) Move() {
	if v, b := r.Outputs["move"]; !b && *v < 1 {
		return
	}

	tp := Pair{r.Pos.i, r.Pos.j}
	switch r.Dir {
	case Up:
		tp.i -= 1
	case Down:
		tp.i += 1
	case Left:
		tp.j -= 1
	case Right:
		tp.j += 1
	}

	if r.Match.PosOk(tp) && r.Match.Grid[tp.i][tp.j].Robot {
		return
	}
}

func (r *Robot) Rotate() {
	// Check if the Outputs exists and the values are at least 1
	rl, rr := false, false
	if v, b := r.Outputs["rotleft"]; b && *v >= 1 {
		rl = true
	}

	if v, b := r.Outputs["rotright"]; b && *v >= 1 {
		rr = true
	}

	// If both rotate outputs are active or off do not rotate, else...
	// Rotate left
	if rl && !rr {
		if r.Dir == 0 {
			r.Dir = 3
		} else {
			r.Dir -= 1
		}
	}

	// Rotate right
	if !rl && rr {
		if r.Dir == 3 {
			r.Dir = 0
		} else {
			r.Dir += 1
		}
	}
}

func (r *Robot) Update() {
	// Set Inputs
	*r.Inputs["bias"] = 1.0
	r.Look()
	r.EvalEnergy()

	// Feed Brain
	r.Brain.Feed()

	// Get Outputs
	r.Rotate()
	r.Move()
	r.Energy -= 2

	if r.Match.GetCell(r.Pos).Energy {
		if r.Energy <= 80 {
			r.Energy += 20
		} else {
			r.Energy = 100
		}
	}

}
