package crunch

import(
  "io"
  "io/ioutil"
  "bufio"
  "log"
  "strings"
  "path"
  "runtime/pprof"
  "os"
  "flag"
)

type Runner struct {
  Opts CrunchOpts
}


func readFile(path string, defaultContent string) string{
  if path != "" {
    content, err := ioutil.ReadFile(path)
    if err != nil {
      log.Fatalf("Cannot find hive template: %s.\nError: %s", path, err)
    }
    return string(content)
  }
  return defaultContent
}

func startCPUProfileIfEnabled(path string){
  if path != "" {
    f, err := os.Create(path)
    if err != nil {
        log.Fatal(err)
    }
    pprof.StartCPUProfile(f)
  }
}

func NewRunner() *Runner {
  return &Runner{ Opts: CrunchOpts{} }
}


func NewRunnerWithOpts(opts CrunchOpts) *Runner {
  return &Runner{ Opts: opts }
}

func (self *Runner) Flags(){
  flag.StringVar(&self.Opts.CpuProfile, "crunch.cpuprofile", "", "Turn on CPU profiling and write to the specified file.")
  flag.StringVar(&self.Opts.PigTemplateFile, "crunch.pigtemplate", "", "Custom Pig template for stub generation.")
  flag.StringVar(&self.Opts.HiveTemplateFile, "crunch.hivetemplate", "", "Custom Hive template for stub generation.")
  flag.StringVar(&self.Opts.StubsOutPath, "crunch.stubs", "", "Generate stubs and output to given path, and exit.")
  flag.BoolVar(&self.Opts.Investigate, "crunch.investigate", false, "Investigation mode - easy output to see processing result.")
}

func (self *Runner) HandleCliWithFlags(row *Row) bool{
  self.Flags()
  flag.Parse()
  return self.HandleCli(row)
}

func (self *Runner) HandleCli(row *Row) bool{
  if self.Opts.StubsOutPath != "" {
    self.GenerateStubs(row)
    return false
  }
  return true
}


func (self *Runner) GenerateStubs(row *Row){
  hiveTempl := readFile(self.Opts.HiveTemplateFile, TMPL_HIVE)
  pigTempl := readFile(self.Opts.PigTemplateFile, TMPL_PIG)
  outPath := "."
  if self.Opts.StubsOutPath != "" {
    outPath = self.Opts.StubsOutPath
  }

  schema := NewSchema()
  hiveTempl = strings.Replace(hiveTempl, "%%schema%%", schema.Hive(row), -1)
  pigTempl = strings.Replace(pigTempl, "%%schema%%", schema.Pig(row), -1)
  pigTempl = strings.Replace(pigTempl, "%%process%%", path.Base(os.Args[0]), -1)

  pigout := path.Join(outPath, "crunch.pig")
  ioutil.WriteFile(pigout, []byte(pigTempl), 0666)
  log.Printf("Generated: %s", pigout)

  hiveout := path.Join(outPath, "crunch.hql")
  log.Printf("Generated: %s", hiveout)
  ioutil.WriteFile(hiveout, []byte(hiveTempl), 0666)
}

func (self *Runner) JsonRowProcessor(row *Row, in io.Reader, out io.Writer){
  startCPUProfileIfEnabled(self.Opts.CpuProfile)
  if self.Opts.CpuProfile != "" {
    defer pprof.StopCPUProfile()
  }

  transform := NewTransformer()

  var writer DataWriter
  writer = NewTsvWriter(out, row)
  if self.Opts.Investigate {
    writer = NewInvestigateWriter(out, row)
  }

  stdin := bufio.NewReader(in)
  for {
    line, err := stdin.ReadString('\n')
    if err !=nil{
      return
    }
    data := transform.FromJson(line)
    reader := NewMapReader(data)
    row.Write(reader, writer)
  }
}


func (self *Runner) CsvRowProcessor(row *Row, in io.Reader, out io.Writer){
  startCPUProfileIfEnabled(self.Opts.CpuProfile)
  if self.Opts.CpuProfile != "" {
    defer pprof.StopCPUProfile()
  }

  transform := NewTransformer()
  var writer DataWriter
  writer = NewTsvWriter(out, row)
  if self.Opts.Investigate {
    writer = NewInvestigateWriter(out, row)
  }

  stdin := bufio.NewReader(in)
  for {
    line, err := stdin.ReadString('\n')
    if err !=nil{
      return
    }
    data := transform.FromCsv(line)
    reader := NewArrayReader(data)
    row.Write(reader, writer)
  }
}



