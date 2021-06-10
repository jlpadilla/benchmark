package postgresql

import (
	"context"
	"fmt"
	"os"
	"sync"

	pgxpool "github.com/jackc/pgx/v4/pgxpool"
	"github.com/jlpadilla/benchmark/pkg/generator"
)

const databaseName = "benchmark"
const insertType = "batch" // "batch or copy"
const batchSize = 1000
const maxConnections = 16

var tables = []string{"resources"}
var WG sync.WaitGroup
var pool *pgxpool.Pool

func init() {
	createPool()
	InitializeDB()
}

// Initialize the database. Drop existing tables and create new tables for this test.
func InitializeDB() {
	// Clear resources table
	for _, table := range tables {
		_, error := pool.Exec(context.Background(), fmt.Sprintf("DROP TABLE %s", table))
		if error != nil {
			fmt.Println("Error dropping table. ", table, error)
		}
		_, err := pool.Exec(context.Background(), fmt.Sprintf("CREATE TABLE %s(UID text PRIMARY KEY, Cluster text, NAME text, DATA JSONB)", table))
		if err != nil {
			fmt.Println("Error creating table ", table, error)
		}
	}
}

// Initializes the connection pool.
func createPool() {
	database_url := "postgres://postgres:dev-pass!@localhost:5432/" + databaseName

	config, _ := pgxpool.ParseConfig(database_url)
	config.MaxConns = maxConnections
	conn, err := pgxpool.ConnectConfig(context.Background(), config)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	pool = conn
}

func ProcessInsert(instance string, insertChan chan *generator.Record) {
	if insertType == "batch" {
		batchInsert(instance, insertChan)
	} else {
		copyInsert(instance, insertChan)
	}
}
