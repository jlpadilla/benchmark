package redisgraph

import (
	"fmt"
	"strings"

	rg "github.com/redislabs/redisgraph-go"
)

func (t *transaction) batchDelete() {
	t.Simulation.WG.Add(1)
	defer t.Simulation.WG.Done()
	conn := Pool.Get()
	defer conn.Close()

	g := rg.Graph{
		Conn: conn,
		Id:   GRAPH_NAME,
	}

	uids := []string{}

	for {
		record, more := <-t.Simulation.DeleteChan
		if more {
			uids = append(uids, fmt.Sprintf("'%s'", record))
		}
		if len(uids) == t.options.BatchSize || (!more && len(uids) > 0) {
			q := fmt.Sprintf("MATCH (n) WHERE n._uid IN [%s] DELETE n", strings.Join(uids, ", "))
			_, err := g.Query(q)

			if err != nil {
				fmt.Println("error: ", err)
			}
			uids = []string{}
		}
		if !more {
			break
		}
	}
}
