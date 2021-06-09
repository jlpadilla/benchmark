package postgresql

import (
	"context"
	"encoding/json"
	"fmt"

	pgx "github.com/jackc/pgx/v4"
	"github.com/jlpadilla/benchmark/pkg/generator"
)

// Process records using batched INSERT requests.
func batchInsert(instance string, insertChan chan *generator.Record) {
	batch := &pgx.Batch{}

	for {
		record := <-insertChan

		// Marshal record.Properties to JSON
		json, err := json.Marshal(record.Properties)
		if err != nil {
			panic(fmt.Sprintf("Error Marshaling json. %v %v", err, json))
		}

		batch.Queue("insert into resources values($1,$2,$3,$4)", record.UID, record.Cluster, record.Name, string(json))

		if batch.Len() == batchSize {
			fmt.Print(".")
			br := pool.SendBatch(context.Background(), batch)
			res, err := br.Exec()
			if err != nil {
				fmt.Println("res: ", res, "  err: ", err, batch.Len())
			}
			br.Close()
			batch = &pgx.Batch{}
		}
	}
}

// Process records in bulk using COPY.
func copyInsert(instance string, insertChan chan *generator.Record) {
	inputRows := make([][]interface{}, batchSize)
	index := 0
	for {
		record := <-insertChan

		// Marshal record.Properties to JSON
		json, err := json.Marshal(record.Properties)
		if err != nil {
			panic(fmt.Sprintf("Error Marshaling json. %v %v", err, json))
		}

		inputRows[index] = []interface{}{record.UID, record.Cluster, record.Name, json}
		index++

		if index == batchSize {
			WG.Add(1)
			go sendUsingCopy(inputRows)

			inputRows = make([][]interface{}, batchSize)
			index = 0
		}
	}
}

// Load records using the COPY command.
func sendUsingCopy(inputRows [][]interface{}) {
	defer WG.Done()
	// start := time.Now()

	// UID text PRIMARY KEY, Cluster text, NAME text, DATA JSONB
	copyCount, err := pool.CopyFrom(context.Background(), pgx.Identifier{tables[0]}, []string{"uid", "cluster", "name", "data"},
		pgx.CopyFromRows(inputRows))

	if err != nil {
		fmt.Printf("Unexpected error for CopyFrom: %v", err)
	} else if int(copyCount) != len(inputRows) {
		fmt.Printf("Expected CopyFrom to return %d copied rows, but got %d", len(inputRows), copyCount)
	}

	fmt.Print(".")
	// fmt.Println("COPY Took:", time.Since(start))
}
