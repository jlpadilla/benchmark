package postgresql

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"

	pgxpool "github.com/jackc/pgx/v4/pgxpool"
	"github.com/jlpadilla/benchmark/pkg/generator"
)

// Global settings
const postgrePW = "dev-pass!"
const databaseName = "benchmark"
const maxConnections = 8

var tables = []string{"resources"}
var pool *pgxpool.Pool

// Transaction settings
type transaction struct {
	// Configurable fields
	batchSize  int
	goRoutines int
	insertType string
	// Internal fields
	InsertChan chan *generator.Record
	UpdateChan chan *generator.Record
	DeleteChan chan string
	WG         *sync.WaitGroup
}

func NewTransaction() *transaction {
	t := &transaction{
		// Configurable fields
		batchSize:  1000,
		goRoutines: 8,
		insertType: "batch",
		// Internal fields
		InsertChan: make(chan *generator.Record, 100),
		UpdateChan: make(chan *generator.Record, 100),
		DeleteChan: make(chan string, 100),
		WG:         &sync.WaitGroup{},
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
		_, err := pool.Exec(context.Background(), fmt.Sprintf("CREATE TABLE %s(UID text PRIMARY KEY, CLUSTER text, NAME text, DATA JSONB)", table))
		if err != nil {
			fmt.Println("Error creating table ", table, error)
		}
	}
}

// Initializes the connection pool.
func createPool() {
	database_url := "postgres://postgres:" + postgrePW + "@localhost:5432/" + databaseName
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
	for i := 0; i < t.goRoutines; i++ {
		if t.insertType == "batch" {
			go t.batchInsert(strconv.Itoa(i))
		} else {
			go t.copyInsert(strconv.Itoa(i))
		}
		go t.batchUpdate()
		go t.batchDelete()
	}
}
