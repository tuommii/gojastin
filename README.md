## gojastin
![CI](https://github.com/tuommii/gojastin/workflows/CI/badge.svg)

Gojastin is a server measuring times between requests: https://timer.miikka.xyz/


![Screenshot](/pic.png)


>  Challenge where i got inspiration for this project was fun but server side seemed far more interesting, so I decided try to do my own. Proof of Concept. WIP. Link to [challenge](https://github.com/hivehelsinki/remote-challs/tree/master/chall03)

## Features
- [x] Live demo
- [x] No external libraries
- [x] Testing with awesome *httptest*
- [x] Testing includes benchmarks
- [x] I Tried to keep code *clean and simple*


## Running it on local

Clone:

```
git clone https://github.com/tuommii/gojastin.git
```

Build and run:

```
make
```


Navigate to http://localhost:3030/


Test:
```
make test
```

Benchmark:
```
make bench
```


Without

?   	miikka.xyz/gojastin/cmd/gojastin	[no test files]
?   	miikka.xyz/gojastin/config	[no test files]
goos: linux
goarch: amd64
pkg: miikka.xyz/gojastin/server
BenchmarkStart
BenchmarkStart-4                 	 5351934	       218 ns/op	      48 B/op	       1 allocs/op
BenchmarkStartAndStop
BenchmarkStartAndStop-4          	 2285290	       522 ns/op	     148 B/op	       4 allocs/op
BenchmarkStartAndHalfStopped
BenchmarkStartAndHalfStopped-4   	 3116977	       381 ns/op	      98 B/op	       3 allocs/op
PASS
ok  	miikka.xyz/gojastin/server	4.706s


With
go test -v -run=XXX -bench . ./... -benchmem
?   	miikka.xyz/gojastin/cmd/gojastin	[no test files]
?   	miikka.xyz/gojastin/config	[no test files]
goos: linux
goarch: amd64
pkg: miikka.xyz/gojastin/server
BenchmarkAllocWithoutPool
BenchmarkAllocWithoutPool-4      	 1792692	       594 ns/op	     143 B/op	       1 allocs/op
BenchmarkAllocWithPool
BenchmarkAllocWithPool-4         	 3741678	       347 ns/op	      91 B/op	       0 allocs/op



BenchmarkStart
BenchmarkStart-4                 	16693008	        72.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkStartAndStop
BenchmarkStartAndStop-4          	 4631438	       256 ns/op	      52 B/op	       2 allocs/op
BenchmarkStartAndHalfStopped
BenchmarkStartAndHalfStopped-4   	 7048154	       185 ns/op	      26 B/op	       1 allocs/op
PASS
ok  	miikka.xyz/gojastin/server	7.717s




## Todo
- [ ] Visitor based rate limiting
- [ ] Test config for bad values
- [ ] Template for "result" and templates to file
- [x] [sync.Pool](https://developer20.com/using-sync-pool/index.html)
(Might be premature optimization)
- [x] Basic styles
- [x] Testing
- [x] Test counter reset
- [x] Test visitors get removed
- [x] Test timerStop()
- [x] Easier way to change times
