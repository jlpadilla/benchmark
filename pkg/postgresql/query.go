package postgresql

import (
	"context"
	"fmt"
	"time"
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
func BenchmarkQueries() string {
	result := ""
	result += executeQueryByUID()
	result += executeQueryByJSONB()
	result += executeQueryAllValues()
	return result
}

func executeQueryByUID() string {
	var name string
	var data string
	start := time.Now()
	err := pool.QueryRow(context.Background(), "SELECT name,data FROM resources WHERE uid=$1", "id-0-1").Scan(&name, &data)
	if err != nil {
		return fmt.Sprintln("Query by UID (primary key):\t\t ERROR:", err)
	}

	return fmt.Sprintln("Query by UID (primary key):\t\t", time.Since(start))
}

func executeQueryByJSONB() string {
	var name string
	var data string
	start := time.Now()
	err := pool.QueryRow(context.Background(), "SELECT name,data FROM resources WHERE data->>'color' = $1", "Blue").Scan(&name, &data)
	if err != nil {
		return fmt.Sprintln("Query by property (JSONB):\t\t ERROR:", err)
	}

	return fmt.Sprintln("Query by property (JSONB):\t\t", time.Since(start))
}

func executeQueryAllValues() string {
	var values string
	start := time.Now()
	err := pool.QueryRow(context.Background(), "SELECT DISTINCT data->'color' AS color from resources").Scan(&values)
	if err != nil {
		return fmt.Sprintln("Query all distinct values of property:\t ERROR:", err)
	}

	return fmt.Sprintln("Query all distinct values of property:\t", time.Since(start))
}
