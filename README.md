Simple in-memory storage with JSON disk persistence.

Store all your data to disk by registering references only once on startup.

### Example

```
package main

type SomeModel struct {
	A int
	B string
}

type MyStorage1 struct {
	ID       int
	SomeList []string
}

type MyStorage2 struct {
	SomeCustModel    SomeModel
	SomeMiscSettings map[string]interface{}
}

func Main() {
	myStore1 := &MyStorage1{
		ID:       1,
		SomeList: make([]string, 0),
	}

	myStore2 := &MyStorage2{
		SomeCustModel:    SomeModel{},
		SomeMiscSettings: make(map[string]interface{}),
	}

	// this reads in from disk and stores in a temporary raw json map
    storage := NewStorage("db.json")
	err := storage.ReadFromDisk()
	if err != nil {
		// check error for file create error or corrupt data
		//  ...
		return
	}

	// register data
	//  this directly manipulates the arguments passed to it and parses the matching keys
	//  from the raw json map read from disk
	storage.RegisterData("id", &myStore1.ID)
	storage.RegisterData("some-list", &myStore1.SomeList)
	storage.RegisterData("struct-all", &myStore2)

	// do stuff...

	// save back to disk. Ideally done periodically
	storage.WriteToDisk()
}
```
