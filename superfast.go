// This file implements the SuperFastHash by Paul Hsieh
// This code is a derivative work of the LGPL code at
// http://www.azillionmonkeys.com/qed/hash.html
// This implementation Copyright (c) 2011 Damian Gryski <damian@gryski.com>
// License: LGPL 2.1 (same terms as original code)

// Much of the framework (tracking tail bytes, etc) is duplicated from murmur3.
// It would be nice to combine the logic somehow.  If we had polymorphism, the
// Write() function would be shared if we rewrote the update() function to take
// the 32 bits in a consistent format

package dgohash

import (
	"hash"
)

type superfast struct {
	h1  uint32  // our hash state
	t   [4]byte // as-yet-unprocessed bytes
	rem int     // how many bytes in t[] are valid
}

func (m *superfast) Size() int      { return 4 }
func (m *superfast) BlockSize() int { return 4 }
func (m *superfast) Reset()         { m.h1 = uint32(0); m.rem = 0 }

// NewSuperFastHash returns a new hash.Hash32 object computing the incremental SuperFastHash
func NewSuperFastHash() hash.Hash32 {
	return new(superfast)
}

// computes new hash state h1 merged with bytes in k1,k2
func (m *superfast) update(k1, k2 uint32) {
	m.h1 += k1
	tmp := (k2 << 11) ^ m.h1
	m.h1 = (m.h1 << 16) ^ tmp
	m.h1 += m.h1 >> 11
}

// virtually identical to murmur3:Write()
func (m *superfast) Write(data []byte) (int, error) {

	length := len(data)

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

		var k1, k2 uint32

		switch need {
		case 1:
			k1 = uint32(m.t[0]) | uint32(m.t[1])<<8
			k2 = uint32(m.t[2]) | uint32(data[0])<<8
		case 2:
			k1 = uint32(m.t[0]) | uint32(m.t[1])<<8
			k2 = uint32(data[0]) | uint32(data[1])<<8
		case 3:
			k1 = uint32(m.t[0]) | uint32(data[0])<<8
			k2 = uint32(data[1]) | uint32(data[2])<<8
		}

		m.update(k1, k2)

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
		k1 := uint32(data[i]) | uint32(data[i+1])<<8
		k2 := uint32(data[i+2]) | uint32(data[i+3])<<8
		m.update(k1, k2)
	}

	// copy the tail for later
	copy(m.t[:rem], data[b:])

	m.rem = rem

	return length, nil
}

func (m *superfast) Sum(b []byte) []byte {
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

// superfast finalize step
func (m *superfast) Sum32() uint32 {

	// copy so as not to change the internal state
	h1 := m.h1

	switch m.rem {
	case 3:
		h1 += uint32(m.t[0]) | uint32(m.t[1])<<8
		h1 ^= h1 << 16
		h1 ^= uint32(m.t[2]) << 18
		h1 += h1 >> 11
		break
	case 2:
		h1 += uint32(m.t[0]) | uint32(m.t[1])<<8
		h1 ^= h1 << 11
		h1 += h1 >> 17
		break
	case 1:
		h1 += uint32(m.t[0])
		h1 ^= h1 << 10
		h1 += h1 >> 1
		break
	}

	// Force "avalanching" of final 127 bits
	h1 ^= h1 << 3
	h1 += h1 >> 5
	h1 ^= h1 << 4
	h1 += h1 >> 17
	h1 ^= h1 << 25
	h1 += h1 >> 6

	return h1
}
