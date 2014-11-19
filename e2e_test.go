package crunch

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

func TestEnds(t *testing.T) {
	makeQuery := func(path string) func(r DataReader) (interface{}, error) {
		parts := strings.Split(path, ".")
		return func(r DataReader) (interface{}, error) { return r.Query(parts) }
	}

	transform := NewTransformer()
	row := NewRow()
	row.FieldWithDefault("ip", "", makeQuery("client.ip"), transform.AsIs)
	row.FieldWithDefault("host", "", makeQuery("client.host"), transform.AsIs)
	row.FieldWithDefault("retries", "", makeQuery("client.retries"), transform.AsIs)
	row.Field("lat", makeQuery("client.lat"), transform.AsIs)
	row.Field("apps_json", makeQuery("apps"), transform.AsJson)
	row.Field("width", makeQuery("screen.width"), transform.AsIs)
	row.Field("height", makeQuery("screen.height"), transform.AsIs)
	row.Field("debug", makeQuery("debug"), transform.AsIs)
	row.FieldWithDefault("ev_yesno:int", "0", makeQuery("action.yesno"), func(val interface{}) string {
		if val.(string) == "yes" {
			return "1"
		}
		return "0"
	})
	row.Feature("getting ip to location", []string{"country", "city"},
		func(r DataReader, row *Row) []string {
			return []string{"israel", "telaviv"}
		})

	file, err := os.Open("test/fixtures/e2e_sanity.in.json") // For read access.
	if err != nil {
		log.Fatalf("Cannot open file %v", err)
	}
	outbytes, err := ioutil.ReadFile("test/fixtures/e2e_sanity.out.tsv")

	expected := string(outbytes)

	if err != nil {
		log.Fatalf("Cannot open file %v", err)
	}

	var out bytes.Buffer
	runner := NewRunner()
	shouldStream := runner.HandleCliWithFlags(row)
	if shouldStream {
		runner.JsonRowProcessor(row, file, &out)
	}
	result := out.String()
	if result != expected {
		t.Fatalf("\nin : [%s]\nout: [%s]\n", result, expected)
	}

}
