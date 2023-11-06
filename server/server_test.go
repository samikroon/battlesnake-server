package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/samikroon/battlesnake-server/battlesnake"
	"github.com/stretchr/testify/assert"
)

type MockSolver struct {
	moveFunc func(state battlesnake.State) string
}

func (m MockSolver) Move(state battlesnake.State) string {
	return m.moveFunc(state)
}

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

	res := w.Result()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, server.info, string(body))
}

func Test_Server_Redoc(t *testing.T) {
	server := Server{
		info: `{"apiVersion":"1"}`,
	}

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	server.redoc(w, r)

	res := w.Result()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, contentTypeHTML, res.Header.Get("Content-Type"))
}

func Test_Server_Start_Success(t *testing.T) {
	server := Server{}

	r := httptest.NewRequest(http.MethodPost, "/start", strings.NewReader(`{}`))
	r.Header.Add("Content-Type", contentTypeJSON)
	w := httptest.NewRecorder()

	server.start(w, r)

	res := w.Result()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func Test_Server_Start_InvalidContentType_Error(t *testing.T) {
	server := Server{}

	r := httptest.NewRequest(http.MethodPost, "/start", nil)
	w := httptest.NewRecorder()

	server.start(w, r)

	res := w.Result()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"error":"invalid Content-Type"}`, string(body))
}

func Test_Server_Start_InvalidBody_Error(t *testing.T) {
	server := Server{}

	r := httptest.NewRequest(http.MethodPost, "/start", nil)
	r.Header.Add("Content-Type", contentTypeJSON)
	w := httptest.NewRecorder()

	server.start(w, r)

	res := w.Result()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	body, err := io.ReadAll(w.Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"error":"invalid JSON body"}`, string(body))
}

func Test_Server_Move_Success(t *testing.T) {
	solver := MockSolver{
		moveFunc: func(state battlesnake.State) string { return "up" },
	}

	server := Server{solver: solver}

	r := httptest.NewRequest(http.MethodPost, "/move", strings.NewReader(`{}`))
	r.Header.Add("Content-Type", contentTypeJSON)
	w := httptest.NewRecorder()

	server.move(w, r)

	res := w.Result()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"move": "up", "shout": "Moving up!"}`, string(body))
}

func Test_Server_Move_InvalidContentType_Error(t *testing.T) {
	server := Server{}

	r := httptest.NewRequest(http.MethodPost, "/move", nil)
	w := httptest.NewRecorder()

	server.move(w, r)

	res := w.Result()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"error":"invalid Content-Type"}`, string(body))
}

func Test_Server_Move_InvalidBody_Error(t *testing.T) {
	server := Server{}

	r := httptest.NewRequest(http.MethodPost, "/move", nil)
	r.Header.Add("Content-Type", contentTypeJSON)
	w := httptest.NewRecorder()

	server.move(w, r)

	res := w.Result()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"error":"invalid JSON body"}`, string(body))
}

func Test_Server_End_Success(t *testing.T) {
	server := Server{}

	r := httptest.NewRequest(http.MethodPost, "/end", strings.NewReader(`{}`))
	r.Header.Add("Content-Type", contentTypeJSON)
	w := httptest.NewRecorder()

	server.end(w, r)

	res := w.Result()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func Test_Server_End_InvalidContentType_Error(t *testing.T) {
	server := Server{}

	r := httptest.NewRequest(http.MethodPost, "/end", nil)
	w := httptest.NewRecorder()

	server.end(w, r)

	res := w.Result()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"error":"invalid Content-Type"}`, string(body))
}

func Test_Server_End_InvalidBody_Error(t *testing.T) {
	server := Server{}

	r := httptest.NewRequest(http.MethodPost, "/end", nil)
	r.Header.Add("Content-Type", contentTypeJSON)
	w := httptest.NewRecorder()

	server.end(w, r)

	res := w.Result()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"error":"invalid JSON body"}`, string(body))
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
