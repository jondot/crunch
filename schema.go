package crunch

import(
  "strings"
)

type Schema struct{
}

func toPigField(field string) string{
  parts := strings.Split(field, " ")
  name := parts[0]
  ftype := "chararray"
  if len(parts) == 2 {
    switch parts[1]{
      case "string":
        ftype = "chararray"
      default:
        ftype = parts[1]
    }
  }

  return strings.Join([]string{name, ftype}, ":")
}

func toHiveField(field string) string{
  parts := strings.Split(field, " ")
  name := parts[0]
  ftype := "string"
  if len(parts) == 2 {
    ftype = parts[1]
  }
  return strings.Join([]string{name, ftype}, " ")
}

func makeSchema(row *Row, converter func(string)string) string {
  fields := []string{}
  for _, field := range(row.Fields){
    fields = append(fields, converter(field.Name))
  }
  return strings.Join(fields, ",\n")
}

func NewSchema() *Schema{
  return &Schema{}
}

func (self *Schema) Hive(row *Row) string{
  return makeSchema(row, toHiveField)
}

func (self *Schema) Pig(row *Row) string{
  return makeSchema(row, toPigField)
}



