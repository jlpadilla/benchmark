package postgresql

import (
	"context"
	"fmt"

	pgx "github.com/jackc/pgx/v4"
)

func (t *transaction) batchDelete() {
	t.WG.Add(1)
	defer t.WG.Done()
	batch := &pgx.Batch{}
	queryTemplate := "DELETE FROM " + tables[0] + " WHERE UID=$1"

	for {
		record, more := <-t.DeleteChan

		if more {
			batch.Queue(queryTemplate, record)
		}

		if batch.Len() == t.batchSize || (!more && batch.Len() > 0) {
			fmt.Print("-")
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
