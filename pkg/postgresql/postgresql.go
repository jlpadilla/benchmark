package postgresql

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	pgx "github.com/jackc/pgx/v4"
	"github.com/jlpadilla/benchmark/pkg/generator"
)

func createConn() *pgx.Conn {
	// start := time.Now()
	database_url := "postgres://postgres:dev-pass!@localhost:5432/benchmark"
	// conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	// fmt.Println("Connection Took:", time.Since(start))
	return conn
}

func init() {
	c := createConn()
	defer c.Close(context.Background())

	// Clear resources table
	tables := []string{"resources"} //, "resources0", "resources1", "resources2", "resources3", "resources4", "resources5", "resources6", "resources7"}
	for _, table := range tables {
		_, error := c.Exec(context.Background(), fmt.Sprintf("DROP TABLE %s", table))
		if error != nil {
			fmt.Println("Error dropping table. ", table, error)
		}
		_, err := c.Exec(context.Background(), fmt.Sprintf("CREATE TABLE %s(UID text PRIMARY KEY, Cluster text, NAME text, DATA JSONB)", table))
		if err != nil {
			fmt.Println("Error creating table ", table, error)
		}
	}

}

func ProcessInsert(instance string, insertChan chan *generator.Record) {
	conn := createConn()
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
			br := conn.SendBatch(context.Background(), batch)
			res, err := br.Exec()
			if err != nil {
				fmt.Println("res: ", res, "  err: ", err, batch.Len())
			}
			br.Close()
			batch = &pgx.Batch{}
		}
	}
}

/* SAMPLE QUERIES //

// SELECT * from resources WHERE data->'counter' = '3';
// SELECT * from resources WHERE data->>'color' = 'Blue' LIMIT 10;
// SELECT * from resources WHERE data->>'color' = 'Blue' AND data->'counter' < '100';

// List all values for a property.
// SELECT DISTINCT data->'color' AS color from resources;

// List all properties
//
*/
func BenchmarkQueries() {
	conn := createConn()
	defer conn.Close(context.Background())

	executeQueryByUID(conn)

	executeQueryByJSONB(conn)

	executeQueryAllValues(conn)
}

func executeQueryByUID(conn *pgx.Conn) {
	var name string
	var data string
	start := time.Now()
	err := conn.QueryRow(context.Background(), "SELECT name,data FROM resources WHERE uid=$1", "id-1-1").Scan(&name, &data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query by UID failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Query by UID (primary key):\t\t", time.Since(start))
}

func executeQueryByJSONB(conn *pgx.Conn) {
	var name string
	var data string
	start := time.Now()
	err := conn.QueryRow(context.Background(), "SELECT name,data FROM resources WHERE data->>'color' = $1", "Blue").Scan(&name, &data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query JSONB property failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Query by property (JSONB):\t\t", time.Since(start))
}

func executeQueryAllValues(conn *pgx.Conn) {
	var values string
	start := time.Now()
	err := conn.QueryRow(context.Background(), "SELECT DISTINCT data->'color' AS color from resources").Scan(&values)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query get all values for JSONB property failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Query all distinct values of property:\t", time.Since(start))
}
