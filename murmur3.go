// This file is an implementation of the murmur3 x86_32 hash function by Austin Appleby 
// This code is translated from the public domain source code at http://code.google.com/p/smhasher/source/browse/trunk/MurmurHash3.cpp
// There are also 128 bit hashes, not included here. (yet...)

package dgohash

import (
	"hash"
	"os"
)

func rotl32(x uint32, r uint8) uint32 {
	return (x << r) | (x >> (32 - r))
}

type murmur3 struct {
	h1     uint32
	t      [4]byte
	length uint32
	rem    int
}

func (m *murmur3) Size() int { return 4 }
func (m *murmur3) Reset()    { m.h1 = uint32(0); m.length = 0; m.rem = 0 }

func NewMurmur3_x86_32() hash.Hash32 {
	return new(murmur3)
}

const c1 = uint32(0xcc9e2d51)
const c2 = uint32(0x1b873593)

func update(h1, k1 uint32) uint32 {
	k1 *= c1
	k1 = rotl32(k1, 15)
	k1 *= c2

	h1 ^= k1
	h1 = rotl32(h1, 13)
	h1 = h1*5 + 0xe6546b64

	return h1
}

func (m *murmur3) Write(data []byte) (int, os.Error) {

	length := len(data)

	m.length += uint32(length)

	if m.rem != 0 {

		need := 4 - m.rem

		if length < need {
			for i := 0; i < len(data); i++ {
				m.t[m.rem] = data[i]
				m.rem++
			}

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

		m.h1 = update(m.h1, k1)

		// we've used up some bytes
		length -= need
		// nothing is left in the tail
		m.rem = 0
		data = data[need:]
	}

	rem := length & 3
	b := length - rem

	for i := 0; i < b; i += 4 {
		k1 := uint32(data[i]) | uint32(data[i+1])<<8 | uint32(data[i+2])<<16 | uint32(data[i+3])<<24
		m.h1 = update(m.h1, k1)
	}

	// copy the tail for later
        // this should probably be unrolled
	for i := 0; i < rem; i++ {
		m.t[i] = data[b+i]
	}

	m.rem = rem

	return length, nil
}

func (m *murmur3) Sum() []byte {
	p := make([]byte, 4)
	h1 := m.Sum32()
	p[0] = byte(h1 >> 24)
	p[1] = byte(h1 >> 16)
	p[2] = byte(h1 >> 8)
	p[3] = byte(h1)
	return p
}

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
