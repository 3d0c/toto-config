# TOTO Configuration Server

## Disclaimer

To simplify testing, deployment and delivery the server uses SQLite as a database backend. It obviously locks server into single instance.  
To scale it into multiple instances SQLite should be replaced by any RDBMS supporting network access. 
It can be easily done by changing `DSN` and `dialect` parameters in config file.

## Benchmarks

To reproduce benchmark run `go test -bench=. -benchmem` inside `toto-config/pkg/apiserver/models` first. It will generate random data and stores it into `/tmp/models-test.db` which will be reused. (NOTE: no clean up here, the part of testing.Main which is in charge of cleanup is commented). 

Expected output: 

```sh
2022/12/05 13:19:39 goose: no migrations to run. current version: 1
Looking for: 2a950a41dbd5a8d429b79ad1d1abf31293e092589ed4e139dbb57463f883a033, 15, 79
goos: darwin
goarch: amd64
pkg: github.com/3d0c/toto-config/pkg/apiserver/models
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkFindBy-12         48506             24463 ns/op            1384 B/op         42 allocs/op
PASS
```

### Apache benchmarks

To reproduce it generate the database first by running `go test` inside `models/` directory

Here is the single curl request verbose output

```sh
~> curl -k -v -XGET -H "X-AppEngine-Country: 89" https://localhost:8443/v1/sku/f3db76df714300f52653e53de8f31e687d13a454d144f05f1460a24819ca5ac5                                             

> GET /v1/sku/f3db76df714300f52653e53de8f31e687d13a454d144f05f1460a24819ca5ac5 HTTP/1.1
> Host: localhost:8443
> User-Agent: curl/7.79.1
> Accept: */*
> X-AppEngine-Country: 89
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Sun, 04 Dec 2022 19:24:59 GMT
< Content-Length: 86
< Content-Type: text/plain; charset=utf-8
< 
{
    "main_sku": "7377fdf723eb4ed2a3a52bb869c162fdc621e373cbfe41dc37bfb26eb84c2031"
}
Total; 0.03784 (sec)
```
The result of Apache Benchmark of hitting the single endpoint:

```sh
ab -c100  -n1000 -k -H "X-AppEngine-Country: 89" https://localhost:8443/v1/sku/f3db76df714300f52653e53de8f31e687d13a454d144f05f1460a24819ca5ac5                                                                                    

Server Software:        
Server Hostname:        localhost
Server Port:            8443
SSL/TLS Protocol:       TLSv1.2,ECDHE-RSA-AES128-GCM-SHA256,2048,128
Server Temp Key:        ECDH X25519 253 bits
TLS Server Name:        localhost

Document Path:          /v1/sku/f3db76df714300f52653e53de8f31e687d13a454d144f05f1460a24819ca5ac5
Document Length:        86 bytes

Concurrency Level:      100
Time taken for tests:   0.291 seconds
Complete requests:      1000
Failed requests:        15
   (Connect: 0, Receive: 0, Length: 15, Exceptions: 0)
Non-2xx responses:      15
Keep-Alive requests:    1000
Total transferred:      225185 bytes
HTML transferred:       84710 bytes
Requests per second:    3431.29 [#/sec] (mean)
Time per request:       29.144 [ms] (mean)
Time per request:       0.291 [ms] (mean, across all concurrent requests)
Transfer rate:          754.56 [Kbytes/sec] received

```

## API

- `GET /v1/sku/{package}`  

  Results:  
  - 200 OK and valid JSON like:
  
    ```json
    {
        "main_sku": "7377fdf723eb4ed2a3a52bb869c162fdc621e373cbfe41dc37bfb26eb84c2031"
    }
    ```
    
  - 404 Not Found
  - 500 Internal server error

## Build and Run

Just run locally for testing purposes:

```sh
go run main.go run --config=configs/dev/local_config.yaml
```

Building:

```sh
# Do static analysis
make lint

# Build binary (will be in build/ directory)
make build
```

## Improvements

- The `ConfigScheme` data structure is hardcoded. For better integration and maintenance experience code generator might be using here, using JSON scheme format.
- GeoTarget using `mmdb` database can be detached into another service.

## Caveats

- `logger` doesn't create destination directory, so it should be done manually!
- TLS configuration is `insecure` by default
