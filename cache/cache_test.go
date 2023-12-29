package cache

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestNewTTLCache(t *testing.T) {
	var cache = New[string, int]()
	if cache == nil {
		t.Error("Expected non-nil cache, got nil")
	}
}

func TestSetAndGet(t *testing.T) {
	var cache = New[string, string]()
	key := "key1"
	value := "value1"

	cache.Set(key, value)
	retrievedValue, exists := cache.Get(key)

	if !exists || retrievedValue != value {
		t.Errorf("Expected %v, got %v", value, retrievedValue)
	}
}

func TestRemove(t *testing.T) {
	var cache = New[string, string]()
	key := "key1"
	value := "value1"

	cache.Set(key, value)
	cache.Remove(key)
	_, exists := cache.Get(key)

	if exists {
		t.Errorf("Expected key %v to be removed", key)
	}
}

func TestPop(t *testing.T) {
	var cache = New[string, string]()
	key := "key1"
	value := "value1"

	cache.Set(key, value)
	retrievedValue, exists := cache.Pop(key)

	if !exists || retrievedValue != value {
		t.Errorf("Expected %v, got %v", value, retrievedValue)
	}

	_, stillExists := cache.Get(key)
	if stillExists {
		t.Errorf("Expected key %v to be removed after Pop", key)
	}
}

func TestTTL(t *testing.T) {
	var cache = New[string, string]()
	key := "key1"
	value := "value1"
	ttl := 50 * time.Millisecond

	cache.Set(key, value, ttl)
	time.Sleep(60 * time.Millisecond)
	_, exists := cache.Get(key)

	if exists {
		t.Errorf("Expected key %v to be expired", key)
	}
}

func TestConcurrentSetAndGet(t *testing.T) {
	cache := New[string, int]()
	var wg sync.WaitGroup
	concurrencyLevel := 100
	keyValuePairs := 1000

	for i := 0; i < concurrencyLevel; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < keyValuePairs; j++ {
				key := "key" + strconv.Itoa(goroutineID*keyValuePairs+j)
				value := goroutineID*keyValuePairs + j
				cache.Set(key, value)

				if retrievedValue, exists := cache.Get(key); !exists || retrievedValue != value {
					t.Errorf("Expected %v, got %v", value, retrievedValue)
				}
			}
		}(i)
	}

	wg.Wait()
}

func BenchmarkConcurrentSet10(b *testing.B) {
	benchmarkConcurrentSet(b, 10)
}

func BenchmarkConcurrentGet10(b *testing.B) {
	benchmarkConcurrentGet(b, 10)
}

func BenchmarkConcurrentRemove10(b *testing.B) {
	benchmarkConcurrentRemove(b, 10)
}

func BenchmarkConcurrentPop10(b *testing.B) {
	benchmarkConcurrentPop(b, 10)
}

func BenchmarkConcurrentSet100(b *testing.B) {
	benchmarkConcurrentSet(b, 100)
}

func BenchmarkConcurrentGet100(b *testing.B) {
	benchmarkConcurrentGet(b, 100)
}

func BenchmarkConcurrentRemove100(b *testing.B) {
	benchmarkConcurrentRemove(b, 100)
}

func BenchmarkConcurrentPop100(b *testing.B) {
	benchmarkConcurrentPop(b, 100)
}

func benchmarkConcurrentSet(b *testing.B, concurrencyLevel int) {
	cache := New[string, int]()
	var wg sync.WaitGroup

	for n := 0; n < b.N; n++ {
		for i := 0; i < concurrencyLevel; i++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()
				key := "keySet" + strconv.Itoa(goroutineID)
				cache.Set(key, goroutineID)
			}(i)
		}
		wg.Wait()
	}
}

func benchmarkConcurrentGet(b *testing.B, concurrencyLevel int) {
	cache := New[string, int]()
	var wg sync.WaitGroup

	// Pre-populate the cache
	for i := 0; i < concurrencyLevel; i++ {
		key := "keyGet" + strconv.Itoa(i)
		cache.Set(key, i)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < concurrencyLevel; i++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()
				key := "keyGet" + strconv.Itoa(goroutineID)
				cache.Get(key)
			}(i)
		}
		wg.Wait()
	}
}

func benchmarkConcurrentRemove(b *testing.B, concurrencyLevel int) {
	cache := New[string, int]()
	var wg sync.WaitGroup

	for n := 0; n < b.N; n++ {
		// Pre-populate the cache
		for i := 0; i < concurrencyLevel; i++ {
			key := "keyRemove" + strconv.Itoa(i)
			cache.Set(key, i)
		}

		b.ResetTimer()
		for i := 0; i < concurrencyLevel; i++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()
				key := "keyRemove" + strconv.Itoa(goroutineID)
				cache.Remove(key)
			}(i)
		}
		wg.Wait()
	}
}

func benchmarkConcurrentPop(b *testing.B, concurrencyLevel int) {
	cache := New[string, int]()
	var wg sync.WaitGroup

	for n := 0; n < b.N; n++ {
		// Pre-populate the cache
		for i := 0; i < concurrencyLevel; i++ {
			key := "keyPop" + strconv.Itoa(i)
			cache.Set(key, i)
		}

		b.ResetTimer()
		for i := 0; i < concurrencyLevel; i++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()
				key := "keyPop" + strconv.Itoa(goroutineID)
				cache.Pop(key)
			}(i)
		}
		wg.Wait()
	}
}
