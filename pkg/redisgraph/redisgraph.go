package redisgraph

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/jlpadilla/benchmark/pkg/generator"
)

// A global redis pool for other parts of this package to use
var Pool *redis.Pool

const (
	POOL_IDLE_TIMEOUT = 60
	POOL_MAX_IDLE     = 10
	POOL_MAX_ACTIVE   = 20
	GRAPH_NAME        = "benchmark"
)

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

// Initializes the pool using functions in this file.
func init() {
	Pool = &redis.Pool{
		MaxIdle:      10,
		MaxActive:    20,
		Dial:         getRedisConnection,
		TestOnBorrow: testRedisConnection,
		Wait:         true,
	}

}

func getRedisConnection() (redis.Conn, error) {
	// var redisConn redis.Conn

	redisConn, err := redis.Dial("tcp", net.JoinHostPort("localhost", "6379"),
		redis.DialUseTLS(false)) // Set this to false when you want to connect to redis via SSH from local laptop
	if err != nil {
		fmt.Println("Error connecting redis using SSH.  Original error: ", err)
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
	for i := 0; i < t.goRoutines; i++ {
		go t.batchInsert(strconv.Itoa(i))
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
