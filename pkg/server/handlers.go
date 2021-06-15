package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jlpadilla/benchmark/pkg/generator"
	"github.com/jlpadilla/benchmark/pkg/postgresql"
	"github.com/jlpadilla/benchmark/pkg/redisgraph"
)

func Generate(w http.ResponseWriter, req *http.Request) {
	opts := ParseQuery(req)
	fmt.Printf("Starting generator with options: %+v\n", opts)
	start := time.Now()

	switch opts.Database {
	case "postgresql":
		sim := postgresql.NewTransaction(opts)
		generator.Generate(opts, sim.Simulation)
		sim.Simulation.WG.Wait()
	case "redisgraph":
		sim := redisgraph.NewTransaction(opts)
		generator.Generate(opts, sim.Simulation)
		sim.Simulation.WG.Wait()
	default:
		fmt.Println("\nDatabase not supported: ", opts.Database)
	}
	fmt.Printf("DONE\n")
	fmt.Fprintf(w, "Options: %+v\nTook:\t\t%v\n", opts, time.Since(start))
}

func Query(w http.ResponseWriter, req *http.Request) {
	opts := ParseQuery(req)
	fmt.Printf("Starting Query with options %+v\n", opts)
	sim := postgresql.NewTransaction(opts)
	result := sim.BenchmarkQueries()
	fmt.Printf("Query results:\n%s", result)
	fmt.Fprintf(w, "Query results:\n%s", result)
}

func ClearDB(w http.ResponseWriter, req *http.Request) {
	postgresql.InitializeDB()
}

func ParseQuery(req *http.Request) generator.Options {
	insert := 1000 // Default to 1000 records.
	insertQuery := req.URL.Query()["insert"]
	if len(insertQuery) > 0 {
		insert, _ = strconv.Atoi(req.URL.Query()["insert"][0])
	}
	update := 0
	updateQuery := req.URL.Query()["update"]
	if len(updateQuery) > 0 {
		update, _ = strconv.Atoi(req.URL.Query()["update"][0])
	}
	delete := 0
	deleteQuery := req.URL.Query()["delete"]
	if len(deleteQuery) > 0 {
		delete, _ = strconv.Atoi(req.URL.Query()["delete"][0])
	}

	database := "postgresql"
	databaseQuery := req.URL.Query()["database"]
	if len(databaseQuery) > 0 {
		database = string(req.URL.Query()["database"][0])
	}

	batchSize := 1000
	batchQuery := req.URL.Query()["batchsize"]
	if len(batchQuery) > 0 {
		batchSize, _ = strconv.Atoi(req.URL.Query()["batchsize"][0])
	}

	goRoutines := 8
	routinesQuery := req.URL.Query()["goroutines"]
	if len(routinesQuery) > 0 {
		goRoutines, _ = strconv.Atoi(req.URL.Query()["goroutines"][0])
	}

	insertType := "batch"

	return generator.Options{
		BatchSize:  batchSize,
		Database:   database,
		GoRoutines: goRoutines,
		InsertType: insertType,
		Insert:     insert,
		Update:     update,
		Delete:     delete,
	}
}
