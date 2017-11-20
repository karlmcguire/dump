# dump
[![GoDoc](https://img.shields.io/badge/api-reference-blue.svg)](https://godoc.org/github.com/karlmcguire/dump)
[![Go Report Card](https://img.shields.io/badge/go%20report-A%2B-green.svg)](https://goreportcard.com/report/github.com/karlmcguire/dump)
[![Coverage](https://img.shields.io/badge/coverage-100%25-ff69b4.svg)](https://gocover.io/karlmcguire/dump)

Simple Go library for in-memory data storage with file persistence.

## persistence

Dumps save to the disk (usually with a ".db" file extension).
There are currently three persistence settings available.

### manually

Using the `dump.PERSIST_MANUAL` constant allows you to control when the dump is saved to disk. 
You can manually save the dump by calling the `*Dump.Save()` function.

```go
... = dump.NewDump(..., dump.PERSIST_MANUAL, ...)
```

### on writes

Using the `dump.PERSIST_WRITES` constant will cause the dump to save to disk when `*Dump.Add()` or `*Dump.Update()` is called.

```go
... = dump.NewDump(..., dump.PERSIST_WRITES, ...)
```

### on an interval

Using the `dump.PERSIST_INTERVAL` constant will cause the dump to save to disk on a timed interval (currently 60 seconds).

```go
... = dump.NewDump(..., dump.PERSIST_INTERVAL, ...)
```

## examples

### creating a dump

```go
users, err := dump.NewDump("users.db", dump.PERSIST_WRITES, dump.Type{"main.User", User{}})
```

### adding an item

```go
// id is assigned to the index location of the item after it is added
id, err := users.Add(&User{Name: "karl"})
```

### getting an item

```go
err := users.View(func(items []dump.Item) error {
    println(items[id].(*User).Name) // will output "karl"
    return nil
})
```

### updating an item

```go
err := users.Update(func(items []dump.Item) error {
    items[id].(*User).Name = "santa"
    return nil
})
```
