package storage

import (
	"os"
	"testing"
)

const filename = "db_test.json"

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

func TestStorage(t *testing.T) {
	var err error
	myStore1 := &MyStorage1{
		ID:       1,
		SomeList: make([]string, 0),
	}

	myStore2 := &MyStorage2{
		SomeCustModel:    SomeModel{},
		SomeMiscSettings: make(map[string]interface{}),
	}

	storage := NewStorage(filename)

	storage.RegisterData("id", &myStore1.ID)
	storage.RegisterData("some-list", &myStore1.SomeList)
	storage.RegisterData("my-store-full", &myStore2)

	myStore1.ID = 1234
	myStore1.SomeList = append(myStore1.SomeList, "a")
	myStore1.SomeList = append(myStore1.SomeList, "b")
	myStore2.SomeCustModel.A = 4321
	myStore2.SomeCustModel.B = "4321s"
	myStore2.SomeMiscSettings["test1"] = 54321
	myStore2.SomeMiscSettings["test2"] = "54321s"

	err = storage.WriteToDisk()
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	// RESET

	myStore3 := &MyStorage1{}
	myStore4 := &MyStorage2{}

	storage2 := NewStorage(filename)
	err = storage2.ReadFromDisk()
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	storage2.RegisterData("id", &myStore3.ID)
	storage2.RegisterData("some-list", &myStore3.SomeList)
	storage2.RegisterData("my-store-full", &myStore4)

	if myStore3.ID != 1234 {
		t.Errorf("expected: %v, result: %v", 1234, myStore3.ID)
	}
	if len(myStore3.SomeList) != 2 {
		t.Errorf("expected: %v, result: %v", 2, len(myStore3.SomeList))
	}
	if myStore3.SomeList[0] != "a" {
		t.Errorf("expected: %v, result: %v", "a", myStore3.SomeList[0])
	}
	if myStore3.SomeList[1] != "b" {
		t.Errorf("expected: %v, result: %v", "b", myStore3.SomeList[1])
	}
	if myStore4.SomeCustModel.A != 4321 {
		t.Errorf("expected: %v, result: %v", 4321, myStore4.SomeCustModel.A)
	}
	if myStore4.SomeCustModel.B != "4321s" {
		t.Errorf("expected: %v, result: %v", "4321s", myStore4.SomeCustModel.B)
	}
	if int(myStore4.SomeMiscSettings["test1"].(float64)) != 54321 {
		t.Errorf("expected: %v, result: %v", 54321, myStore4.SomeMiscSettings["test1"])
	}
	if myStore4.SomeMiscSettings["test2"] != "54321s" {
		t.Errorf("expected: %v, result: %v", "54321s", myStore4.SomeMiscSettings["test2"])
	}

	os.Remove(filename)
}
