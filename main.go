package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jlpadilla/benchmark/pkg/generator"
	"github.com/jlpadilla/benchmark/pkg/postgresql"
)

var targetDb = "postgresql"

func main() {

	if len(os.Args) < 2 {
		fmt.Println("To use as a standalone go program pass [numRecords].")
		fmt.Println("Example: \"go run main.go [numRecords]\"")

		fmt.Println("\nStarting web server at localhost:8090")
		fmt.Println("\nSample curl commands:")
		fmt.Println("\tcurl 'localhost:8090/generate?records=100000'")
		fmt.Println("\tcurl localhost:8090/query")
		fmt.Println("\tcurl localhost:8090/clear")
		fmt.Println("")
		startHttpServer()
	} else {
		_, numRecords := readInputs()
		fmt.Println("Running benchmark with:")
		fmt.Println("\tDatabase: ", targetDb)
		fmt.Println("\tRecords : ", numRecords)

		// Create a channel to send resources from the generator to the db insert.
		insertChan := make(chan *generator.Record, 100)

		// Start routine to process records from channel into the target database.
		startPostgre(insertChan)

		// Start generating records.
		start := time.Now()
		generator.Generate(numRecords, insertChan)
		postgresql.WG.Wait()
		fmt.Printf("\nInsert %d records took: %s", numRecords, time.Since(start))

		// Benchmark queries.
		fmt.Println("\nWaiting 10 seconds before running queries.")
		time.Sleep(10 * time.Second)
		result := postgresql.BenchmarkQueries()
		fmt.Printf("\nQuery benchmark results:\n%s", result)
	}
}

func readInputs() (targetDb string, numRecords int) {
	numRecords, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Must pass number of records to benchmark.")
		fmt.Println("USAGE: go run main.go [numRecords]")
		os.Exit(1)
	}

	return "postgresql", numRecords
}

func startPostgre(insertChan chan *generator.Record) {
	for i := 0; i < 8; i++ {
		go postgresql.ProcessInsert(strconv.Itoa(i), insertChan)
	}
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
	records := 1000 // Default to 1000 records.
	recordsQuery := req.URL.Query()["records"]
	if len(recordsQuery) > 0 {
		records, _ = strconv.Atoi(req.URL.Query()["records"][0])
	}
	insertChan := make(chan *generator.Record, 100)
	startPostgre(insertChan)

	start := time.Now()
	generator.Generate(records, insertChan)
	postgresql.WG.Wait()
	fmt.Fprintf(w, "Database:\t%s\nRecords:\t%d\nTook:\t\t%v\n", targetDb, records, time.Since(start))
}

func query(w http.ResponseWriter, req *http.Request) {
	result := postgresql.BenchmarkQueries()
	fmt.Printf("Query results:\n%s", result)
	fmt.Fprintf(w, "Query results:\n%s", result)
}

func clearDB(w http.ResponseWriter, req *http.Request) {
	postgresql.InitializeDB()
}
