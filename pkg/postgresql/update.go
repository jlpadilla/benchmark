package postgresql

import (
	"context"
	"encoding/json"
	"fmt"

	pgx "github.com/jackc/pgx/v4"
)

func (t *transaction) batchUpdate() {
	t.Simulation.WG.Add(1)
	defer t.Simulation.WG.Done()
	batch := &pgx.Batch{}
	queryTemplate := "UPDATE " + tables[0] + " SET CLUSTER=$2, DATA=$3 WHERE UID=$1"

	for {
		record, more := <-t.Simulation.UpdateChan

		if more {
			// Marshal record.Properties to JSON
			json, err := json.Marshal(record.Properties)
			if err != nil {
				panic(fmt.Sprintf("Error Marshaling json. %v %v", err, json))
			}

			batch.Queue(queryTemplate, record.UID, record.Cluster, string(json))
		}

		if batch.Len() == t.options.BatchSize || (!more && batch.Len() > 0) {
			fmt.Print(".")
			br := pool.SendBatch(context.Background(), batch)
			res, err := br.Exec()
			if err != nil {
				fmt.Println("res: ", res, "  err: ", err, batch.Len())
			}
			br.Close()
			batch = &pgx.Batch{}
		}
		if !more {
			break
		}
	}
}
