// Package dgohash implements a number of well-known string hashing functions.
// They all conform to the hash.Hash32 interface.
// Copyright (c) 2011 Damian Gryski <damian@gryski.com>
// Licensed under the GPLv3, or at your option any later version.
package dgohash

import (
	"hash"
)

type stringHash32 struct {
	h uint32
}

func (sh *stringHash32) Size() int     { return 4 }
func (sh *stringHash32) Sum32() uint32 { return sh.h }
func (sh *stringHash32) Reset()        { sh.h = 0 }
func (sh *stringHash32) Sum() []byte {
	p := make([]byte, 4)
	p[0] = byte(sh.h >> 24)
	p[1] = byte(sh.h >> 16)
	p[2] = byte(sh.h >> 8)
	p[3] = byte(sh.h)
	return p
}

type javaStringHash32 struct {
	stringHash32
}

// NewJava32 returns a new hash.Hash32 object, computing Java's string.hashCode() algorithm
func NewJava32() hash.Hash32 {
	return new(javaStringHash32)
}

func (sh *javaStringHash32) Write(b []byte) (int, error) {
	for _, c := range b {
		sh.h = 31*sh.h + uint32(c)
	}
	return len(b), nil
}

type djb2StringHash32 struct {
	stringHash32
}

// NewDjb32 returns a new hash.Hash32 object, computing Daniel J. Bernstein's hash
func NewDjb32() hash.Hash32 {
	sh := new(djb2StringHash32)
	sh.Reset()
	return sh
}

func (sh *djb2StringHash32) Reset() { sh.h = 5381 }

func (sh *djb2StringHash32) Write(b []byte) (int, error) {
	for _, c := range b {
		sh.h = 33*sh.h + uint32(c)
	}
	return len(b), nil
}

type elf32StringHash32 struct {
	stringHash32
}

// NewElf32 returns a new hash.Hash32 object computing the ELF32 symbol hash
func NewElf32() hash.Hash32 {
	return new(elf32StringHash32)
}

func (sh *elf32StringHash32) Write(b []byte) (int, error) {

	for _, c := range b {
		sh.h = (sh.h << 4) + uint32(c)
		g := sh.h & 0xf0000000
		if g != 0 {
			sh.h ^= g >> 24
			sh.h &= ^g
		}
	}

	return len(b), nil
}

type sdbmStringHash32 struct {
	stringHash32
}

// NewSDBM32 returns a new hash.Hash32 object, computing the string hash function from SDBM
func NewSDBM32() hash.Hash32 {
	return new(sdbmStringHash32)
}

func (sh *sdbmStringHash32) Write(b []byte) (int, error) {
	for _, c := range b {
		sh.h = uint32(c) + (sh.h << 6) + (sh.h << 16) - sh.h
	}
	return len(b), nil
}

type sqlite3StringHash32 struct {
	stringHash32
}

// NewSQLite32 returns a new hash.Hash32 object, computing the string hash function from SQLite3
func NewSQLite32() hash.Hash32 {
	return new(sqlite3StringHash32)
}

func (sh *sqlite3StringHash32) Write(b []byte) (int, error) {
	for _, c := range b {
		sh.h = (sh.h << 3) ^ sh.h ^ uint32(c)
	}
	return len(b), nil
}

type jenkinsStringHash32 struct {
	stringHash32
}

// NewJenkinsStringHash32 returns a new hash.Hash32 object, computing the Robert Jenkins' one-at-a-time string hash function
func NewJenkins32() hash.Hash32 {
	return new(jenkinsStringHash32)
}

func (sh *jenkinsStringHash32) Write(b []byte) (int, error) {
	for _, c := range b {
		sh.h += uint32(c)
		sh.h += (sh.h << 10)
		sh.h ^= (sh.h >> 6)
	}
	return len(b), nil
}

func (sh *jenkinsStringHash32) Sum32() uint32 {
	h := sh.h // copy so we don't mess with the internal state

	// Jenkins' finalize
	h += (h << 3)
	h ^= (h >> 11)
	h += (h << 15)

	return h
}

// This is a duplicate of stringHash32 Sum(), above, but is needed because otherwise a call to Sum() will call stringHash32.Sum32(), and not the Jenkins' finalize
func (sh *jenkinsStringHash32) Sum() []byte {

	h := sh.Sum32()

	p := make([]byte, 4)
	p[0] = byte(h >> 24)
	p[1] = byte(h >> 16)
	p[2] = byte(h >> 8)
	p[3] = byte(h)
	return p
}
