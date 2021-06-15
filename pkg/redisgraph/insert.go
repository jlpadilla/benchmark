package redisgraph

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"

	rg "github.com/redislabs/redisgraph-go"
)

func (t *transaction) batchInsert(instance string) {
	t.WG.Add(1)
	defer t.WG.Done()
	conn := Pool.Get()
	defer conn.Close()

	g := rg.Graph{
		Conn: conn,
		Id:   GRAPH_NAME,
	}

	resourceStrings := []string{}

	for {
		record, more := <-t.InsertChan

		if more {
			encodedProps, err := encodeProperties(record.Properties)
			if err != nil {
				fmt.Println("Cannot encode resource ", record.UID, ", excluding it from insertion: ", err)
				continue
			}
			propStrings := []string{}
			for k, v := range encodedProps {
				switch typed := v.(type) { // At this point it's either string or int64. Need to wrap in quotes if it's string
				case int64:
					propStrings = append(propStrings, fmt.Sprintf("%s:%d", k, typed)) // e.g. key>:<value>
				default:
					propStrings = append(propStrings, fmt.Sprintf("%s:'%s'", k, typed)) // e.g. <key>:'<value>'
				}
			}
			resource := fmt.Sprintf("(:%s {_uid:'%s', %s})", record.Properties["kind"], record.UID, strings.Join(propStrings, ", ")) // e.g. (:Pod {_uid: 'abc123', prop1:5, prop2:'cheese'})

			resourceStrings = append(resourceStrings, resource)
		}
		if len(resourceStrings) == t.batchSize || !more {
			q := fmt.Sprintf("%s %s", "CREATE", strings.Join(resourceStrings, ", "))
			_, err := g.Query(q)

			if err != nil {
				fmt.Println("error: ", err)
			}
			resourceStrings = []string{}
		}
		if !more {
			break
		}

	}

}

func encodeProperties(props map[string]interface{}) (map[string]interface{}, error) {
	res := make(map[string]interface{}, len(props))
	for k, v := range props {
		// Get all the rg props for this property.
		partial, err := encodeProperty(k, v)
		if err != nil { // if anything went wrong just log a warning and skip it
			// glog.Warning("Skipping property ", k, " on resource ", r.UID, ": ", err)
			continue
		}

		// Merge all the props that came out into the larger map.
		for pk, pv := range partial {
			res[pk] = pv
		}
	}

	if len(res) == 0 {
		return nil, errors.New("No valid redisgraph properties found")
	}
	return res, nil
}

// Outputs all the redisgraph properties that come out of a given property on a resource.
// Outputs exclusively in our supported types: string, []string, map[string]string, and int64.
func encodeProperty(key string, value interface{}) (map[string]interface{}, error) {

	// Sanitize value
	if value == nil || value == "" { // value == "" is false for anything not a string
		return nil, errors.New("Empty Value")
	}

	res := make(map[string]interface{})

	// Switch over all the default json.Unmarshal types. These are the only possible types that could be in the map. For each, we go through and convert to what we want them to be.
	// Useful doc regarding default types: https://golang.org/pkg/encoding/json/#Unmarshal
	switch typedVal := value.(type) {
	case string:
		if key == "kind" { // we lowercase the kind.
			res[key] = strings.ToLower(typedVal) //sanitizeValue(typedVal))
		} else {
			res[key] = typedVal //sanitizeValue(typedVal)
		}

	case []interface{}:
		// RedisGraph 1.0.15 doesn't support a list of properties. As a workaround to this limitation
		// we are encoding a list of values in a single string.
		elementStrings := make([]string, 0, len(typedVal))
		for _, e := range typedVal {
			elementString := fmt.Sprintf("%v", e)
			elementStrings = append(elementStrings, elementString)
		}
		sort.Strings(elementStrings)                  // Sotring to make comparisons more predictable
		res[key] = strings.Join(elementStrings, ", ") // sanitizeValue(strings.Join(elementStrings, ", ")) // e.g. val1, val2, val3

	case map[string]interface{}:
		// RedisGraph 1.0.15 doesn't support a list of properties. As a workaround to this limitation
		// we are encoding the labels in a single string.
		if key == "label" {
			labelStrings := make([]string, 0, len(typedVal))
			for key, value := range typedVal {
				labelString := fmt.Sprintf("%s=%s", key, value)
				labelStrings = append(labelStrings, labelString)
			}
			sort.Strings(labelStrings)                  // Sotring to make comparisons more predictable
			res[key] = strings.Join(labelStrings, "; ") // sanitizeValue(strings.Join(labelStrings, "; ")) // e.g. key1=val1; key2=val2; key3=val3
		}

	case int64:
		res[key] = typedVal
	case float64: // As of 4/15/2019 we don't have any numerical properties that aren't ints.
		res[key] = int64(typedVal)
	case bool: // as of 4/2/2019 redisgraph does not support bools so we convert to string.
		if typedVal {
			res[key] = "true"
		} else {
			res[key] = "false"
		}
	default:
		return nil, fmt.Errorf("Property type unsupported: %s %v", reflect.TypeOf(typedVal), typedVal)
	}

	if len(res) == 0 {
		return nil, errors.New("No valid redisgraph properties found")
	}

	return res, nil
}
