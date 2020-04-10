![Gopher](/soldering.svg)

![CI](https://github.com/tuommii/gojastin/workflows/CI/badge.svg)

## gojastin
Gojastin is a server for measuring times between requests: https://timer.miikka.xyz/

>  The Challenge where i got a inspiration for this project was fun, but validation for solutions seemed far more interesting. So I decided try to do my own *Proof of Concept*. You can read the subject [here](https://github.com/hivehelsinki/remote-challs/tree/master/chall03).

## Features
- [x] [Live demo](https://timer.miikka.xyz/)
- [x] Uses sync.Pool
- [x] Benchmark testing
- [x] Testing with awesome *httptest*
- [x] Github actions
- [x] No external libraries
- [x] Useful Makefile
- [x] Hosted on DigitalOcean behind nginx


## Bechmarks
| Function | Iterations | ns/op | B/op | allocs/op |
|---|--:|--:|--:|--:|
|[without sync.Loop](https://github.com/tuommii/gojastin/blob/02dbae4ad50f6fe8d68dd62a585b9e58bbc69760/server/visitor.go#L29)| 5351934 | 218 | 48 | 1 |
|[with sync.Loop](https://github.com/tuommii/gojastin/blob/21ad33431767dfb9b4c9a6d8b9f63c9f720f66e2/server/visitor.go#L29)|  15951188 | 75 | 0 | 0|



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


## Screenshot

![Screenshot](/pic.png)
