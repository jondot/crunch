package main
import(
  "strings"
  "../../crunch"
)

func makeQuery(path string) func(crunch.DataReader)(interface{},error){
  parts := strings.Split(path, ".")
  return func(r crunch.DataReader)(interface{}, error){ return r.Query(parts) }
}


func main(){
  transform := crunch.NewTransformer()
  row := crunch.NewRow()
  row.FieldWithDefault("ev_tshd", "", makeQuery("head.timestamp"), transform.AsIs)
  row.FieldWithDefault("ev_ts", "", makeQuery("action.timestamp"), transform.AsIs)
  row.FieldWithDefault("ev_json", "", makeQuery("action"), transform.AsJson)
  row.FieldWithDefault("ev_yesno:int", "0", makeQuery("action.yesno"), func(val interface{})string{
    if val.(string) == "yes" {
      return "1"
    }
    return "0"
  })
  row.FieldWithValue("ev_smp int", "1.0")
  row.FieldWithDefault("ev_source", "", makeQuery("action.source"), transform.AsIs)
  row.Feature("getting ip to location", []string{"country", "city"},
    func(r crunch.DataReader, row *crunch.Row)[]string{
      return []string{ "israel", "telaviv" }
    })

  crunch.ProcessJson(row)
}
