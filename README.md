# store
Simple Go library for in-memory data storage with file persistence.

## example

### adding an item

```go
package main

import (
    "github.com/karlmcguire/store"
)

type User struct {
    Name string
}

// implements store.Item interface
func (u *User) EncodeJson() []byte {
    ...
}

func main() {
    users := store.NewStore("users.db", store.Type{"main.User", &User{}})

    // id can be used to access the item later
    var id int = users.Add(&User{"karl"})

    ...
}
```

### getting an item

```go
package main

import (
    "github.com/karlmcguire/store"
)

type User struct {
    Name string
}

// implements store.Item interface
func (u *User) EncodeJson() []byte {
    ...
}

func main() {
    users := store.NewStore("users.db", store.Type{"main.User", &User{}})

    // read file 
    if err := users.DecodeFile(); err != nil {
        panic(err) 
    }

    // destination object
    var user *User
    // location (see above example)
    var id int = 0

    users.Read(func(items []store.Item) {
        user = items[id].(*User)
    })

    ...
}
```

### updating an item

```go
package main

import (
    "github.com/karlmcguire/store"
)

type User struct {
    Name string
}

// implements store.Item interface
func (u *User) EncodeJson() []byte {
    ...
}

func main() {
    users := store.NewStore("users.db", store.Type{"main.User", &User{}})

    // read file 
    if err := users.DecodeFile(); err != nil {
        panic(err) 
    }
    
    // location (see adding example)
    var id int = 0
    
    if err := users.Write(func(items []store.Item) error {
        items[id].(*User).Name = "joe"
        return nil
    }); err != nil {
        panic(err)
    }

    ...
}
```

### json representation of store

```go
```

### write store to file

```go
package main

import (
    "github.com/karlmcguire/store"
)

type User struct {
    Name string
}

// implements store.Item interface
func (u *User) EncodeJson() []byte {
    ...
}

func main() {
    users := store.NewStore("users.db", store.Type{"main.User", &User{}})

    users.Add(&User{"karl"})

    // write to 'users.db' file 
    if err := users.EncodeFile(); err != nil {
        panic(err) 
    }
    
    ...
}
```

### read store from file

```go
package main

import (
    "github.com/karlmcguire/store"
)

type User struct {
    Name string
}

// implements store.Item interface
func (u *User) EncodeJson() []byte {
    ...
}

func main() {
    users := store.NewStore("users.db", store.Type{"main.User", &User{}})

    // read from 'users.db' file 
    if err := users.DecodeFile(); err != nil {
        panic(err) 
    }
    
    ...
}
```
