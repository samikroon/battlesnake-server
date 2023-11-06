package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/samikroon/battlesnake-server/battlesnake"
	"github.com/samikroon/battlesnake-server/server/resources"
)

const (
	contentTypeJSON = "application/json"
	contentTypeHTML = "text/html"
)

var (
	errInvalidJSONBody    = errors.New("invalid JSON body")
	errInvalidContentType = errors.New("invalid Content-Type")
)

type errResponse struct {
	Error string `json:"error"`
}

type ServerInfo struct {
	ApiVersion string `json:"apiversion"`
	Author     string `json:"author"`
	Color      string `json:"color"`
	Head       string `json:"head"`
	Tail       string `json:"tail"`
	Version    string `json:"version"`
}

type Server struct {
	info          string
	listenAddress string
	solver        battlesnake.Solver
	server        *http.Server
}

// NewServer creates a new Server.
func NewServer(listenAddress string, info ServerInfo, solver battlesnake.Solver) (*Server, error) {
	infoBytes, err := json.Marshal(info)
	if err != nil {
		return nil, fmt.Errorf("json marshal server info: %w", err)
	}
	return &Server{
		info:          string(infoBytes),
		listenAddress: listenAddress,
		solver:        solver,
	}, nil
}

// Run sets up the router and server, after this it starts the server.
func (s *Server) Run() error {
	router := mux.NewRouter()
	router.HandleFunc("/", s.home).Methods(http.MethodGet)
	router.HandleFunc("/docs", s.redoc).Methods(http.MethodGet)
	router.HandleFunc("/openapi.json", s.openapi).Methods(http.MethodGet)
	router.HandleFunc("/start", s.start).Methods(http.MethodPost)
	router.HandleFunc("/move", s.move).Methods(http.MethodPost)
	router.HandleFunc("/end", s.end).Methods(http.MethodPost)

	s.server = &http.Server{
		Addr:    s.listenAddress,
		Handler: router,
	}

	return s.server.ListenAndServe()
}

// Shutdown tries to gracefully shutdown the server, on failure, an error is returned.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return s.server.Shutdown(ctx)
}

// home writes the info JSON string initialized in the server to the ResponseWriter.
func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", contentTypeJSON)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, s.info)
}

// redoc writes the docs HTML to the ResponseWriter.
func (s *Server) redoc(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", contentTypeHTML)
	w.WriteHeader(http.StatusOK)
	w.Write(resources.Redoc)
}

// openapi writes the OpenAPI JSON spec to the ResponseWriter.
func (s *Server) openapi(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", contentTypeJSON)
	w.WriteHeader(http.StatusOK)
	w.Write(resources.OpenApiSpec)
}

// start does not have a lot of utility, it validates the request and responds.
// In case of a validation error, a bad request respoinse is sent (400).
// In case of successful validation, a succes response is sent (200).
func (s *Server) start(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != contentTypeJSON {
		handleErrResponse(w, http.StatusBadRequest, errInvalidContentType)
		return
	}

	var state battlesnake.State
	if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
		handleErrResponse(w, http.StatusBadRequest, errInvalidJSONBody)
		return
	}

	log.Printf(
		"INFO Started game with size %dx%d (id=%s)",
		state.Board.Width, state.Board.Height, state.Game.ID,
	)

	w.WriteHeader(http.StatusOK)
}

// move handles requests for a next move in the game state based on the sent state.
// In case of a validation error, a bad request respoinse is sent (400).
// In case of successful validation, a move is calculated and a succes response is sent (200).
// The response body contains the move that the client should execute for the Battlesnake.
func (s *Server) move(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != contentTypeJSON {
		handleErrResponse(w, http.StatusBadRequest, errInvalidContentType)
		return
	}

	var state battlesnake.State
	if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
		handleErrResponse(w, http.StatusBadRequest, errInvalidJSONBody)
		return
	}

	move := s.solver.Move(state)

	responseBody := fmt.Sprintf(`{"move": "%s", "shout": "Moving %s!"}`, move, move)

	w.Header().Add("Content-Type", contentTypeJSON)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, responseBody)
}

// end does not have a lot of utility, it validates the request and responds.
// In case of a validation error, a bad request respoinse is sent (400).
// In case of successful validation, a succes response is sent (200).
func (s *Server) end(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != contentTypeJSON {
		handleErrResponse(w, http.StatusBadRequest, errInvalidContentType)
		return
	}

	var state battlesnake.State
	if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
		handleErrResponse(w, http.StatusBadRequest, errInvalidJSONBody)
		return
	}

	log.Printf(
		"INFO Ended game with size %dx%d after %d turns (id=%s)",
		state.Board.Width, state.Board.Height, state.Turn, state.Game.ID,
	)

	w.WriteHeader(http.StatusOK)
}

// handleErrResponse writes an error response to the response writer
// with the given status code and error.
func handleErrResponse(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)

	b, err := json.Marshal(errResponse{Error: err.Error()})
	if err != nil {
		log.Printf("ERROR marshal error response: %s", err)
		return
	}

	w.Header().Add("Content-Type", contentTypeJSON)
	w.Write(b)
}
