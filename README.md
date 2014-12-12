A hash ring is used to evenly distribute keys among a cluster of servers without needing to rebalance. The goal is to have a fraction of keys map to a different server when a server is added or removed. [David Karger](http://people.csail.mit.edu/karger/) is credited with coming up with the technique *consistent hashing*.

1. hash each server ip address a couple hundred times
2. each unsigned int points to the server, so the server is represented by multiple pointers on the *continuum* from 0 to 2^64-1
3. assign keys to servers by hashing the key, and finding the next pointer on the continuum

This is a small implementation in the spirit of learning by doing. It works.

```go
// create a new continuum
circle := NewContinuum(200)

// keep track of distribution
counts := make(map[string]int)

// add four hosts
for port := 8000; port < 8005; port++ {
  host := fmt.Sprintf("localhost:%d", port)
  circle.Add(host)
  counts[host] = 0
}

// keep track of picks
picks := make(map[string]string)

// pick hosts for cache keys
for i := 0; i < 1000; i++ {
  key := "cache key " + string(i)
  host := circle.Get(key)
  picks[key] = host
  counts[host]++
}

// print out distribution
fmt.Println("Host and count.")
for host, count := range counts {
  fmt.Printf("%s: %d\n", host, count)
}

// remove a host
circle.Remove("localhost:8004")

// count keys that resolve to a new host
misses := 0
for i := 0; i < 1000; i++ {
  key := "cache key " + string(i)
  host := circle.Get(key)
  if picks[key] != host {
    misses++
  }
}

// print out cache misses
fmt.Printf("%d misses\n", misses)
```

And the output.

```
Host and count.
localhost:8000: 195
localhost:8001: 200
localhost:8002: 190
localhost:8003: 199
localhost:8004: 216
216 misses
```

The distribution is pretty even. The md5 hashing makes the one-letter differences significant. After removing the `localhost:8004` host, only the 216 cache keys it had were redistributed. 
