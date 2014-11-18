package crunch

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type Transformer struct {
}

func NewTransformer() *Transformer {
	return &Transformer{}
}

func (self *Transformer) AsIs(data interface{}) string {
	return fmt.Sprintf("%v", data)
}

func (self *Transformer) FromCsv(line string) []string {
	arr, err := csv.NewReader(strings.NewReader(line)).Read()
	if err != nil {
		log.Fatalf("Error parsing CSV line: [%s].\nError: %s.", line, err)
	}
	return arr
}

func (self *Transformer) FromJson(line string) map[string]interface{} {
	data := map[string]interface{}{}
	dec := json.NewDecoder(strings.NewReader(line))
	dec.UseNumber()
	dec.Decode(&data)
	return data
}

func (self *Transformer) AsJson(data interface{}) string {
	b, err := json.Marshal(data)
	if err != nil {
		panic("cant marshall")
	}
	return string(b)
}

func Clone(src map[string]interface{}) map[string]interface{} {
	dst := map[string]interface{}{}
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func (self *Transformer) Explode(data map[string]interface{}, fromArray string, asItem string) []map[string]interface{} {
	exploded := []map[string]interface{}{}

	// cast the array field into array of objects
	for _, item := range data[fromArray].([]interface{}) {
		obj := Clone(data)
		delete(obj, fromArray)
		obj[asItem] = item
		exploded = append(exploded, obj)
	}

	return exploded
}
