package main

import (
	"../../crunch"
	"flag"
	"os"
	"strings"
)

func makeQuery(path string) func(crunch.DataReader) (interface{}, error) {
	parts := strings.Split(path, ".")
	return func(r crunch.DataReader) (interface{}, error) { return r.Query(parts) }
}

func main() {
	transform := crunch.NewTransformer()
	row := crunch.NewRow()
	row.FieldWithValue("ev_id", "xxx-id")
	row.FieldWithDefault("ev_tshd", "", makeQuery("head.timestamp"), transform.AsIs)
	row.FieldWithDefault("ev_ts", "", makeQuery("action.timestamp"), transform.AsIs)
	row.FieldWithValue("ev_smp int", "1.0")
	row.FieldWithDefault("ev_action", "", makeQuery("action.action"), transform.AsIs)
	row.Feature("my feature", []string{"one", "two", "three"},
		func(r crunch.DataReader, row *crunch.Row) []string {
			return []string{"1", "2", "3"}
		})

	//
	// Specify the kind of writer to use: TSV.
	//
	writer := crunch.NewTsvWriter(os.Stdout, row)

	//
	// Set up a runner with Crunch's flags, and add our custom flags
	//
	runner := crunch.NewRunner()
	runner.Flags()
	myflag := ""
	flag.StringVar(&myflag, "foobar", "", "Foo bar the baz.")
	flag.Parse()

	// Note: A more customized processor can be had by breaking up HandleCli further (calling Schema.Pig(row) or Schema.Hive(row), or runner.GenerateStubs(row))
	shouldStream := runner.HandleCli(row)

	if shouldStream {
		/*
		   Run a custom JSON processor (a plain text processor is also avail.).
		   - Use transform operations from Transformer such as: explode
		   - Plug in an explicit DataReader: MapReader
		   - Plug in our DataWriter from earlier on: TsvWriter
		   - Explicitly call `Row.Write` due to our own custom workflow
		*/
		runner.JsonCustomProcessor(os.Stdin, func(data map[string]interface{}) {
			exploded := transform.Explode(data, "actions", "action")
			for _, exp := range exploded {
				reader := crunch.NewMapReader(exp)
				row.Write(reader, writer)
			}
		})
	}
}
