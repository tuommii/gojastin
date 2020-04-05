## gojastin
![CI](https://github.com/tuommii/gojastin/workflows/CI/badge.svg)

Gojastin is a server measuring times between requests: https://timer.miikka.xyz/

>  Challenge where i got inspiration for this project was fun, but validation for solutions seemed far more interesting. So I decided try to do my own *Proof of Concept*. You can read challenge subject [here](https://github.com/hivehelsinki/remote-challs/tree/master/chall03). WIP. 

## Features
- [x] [Live demo](https://timer.miikka.xyz/)
- [x] Uses sync.Pool
- [x] Benchmark testing
- [x] Testing with awesome *httptest*
- [x] Github actions
- [x] No external libraries
- [x] Useful Makefile

## Bechmarks
| Function | Iterations | ns/op | B/op | allocs/op |
|---|--:|--:|--:|--:|
|[without sync.Loop](https://github.com/tuommii/gojastin/blob/21ad33431767dfb9b4c9a6d8b9f63c9f720f66e2/server/visitor.go#L29)| 5351934 | 218 | 48 | 1 |
|[with sync.Loop](https://github.com/tuommii/gojastin/blob/f9cdfa646ed7693d3210a17291abb2a0efd84886/server/visitor.go#L29)|  15951188 | 75 | 0 | 0|



More testes by

## Running it on local

Clone:

```
git clone https://github.com/tuommii/gojastin.git
```

Test:
```
make test
```

Benchmark:
```
make bench
```

Build and run:

```
make
```


Navigate to http://localhost:3030/



## Todo
- [ ] Test config for bad values
- [ ] Template for "result" and templates to file
- [x] sync.Pool (Might be premature optimization)
- [x] Basic styles
- [x] Testing
- [x] Test counter reset
- [x] Test visitors get removed
- [x] Test timerStop()
- [x] Easier way to change times

## Screenshot

![Screenshot](/pic.png)
