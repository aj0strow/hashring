package hashring

import (
  "testing"
  "fmt"
)

func Testhash(t *testing.T) {
  s := "127.0.0.1"
  if hash(s) != 18377291947331486790 {
    t.Errorf("Hash function is incorrect.")
  }
}

func TestAdd(t *testing.T) {
  c := NewContinuum(1)
  c.Add("localhost")
  if _, ok := c.values["localhost"]; !ok {
    t.Errorf("Value not added.")
  }
}

func TestAddCountLookups(t *testing.T) {
  c := NewContinuum(3)
  c.Add("localhost")
  c.Add("localhost")
  if len(c.lookups) != 3 {
    t.Errorf("Wrong count of lookups added.")
  }
}

func TestRemove(t *testing.T) {
  c := NewContinuum(1)
  c.Add("localhost")
  c.Remove("localhost")
  if _, ok := c.values["localhost"]; ok {
    t.Errorf("Value not removed.")
  }
}

func TestGet(t *testing.T) {
  c := NewContinuum(1)
  c.Add("127.0.0.1")
  if c.Get("key") != "127.0.0.1" {
    t.Errorf("Get value failed.")
  }
}

func TestGetConsistency(t *testing.T) {
  c := NewContinuum(10)
  for i := 0; i < 100; i++ {
    c.Add("0.0.0." + string(i))
  }
  for i := 0; i < 100; i++ {
    key := "cache key " + string(i)
    value := c.Get(key)
    if value != c.Get(key) {
      t.Errorf("Inconsistent")
    }
  }
}

func TestStuff(t *testing.T) {
  // create a new continuum

  circle := NewContinuum(200)
  counts := make(map[string]int)
  for port := 8000; port < 8005; port++ {
    host := fmt.Sprintf("localhost:%d", port)
    circle.Add(host)
    counts[host] = 0
  }

  picks := make(map[string]string)
  for i := 0; i < 1000; i++ {
    key := "cache key " + string(i)
    host := circle.Get(key)
    picks[key] = host
    counts[host]++
  }
  fmt.Println("Host and count.")
  for host, count := range counts {
    fmt.Printf("%s: %d\n", host, count)
  }

  circle.Remove("localhost:8004")
  misses := 0

  for i := 0; i < 1000; i++ {
    key := "cache key " + string(i)
    host := circle.Get(key)
    if picks[key] != host {
      misses++
    }
  }

  fmt.Printf("%d misses\n", misses)  
}
