# ilog

`ilog` provides an interface (LoggerInterface) that specifies three functions:
* Init() error
* Info(string)
* Error(string)

This repository implements three backends for this interface:

* EmptyLogger- Logging is off in this case
* SimpleLogger- Unbuffered writes to stderr or another file descriptor
* ZapWrap- A backend of uber-go/zap
* TestLogger- A backend that calls t.Log or b.Log (needs to be initialized for every test)

You can have a global default LoggerInterface or a LoggerInterface specific to some object, and set the backend to EmptyLogger for sub-nanosecond cancellation of logging, or whatever you want.

Use the `ilog_test.go` file as an example of how to set up loggers.

## Benchmarks


| Name                                            |Iterations   |Speed			 |Memory	 |Allocs      |
|:----------------------------------------------- | -----------:| ----------:| -------:| ----------:|
| BenchmarkLogger/Benchmark_empty_logger					|2000000000	  |0.64 ns/op	 |0 B/op	 |0 allocs/op	|
| BenchmarkLogger/Benchmark_simple_logger					|1000000			|1166 ns/op	 |0 B/op	 |0 allocs/op	|
| BenchmarkLogger/Benchmark_zap_production_logger	|5000000			|308 ns/op	 |2 B/op	 |0 allocs/op	|
| BenchmarkLogger/Benchmark_zap_sugared_logger		|2000000			|611 ns/op	 |50 B/op	 |2 allocs/op	|
