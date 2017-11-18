# store
Simple Go library for in-memory data storage with file persistence.

## example

### adding an item

```go
// id can be used to access the item later 
var id int = users.Add(&User{"karl"})
```

### getting an item

```go
var user &User

users.Read(func(items []store.Item) {
    user = items[id].(*User)
})
```

### updating an item

```go
var err error = users.Write(func(items []store.Item) error {
    items[id].(*User).Name = "joe"
    return nil
})
```

### json representation of store

```go
var data []byte = users.EncodeJson()
```

### write store to file

```go
var users *store.Store = store.NewStore("users.db")

gob.RegisterName("main.User", &User{})

// writes to 'users.db' file
var err error = users.EncodeFile()
```

### read store from file

```go
var users *store.Store = store.NewStore("users.db")

gob.RegisterName("main.User", &User{})

// reads from 'users.db' file
var err error = users.DecodeFile()
```
