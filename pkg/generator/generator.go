package generator

import (
	"fmt"
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
			Name:    fmt.Sprintf("%s%d", "name-", i),
			Kind:    "my-kind",
			Cluster: "my-cluster",
			Properties: map[string]interface{}{
				// "name":       fmt.Sprintf("%s%d", "name-", i),
				// "_uid":       fmt.Sprintf("%s%d", "id-", i),
				"_rbac":      fmt.Sprintf("%s%d", "rbac-", i%50),
				"property0":  "value0",
				"property1":  "value1",
				"property2":  "value2",
				"property3":  "value3",
				"property4":  "value4",
				"property5":  "value5",
				"property6":  "value6",
				"property7":  "value7",
				"property8":  "value8",
				"property9":  "value9",
				"property10": "value10",
				"property11": "value11",
				"property12": "value12",
				"property13": "value13",
				"property14": "value14",
				"property15": "value15",
				"property16": "value16",
				"property17": "value17",
				"property18": "value18",
				"property19": "value19",
				"property20": "value20",
			},
		}

		insertChan <- &record
	}
}
