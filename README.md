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
