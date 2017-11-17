# store
Simple Go library for in-memory data storage with file persistence.

## example

### adding an item

```go
type User struct {
    Name string 
}

// implementing store.Item interface
func (u *User) EncodeJson() []byte {
    ...
}

func main() {
    var users *store.Store
    users = store.NewStore("users.db")
    
    var id int 
    id = users.Add(&Users{"karl"})
    
    users.Read(func(items []store.Items) {
        println(items[id].(*User).Name) // will output "karl"
    })
}
```
