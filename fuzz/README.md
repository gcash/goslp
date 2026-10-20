## Fuzzer for goslp

Fuzzing helps find bugs which may have been missed by unit tests.  A differential fuzzer has been implemented for comparing the goslp v1parser and slp-validate javascript library.

### Running

The v1parser has been setup in `../v1parser/v1_parser_fuzz.go` in order to use the `dvyukov/go-fuzz` fuzz tool.  That file builds only with the `gofuzz` build tag, which go-fuzz-build sets, so it stays out of programs that import v1parser.  A corpus file `corpus.tar.gz` has been included which was generated from previous fuzzing campaigns.

**The following steps will allow you to run the fuzzer:**

1. Start the http server used for getting the js parsing result

```
$ cd ./fuzz
$ npm i
$ node ./server.js
```

2. Extract initial seed corpus to v1parser directory

```
$ cd ../v1parser
$ tar -xzvf ../fuzz/corpus.tar.gz
```

3. Build and start fuzzer

```
$ go install github.com/dvyukov/go-fuzz/go-fuzz@latest github.com/dvyukov/go-fuzz/go-fuzz-build@latest
$ go get github.com/dvyukov/go-fuzz/go-fuzz-dep
$ go-fuzz-build
$ go-fuzz --procs 4
```

go-fuzz-build needs go-fuzz-dep in go.mod, so `go get` adds it there and to go.sum.  Don't commit that change.