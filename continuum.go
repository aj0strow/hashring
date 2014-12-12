package hashring

import (
  "crypto/md5"
  "encoding/binary"
  "sort"
)

// actual code

type Continuum struct {
  count int
  values map[string]bool
  lookups map[uint64]string
  keys uint64slice
}

func NewContinuum(count int) *Continuum {
  return &Continuum{
    count: count,
    values: make(map[string]bool),
    lookups: make(map[uint64]string),
    keys: make(uint64slice, 0, 0),
  }
}

func (c *Continuum) Add(value string) {
  if _, ok := c.values[value]; ok {
    return
  }
  c.values[value] = true

  keys := make(uint64slice, len(c.keys) + c.count)
  copy(keys, c.keys)
  for i, key := range c.getKeys(value) {
    c.lookups[key] = value
    keys[len(c.keys) + i] = key
  }
  sort.Sort(keys)
  c.keys = keys
}

func (c *Continuum) Remove(value string) {
  if _, ok := c.values[value]; !ok {
    return
  }
  delete(c.values, value)

  keys := make(uint64slice, len(c.keys) - c.count)
  rm := make(map[uint64]bool)
  for _, key := range c.getKeys(value) {
    delete(c.lookups, key)
    rm[key] = true
  }
  i := 0
  for _, key := range c.keys {
    if _, ok := rm[key]; !ok {
      keys[i] = key
      i++
    }
  }
  c.keys = keys
}

func (c *Continuum) Get(input string) string {
  num := hash(input)

  // TODO: maybe slow, try binary search or guess and correct
  // i := int(num / (MAX / uint64(len(c.keys))))

  for _, key := range c.keys {
    if num < key {
      return c.lookups[key]
    }
  }

  return c.lookups[c.keys[0]]
}

// generate hash lookups

func (c *Continuum) getKeys(value string) []uint64 {
  keys := make(uint64slice, c.count)
  for i := range keys {
    keys[i] = hash(value + string(i))
  }
  return keys
}

func hash(v string) uint64 {
  var bites [16]byte = md5.Sum([]byte(v))
  var num [8]byte
  for i := range num {
    num[i] = bites[i] ^ bites[i + 8]
  }
  return binary.LittleEndian.Uint64(num[:])
}

// sort int64 slices

type uint64slice []uint64

func (slice uint64slice) Len() int {
  return len(slice)
}

func (slice uint64slice) Swap(i, j int) {
  slice[i], slice[j] = slice[j], slice[i]
}

func (slice uint64slice) Less(i, j int) bool {
  return slice[i] < slice[j]
}

