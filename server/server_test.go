package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/samikroon/battlesnake-server/battlesnake"
	"github.com/stretchr/testify/assert"
)

func Test_NewServer(t *testing.T) {
	solver, err := battlesnake.NewSolver(battlesnake.SOLVER_RANDOM)
	assert.NoError(t, err)

	server, err := NewServer(":8080", ServerInfo{}, solver)
	assert.NoError(t, err)
	assert.IsType(t, Server{}, *server)
}

func Test_Server_Home(t *testing.T) {
	server := Server{
		info: `{"apiVersion":"1"}`,
	}

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	server.home(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	res, err := io.ReadAll(w.Body)
	assert.NoError(t, err)
	assert.Equal(t, server.info, string(res))
}

func Test_Server_Start_Success(t *testing.T) {
	server := Server{}

	r := httptest.NewRequest(http.MethodPost, "/start", strings.NewReader(`{}`))
	w := httptest.NewRecorder()

	server.start(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func Test_Server_Start_Failure(t *testing.T) {
	server := Server{}

	r := httptest.NewRequest(http.MethodPost, "/start", nil)
	w := httptest.NewRecorder()

	server.start(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	res, err := io.ReadAll(w.Body)

	assert.NoError(t, err)
	assert.Equal(t, `{"error":"invalid JSON body"}`, string(res))
}

func Test_Server_Move_Success(t *testing.T) {
	solver, err := battlesnake.NewSolver(battlesnake.SOLVER_RANDOM)
	assert.NoError(t, err)

	server := Server{solver: solver}

	r := httptest.NewRequest(http.MethodPost, "/move", strings.NewReader(`{}`))
	w := httptest.NewRecorder()

	server.move(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	res, err := io.ReadAll(w.Body)
	assert.NoError(t, err)

	re := regexp.MustCompile(
		`{"move": "(up|right|down|left)", "shout": "Moving (up|right|down|left)!"}`,
	)
	assert.True(t, re.Match(res))
}

func Test_Server_Move_Failure(t *testing.T) {
	server := Server{}

	r := httptest.NewRequest(http.MethodPost, "/move", nil)
	w := httptest.NewRecorder()

	server.move(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	res, err := io.ReadAll(w.Body)

	assert.NoError(t, err)
	assert.Equal(t, `{"error":"invalid JSON body"}`, string(res))
}

func Test_Server_End_Success(t *testing.T) {
	server := Server{}

	r := httptest.NewRequest(http.MethodPost, "/end", strings.NewReader(`{}`))
	w := httptest.NewRecorder()

	server.end(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func Test_Server_End_Failure(t *testing.T) {
	server := Server{}

	r := httptest.NewRequest(http.MethodPost, "/end", nil)
	w := httptest.NewRecorder()

	server.end(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	res, err := io.ReadAll(w.Body)

	assert.NoError(t, err)
	assert.Equal(t, `{"error":"invalid JSON body"}`, string(res))
}

func Test_Server_Run_and_Shutdown(t *testing.T) {
	server := Server{listenAddress: ":8080"}
	go func() {
		err := server.Run()
		assert.ErrorIs(t, err, http.ErrServerClosed)
	}()
	time.Sleep(time.Millisecond)
	err := server.Shutdown()
	assert.NoError(t, err)
}
