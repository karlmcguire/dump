# dump
[![GoDoc](https://img.shields.io/badge/api-reference-blue.svg)](https://godoc.org/github.com/karlmcguire/dump)
[![Go Report Card](https://img.shields.io/badge/go%20report-A%2B-green.svg)](https://goreportcard.com/report/github.com/karlmcguire/dump)
[![Coverage](https://img.shields.io/badge/coverage-100%25-ff69b4.svg)](https://gocover.io/karlmcguire/dump)

Simple Go library for in-memory data storage with file persistence.

## examples

### creating a dump

```go
users, err := dump.NewDump("users.db", dump.PERSIST_WRITES, dump.Type{"main.User", User{}})
```
