// Package dgohash implements a number of well-known string hashing functions.
// They all conform to the hash.Hash32 interface.
// Copyright (c) 2011 Damian Gryski <damian@gryski.com>
// Licensed under the GPLv3, or at your option any later version.
package dgohash

import (
	"hash"
)

type javaStringHash32 uint32

// NewJava32 returns a new hash.Hash32 object, computing Java's string.hashCode() algorithm
func NewJava32() hash.Hash32                { sh := javaStringHash32(0); sh.Reset(); return &sh }
func (sh *javaStringHash32) Size() int      { return 4 }
func (sh *javaStringHash32) BlockSize() int { return 1 }
func (sh *javaStringHash32) Sum32() uint32  { return uint32(*sh) }
func (sh *javaStringHash32) Reset()         { *sh = javaStringHash32(0) }
func (sh *javaStringHash32) Sum(b []byte) []byte {
	v := uint32(*sh)
	return append(b, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

func (sh *javaStringHash32) Write(b []byte) (int, error) {
	h := uint32(*sh)
	for _, c := range b {
		h = 31*h + uint32(c)
	}
	*sh = javaStringHash32(h)
	return len(b), nil
}

type djb2StringHash32 uint32

// NewDjb32 returns a new hash.Hash32 object, computing Daniel J. Bernstein's hash
func NewDjb32() hash.Hash32                 { sh := djb2StringHash32(0); sh.Reset(); return &sh }
func (sh *djb2StringHash32) Size() int      { return 4 }
func (sh *djb2StringHash32) BlockSize() int { return 1 }
func (sh *djb2StringHash32) Sum32() uint32  { return uint32(*sh) }
func (sh *djb2StringHash32) Reset()         { *sh = djb2StringHash32(5381) }
func (sh *djb2StringHash32) Sum(b []byte) []byte {
	v := uint32(*sh)
	return append(b, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

func (sh *djb2StringHash32) Write(b []byte) (int, error) {
	h := uint32(*sh)
	for _, c := range b {
		h = 33*h + uint32(c)
	}
	*sh = djb2StringHash32(h)
	return len(b), nil
}

type djb2aStringHash32 uint32

// NewDjb32a returns a new hash.Hash32 object, computing a variant of Daniel J. Bernstein's hash that uses xor instead of +
func NewDjb32a() hash.Hash32                 { sh := djb2aStringHash32(0); sh.Reset(); return &sh }
func (sh *djb2aStringHash32) Size() int      { return 4 }
func (sh *djb2aStringHash32) BlockSize() int { return 1 }
func (sh *djb2aStringHash32) Sum32() uint32  { return uint32(*sh) }
func (sh *djb2aStringHash32) Reset()         { *sh = djb2aStringHash32(5381) }
func (sh *djb2aStringHash32) Sum(b []byte) []byte {
	v := uint32(*sh)
	return append(b, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

func (sh *djb2aStringHash32) Write(b []byte) (int, error) {
	h := uint32(*sh)
	for _, c := range b {
		h = 33*h ^ uint32(c)
	}
	*sh = djb2aStringHash32(h)
	return len(b), nil
}

type elf32StringHash32 uint32

// NewElf32 returns a new hash.Hash32 object computing the ELF32 symbol hash
func NewElf32() hash.Hash32                  { sh := elf32StringHash32(0); sh.Reset(); return &sh }
func (sh *elf32StringHash32) Size() int      { return 4 }
func (sh *elf32StringHash32) BlockSize() int { return 1 }
func (sh *elf32StringHash32) Sum32() uint32  { return uint32(*sh) }
func (sh *elf32StringHash32) Reset()         { *sh = elf32StringHash32(0) }
func (sh *elf32StringHash32) Sum(b []byte) []byte {
	v := uint32(*sh)
	return append(b, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

func (sh *elf32StringHash32) Write(b []byte) (int, error) {
	h := uint32(*sh)
	for _, c := range b {
		h = (h << 4) + uint32(c)
		g := h & 0xf0000000
		if g != 0 {
			h ^= g >> 24
			h &= ^g
		}
	}
	*sh = elf32StringHash32(h)
	return len(b), nil
}

type sdbmStringHash32 uint32

// NewSDBM32 returns a new hash.Hash32 object, computing the string hash function from SDBM
func NewSDBM32() hash.Hash32                { sh := sdbmStringHash32(0); sh.Reset(); return &sh }
func (sh *sdbmStringHash32) Size() int      { return 4 }
func (sh *sdbmStringHash32) BlockSize() int { return 1 }
func (sh *sdbmStringHash32) Sum32() uint32  { return uint32(*sh) }
func (sh *sdbmStringHash32) Reset()         { *sh = sdbmStringHash32(0) }
func (sh *sdbmStringHash32) Sum(b []byte) []byte {
	v := uint32(*sh)
	return append(b, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

func (sh *sdbmStringHash32) Write(b []byte) (int, error) {
	h := uint32(*sh)
	for _, c := range b {
		h = uint32(c) + (h << 6) + (h << 16) - h
	}
	*sh = sdbmStringHash32(h)
	return len(b), nil
}

type sqlite3StringHash32 uint32

// NewSQLite32 returns a new hash.Hash32 object, computing the string hash function from SQLite3
func NewSQLite32() hash.Hash32                 { sh := sqlite3StringHash32(0); sh.Reset(); return &sh }
func (sh *sqlite3StringHash32) Size() int      { return 4 }
func (sh *sqlite3StringHash32) BlockSize() int { return 1 }
func (sh *sqlite3StringHash32) Sum32() uint32  { return uint32(*sh) }
func (sh *sqlite3StringHash32) Reset()         { *sh = sqlite3StringHash32(0) }
func (sh *sqlite3StringHash32) Sum(b []byte) []byte {
	v := uint32(*sh)
	return append(b, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

func (sh *sqlite3StringHash32) Write(b []byte) (int, error) {
	h := uint32(*sh)
	for _, c := range b {
		h = (h << 3) ^ h ^ uint32(c)
	}
	*sh = sqlite3StringHash32(h)
	return len(b), nil
}

type jenkinsStringHash32 uint32

// NewJenkins32 returns a new hash.Hash32 object, computing the Robert Jenkins' one-at-a-time string hash function
func NewJenkins32() hash.Hash32 {
	var s = jenkinsStringHash32(0)
	return &s
}
func (sh *jenkinsStringHash32) Size() int      { return 4 }
func (sh *jenkinsStringHash32) BlockSize() int { return 1 }
func (sh *jenkinsStringHash32) Reset()         { *sh = jenkinsStringHash32(0) }

func (sh *jenkinsStringHash32) Write(b []byte) (int, error) {
	h := uint32(*sh)
	for _, c := range b {
		h += uint32(c)
		h += (h << 10)
		h ^= (h >> 6)
	}
	*sh = jenkinsStringHash32(h)
	return len(b), nil
}

func (sh *jenkinsStringHash32) Sum32() uint32 {
	h := uint32(*sh) // copy so we don't mess with the internal state

	// Jenkins' finalize
	h += (h << 3)
	h ^= (h >> 11)
	h += (h << 15)

	return h
}

func (sh *jenkinsStringHash32) Sum(b []byte) []byte {
	v := sh.Sum32()
	return append(b, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}
