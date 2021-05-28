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
	targetdb := os.Args[1]
	if targetdb == "" {
		fmt.Println("usage: go run main.go [targetdb] [numRecords]")
		panic("Must pass target database.")
	}
	fmt.Println("Target database: ", targetdb)

	numRecords, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("usage: go run main.go [targetdb] [numRecords]")
		panic("Must pass number of records to profile.")
	}
	fmt.Println("Records to add : ", numRecords)

	// Create a channel to send resources from the generator to the db insert.
	insertChan := make(chan *generator.Record)

	for i := 0; i < 6; i++ {
		go postgresql.ProcessInsert(strconv.Itoa(i), insertChan)
	}

	start := time.Now()
	for i := 0; i < 3; i++ {
		go generator.Generate(strconv.Itoa(i), numRecords/4, insertChan)
	}
	generator.Generate("4", numRecords/4, insertChan)

	fmt.Println("\n Took:", time.Now().Sub(start))
}
