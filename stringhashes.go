package dgohash

import "os" // for os.Error

// TODO: add mumblehash3
// TODO: add superfasthash
// TODO: check signedness of characters in hash functions
// TODO: add test cases for the hashes

type StringHash32 struct {
	h uint32
}

func (sh *StringHash32) Size() int     { return 4 }
func (sh *StringHash32) Sum32() uint32 { return sh.h }
func (sh *StringHash32) Reset()        { sh.h = 0 }
func (sh *StringHash32) Sum() []byte {
	p := make([]byte, 4)
	p[0] = byte(sh.h >> 24)
	p[1] = byte(sh.h >> 16)
	p[2] = byte(sh.h >> 8)
	p[3] = byte(sh.h)
	return p
}

// default java string.hashCode() algorithm
type JavaStringHash32 struct {
	StringHash32
}

func NewJavaStringHash32() *JavaStringHash32 {
	return new(JavaStringHash32)
}

func (j *JavaStringHash32) Write(b []byte) (int, os.Error) {
	for _, c := range b {
		j.h = 31*j.h + uint32(c)
	}
	return len(b), nil
}

// Bernstein hash: used in Elf32, Glib, ...
type Djb2StringHash32 struct {
	StringHash32
}

func NewDjb2StringHash32() *Djb2StringHash32 {
	sh := new(Djb2StringHash32)
	sh.Reset()
	return sh
}

func (sh *Djb2StringHash32) Reset() { sh.h = 5381 }

func (j *Djb2StringHash32) Write(b []byte) (int, os.Error) {
	for _, c := range b {
		j.h = 33*j.h + uint32(c)
	}
	return len(b), nil
}

// hash function from sdbm
type SDBMStringHash32 struct {
	StringHash32
}

func NewSDBMStringHash32() *SDBMStringHash32 {
	return new(SDBMStringHash32)
}

func (j *SDBMStringHash32) Write(b []byte) (int, os.Error) {
	for _, c := range b {
		j.h = uint32(c) + (j.h << 6) + (j.h << 16) - j.h
	}
	return len(b), nil
}

type SQLite3StringHash32 struct {
	StringHash32
}

func NewSQLite3StringHash32() *SQLite3StringHash32 {
	return new(SQLite3StringHash32)
}

func (sh *SQLite3StringHash32) Write(b []byte) (int, os.Error) {
	for _, c := range b {
		sh.h = (sh.h << 3) ^ sh.h ^ uint32(c)
	}
	return len(b), nil
}

type JenkinsStringHash32 struct {
	StringHash32
}

func NewJenkinsStringHash32() *JenkinsStringHash32 {
	return new(JenkinsStringHash32)
}

func (sh *JenkinsStringHash32) Write(b []byte) (int, os.Error) {
	for _, c := range b {
		sh.h += uint32(c)
		sh.h += (sh.h << 10)
		sh.h ^= (sh.h >> 6)
	}
	return len(b), nil
}

func (sh *JenkinsStringHash32) Sum32() uint32 {
	h := sh.h // copy so we don't mess with the internal state

        // Jenkins' finalize
	h += (h << 3)
	h ^= (h >> 11)
	h += (h << 15)

	return h
}

func (sh *JenkinsStringHash32) Sum() []byte {

	h := sh.Sum32()

	p := make([]byte, 4)
	p[0] = byte(h >> 24)
	p[1] = byte(h >> 16)
	p[2] = byte(h >> 8)
	p[3] = byte(h)
	return p
}
