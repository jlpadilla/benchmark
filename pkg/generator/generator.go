package generator

import (
	"fmt"
	"math/rand"
	"reflect"
	"sync"

	"github.com/brianvoe/gofakeit/v6"
)

type Record struct {
	UID        string
	Cluster    string
	Name       string
	Properties map[string]interface{}
}

// Keep track of existing record UID
var recordsMap = make(map[string]bool)
var recordCounter = 0
var mux sync.Mutex

func newRecord(index int, uid string) Record {
	mux.Lock()
	defer mux.Unlock()
	if uid == "" {
		uid = fmt.Sprintf("id-%d", recordCounter)
		recordCounter++
	}

	return Record{
		UID:     uid,
		Name:    fmt.Sprintf("name-%d", index),
		Cluster: fmt.Sprintf("cluster-%d", gofakeit.Number(0, 100)),
		Properties: map[string]interface{}{
			"name":    fmt.Sprintf("name-%d", index),
			"kind":    gofakeit.Color(),
			"counter": index,
			"number":  gofakeit.Number(1, 9999),
			"bool":    gofakeit.Bool(),
			"beer":    gofakeit.BeerName(),
			"car":     gofakeit.CarModel(),
			"color":   gofakeit.Color(),
			"city":    gofakeit.City(),
			// "map":     map[string]string{"key1": "value1", "key2": "value2"},
			// "list":    []string{"a", "b", "c"},
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
}

func Generate(opts Options, sim Simulation) {
	var wg sync.WaitGroup
	wg.Add(3)
	go addRecords(opts.Insert, sim.InsertChan, &wg)
	go updateRecords(opts.Update, sim.UpdateChan, &wg)
	go deleteRecords(opts.Delete, sim.DeleteChan, &wg)

	wg.Wait()

	close(sim.InsertChan)
	close(sim.UpdateChan)
	close(sim.DeleteChan)
}

func addRecords(insert int, insertChan chan *Record, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < insert; i++ {
		record := newRecord(i, "")
		mux.Lock()
		recordsMap[record.UID] = true
		mux.Unlock()
		insertChan <- &record
	}
}

func updateRecords(update int, updateChan chan *Record, wg *sync.WaitGroup) {
	defer wg.Done()
	if update > len(recordsMap) {
		fmt.Println("Can't update. Not enough records.")
		return
	}
	for i := 0; i < update; i++ {
		recordID := getRandomRecordID()
		record := newRecord(i, recordID)
		updateChan <- &record
	}
}

func deleteRecords(del int, deleteChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	if del > len(recordsMap) {
		fmt.Println("Can't delete. Not enough records.")
		return
	}
	for i := 0; i < del; i++ {
		recordID := getRandomRecordID()
		mux.Lock()
		delete(recordsMap, recordID)
		mux.Unlock()
		deleteChan <- recordID
	}
}

func getRandomRecordID() string {
	mux.Lock()
	defer mux.Unlock()
	records := reflect.ValueOf(recordsMap).MapKeys()
	return records[rand.Intn(len(recordsMap))].Interface().(string)
}
