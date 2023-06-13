package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type HashFunc func(data []byte) uint32

type ConsistentHash struct {

	// hash use to storage hash algorithm,default use crc32 algorithm
	hash HashFunc
	// a ring of hash,must be sorted
	ring []int
	// the multiple of virtual node
	replicas int
	// reflect virtual nodes into real nodes
	hashmap map[int]string
}

func NewConsistentHash(multiple int, fn HashFunc) *ConsistentHash {
	var consistentHash = &ConsistentHash{
		hash:     fn,
		replicas: multiple,
		hashmap:  make(map[int]string),
	}
	if fn == nil {
		consistentHash.hash = crc32.ChecksumIEEE
	}
	return consistentHash
}

func (c *ConsistentHash) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < c.replicas; i++ {
			value := int(c.hash(([]byte)(strconv.Itoa(i) + key)))
			c.ring = append(c.ring, value)
			c.hashmap[value] = key
		}
	}
	sort.Ints(c.ring)
}

// Get to searches which node should be chosen to storage the key
func (c *ConsistentHash) Get(Key string) string {

	value := int(c.hash([]byte(Key)))
	var i int
	// TODO can be replaced by binary search
	for i = 0; i < len(c.ring); i++ {
		if c.ring[i] > value {
			break
		}
	}
	return c.hashmap[c.ring[i%len(c.ring)]]
}
func (c *ConsistentHash) IsEmpty() bool {
	if len(c.ring) == 0 {
		return true
	}
	return false
}
