# gojastin
![Build & Test Go](https://github.com/tuommii/gojastin/workflows/Build%20&%20Test%20Go/badge.svg?branch=master)

Gojastin is a server measuring times between requests. Live: http://timer.miikka.xyz/

[Challenge](https://github.com/hivehelsinki/remote-challs/tree/master/chall03) where i got inspiration for this project was fun but server side seemed far more interesting so I decided try to do my own. Proof of Concept. WIP.

## Run
`make`

## Test
`make test`

### Benchmark
```make bench```

## Todo
- [x] Basic styles
- [ ] Visitor based rate limiting
- [ ] [sync.Pool](https://developer20.com/using-sync-pool/index.html)
- [ ] Template for "result" and templates to file
- [x] Testing
- [x] Test counter reset
- [x] Test visitors get removed
- [ ] Test config for bad values
- [x] Test timerStop()
- [x] Easier way to change times
