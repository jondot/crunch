package crunch

type DataReader interface {
	Query([]string) (interface{}, error)
	QueryIndex(i int) (interface{}, error)
}

type DataWriter interface {
	Field(string)
	End()
}

type CrunchOpts struct {
	CpuProfile       string
	PigTemplateFile  string
	HiveTemplateFile string
	StubsOutPath     string
	Investigate      bool
}
