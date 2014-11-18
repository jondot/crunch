package main
import(
  "../../crunch"
)

func makeQuery(rownumber int) func(crunch.DataReader)(interface{},error){
  return func(r crunch.DataReader)(interface{}, error){ return r.QueryIndex(rownumber) }
}


func main(){
  transform := crunch.NewTransformer()
  row := crunch.NewRow()
  row.FieldWithDefault("ev_tshd", "", makeQuery(0), transform.AsIs)
  row.FieldWithDefault("ev_ts", "", makeQuery(1), transform.AsIs)
  row.FieldWithDefault("ev_json", "", makeQuery(2), transform.AsJson)
  row.FieldWithDefault("ev_yesno:int", "0", makeQuery(3), func(val interface{})string{
    if val.(string) == "yes" {
      return "1"
    }
    return "0"
  })
  row.FieldWithValue("ev_smp int", "1.0")
  row.FieldWithDefault("ev_source", "", makeQuery(4), transform.AsIs)
  row.Feature("getting ip to location", []string{"country", "city"},
    func(r crunch.DataReader, row *crunch.Row)[]string{
      return []string{ "israel", "telaviv" }
    })

  crunch.ProcessCsv(row)
}
