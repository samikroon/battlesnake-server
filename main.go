package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/samikroon/battlesnake-server/battlesnake"
	"github.com/samikroon/battlesnake-server/server"
)

// getEnvString retrieves an environment variable by key.
// If the environment variable is not set, the defaultVal is returned.
func getEnvString(key string, defaultVal string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultVal
	}
	return val
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("WARNING Could not load env, continuing with system variables.")
	}

	serverInfo := server.ServerInfo{
		ApiVersion: getEnvString("API_VERSION", "1"),
		Author:     getEnvString("AUTHOR", "Sam Kroon"),
		Color:      getEnvString("COLOR", "#4287f5"),
		Head:       getEnvString("HEAD", "default"),
		Tail:       getEnvString("TAIL", "default"),
		Version:    getEnvString("VERSION", "0.0.1"),
	}
	listenAddress := getEnvString("LISTEN_ADDRESS", "127.0.0.1:8080")
	solverType := getEnvString("SOLVER", battlesnake.SOLVER_RANDOM)

	solver, err := battlesnake.NewSolver(solverType)
	if err != nil {
		log.Fatalf("ERROR Create solver with type %s: %s", solverType, err)
	}

	server, err := server.NewServer(listenAddress, serverInfo, solver)
	if err != nil {
		log.Fatalf("ERROR New server: %s", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ERROR Run server: %s", err)
		}
	}()
	log.Println("INFO Server started! Listening at: ", listenAddress)

	<-done
	log.Println("INFO Server stopped!")
	if err := server.Shutdown(); err != nil {
		log.Fatalf("ERROR Shutdown server: %s", err)
	}
}
