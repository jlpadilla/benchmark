package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jlpadilla/benchmark/pkg/generator"
	"github.com/jlpadilla/benchmark/pkg/postgresql"
	"github.com/jlpadilla/benchmark/pkg/redisgraph"
	"github.com/jlpadilla/benchmark/pkg/server"
)

var targetDb = "postgresql"

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
		_, insert := readInputs()
		fmt.Println("Running benchmark with:")
		fmt.Println("\tDatabase: ", targetDb)
		fmt.Println("\tInsert : ", insert)

		sim := postgresql.NewTransaction()
		// sim := redisgraph.NewTransaction()

		// Start generating records.
		start := time.Now()
		generator.Generate(insert, 0, 0, sim.InsertChan, sim.InsertChan, sim.DeleteChan)
		sim.WG.Wait()
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
	http.HandleFunc("/clear", clearDB)
	http.HandleFunc("/generate", generate)
	http.HandleFunc("/query", query)
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		fmt.Println("Error starting http server.")
	}
}

func generate(w http.ResponseWriter, req *http.Request) {
	opts := server.ParseQuery(req)
	fmt.Printf("Starting with options: %+v", opts)
	start := time.Now()

	switch opts.Database {
	case "postgresql":
		sim := postgresql.NewTransaction()
		generator.Generate(opts.Insert, opts.Update, opts.Delete, sim.InsertChan, sim.UpdateChan, sim.DeleteChan)
		sim.WG.Wait()
	case "redisgraph":
		sim := redisgraph.NewTransaction()
		generator.Generate(opts.Insert, opts.Update, opts.Delete, sim.InsertChan, sim.UpdateChan, sim.DeleteChan)
		sim.WG.Wait()
	default:
		fmt.Println("\nDatabase not supported: ", opts.Database)
	}
	fmt.Printf("DONE\n")
	fmt.Fprintf(w, "Options: %+v\nTook:\t\t%v\n", opts, time.Since(start))
}

func query(w http.ResponseWriter, req *http.Request) {
	sim := postgresql.NewTransaction()
	result := sim.BenchmarkQueries()
	fmt.Printf("Query results:\n%s", result)
	fmt.Fprintf(w, "Query results:\n%s", result)
}

func clearDB(w http.ResponseWriter, req *http.Request) {
	postgresql.InitializeDB()
}
