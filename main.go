package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jlpadilla/benchmark/pkg/generator"
	"github.com/jlpadilla/benchmark/pkg/postgresql"
)

func main() {
	targetDb, numRecords := readInputs()
	fmt.Println("Running benchmark with:")
	fmt.Println("\tDatabase: ", targetDb)
	fmt.Println("\tRecords : ", numRecords)

	// Create a channel to send resources from the generator to the db insert.
	insertChan := make(chan *generator.Record, 100)

	// Reads the records and inserts in the target database.
	startPostgre(insertChan)

	start := time.Now()

	// Generate records.
	generateRecords(numRecords, insertChan)

	postgresql.WG.Wait()
	fmt.Println("\nLoad DB took:", time.Since(start))

	fmt.Println("Waiting 10 seconds, then benchmarking queries")
	time.Sleep(10 * time.Second)
	postgresql.BenchmarkQueries()

	fmt.Println("\nDONE. Goodbye.")
}

func readInputs() (targetDb string, numRecords int) {
	if len(os.Args) < 2 {
		fmt.Println("Must pass number of records to benchmark.")
		fmt.Println("USAGE: go run main.go [numRecords]")
		os.Exit(1)
	}
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
	routines := 1
	for i := 1; i < routines; i++ {
		go generator.Generate(strconv.Itoa(i), numRecords/routines, insertChan)
	}
	generator.Generate(strconv.Itoa(routines), numRecords/routines, insertChan)
}
