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
	fmt.Printf("Starting with options: %+v", opts)
	start := time.Now()

	switch opts.Database {
	case "postgresql":
		sim := postgresql.NewTransaction()
		generator.Generate(opts.Insert, opts.Update, opts.Delete, sim.InsertChan, sim.UpdateChan, sim.DeleteChan)
		sim.WG.Wait()
	case "redisgraph":
		sim := redisgraph.NewTransaction(opts)
		generator.Generate(opts.Insert, opts.Update, opts.Delete, sim.InsertChan, sim.UpdateChan, sim.DeleteChan)
		sim.WG.Wait()
	default:
		fmt.Println("\nDatabase not supported: ", opts.Database)
	}
	fmt.Printf("DONE\n")
	fmt.Fprintf(w, "Options: %+v\nTook:\t\t%v\n", opts, time.Since(start))
}

func Query(w http.ResponseWriter, req *http.Request) {
	sim := postgresql.NewTransaction()
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
	goRoutines := 8
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
