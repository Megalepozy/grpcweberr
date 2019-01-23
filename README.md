### Package for easy creation & reading of GRPC errors with added data for web clients

GRPCWebError was built to solve the problem of how to transfer descriptive errors to web clients through microservices 
which communicate using GRPC, in order to achieve that there are 2 steps which need to be taken:

The first step is where an error is captured (source of error can be from a system/user/3rd side library whatever...), there the New 
func should be used to create the wanted error (mind u that the resulting error is just a plain error which avoid the 
need to do any type inference later), for example:

```go
if 1 != 0 {
    return grpcweberr.New("lcaxafdv23d9d2s", codes.InvalidArgument, 422, "Received invalid values")
}
```

While the error is being transferred through the system it can just be handled like a simple error, usually just return it.
  
When its time to prepare a response to the web client the following functions can be used to get the needed data:
```go
errID := grpcweberr.GetErrorID(err)
httpStatus := grpcweberr.GetHTTPStatus(err)
errMsg := grpcweberr.GetUserErrorMessage(err)
```

There is also a way to track and log specific errors which travel through the microservices, for that, first append a 
logTracingID to the error:
```go
if 1 != 0 {
    err := grpcweberr.New("lcaxafdv23d9d2s", codes.InvalidArgument, 422, "Received invalid values")
    return grpcweberr.AddLogTracingID(tracingID, err)
}
```
Where tracingID is some sort of error id (string), I use it to track errors which I log at StackDriver, I created for that 
another pkg [sdlog](https://github.com/Megalepozy/sdlog) which is not documented... notify me if u want me to document it.

Later at specific points where u may want to log an error u can do:
```go
if logID := grpcweberr.GetLogTracingID(err); logID != "" {
    // some logging code... here I use sdlog but its up to u
    sdlog.New().Info("Error of Info lvl from X service", sdlog.AddLogTracingID(logID), sdlog.Lbl("err", err))
}
```
