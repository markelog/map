// Package list create hash instance and checks value presence
package list

import (
	"bitbucket.org/creachadair/cityhash"
)

// List settings
type List struct {
	list map[uint64]bool
}

// New returns a instance of the list
func New() *List {
	return &List{
		list: make(map[uint64]bool),
	}
}

// Add byte list to the instance with the hash
func (items *List) Add(text []byte) {
	hash := cityhash.Hash64(text)

	items.list[hash] = true
}

// Has checks if these bytes already present in list
func (items List) Has(text []byte) (has bool) {
	hash := cityhash.Hash64(text)
	_, has = items.list[hash]

	return
}
