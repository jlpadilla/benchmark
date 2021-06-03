package generator

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v6"
)

type Record struct {
	UID        string
	Cluster    string
	Kind       string
	Name       string
	Properties map[string]interface{}
}

func Generate(instance string, numRecords int, insertChan chan *Record) {

	for i := 0; i <= numRecords; i++ {
		record := Record{
			UID:     fmt.Sprintf("id-%s-%d", instance, i),
			Name:    fmt.Sprintf("name-%d", i),
			Kind:    gofakeit.Color(),
			Cluster: gofakeit.City(),
			Properties: map[string]interface{}{
				// "name":       fmt.Sprintf("%s%d", "name-", i),
				// "_uid":       fmt.Sprintf("%s%d", "id-", i),
				"_rbac":   fmt.Sprintf("%s%d", "rbac-", i%50),
				"name":    fmt.Sprintf("name-%d", i),
				"counter": i,
				"number":  gofakeit.Number(1, 9999),
				"bool":    gofakeit.Bool(),
				"beer":    gofakeit.BeerName(),
				"car":     gofakeit.CarModel(),
				"color":   gofakeit.Color(),
				"city":    gofakeit.City(),
				"map":     map[string]string{"key1": "value1", "key2": "value2"},
				"list":    []string{"a", "b", "c"},
				// "property10": "value10",
				// "property11": "value11",
				// "property12": "value12",
				// "property13": "value13",
				// "property14": "value14",
				// "property15": "value15",
				// "property16": "value16",
				// "property17": "value17",
				// "property18": "value18",
				// "property19": "value19",
				// "property20": "value20",
			},
		}

		insertChan <- &record
	}
}
