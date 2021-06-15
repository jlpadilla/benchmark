package server

import (
	"net/http"
	"strconv"
)

type options struct {
	Database string
	Insert   int
	Update   int
	Delete   int
}

func ParseQuery(req *http.Request) options {
	insert := 1000 // Default to 1000 records.
	insertQuery := req.URL.Query()["insert"]
	if len(insertQuery) > 0 {
		insert, _ = strconv.Atoi(req.URL.Query()["insert"][0])
	}
	update := 0
	updateQuery := req.URL.Query()["update"]
	if len(updateQuery) > 0 {
		update, _ = strconv.Atoi(req.URL.Query()["update"][0])
	}
	delete := 0
	deleteQuery := req.URL.Query()["delete"]
	if len(deleteQuery) > 0 {
		delete, _ = strconv.Atoi(req.URL.Query()["delete"][0])
	}

	database := "postgresql"
	databaseQuery := req.URL.Query()["database"]
	if len(databaseQuery) > 0 {
		database = string(req.URL.Query()["database"][0])
	}

	return options{
		Database: database,
		Insert:   insert,
		Update:   update,
		Delete:   delete,
	}
}
