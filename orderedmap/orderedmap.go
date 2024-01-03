package orderedmap

import "sync"

// OrderedMap is a concurrent safe generic map that maintains the insertion order of keys.
type OrderedMap[K comparable, V any] struct {
	sync.RWMutex
	keys   []K
	values map[K]V
}

// New creates a new instance of an OrderedMap.
func New[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		keys:   make([]K, 0),
		values: make(map[K]V),
	}
}

// Set adds or updates the key-value pair in the map.
func (o *OrderedMap[K, V]) Set(key K, value V) {
	o.Lock()
	defer o.Unlock()

	if _, exists := o.values[key]; !exists {
		o.keys = append(o.keys, key)
	}
	o.values[key] = value
}

// Get retrieves a value for a given key from the map.
func (o *OrderedMap[K, V]) Get(key K) (V, bool) {
	o.RLock()
	defer o.RUnlock()

	val, exists := o.values[key]
	return val, exists
}

// Delete removes a key-value pair from the map.
func (o *OrderedMap[K, V]) Delete(key K) {
	o.Lock()
	defer o.Unlock()

	if _, exists := o.values[key]; exists {
		delete(o.values, key)
		// Remove the key from the keys slice
		for i, k := range o.keys {
			if k == key {
				o.keys = append(o.keys[:i], o.keys[i+1:]...)
				break
			}
		}
	}
}

// Iterate returns a slice of keys in their insertion order.
func (o *OrderedMap[K, V]) Iterate() []K {
	o.RLock()
	defer o.RUnlock()

	return o.keys
}
