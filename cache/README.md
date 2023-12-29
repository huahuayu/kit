# Cache

The `cache` pkg provides a generic in-memory key-value cache with optional Time-To-Live (TTL) support in Go.

## Features

- Generic key-value cache: The cache can store any type of key-value pairs.
- Optional TTL support: Each key-value pair can have an optional TTL, after which the pair is automatically removed from the cache.
- Thread-safe: The cache uses a `sync.RWMutex` to ensure that it can be safely used from multiple goroutines.

## Usage

First, import the `cache` package:

```go
import "github.com/huahuayu/kit/cache"
```

Create a new cache:

```go
c := cache.New()
```

```go
c.Set("key", "value")
```

Retrieve a value from the cache:

```go
value, found := c.Get("key")
```

Remove a key-value pair from the cache:

```go
c.Remove("key")
```

Pop a value from the cache:

```go
value, found := c.Pop("key")
```

## TTL

Optionally, you can set a TTL for each key-value pair:

```go
c := cache.New(5 * time.Minute)
c.Set("key", "value", 5 * time.Minute)
```