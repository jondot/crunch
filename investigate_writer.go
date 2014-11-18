package crunch

import (
  "io"
  "fmt"
)


type InvestigateWriter struct{
  writer io.Writer
  row *Row
  index int
}


func NewInvestigateWriter(w io.Writer, row *Row) *InvestigateWriter{
  return &InvestigateWriter{ writer: w, row: row}
}

func (self *InvestigateWriter) Field(val string){
  f := self.row.Fields[self.index]
  io.WriteString(self.writer, fmt.Sprintf("- %s\n", f.Name))
  io.WriteString(self.writer, fmt.Sprintf("[%s]\n", val))
  self.index++
}

func (self *InvestigateWriter) End(){
  io.WriteString(self.writer, "----------\n")
  self.index = 0
}


