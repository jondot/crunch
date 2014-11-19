package main

import (
	"../../crunch"
	"strings"
  "github.com/oschwald/geoip2-golang"
  "log"
  "net"
)

func makeQuery(path string) func(crunch.DataReader) (interface{}, error) {
	parts := strings.Split(path, ".")
	return func(r crunch.DataReader) (interface{}, error) { return r.Query(parts) }
}

func main() {
  db, err := geoip2.Open("GeoLite2-City.mmdb")
  if err != nil {
          log.Fatal(err)
  }
  defer db.Close()
  transform := crunch.NewTransformer()

	row := crunch.NewRow()
	row.FieldWithDefault("ip", "", makeQuery("x-forwarded-for"), transform.AsIs)
	row.FieldWithDefault("ev_ts", "", makeQuery("head.timestamp"), transform.AsIs)
	row.FieldWithDefault("ev_json", "", makeQuery("action"), transform.AsJson)
	row.FieldWithDefault("ev_yesno:int", "0", makeQuery("action.yesno"), func(val interface{}) string {
		if val.(string) == "yes" {
			return "1"
		}
		return "0"
	})
	row.FieldWithValue("ev_smp int", "1.0")
	row.FieldWithDefault("ev_source", "", makeQuery("action.source"), transform.AsIs)
	row.Feature("getting ip to location", []string{"country", "timezone"},
		func(r crunch.DataReader, row *crunch.Row) []string {
        ip := net.ParseIP(row.GetFieldValue("ip"))
        record, err := db.City(ip)
        if err != nil {
          log.Fatal(err)
          return []string{"N/A", "N/A"}
        }
			return []string{record.Country.IsoCode, record.Location.TimeZone}
		})

	crunch.ProcessJson(row)
}
