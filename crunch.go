package crunch

import (
	"os"
)

func ProcessCsv(row *Row) {
	runner := NewRunner()
	shouldStream := runner.HandleCliWithFlags(row)
	if shouldStream {
		runner.CsvRowProcessor(row, os.Stdin, os.Stdout)
	}
}

func ProcessJson(row *Row) {
	runner := NewRunner()
	shouldStream := runner.HandleCliWithFlags(row)
	if shouldStream {
		runner.JsonRowProcessor(row, os.Stdin, os.Stdout)
	}
}
