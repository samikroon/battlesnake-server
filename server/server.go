package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/samikroon/battlesnake-server/battlesnake"
)

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

func (s *Server) Run() error {
	router := mux.NewRouter()
	router.HandleFunc("/", s.home).Methods(http.MethodGet)
	router.HandleFunc("/start", s.start).Methods(http.MethodPost)
	router.HandleFunc("/move", s.move).Methods(http.MethodPost)
	router.HandleFunc("/end", s.end).Methods(http.MethodPost)

	s.server = &http.Server{
		Addr:    s.listenAddress,
		Handler: router,
	}

	return s.server.ListenAndServe()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, s.info)
}

func (s *Server) start(w http.ResponseWriter, r *http.Request) {
	var state battlesnake.State
	if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"error":"invalid JSON body"}`)
		return
	}

	log.Printf(
		"INFO Started game with size %dx%d (id=%s)",
		state.Board.Width, state.Board.Height, state.Game.ID,
	)

	w.WriteHeader(http.StatusOK)
}

func (s *Server) move(w http.ResponseWriter, r *http.Request) {
	var state battlesnake.State
	if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"error":"invalid JSON body"}`)
		return
	}

	move := s.solver.Move(state)

	responseBody := fmt.Sprintf(
		`{"move": "%s", "shout": "Moving %s!"}`,
		move, move,
	)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, responseBody)
}

func (s *Server) end(w http.ResponseWriter, r *http.Request) {
	var state battlesnake.State
	if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"error":"invalid JSON body"}`)
		return
	}

	log.Printf(
		"INFO Ended game with size %dx%d after %d turns (id=%s)",
		state.Board.Width, state.Board.Height, state.Turn, state.Game.ID,
	)

	w.WriteHeader(http.StatusOK)
}
