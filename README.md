### Package for easy creation & reading of GRPC errors with added data for web clients

GRPCWebError was built to solve the problem of how to transfer descriptive errors to web clients through microservices 
which communicate using GRPC, in order to achieve that there are 2 steps which need to be taken:

The first step is where an error is created:

```go
if 1 != 0 {
    gwe := grpcweberr.New()
    return gwe.New(codes.InvalidArgument, 422, "Received invalid values")
}
```

While the error is being transferred through the system it can just be handled like a simple error.
  
When its time to prepare a response to the web client the following functions can be used to retrieve the error data:
```go
gwe := grpcweberr.New()
httpStatus := gwe.GetHTTPStatus(err)
errMsg := gwe.GetUserErrorMessage(err)
```

There is also a way to track and log specific errors which travel through the microservices, for that, first append a 
logTracingID to the error:
```go
if 1 != 0 {
    gwe := grpcweberr.New()
    err := gwe.New(codes.InvalidArgument, 422, "Received invalid values")
    return gwe.AddLogTracingID(tracingID, err)
}
```
Where tracingID is some string which can be used to track errors across the microservices. 

Later at specific points where u may want to log an error:
```go
gwe := grpcweberr.New()
if logID := gwe.GetLogTracingID(err); logID != "" {
    // some logging code... here I use sdlog but its up to u
    sdlog.New().Info("Error of Info lvl from X service", sdlog.AddLogTracingID(logID), sdlog.Lbl("err", err))
}
```
