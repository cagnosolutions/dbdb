package main

import (
	"math/rand"
	"testing"

	"github.com/cagnosolutions/dbdb"
)

var ds *dbdb.DataStore = dbdb.NewDataStore()

var stores = []string{
	"items",
	"clients",
	"quotes",
	"orders",
	"invoices",
}

type Item struct {
	Sku        int64
	Name, Desc string
	Price      float32
	Count      int64
	Restock    bool
}

type Client struct {
	Name                 []string
	Email, Phone, Gender string
}

type Order struct {
	Id    uint64
	PO    string
	Items []int64
	Total float32
}

type Quote struct {
	Id       uint64
	Name     string
	Estimate float32
}

var vals = []interface{}{
	Item{38129212, "Mustard", "Condement", 5.99, 1248, false},
	Client{[]string{"John", "Doe"}, "jdoe@example.com", "111-555-6666", "M"},
	Item{84932080, "Swimsuit", "Swimwear/Clothing", 57.25, 289, false},
	Order{767, "PO#123128", []int64{38129212, 84932080, 12892819}, 478.99},
	Quote{2329, "Some Odd Job Estimate", 383.78},
	Item{12892819, "Condoms", "N/A", 3.75, 2, true},
	Client{[]string{"Ed", "Gommel"}, "egom@example.com", "222-333-7777", "M"},
	Client{[]string{"Francheska", "Jude"}, "fjude@example.com", "333-999-1234", "F"},
}

func TestAddStore(t *testing.T) {
	var count int
	for _, store := range stores {
		ds.AddStore(store)
		count++
	}
	if count != 5 {
		t.Errorf("TestGetStore() -> count != 5, was %d\n", count)
	}
}

func TestGetStore(t *testing.T) {
	var count int
	for _, store := range stores {
		if _, ok := ds.GetStore(store); ok {
			count++
		}
	}
	if count != 5 {
		t.Errorf("TestGetStore() -> count != 5, was %d\n", count)
	}
}

func TestAdd(t *testing.T) {
	for _, store := range stores {
		for _, val := range vals {
			id := ds.Add(store, val)
			if id == 0 {
				t.Errorf("TestAdd() -> id == 0\n")
			}
		}
	}
}

func TestGet(t *testing.T) {
	for _, store := range stores {
		for i, _ := range vals {
			if doc := ds.Get(store, uint64(i+1)); doc == nil {
				t.Errorf("TestGet() -> doc == nil\n")
			}
		}
	}
}

func TestSet(t *testing.T) {
	for _, store := range stores {
		for i, val := range vals {
			ds.Set(store, uint64(i+1), val)
		}
	}
}

func TestDel(t *testing.T) {
	for _, store := range stores {
		for i, _ := range vals {
			ds.Del(store, uint64(i+1))
		}
	}
}

func TestDelStore(t *testing.T) {
	var count int
	for _, store := range stores {
		if _, ok := ds.GetStore(store); ok {
			count++
		}
	}
	for _, store := range stores {
		ds.DelStore(store)
		count--
	}
	if count != 0 {
		t.Errorf("TestGetStore() -> count != 0, was %d\n", count)
	}
}

func TestInsert100ish(t *testing.T) {
	store := "100"
	ds.AddStore(store)
	for i := 0; i < 101; i++ {
		if id := ds.Add(store, vals[rand.Intn(7)]); id == 0 {
			t.Errorf("TestInsert100ish() -> id == 0\n")
		}
	}
}

/*
func TestInsert1MillionPlus(t *testing.T) {
	store := "1048576"
	ds.AddStore(store)
	for i := 0; i < 1<<20; i++ {
		if id := ds.Add(store, vals[rand.Intn(7)]); id == 0 {
			t.Errorf("TestInsert1MillionPlus() -> id == 0\n")
		}
	}
}
*/
