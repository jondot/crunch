
![](/media/logo.png)

A fast to iterate, fast to run, Go based toolkit for ETL and feature extraction on Hadoop.

Use [crunch-starter](https://github.com/jondot/crunch-starter) for a boilerplate project to kickstart a production
setup.


## Quick Start

Crunch is optimized to be a big-bang-for-the-buck libary, yet almost
every aspect is extensible.

Let's say you have a log of semi-structured and deeply nested JSON. Each
line contains a record.

You would like to:

1. Parse JSON records
2. Extract fields
3. Optionally cleanup/process fields
4. Optionally extract features - run custom code on field values and
   output the result as new field(s)

![](/media/crunch.gif)


So here's a detailed view:

```go
// Describe your row
transform := crunch.NewTransformer()
row := crunch.NewRow()
// Use "field_name type". Types are Hive types.
row.FieldWithValue("ev_smp int", "1.0")
// If no type given, assume 'string'
row.FieldWithDefault("ip", "0.0.0.0", makeQuery("head.x-forwarded-for"), transform.AsIs)
row.FieldWithDefault("ev_ts", "", makeQuery("action.timestamp"), transform.AsIs)
row.FieldWithDefault("ev_source", "", makeQuery("action.source"), transform.AsIs)
row.Feature("doing ip to location", []string{"country", "city"},
  func(r crunch.DataReader, row *crunch.Row)[]string{
    // call your "standard" Go code for doing ip2location
    return ip2location(row["ip"])
  })

// By default, will build a hadoop-compatible streamer process that understands json: (stdin[JSON] to stdout[TSV])
// Also will plug-in Crunch's CLI utility functions (use -help)
crunch.ProcessJson(row)
```

Build your processor

```
$ go build my_processor.go
```

Generate a Pig driver that uses `my_processor`, and a Hive table
creation DDL.

```
$ ./my_processor -crunch.stubs="."
```

You can now ship your binary and scripts (crunch.hql, crunch.pig) to
your cluster.

In your cluster, you can now setup your table with Hive and run an ETL job with Pig:

```
$ hive -f crunch.hql
$ pig -stop_on_failure --param inurl=s3://inbucket/logs/dt=20140304 --param outurl=s3://outbucket/success/dt=20140304 crunch.pig
```

## Row Setup

The row setup is the most important part of the processor.

Make a row:

```go
transform := crunch.NewTransformer()
row := crunch.NewRow()
```

And start describing fields in it:

```Go
row.FieldWithDefault("name type", "default-value", <lookup function>, <transform function>)
```

A field description is:

* A `name type` pair, where types are Hive types.
* A default value (for `FieldWithDefault`, there are variants of this -- see the API docs)
* A lookup function (the 'Extract' part of ETL) - see one in the
  example processor. It outputs an `interface{}`
* A transform function, which eventually should represent that
  `interface{}` as a string type but its contents can changed based on semantics (JSON, int values, dates, etc).



## The Processor
Crunch comes with a built in processor rig, that packs its API into
a ready-made processor:

```go
crunch.ProcessJson(row)
```
This processor reads JSON and outputs Hadoop-streaming TSV that is compatible with [Pig STREAM](https://pig.apache.org/docs/r0.11.1/basic.html#STREAM) (which we use later), based on your row description and functions.

It also injects the following commands into your binary:

```
$ ./simple_processor -help
Usage of ./simple_processor:
  -crunch.cpuprofile="": Turn on CPU profiling and write to the specified file.
  -crunch.hivetemplate="": Custom Hive template for stub generation.
  -crunch.pigtemplate="": Custom Pig template for stub generation.
  -crunch.stubs="": Generate stubs and output to given path, and exit.
```

## Building a binary

Since go packs all dependencies into your binary, this makes a great
delivery package to hadoop.

Simply take a starter processor from `/examples` and build your processor based on it. Then build it:

```
$ go build simple_processor.go
$ ./simple_processor -crunch.stubs="."
Generated crunch.pig
Generated crunch.hql
```

The resulting binary should be ready for action, using Pig (see next
section)

## Generating Pig and Hive stubs

Crunch injects useful commands into your processor, one of them supports
script generation to create your Hive table, and your Pig job.

```
$ ./simple_processor -crunch.stubs="."
Generated crunch.pig
Generated crunch.hql
```

You can use your own templates with the `-crunch.hivetemplate` and `-crunch.pigtemplate` flags, as long as you include a `%%schema%%` (and `%%process%%` for the pig script) special pragma so that Crunch will replace it with the actual Pig or Hive schema.

## Extending Crunch

[this section is WIP]

Crunch is packaged into use-cases accessible from the crunch package, `crunch.ProcessJson` to name one.

However beneath the usecase facade, lies an extensible API which lets
you have any kind of granularity over using Crunch.

Some detailed examples can be seen in `/examples/detailed_processor.go`.



