# mincache

A very minimal in-memory cache for Go.

## Installation

```bash
go get github.com/NoBypass/mincache
```

## Usage

```go
package main

import (
	"github.com/NoBypass/mincache"
	"time"
)

func main() {
	c := mincache.New()
	defer c.Close()   // Close the cache when you're done with it

	c.Set("key", "value", mincache.ExpireIn(5*time.Minute)) // Set a key with a TTL of 5 minutes
	c.Get("key")         // Get the value of the key
	c.Delete("key")       // Delete the key manually

	c.Set("key", 123, mincache.ExpireAt(time.Now().Add(24*time.Hour))) // Overwrite the key with a new value and a new TTL

	c.Set(123, "value")    // You can use any type as a key or value

	safeCache := mincache.NewSafe[int, string]() // Create a typesafe cache
	safeCache.Set(123, "hello world")
	safeCache.Get(123)    // Returns "hello world". No type casting needed
}
```

## How it works

This library is supposed to be as simple as possible. If you want any functionality besides setting, getting and
deleting keys, something like [go-cache](https://github.com/patrickmn/go-cache) might be what you're looking for.

Each instance of the cache will have its own goroutine running in the background (which can be stopped by calling
`Cache.Close()`). This goroutine will make sure that keys with a TTL are deleted when their time is up.