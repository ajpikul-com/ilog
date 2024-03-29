# ilog

Notes: Some changes were made. We now have Warn and Debug level interfaces that can be turned on and off.
In simple logger, it's a Level(DEBUG|INFO|WARN|ERROR|NONE) function, in zap, it's a structure member Level
I realized in zap we can't use most of the sugar logger interface with this structure, and I'm leaning towards just using
zap entirely. Also, tests don't implement any of this yet.

I generally like it, but I want better stack trace processing. Also, need to set levels in consistent way.

`ilog` provides an interface (LoggerInterface) that specifies three functions:
* Init() error
* Info(string)
* Error(string)
* Debug(string)
* Warn(string)

You must choose (or write) a backend:

* EmptyLogger- Logging is off in this case
```
defaultLogger = &ilog.EmptyLogger{}
defaultLogger.Init()
```
* SimpleLogger- Unbuffered writes to stderr or another file descriptor
```
defaultLogger = &ilog.SimpleLogger{} 
defaultLogger.Init()
defaultLogger.Level(ilog.DEBUG|INFO|WARN|ERROR)
```
* ZapWrap- A backend of uber-go/zap
```
defaultLogger = &ilog.ZapWrap{}
type ilog.ZapWrap struct { 
 Sugar bool // If true support printf and key-value pairs
 Paths []string // Out, default is stderr
 Level int // One of ilog.DEBUG|INFO|WARN|ERROR
 ZapLogger *zap.Logger // Access after init
 SugarLogger *zap.SugaredLogger // Access after init
}
```
* TestLogger- A backend that calls t.Log or b.Log (has seperate constructor which accepts as an argument the test being run)

Use the `ilog_test.go` file as an example of how to set up loggers. Or TODO: write readme for each one

## Use

### In Application Example
```
newLogger := new(ilog.ZapWrap)
err := newLogger.Init()
if err != nil {
	panic(err)
}
myModule.SetDefaultLogger(newLogger)
```

### In-Module Example
```
// Set a global logger for the library
var defaultLogger ilog.LoggerInterface

// Establish a default logger
func init() {
	if defaulLogger == nil {
		defaultLogger = new(ilog.EmptyLogger)
	}
}

// Allow calling program to change default logger
func SetDefaultLogger(newLogger ilog.LoggerInterface) {
	defaultLogger = newLogger
	defaultLogger.Info("Default Logger Set")
}
```

## Benchmarks


| Name                                            |Iterations   |Speed			 |Memory	 |Allocs      |
|:----------------------------------------------- | -----------:| ----------:| -------:| ----------:|
| BenchmarkLogger/Benchmark_empty_logger					|2000000000	  |0.64 ns/op	 |0 B/op	 |0 allocs/op	|
| BenchmarkLogger/Benchmark_simple_logger					|1000000			|1166 ns/op	 |0 B/op	 |0 allocs/op	|
| BenchmarkLogger/Benchmark_zap_production_logger	|5000000			|308 ns/op	 |2 B/op	 |0 allocs/op	|
| BenchmarkLogger/Benchmark_zap_sugared_logger		|2000000			|611 ns/op	 |50 B/op	 |2 allocs/op	|

