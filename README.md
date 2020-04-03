# gojastin

Gojastin is a server measuring times between requests. Live: http://timer.miikka.xyz/

[Challenge](https://github.com/hivehelsinki/remote-challs/tree/master/chall03) where i got inspiration was fun but server side seemed far more interesting so I decided try to do my own. Proof of Concept. WIP.

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
- [ ] Template for "result"
- [x] Testing
- [ ] Test counter reset
- [ ] Test visitors get removed
- [x] Easier way to change times
