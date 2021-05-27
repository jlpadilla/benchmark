package postgresql

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

func Start(numRecords int) {
	fmt.Println("Profiling postgresql. Records: ", numRecords)

	database_url := "postgres://postgres:dev-pass!@localhost:5432/jorge-demo"
	// conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var name string
	var age int64
	err = conn.QueryRow(context.Background(), "select name, age from company where id=$1", 1).Scan(&name, &age)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(name, age)
}
