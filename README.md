# store
Simple Go library for in-memory data storage with file persistence.

## example

### adding an item

```go
type User struct {
    Name string 
}

func (u *User) EncodeJson() []byte {
    ...
}

func main() {
    users := store.NewStore("users.db")

    id := users.Add(&User{"karl"})

    users.Read(func(items []store.Items) {
        println(items[id].(*User).Name) // will output "karl"
    })
}
```
