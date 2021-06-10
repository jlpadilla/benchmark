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

var numRecords = 100000
var targetDb = "postgresql"
var generateCounter = 0

func main() {

	if len(os.Args) < 2 {
		fmt.Println("No arguments were passed, starting web server.")
		fmt.Println("To use as a standalone go program pass the number of records to simulate like \"go run main.go [numRecords]\"")
		startHttpServer()
	} else {

		// Create a channel to send resources from the generator to the db insert.
		insertChan := make(chan *generator.Record, 100)

		// Reads the records and inserts in the target database.
		startPostgre(insertChan)

		targetDb, numRecords = readInputs()
		fmt.Println("Running benchmark with:")
		fmt.Println("\tDatabase: ", targetDb)
		fmt.Println("\tRecords : ", numRecords)

		// Generate records.
		generateRecords(numRecords, insertChan)

		fmt.Println("Waiting 10 seconds, then benchmarking queries")
		time.Sleep(10 * time.Second)
		postgresql.BenchmarkQueries()
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

func generateRecords(numRecords int, insertChan chan *generator.Record) {
	start := time.Now()

	// NOTE: I experimented with multiple generate routines, but it didn't make a difference.
	generator.Generate(strconv.Itoa(generateCounter), numRecords, insertChan)
	generateCounter++
	postgresql.WG.Wait()
	fmt.Println("\nInserting records took:", time.Since(start))
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
	records, _ := strconv.Atoi(req.URL.Query()["records"][0])

	insertChan := make(chan *generator.Record, 100)
	startPostgre(insertChan)

	start := time.Now()
	generateRecords(records, insertChan)

	fmt.Fprintf(w, "Database:\t%s\nRecords:\t%d\nTook:\t\t%v\n", targetDb, numRecords, time.Since(start))
}

func query(w http.ResponseWriter, req *http.Request) {
	result := postgresql.BenchmarkQueries()
	fmt.Printf("Query results:\n%s", result)
	fmt.Fprintf(w, "Query results:\n%s", result)
}

func clearDB(w http.ResponseWriter, req *http.Request) {
	postgresql.InitializeDB()
}
