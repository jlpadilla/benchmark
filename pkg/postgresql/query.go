package postgresql

import (
	"context"
	"fmt"
	"os"
	"time"

	pgx "github.com/jackc/pgx/v4"
)

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
