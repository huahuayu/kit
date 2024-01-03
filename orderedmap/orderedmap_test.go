package orderedmap

import (
	"sync"
	"testing"
)

func TestSetAndGet(t *testing.T) {
	om := New[string, int]()
	om.Set("key1", 10)
	om.Set("key2", 20)

	if value, exists := om.Get("key1"); !exists || value != 10 {
		t.Errorf("Get(\"key1\") = %v, %v, want 10, true", value, exists)
	}

	if value, exists := om.Get("key2"); !exists || value != 20 {
		t.Errorf("Get(\"key2\") = %v, %v, want 20, true", value, exists)
	}
}

func TestDelete(t *testing.T) {
	om := New[string, int]()
	om.Set("key1", 10)
	om.Delete("key1")

	if _, exists := om.Get("key1"); exists {
		t.Errorf("Delete(\"key1\") failed, key1 still exists")
	}
}

func TestIterate(t *testing.T) {
	om := New[string, int]()
	om.Set("key1", 10)
	om.Set("key2", 20)
	expectedKeys := []string{"key1", "key2"}
	keys := om.Iterate()

	if len(keys) != len(expectedKeys) {
		t.Errorf("Iterate() returned %v, want %v", keys, expectedKeys)
		return
	}
	for i, key := range keys {
		if key != expectedKeys[i] {
			t.Errorf("Iterate() returned %v, want %v", keys, expectedKeys)
			return
		}
	}
}

func TestConcurrency(t *testing.T) {
	om := New[int, int]()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			om.Set(i, i*10)
		}(i)
	}
	wg.Wait()

	for i := 0; i < 100; i++ {
		if value, exists := om.Get(i); !exists || value != i*10 {
			t.Errorf("Get(%d) = %d, %v, want %d, true", i, value, exists, i*10)
		}
	}
}
