package redisgraph

// import (
// 	rg "github.com/redislabs/redisgraph-go"
// )

func (t *transaction) batchUpdate() {
	t.WG.Add(1)
	defer t.WG.Done()
	conn := Pool.Get()
	defer conn.Close()

	// g := rg.Graph{
	// 	Conn: conn,
	// 	Id:   GRAPH_NAME,
	// }

	// resourceStrings := []string{}

	for {
		_, more := <-t.DeleteChan

		// resourceStrings = append(resourceStrings, record)

		// if len(uids) == t.batchSize || !more {
		// 	q := fmt.Sprintf("MATCH n WHERE n IN [%s] DELETE n", strings.Join(uids, ", "))
		// 	_, err := g.Query(q)

		// 	if err != nil {
		// 		fmt.Println("error: ", err)
		// 	}
		// 	uids = []string{}
		// }
		if !more {
			break
		}
	}
}
