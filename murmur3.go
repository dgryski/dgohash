// This file is an implementation of the murmur3 x86_32 hash function by Austin Appleby
// The code is translated from the public domain source code at http://code.google.com/p/smhasher/source/browse/trunk/MurmurHash3.cpp
// This implementation Copyright (c) 2011 Damian Gryski <damian@gryski.com>
// Licensed under the GPLv3, or at your option any later version

// Murmur3 also has 128-bit hashes, which are not (yet) included here.

package dgohash

import (
	"hash"
)

// rotate x left by r bits
func rotl32(x uint32, r uint8) uint32 {
	return (x << r) | (x >> (32 - r))
}

type murmur3 struct {
	h1     uint32  // our hash state
	length uint32  // current bytes written so far (needed for finalize)
	t      [4]byte // as-yet-unprocessed bytes
	rem    int     // how many bytes in t[] are valid
}

func (m *murmur3) Size() int      { return 4 }
func (m *murmur3) BlockSize() int { return 4 }
func (m *murmur3) Reset()         { m.h1 = uint32(0); m.length = 0; m.rem = 0 }

// NewMurmur3_x86_32 returns a new hash.Hash32 object computing the Murmur3 x86 32-bit hash
func NewMurmur3_x86_32() hash.Hash32 {
	return new(murmur3)
}

const c1 = uint32(0xcc9e2d51)
const c2 = uint32(0x1b873593)

// computes new hash state h1 merged with bytes in k1
func (m *murmur3) update(k1 uint32) {
	k1 *= c1
	k1 = rotl32(k1, 15)
	k1 *= c2

	m.h1 ^= k1
	m.h1 = rotl32(m.h1, 13)
	m.h1 = m.h1*5 + 0xe6546b64
}

func (m *murmur3) Write(data []byte) (int, error) {

	length := len(data)

	m.length += uint32(length)

	// Since the hash actually processes uint32s, but we allow []byte to be
	// Written, we have to keep track of the tail bytes that haven't yet
	// been processed, and do that on next round if we can scrounge
	// together a uint32.  If they're not merged here, they're pulled in
	// during the finalize step
	if m.rem != 0 {

		need := 4 - m.rem

		if length < need {
			copy(m.t[m.rem:], data[:length])
			m.rem += length

			return length, nil
		}

		var k1 uint32

		switch need {
		case 1:
			k1 = uint32(m.t[0]) | uint32(m.t[1])<<8 | uint32(m.t[2])<<16 | uint32(data[0])<<24
		case 2:
			k1 = uint32(m.t[0]) | uint32(m.t[1])<<8 | uint32(data[0])<<16 | uint32(data[1])<<24
		case 3:
			k1 = uint32(m.t[0]) | uint32(data[0])<<8 | uint32(data[1])<<16 | uint32(data[2])<<24
		}

		m.update(k1)

		// we've used up some bytes
		length -= need
		// nothing is left in the tail
		m.rem = 0
		data = data[need:]
	}

	// figure out the length of the tail, and round down b
	rem := length & 3
	b := length - rem

	for i := 0; i < b; i += 4 {
		k1 := uint32(data[i]) | uint32(data[i+1])<<8 | uint32(data[i+2])<<16 | uint32(data[i+3])<<24
		m.update(k1)
	}

	// copy the tail for later
	copy(m.t[:rem], data[b:])

	m.rem = rem

	return length, nil
}

func (m *murmur3) Sum(b []byte) []byte {
	h1 := m.Sum32()
	p := make([]byte, 4)
	p[0] = byte(h1 >> 24)
	p[1] = byte(h1 >> 16)
	p[2] = byte(h1 >> 8)
	p[3] = byte(h1)
	if b == nil {
		return p
	}
	return append(b, p...)
}

// murmur3 finalize step
func (m *murmur3) Sum32() uint32 {

	k1 := uint32(0)

	// copy so as not to change the internal state
	h1 := m.h1

	switch m.rem {
	case 3:
		k1 ^= uint32(m.t[2]) << 16
		fallthrough
	case 2:
		k1 ^= uint32(m.t[1]) << 8
		fallthrough
	case 1:
		k1 ^= uint32(m.t[0])
		k1 *= c1
		k1 = rotl32(k1, 15)
		k1 *= c2
		h1 ^= k1
	}

	h1 ^= m.length

	h1 ^= h1 >> 16
	h1 *= uint32(0x85ebca6b)
	h1 ^= h1 >> 13
	h1 *= uint32(0xc2b2ae35)
	h1 ^= h1 >> 16

	return h1
}
