package battlesnake_test

import (
	"slices"
	"testing"

	"github.com/samikroon/battlesnake-server/battlesnake"
	"github.com/stretchr/testify/assert"
)

func Test_NewSolver(t *testing.T) {
	_, err := battlesnake.NewSolver(battlesnake.SOLVER_RANDOM)
	assert.NoError(t, err)

	_, err = battlesnake.NewSolver(battlesnake.SOLVER_SAFE)
	assert.NoError(t, err)

	_, err = battlesnake.NewSolver("not exists")
	assert.ErrorIs(t, err, battlesnake.ErrSolverNotFound)
}

func Test_RandSolver_Move(t *testing.T) {
	solver, err := battlesnake.NewSolver(battlesnake.SOLVER_RANDOM)
	assert.NoError(t, err)

	move := solver.Move(battlesnake.State{})

	assert.True(t, slices.Contains(battlesnake.Moves, move))
}

func Test_SafeSolver_Move_Possible(t *testing.T) {
	solver, err := battlesnake.NewSolver(battlesnake.SOLVER_SAFE)
	assert.NoError(t, err)

	move := solver.Move(battlesnake.State{
		Board: battlesnake.Board{
			Height: 4,
			Width:  4,
		},
		You: battlesnake.Snake{
			Body: []battlesnake.Point{
				{1, 3}, {1, 2}, {1, 1}, {1, 0}, {0, 0},
			},
			Head: battlesnake.Point{1, 3},
		},
	})

	assert.True(t, slices.Contains(battlesnake.Moves, move))
}

func Test_SafeSolver_Move_Impossible(t *testing.T) {
	solver, err := battlesnake.NewSolver(battlesnake.SOLVER_SAFE)
	assert.NoError(t, err)

	move := solver.Move(battlesnake.State{
		Board: battlesnake.Board{
			Height: 2,
			Width:  2,
		},
		You: battlesnake.Snake{
			Body: []battlesnake.Point{
				{1, 1}, {1, 0}, {0, 0}, {0, 1},
			},
			Head: battlesnake.Point{1, 1},
		},
	})

	assert.Equal(t, battlesnake.MOVE_UP, move)
}
