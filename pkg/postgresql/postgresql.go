package postgresql

import (
	"context"
	"fmt"
	"os"
	"strconv"

	pgxpool "github.com/jackc/pgx/v4/pgxpool"
	"github.com/jlpadilla/benchmark/pkg/generator"
)

// Global settings
const dbHost = "localhost"
const dbPort = "5432"
const postgrePW = "dev-pass!"
const databaseName = "benchmark"
const maxConnections = 8

var tables = []string{"resources"}
var pool *pgxpool.Pool

// Transaction settings
type transaction struct {
	// Configurable fields
	options generator.Options
	// Internal fields
	Simulation generator.Simulation
}

func NewTransaction(options generator.Options) *transaction {
	t := &transaction{
		// Configurable fields
		options: options,
		// Internal fields
		Simulation: generator.NewSimulation(),
	}
	t.startConnectors()
	return t
}

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
		_, err := pool.Exec(context.Background(), fmt.Sprintf("CREATE TABLE %s(UID text PRIMARY KEY, CLUSTER text, DATA JSONB)", table))
		if err != nil {
			fmt.Println("Error creating table ", table, error)
		}
	}

	// Use 100 tables
	// for i := 0; i < 100; i++ {
	// 	table := fmt.Sprintf("cluster_%d", i)
	// 	_, error := pool.Exec(context.Background(), fmt.Sprintf("DROP TABLE %s", table))
	// 	if error != nil {
	// 		fmt.Println("Error dropping table. ", table, error)
	// 	}
	// 	_, err := pool.Exec(context.Background(), fmt.Sprintf("CREATE TABLE %s(UID text PRIMARY KEY, CLUSTER text, DATA JSONB)", table))
	// 	if err != nil {
	// 		fmt.Println("Error creating table ", table, error)
	// 	}
	// }
}

// Initializes the connection pool.
func createPool() {
	database_url := "postgres://postgres:" + postgrePW + "@" + dbHost + ":" + dbPort + "/" + databaseName
	fmt.Println("Connecting to PostgreSQL at: ", database_url)
	config, _ := pgxpool.ParseConfig(database_url)
	config.MaxConns = maxConnections
	conn, err := pgxpool.ConnectConfig(context.Background(), config)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	pool = conn
}

func (t *transaction) startConnectors() {
	for i := 0; i < t.options.GoRoutines; i++ {
		if t.options.InsertType == "batch" {
			go t.batchInsert(strconv.Itoa(i))
		} else {
			go t.copyInsert(strconv.Itoa(i))
		}
		go t.batchUpdate()
		go t.batchDelete()
	}
}
