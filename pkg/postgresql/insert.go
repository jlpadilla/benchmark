package postgresql

import (
	"context"
	"encoding/json"
	"fmt"

	pgx "github.com/jackc/pgx/v4"
	"github.com/jlpadilla/benchmark/pkg/generator"
)

// Sends INSERT commands in a batch request.
func batchInsert(instance string, insertChan chan *generator.Record) {
	// conn := createConn()
	batch := &pgx.Batch{}

	for {
		record := <-insertChan

		// Marshal record.Properties to JSON
		json, err := json.Marshal(record.Properties)
		if err != nil {
			panic(fmt.Sprintf("Error Marshaling json. %v %v", err, json))
		}

		// batch.Queue("insert into resources values($1,$2,$3,$4)", record.UID, record.Cluster, record.Kind, record.Name)
		// batch.Queue(fmt.Sprintf("insert into resources%s values($1,$2,$3,$4,$5)", instance), record.UID, record.Cluster, record.Kind, record.Name, string(json))
		batch.Queue("insert into resources values($1,$2,$3,$4)", record.UID, record.Cluster, record.Name, string(json))

		if batch.Len()%300 == 0 {
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

// Inserts records in bulk using the COPY command.
func copyInsert(instance string, insertChan chan *generator.Record) {
	var inputRows = make([][]interface{}, 0)
	for {
		record := <-insertChan

		// Marshal record.Properties to JSON
		json, err := json.Marshal(record.Properties)
		if err != nil {
			panic(fmt.Sprintf("Error Marshaling json. %v %v", err, json))
		}

		row := []interface{}{record.UID, record.Cluster, record.Name, json}
		inputRows = append(inputRows, row)

		if len(inputRows) == 1000 {
			WG.Add(1)
			fmt.Println("Connections  MAX: ", pool.Stat().MaxConns(), "  CURRENT: ", pool.Stat().TotalConns())
			sendUsingCopy(inputRows)
			// WG.Done()
			inputRows = make([][]interface{}, 0)
		}
	}
}

func sendUsingCopy(inputRows [][]interface{}) {
	defer WG.Done()
	// start := time.Now()
	// fmt.Println("Conn Took:", time.Since(start))

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
