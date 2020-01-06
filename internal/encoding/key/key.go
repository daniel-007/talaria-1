// Copyright 2019 Grabtaxi Holdings PTE LTE (GRAB), All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file

package key

import (
	"encoding/binary"
	"math"
	"sync/atomic"
	"time"

	"github.com/spaolacci/murmur3"
)

const size = 16

// Next is used as a sequence number, it's okay to overflow.
var next uint32

// New generates a new key for the storage.
func New(eventName string, tsi time.Time) []byte {
	out := make([]byte, size)
	binary.BigEndian.PutUint32(out[0:4], murmur3.Sum32([]byte(eventName)))
	binary.BigEndian.PutUint64(out[4:12], uint64(tsi.Unix()))
	binary.BigEndian.PutUint32(out[12:16], atomic.AddUint32(&next, 1))
	return out
}

// HashOf returns the hash value of the key
func HashOf(k []byte) uint32 {
	return binary.BigEndian.Uint32(k[0:4])
}

// PrefixOf a common prefix between two keys (common leading bytes) which is
// then used as a prefix for Badger to narrow down SSTables to traverse.
func PrefixOf(seek, until []byte) []byte {
	var prefix []byte

	// Calculate the minimum length
	length := len(seek)
	if len(until) < length {
		length = len(until)
	}

	// Iterate through the bytes and append common ones
	for i := 0; i < length; i++ {
		if seek[i] != until[i] {
			break
		}
		prefix = append(prefix, seek[i])
	}
	return prefix
}

// First returns the smallest possible key
func First() []byte {
	return make([]byte, size)
}

// Last returns the largest possible key
func Last() []byte {
	out := make([]byte, size)
	for i := 0; i < size; i++ {
		out[i] = math.MaxUint8
	}
	return out
}
