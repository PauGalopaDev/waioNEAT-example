package main

import (
	"math/rand"

	waio "github.com/PauGalopaDev/waioNEAT"
)

func SliceRemove[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

type Pair struct {
	i int
	j int
}

func (x *Pair) Add(y Pair) {
	x.i += y.i
	x.j += y.j
}

type Cell struct {
	Robot  bool
	Energy bool
}

type Match struct {
	Grid    [][]Cell
	Energy  int
	Robots  []*Robot
	Genomes []*waio.Genome
}

func MakeMatch(size int, chance float64, g *waio.Genome, nb int) *Match {
	result := &Match{Grid: make([][]Cell, size)}
	result.Genomes = make([]*waio.Genome, nb)
	result.Robots = make([]*Robot, 0, nb)
	energy := false

	// Create Grid Cells and spawn energy
	for i := range result.Grid {
		result.Grid[i] = make([]Cell, size)
		for j := range result.Grid[i] {
			if rand.Float64() < chance {
				energy = true
			}
			result.Grid[i][j] = Cell{
				Energy: energy,
			}
		}
	}

	// Spwan the robots
	// Make the first robot, with the given genome
	result.Robots = append(result.Robots, MakeRobot(g))
	currentRobot := result.Robots[len(result.Robots)-1]
	currentRobot.Match = result

	// Set its position
	currentRobot.Pos.i = rand.Intn(size)
	currentRobot.Pos.j = rand.Intn(size)

	// Set the grid cell robot marker as true
	result.GetCell(currentRobot.Pos).Robot = true

	for i := 0; i < nb; i += 1 {
		genome := &waio.Genome{}
		*genome = *g
		if c := rand.Float64(); c <= 0.5 {
			genome.MutateAddNode()
		} else {
			genome.MutateAddSynapse()
		}
		// Make the next robot, with the given genome
		result.Robots = append(result.Robots, MakeRobot(genome))
		currentRobot := result.Robots[len(result.Robots)-1]
		currentRobot.Match = result

		// Set its position
		currentRobot.Pos.i = rand.Intn(size)
		currentRobot.Pos.j = rand.Intn(size)

		for result.GetCell(currentRobot.Pos).Robot {
			currentRobot.Pos.i = rand.Intn(size)
			currentRobot.Pos.j = rand.Intn(size)
		}

		// Set the grid cell robot marker as true
		result.GetCell(currentRobot.Pos).Robot = true
	}

	return result
}

func (m *Match) PosOk(p Pair) bool {
	return p.i >= 0 && p.j >= 0 && p.i < m.Rows() && p.j < m.Cols()
}

func (m *Match) Rows() int {
	return len(m.Grid)
}

func (m *Match) Cols() int {
	return len(m.Grid[0])
}

func (m *Match) GetCell(p Pair) *Cell {
	return &m.Grid[p.i][p.j]
}

func (m *Match) Update() {
	// For each robot

	for i := len(m.Robots) - 1; i >= 0; i-- {
		r := m.Robots[i]
		lastPos := r.Pos
		r.Update()

		// If there is energy on the current position, remove form cells
		m.GetCell(r.Pos).Energy = false

		// If the energy depleted
		if r.Energy > 0 {
			// Update cell robot presence
			m.GetCell(lastPos).Robot = false
			m.GetCell(r.Pos).Robot = true
		} else {
			// remove the robot form the cells
			m.GetCell(lastPos).Robot = false

			// Save its genome
			m.Genomes = append(m.Genomes, r.Genome)

			// Remove robot from the list
			m.Robots = SliceRemove(m.Robots, i)
		}
		i--
	}

	if len(m.Robots) == 1 {
		m.Genomes = append(m.Genomes, m.Robots[0].Genome)
		m.Robots = SliceRemove(m.Robots, 0)
	}
}
