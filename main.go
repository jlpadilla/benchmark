package main

import (
	"fmt"
	"os"
	"strconv"

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

	postgresql.Start(numRecords)
}
