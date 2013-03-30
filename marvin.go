// This file is an implementation of the InternalMarvin32HashString hash function.
// The code is based on the implementation found at https://github.com/floodyberry/Marvin32

// This implementation Copyright (c) 2011 Damian Gryski <damian@gryski.com>
// Licensed under the GPLv3, or at your option any later version

package dgohash

import (
	"hash"
)

type marvin struct {
	seed   uint64
	lo, hi uint32
	t      [4]byte // as-yet-unprocessed bytes
	rem    int     // how many bytes in t[] are valid
}

func (st *marvin) update(v uint32) {
	st.lo += v
	st.hi ^= st.lo
	st.lo = rotl32(st.lo, 20) + st.hi
	st.hi = rotl32(st.hi, 9) ^ st.lo
	st.lo = rotl32(st.lo, 27) + st.hi
	st.hi = rotl32(st.hi, 19)
}

// NewMarvin32 returns a new hash.Hash32 object computing Microsoft's InternalMarvin32HashString seeded hash.
func NewMarvin32(seed uint64) hash.Hash32 {
	m := new(marvin)
	m.seed = seed
	m.Reset()
	return m
}

func (m *marvin) Size() int      { return 4 }
func (m *marvin) BlockSize() int { return 4 }
func (m *marvin) Reset()         { m.lo = uint32(m.seed); m.hi = uint32(m.seed >> 32); m.rem = 0 }

func (m *marvin) Write(data []byte) (int, error) {

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

func (m *marvin) Sum(b []byte) []byte {
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

// marvin finalize step
func (m *marvin) Sum32() uint32 {

	/* pad the final 0-3 bytes with 0x80 */
	final := uint32(0x80)

	// copy so as not to change the internal state
	m_tmp := *m

	switch m_tmp.rem {

	case 3:
		final = (final << 8) | uint32(m_tmp.t[2])
		fallthrough
	case 2:
		final = (final << 8) | uint32(m_tmp.t[1])
		fallthrough
	case 1:
		final = (final << 8) | uint32(m_tmp.t[0])
	}

	m_tmp.update(final)
	m_tmp.update(0)

	return m_tmp.lo ^ m_tmp.hi
}
