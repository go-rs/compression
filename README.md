# Response Payload Compressor
Added gzip compression mechanism for response payload.

## How to use?

```go
var api rest.API

// rest.API pointer
// use value between level -2 to 9 OR gzip constants
h := compression.Handler(&api, gzip.BestCompression)

http.ListenAndServe(":8080", h)
```