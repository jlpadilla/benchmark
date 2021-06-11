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
		fmt.Println("To use as a standalone go program pass number of records to insert.")
		fmt.Println("Example: \"go run main.go [insert]\"")

		fmt.Println("\nStarting web server at localhost:8090")
		fmt.Println("\nSample curl commands:")
		fmt.Println("\tcurl 'localhost:8090/generate?insert=100000&update=100&delete=100'")
		fmt.Println("\tcurl localhost:8090/query")
		fmt.Println("\tcurl localhost:8090/clear")
		fmt.Println("")
		startHttpServer()
	} else {
		_, insert := readInputs()
		fmt.Println("Running benchmark with:")
		fmt.Println("\tDatabase: ", targetDb)
		fmt.Println("\tInsert : ", insert)

		postgre := postgresql.NewTransaction()

		// Start generating records.
		start := time.Now()
		generator.Generate(insert, 0, 0, postgre.InsertChan, postgre.InsertChan, postgre.DeleteChan)
		postgre.WG.Wait()
		fmt.Printf("\nInsert %d records took: %s", insert, time.Since(start))

		// Benchmark queries.
		fmt.Println("\nWaiting 5 seconds before running queries.")
		time.Sleep(5 * time.Second)
		result := postgresql.BenchmarkQueries()
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

	postgre := postgresql.NewTransaction()

	start := time.Now()
	generator.Generate(insert, update, delete, postgre.InsertChan, postgre.UpdateChan, postgre.DeleteChan)
	postgre.WG.Wait()
	fmt.Printf("DONE\n")
	fmt.Fprintf(w, "Database:\t%s\nInsert:\t\t%d\nUpdate:\t\t%d\nDelete:\t\t%d\nTook:\t\t%v\n", targetDb, insert, update, delete, time.Since(start))
}

func query(w http.ResponseWriter, req *http.Request) {
	result := postgresql.BenchmarkQueries()
	fmt.Printf("Query results:\n%s", result)
	fmt.Fprintf(w, "Query results:\n%s", result)
}

func clearDB(w http.ResponseWriter, req *http.Request) {
	postgresql.InitializeDB()
}
