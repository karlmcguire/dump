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

users.Read(func(items []store.Item) error {
    user = items[id].(*User)
})
```

### json representation of store

```go
var data []byte = users.EncodeJson()
```

### write store to file

```go
users := store.NewStore("users.db")

gob.RegisterName("main.User", &User{})

users.EncodeFile() // writes to 'users.db' file
```

### read store from file

```go

users := store.NewStore("users.db")

gob.RegisterName("main.User", &User{})

users.DecodeFile() // reads from 'users.db' file
```
