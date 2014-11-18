package crunch

import (
  "errors"
)


type MapReader struct{
  jq *JsonQuery
}

func NewMapReader(data map[string]interface{}) *MapReader{
  return &MapReader{jq: NewQuery(data)}
}

func (self *MapReader) QueryIndex(i int) (interface{}, error){
  return nil, errors.New("Index not supported with MapReader, use Query")
}

func (self *MapReader) Query(parts []string) (interface{}, error){
  i, err := self.jq.Interface(parts...)
  if err!=nil {
    return nil, err
  }

  return i, nil
}

