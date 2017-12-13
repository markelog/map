// Package list create hash instance and checks value presence
package list

import (
	"sync"

	"bitbucket.org/creachadair/cityhash"
)

// List settings
type List struct {
	mutex *sync.RWMutex
	list  map[uint64]bool
}

// New returns a instance of the list
func New() *List {
	return &List{
		mutex: &sync.RWMutex{},
		list:  make(map[uint64]bool),
	}
}

// Add byte list to the instance with the hash
func (me *List) Add(text []byte) {
	hash := cityhash.Hash64(text)

	me.mutex.Lock()
	me.list[hash] = true
	me.mutex.Unlock()
}

// Has checks if these bytes already present in list
func (me List) Has(text []byte) (has bool) {
	hash := cityhash.Hash64(text)

	me.mutex.RLock()
	_, has = me.list[hash]
	me.mutex.RUnlock()

	return
}
