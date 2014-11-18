package crunch

type field struct {
  Name string
  Default string
  Value string
  Extract func(r DataReader)(interface{}, error)
  Transform func(interface{})string
  HasDefault bool
  HasValue bool
  IsComputed bool
}
func (self *field) FetchValue(r DataReader) (string, error){
    if self.HasValue {
      return self.Value, nil
    }

    extracted, err := self.Extract(r)
    if err != nil {
      if self.HasDefault {
        return self.Default, nil
      }
      return "", err
    }
    return self.Transform(extracted), nil
}

type feature struct {
  Description string
  OutFields []string
  Process func(DataReader, *Row)[]string
}

type Row struct {
  Fields []field
  features []feature
}

func NewRow() *Row{
  return &Row{}
}

func (self *Row) Size() int {
  return len(self.Fields)
}

func (self *Row) Write(r DataReader, w DataWriter) error {
  for _, field := range self.Fields {
    if field.IsComputed {
      continue
    }

    val, err := field.FetchValue(r)
    if err != nil{
      return err
    }

    w.Field(val)
  }

  // materialize all computed field by features
  computedCache := map[string]string {}

  for _, feature := range self.features {
    results := feature.Process(r, self)
    for i, result := range results {
      computedCache[feature.OutFields[i]] = result
    }
  }


  // hydrate computed fields
  // pull all values from computed fields
  for _, field := range self.Fields {
    if !field.IsComputed {
      continue
    }

    val := computedCache[field.Name]
    w.Field(val)
  }

  w.End()
  return nil
}



func (self *Row) Field(decl string, extract func(r DataReader)(interface{}, error), transform func(interface{})(string)){
  self.Fields = append(self.Fields, field{
    Name:decl,
    Extract: extract,
    Transform: transform,
  })
}

func (self *Row) FieldWithDefault(decl string, defval string, extract func(r DataReader)(interface{}, error),  transform func(interface{})(string)){
  self.Fields = append(self.Fields, field{
    Name:decl,
    Extract: extract,
    Transform: transform,
    Default: defval,
    HasDefault:true,
  })
}

func (self *Row) FieldWithValue(decl string, val string){
  self.Fields = append(self.Fields, field{
    Name:decl,
    Value: val,
    HasValue: true,
  })
}

func (self *Row) Feature(desc string, decls []string, process func(DataReader, *Row)[]string){
  for _, decl := range decls{
    self.Fields = append(self.Fields, field{
      Name: decl,
      IsComputed: true,
    })
  }
  self.features = append(self.features, feature{
    Description: desc,
    OutFields: decls,
    Process: process,
  })
}

