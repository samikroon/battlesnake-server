package battlesnake

import (
	"errors"
	"math/rand"
	"time"
)

const (
	SOLVER_RANDOM = "random"
	SOLVER_SAFE   = "safe"
)

var (
	ErrSolverNotFound = errors.New("solver not found")
)

type Solver interface {
	Move(state State) string
}

// NewSolver creates a new solver based on the given type.
// If the given type is not implemented in this server, an error is returned.
func NewSolver(solverType string) (Solver, error) {
	switch solverType {
	case SOLVER_RANDOM:
		src := rand.NewSource(time.Now().UnixNano())
		return randSolver{
			random: rand.New(src),
		}, nil
	case SOLVER_SAFE:
		return safeMoveSolver{}, nil
	default:
		return nil, ErrSolverNotFound
	}
}

type randSolver struct {
	random *rand.Rand
}

// Move returns a random move.
func (r randSolver) Move(state State) string {
	return Moves[r.random.Intn(4)]
}

type safeMoveSolver struct{}

// Move returns the first safe move it can find, if no safe move is found, up is returned.
func (s safeMoveSolver) Move(state State) string {
	for moveName, move := range MovesMap {
		if !isCollision(state, move) {
			return moveName
		}
	}
	return MOVE_UP // no safe moves
}

// isCollision checks if the a move results in a collision with a wall or the snake itself.
func isCollision(state State, move Point) bool {
	headX := state.You.Head.X + move.X
	headY := state.You.Head.Y + move.Y

	if headX < 0 || headY < 0 || headX >= state.Board.Width || headY >= state.Board.Height { // wall collision
		return true
	}

	for i := 1; i < len(state.You.Body); i++ { // self collision
		if headX == state.You.Body[i].X && headY == state.You.Body[i].Y {
			return true
		}
	}

	return false
}
