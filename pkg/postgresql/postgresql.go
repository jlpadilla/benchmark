package postgresql

import (
	"context"
	"fmt"
	"os"

	pgx "github.com/jackc/pgx/v4"
	"github.com/jlpadilla/benchmark/pkg/generator"
)

func createConn() *pgx.Conn {
	database_url := "postgres://postgres:dev-pass!@localhost:5432/jorge-demo"
	// conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return conn
}

func init() {
	c := createConn()
	defer c.Close(context.Background())

	// Clear resources table
	_, error := c.Exec(context.Background(), "DROP TABLE resources")
	if error != nil {
		fmt.Println("Error dropping table RESOURCES. ", error)
	}
	_, err := c.Exec(context.Background(), "CREATE TABLE resources(UID text PRIMARY KEY, NAME text, KIND text, Cluster text)")
	if err != nil {
		fmt.Println("Error creating table RESOURCES.")
	}

}

func ProcessInsert(instance string, insertChan chan *generator.Record) {
	conn := createConn()
	batch := &pgx.Batch{}

	for {
		record := <-insertChan

		batch.Queue("insert into resources(UID,Cluster,Kind,NAME) values($1,$2,$3,$4)", record.UID, record.Cluster, record.Kind, record.Name)

		if batch.Len()%250 == 0 {
			// fmt.Println("Sending batch from instance: ", instance)
			fmt.Print(".")
			br := conn.SendBatch(context.Background(), batch)
			br.Close()
			batch = &pgx.Batch{}
		}
	}
}

func QueryRecord() {
	// conn := createConn()
	// defer conn.Close(context.Background())
	// var name string
	// var age int64
	// err = conn.QueryRow(context.Background(), "select name, age from company where id=$1", 1).Scan(&name, &age)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	// 	os.Exit(1)
	// }

	// fmt.Println(name, age)
}
