package redisgraph

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/jlpadilla/benchmark/pkg/generator"
	rg "github.com/redislabs/redisgraph-go"
)

// A global redis pool for other parts of this package to use
var Pool *redis.Pool

const (
	POOL_IDLE_TIMEOUT = 60
	POOL_MAX_IDLE     = 20
	POOL_MAX_ACTIVE   = 30
	GRAPH_NAME        = "benchmark"
)

// Transaction settings
type transaction struct {
	options    generator.Options
	InsertChan chan *generator.Record
	UpdateChan chan *generator.Record
	DeleteChan chan string
	WG         *sync.WaitGroup
}

func NewTransaction(options generator.Options) *transaction {
	t := &transaction{
		// Configurable fields
		options: options,
		// Internal fields
		InsertChan: make(chan *generator.Record, 100),
		UpdateChan: make(chan *generator.Record, 100),
		DeleteChan: make(chan string, 100),
		WG:         &sync.WaitGroup{},
	}

	t.startConnectors()
	return t
}

// Initializes the pool using functions in this file.
func init() {
	Pool = &redis.Pool{
		MaxIdle:      10,
		MaxActive:    20,
		Dial:         getRedisConnection,
		TestOnBorrow: testRedisConnection,
		Wait:         true,
	}
	InitializeDB()
}

// Initialize the database. Drop existing tables and create new tables for this test.
func InitializeDB() {
	conn := Pool.Get()
	defer conn.Close()

	g := rg.Graph{
		Conn: conn,
		Id:   GRAPH_NAME,
	}
	_, err := g.Query("MATCH (n) DELETE n")
	if err != nil {
		fmt.Println("Error initializing Redisgraph. ", err)
	}
}

func getRedisConnection() (redis.Conn, error) {
	redisConn, err := redis.Dial("tcp", net.JoinHostPort("localhost", "6379"))
	if err != nil {
		fmt.Println("Error connecting Redisgraph.  Original error: ", err)
		return nil, err
	}
	return redisConn, nil
}

// Used by the pool to test if redis connections are still okay. If they have been idle for less than a minute, just assumes they are okay. If not, calls PING.
func testRedisConnection(c redis.Conn, t time.Time) error {
	if time.Since(t) < POOL_IDLE_TIMEOUT*time.Second {
		return nil
	}
	_, err := c.Do("PING")
	return err
}

func (t *transaction) startConnectors() {
	for i := 0; i < t.options.GoRoutines; i++ {
		go t.batchInsert(strconv.Itoa(i))
		go t.batchUpdate()
		go t.batchDelete()
	}
	// for i := 0; i < t.goRoutines; i++ {
	// 	if t.insertType == "batch" {
	// 		go t.batchInsert(strconv.Itoa(i))
	// 	} else {
	// 		go t.copyInsert(strconv.Itoa(i))
	// 	}
	// 	go t.batchUpdate()
	// 	go t.batchDelete()
	// }
}
