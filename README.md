## gojastin
![CI](https://github.com/tuommii/gojastin/workflows/CI/badge.svg)

Gojastin is a server measuring times between requests: http://timer.miikka.xyz/


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

## Todo
- [ ] Visitor based rate limiting
- [ ] Test config for bad values
- [ ] Template for "result" and templates to file
- [ ] [sync.Pool](https://developer20.com/using-sync-pool/index.html)
(Might be premature optimization)
- [x] Basic styles
- [x] Testing
- [x] Test counter reset
- [x] Test visitors get removed
- [x] Test timerStop()
- [x] Easier way to change times
