package libgrequest

import (
	"sync"
)

// SafeCache is safe to use concurrently.
type SafeCache struct {
	v   map[string]int
	mux sync.Mutex
}

// Contains :
func (c *SafeCache) Contains(s string) bool {
	c.mux.Lock()
	defer c.mux.Unlock()
	_, found := c.v[s]
	return found
}

// Inc :
func (c *SafeCache) Inc(s string) {
	c.mux.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	c.v[s]++
	c.mux.Unlock()

}