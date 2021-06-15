package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jlpadilla/benchmark/pkg/generator"
	"github.com/jlpadilla/benchmark/pkg/postgresql"
	"github.com/jlpadilla/benchmark/pkg/server"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("To use as a standalone go program pass number of records to insert.")
		fmt.Println("Example: \"go run main.go [insert]\"")

		fmt.Println("\nStarting web server at localhost:8090")
		fmt.Println("\nSample curl commands:")
		fmt.Println("\tcurl 'localhost:8090/generate?database=postgresql&insert=100000&update=100&delete=100'")
		fmt.Println("\tcurl localhost:8090/query")
		fmt.Println("\tcurl localhost:8090/clear")
		fmt.Println("")
		startHttpServer()
	} else {
		targetDb, insert := readInputs()
		fmt.Println("Running benchmark with:")
		fmt.Println("\tDatabase: ", targetDb)
		fmt.Println("\tInsert : ", insert)

		opts := generator.Options{Database: targetDb, Insert: insert, Update: 0, Delete: 0, BatchSize: 1000, GoRoutines: 8, InsertType: "batch"}
		start := time.Now()

		sim := postgresql.NewTransaction(opts)
		// sim := redisgraph.NewTransaction()

		// Start generating records.

		generator.Generate(opts, sim.Simulation)
		sim.Simulation.WG.Wait()
		fmt.Printf("\nInsert %d records took: %s", insert, time.Since(start))

		// Benchmark queries.
		fmt.Println("\nWaiting 5 seconds before running queries.")
		time.Sleep(5 * time.Second)
		result := sim.BenchmarkQueries()
		fmt.Printf("\nQuery benchmark results:\n%s", result)
	}
}

func readInputs() (targetDb string, numRecords int) {
	numRecords, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Must pass number of records to benchmark.")
		fmt.Println("USAGE: go run main.go [insert]")
		os.Exit(1)
	}

	return "postgresql", numRecords
}

//
// HTTP Server
//

func startHttpServer() {
	http.HandleFunc("/clear", server.ClearDB)
	http.HandleFunc("/generate", server.Generate)
	http.HandleFunc("/query", server.Query)
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		fmt.Println("Error starting http server.")
	}
}
