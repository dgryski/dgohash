// Unit tests for hash functions
// Copyright (c) 2011 Damian Gryski <damian@gryski.com>
// Licensed under the GPLv3, or at your option any later version.
package dgohash

import (
	"encoding/binary"
	"hash"
	"testing"
)

type _Golden struct {
	out uint32
	in  string
}

// These tables were all generated from reference C implementations of the associated hashes.

var golden_java = []_Golden{
	{0x00000000, ""},
	{0x00000061, "a"},
	{0x00000c21, "ab"},
	{0x00017862, "abc"},
	{0x002d9442, "abcd"},
	{0x0584f463, "abcde"},
	{0xab199863, "abcdef"},
	{0xb8197464, "abcdefg"},
	{0x4b151884, "abcdefgh"},
	{0x178df865, "abcdefghi"},
	{0xda3114a5, "abcdefghij"},
	{0x507cbe5d, "Discard medicine more than two years old."},
	{0xcf8332bc, "He who has a shady past knows that nice guys finish last."},
	{0x94ddaa0e, "I wouldn't marry him with a ten foot pole."},
	{0xd1a67f32, "Free! Free!/A trip/to Mars/for 900/empty jars/Burma Shave"},
	{0x29e1993d, "The days of the digital watch are numbered.  -Tom Stoppard"},
	{0x46b8e871, "Nepal premier won't resign."},
	{0x80a347dc, "For every action there is an equal and opposite government program."},
	{0xb560b45d, "His money is twice tainted: 'taint yours and 'taint mine."},
	{0x123c79c6, "There is no reason for any individual to have a computer in their home. -Ken Olsen, 1977"},
	{0x3f1ff283, "It's a tiny change to the code and not completely disgusting. - Bob Manchek"},
	{0xbf045f20, "size:  a.out:  bad magic"},
	{0x30642382, "The major problem is with sendmail.  -Mark Horton"},
	{0xf11f3607, "Give me a rock, paper and scissors and I will move the world.  CCFestoon"},
	{0xb68626c4, "If the enemy is within range, then so are you."},
	{0x872d8aba, "It's well we cannot hear the screams/That we create in others' dreams."},
	{0xd68213e8, "You remind me of a TV show, but that's all right: I watch it anyway."},
	{0xd55e6f3e, "C is as portable as Stonehedge!!"},
	{0xb34d3565, "Even if I could be Shakespeare, I think I should still choose to be Faraday. - A. Huxley"},
	{0x1f5a0d48, "The fugacity of a constituent in a mixture of gases at a given temperature is proportional to its mole fraction.  Lewis-Randall Rule"},
	{0xda3df8dd, "How can you write a big system without C++?  -Paul Glick"},
}

var golden_djb32 = []_Golden{
	{0x00001505, ""},
	{0x0002b606, "a"},
	{0x00597728, "ab"},
	{0x0b885c8b, "abc"},
	{0x7c93ee4f, "abcd"},
	{0x0f11b894, "abcde"},
	{0xf148cb7a, "abcdef"},
	{0x1a623b21, "abcdefg"},
	{0x66a99fa9, "abcdefgh"},
	{0x3bdd9532, "abcdefghi"},
	{0xb7903bdc, "abcdefghij"},
	{0xa61e3ba6, "Discard medicine more than two years old."},
	{0xf9827a7b, "He who has a shady past knows that nice guys finish last."},
	{0xa68ea4c5, "I wouldn't marry him with a ten foot pole."},
	{0xe31c5f19, "Free! Free!/A trip/to Mars/for 900/empty jars/Burma Shave"},
	{0xbaef90a4, "The days of the digital watch are numbered.  -Tom Stoppard"},
	{0xa2e15ace, "Nepal premier won't resign."},
	{0x3dd4f3e1, "For every action there is an equal and opposite government program."},
	{0xeffef6c6, "His money is twice tainted: 'taint yours and 'taint mine."},
	{0xbfd5d7e7, "There is no reason for any individual to have a computer in their home. -Ken Olsen, 1977"},
	{0x14a6762e, "It's a tiny change to the code and not completely disgusting. - Bob Manchek"},
	{0x9dc2ebc3, "size:  a.out:  bad magic"},
	{0x2fc35375, "The major problem is with sendmail.  -Mark Horton"},
	{0xbd0267c8, "Give me a rock, paper and scissors and I will move the world.  CCFestoon"},
	{0x682419cf, "If the enemy is within range, then so are you."},
	{0x82f44aeb, "It's well we cannot hear the screams/That we create in others' dreams."},
	{0x41db5feb, "You remind me of a TV show, but that's all right: I watch it anyway."},
	{0xa3b3be6d, "C is as portable as Stonehedge!!"},
	{0x42b489b4, "Even if I could be Shakespeare, I think I should still choose to be Faraday. - A. Huxley"},
	{0x57e38ab3, "The fugacity of a constituent in a mixture of gases at a given temperature is proportional to its mole fraction.  Lewis-Randall Rule"},
	{0x9f8d455a, "How can you write a big system without C++?  -Paul Glick"},
}

var golden_elf32 = []_Golden{
	{0x00000000, ""},
	{0x00000061, "a"},
	{0x00000672, "ab"},
	{0x00006783, "abc"},
	{0x00067894, "abcd"},
	{0x006789a5, "abcde"},
	{0x06789ab6, "abcdef"},
	{0x0789aba7, "abcdefg"},
	{0x089abaa8, "abcdefgh"},
	{0x09abaa69, "abcdefghi"},
	{0x0abaa66a, "abcdefghij"},
	{0x0ab8c77e, "Discard medicine more than two years old."},
	{0x0c2895ee, "He who has a shady past knows that nice guys finish last."},
	{0x0d88846e, "I wouldn't marry him with a ten foot pole."},
	{0x00f84415, "Free! Free!/A trip/to Mars/for 900/empty jars/Burma Shave"},
	{0x0ffe12f4, "The days of the digital watch are numbered.  -Tom Stoppard"},
	{0x0ce8fd4e, "Nepal premier won't resign."},
	{0x0db274ae, "For every action there is an equal and opposite government program."},
	{0x00bd1fee, "His money is twice tainted: 'taint yours and 'taint mine."},
	{0x0c80df37, "There is no reason for any individual to have a computer in their home. -Ken Olsen, 1977"},
	{0x0b49043b, "It's a tiny change to the code and not completely disgusting. - Bob Manchek"},
	{0x04724b83, "size:  a.out:  bad magic"},
	{0x02955e6e, "The major problem is with sendmail.  -Mark Horton"},
	{0x035111fe, "Give me a rock, paper and scissors and I will move the world.  CCFestoon"},
	{0x0a07b02e, "If the enemy is within range, then so are you."},
	{0x0c2c655e, "It's well we cannot hear the screams/That we create in others' dreams."},
	{0x0e8fc43e, "You remind me of a TV show, but that's all right: I watch it anyway."},
	{0x02450da1, "C is as portable as Stonehedge!!"},
	{0x03568a09, "Even if I could be Shakespeare, I think I should still choose to be Faraday. - A. Huxley"},
	{0x0aa09cd5, "The fugacity of a constituent in a mixture of gases at a given temperature is proportional to its mole fraction.  Lewis-Randall Rule"},
	{0x0810f11b, "How can you write a big system without C++?  -Paul Glick"},
}

var golden_sdbm = []_Golden{
	{0x00000000, ""},
	{0x00000061, "a"},
	{0x00611841, "ab"},
	{0x3025f862, "abc"},
	{0xd1ba2082, "abcd"},
	{0xbd500063, "abcde"},
	{0x971318c3, "abcdef"},
	{0x46761864, "abcdefg"},
	{0x6f740104, "abcdefgh"},
	{0x6e904065, "abcdefghi"},
	{0x75e4d945, "abcdefghij"},
	{0x046d355d, "Discard medicine more than two years old."},
	{0x718c9e9c, "He who has a shady past knows that nice guys finish last."},
	{0x14c663ae, "I wouldn't marry him with a ten foot pole."},
	{0xf21ea712, "Free! Free!/A trip/to Mars/for 900/empty jars/Burma Shave"},
	{0x2ab38c1d, "The days of the digital watch are numbered.  -Tom Stoppard"},
	{0x354c9f71, "Nepal premier won't resign."},
	{0x8b82905c, "For every action there is an equal and opposite government program."},
	{0x2157591d, "His money is twice tainted: 'taint yours and 'taint mine."},
	{0xdda5cb46, "There is no reason for any individual to have a computer in their home. -Ken Olsen, 1977"},
	{0x87619563, "It's a tiny change to the code and not completely disgusting. - Bob Manchek"},
	{0x2dfefd80, "size:  a.out:  bad magic"},
	{0x541955e2, "The major problem is with sendmail.  -Mark Horton"},
	{0xe8d7cbc7, "Give me a rock, paper and scissors and I will move the world.  CCFestoon"},
	{0x434fdd24, "If the enemy is within range, then so are you."},
	{0x3bd7247a, "It's well we cannot hear the screams/That we create in others' dreams."},
	{0x777bd008, "You remind me of a TV show, but that's all right: I watch it anyway."},
	{0x60ac769e, "C is as portable as Stonehedge!!"},
	{0x65db3345, "Even if I could be Shakespeare, I think I should still choose to be Faraday. - A. Huxley"},
	{0x3a182aa8, "The fugacity of a constituent in a mixture of gases at a given temperature is proportional to its mole fraction.  Lewis-Randall Rule"},
	{0x2129ea9d, "How can you write a big system without C++?  -Paul Glick"},
}

var golden_sqlite = []_Golden{
	{0x00000000, ""},
	{0x00000061, "a"},
	{0x0000030b, "ab"},
	{0x00001b30, "abc"},
	{0x0000c2d4, "abcd"},
	{0x0006d411, "abcde"},
	{0x003074ff, "abcdef"},
	{0x01b3d360, "abcdefg"},
	{0x0c2d4808, "abcdefgh"},
	{0x6d470821, "abcdefghi"},
	{0x077f4943, "abcdefghij"},
	{0xbfd88981, "Discard medicine more than two years old."},
	{0x3d7ed466, "He who has a shady past knows that nice guys finish last."},
	{0x1d05fae6, "I wouldn't marry him with a ten foot pole."},
	{0x68662562, "Free! Free!/A trip/to Mars/for 900/empty jars/Burma Shave"},
	{0x92d13743, "The days of the digital watch are numbered.  -Tom Stoppard"},
	{0x0c85b42d, "Nepal premier won't resign."},
	{0xb4edacc0, "For every action there is an equal and opposite government program."},
	{0xe6412c11, "His money is twice tainted: 'taint yours and 'taint mine."},
	{0x0ff516f4, "There is no reason for any individual to have a computer in their home. -Ken Olsen, 1977"},
	{0x575c3671, "It's a tiny change to the code and not completely disgusting. - Bob Manchek"},
	{0xd6f1fcda, "size:  a.out:  bad magic"},
	{0x07ef9f3a, "The major problem is with sendmail.  -Mark Horton"},
	{0x08f40b55, "Give me a rock, paper and scissors and I will move the world.  CCFestoon"},
	{0x16c45448, "If the enemy is within range, then so are you."},
	{0x62c3718a, "It's well we cannot hear the screams/That we create in others' dreams."},
	{0x898c6d90, "You remind me of a TV show, but that's all right: I watch it anyway."},
	{0xe13a717a, "C is as portable as Stonehedge!!"},
	{0x71f77adf, "Even if I could be Shakespeare, I think I should still choose to be Faraday. - A. Huxley"},
	{0x658886ca, "The fugacity of a constituent in a mixture of gases at a given temperature is proportional to its mole fraction.  Lewis-Randall Rule"},
	{0xb1d7f9e5, "How can you write a big system without C++?  -Paul Glick"},
}

var golden_jenkins = []_Golden{
	{0x00000000, ""},
	{0xca2e9442, "a"},
	{0x45e61e58, "ab"},
	{0xed131f5b, "abc"},
	{0xcd8b6206, "abcd"},
	{0xb98559fc, "abcde"},
	{0x0161526f, "abcdef"},
	{0x4ac70178, "abcdefg"},
	{0x44d2d3e1, "abcdefgh"},
	{0xc8b4ca7d, "abcdefghi"},
	{0x7031289d, "abcdefghij"},
	{0x41454415, "Discard medicine more than two years old."},
	{0xd6995686, "He who has a shady past knows that nice guys finish last."},
	{0xd77ff8d6, "I wouldn't marry him with a ten foot pole."},
	{0x353105d6, "Free! Free!/A trip/to Mars/for 900/empty jars/Burma Shave"},
	{0x2599d1ab, "The days of the digital watch are numbered.  -Tom Stoppard"},
	{0x402a21c2, "Nepal premier won't resign."},
	{0xcc522896, "For every action there is an equal and opposite government program."},
	{0xa869b6fb, "His money is twice tainted: 'taint yours and 'taint mine."},
	{0x1a8b3dcd, "There is no reason for any individual to have a computer in their home. -Ken Olsen, 1977"},
	{0x660a13c1, "It's a tiny change to the code and not completely disgusting. - Bob Manchek"},
	{0x7878a798, "size:  a.out:  bad magic"},
	{0x66e9dba8, "The major problem is with sendmail.  -Mark Horton"},
	{0xbc1b46f0, "Give me a rock, paper and scissors and I will move the world.  CCFestoon"},
	{0x4f6762bf, "If the enemy is within range, then so are you."},
	{0x183959f7, "It's well we cannot hear the screams/That we create in others' dreams."},
	{0x6aff9b36, "You remind me of a TV show, but that's all right: I watch it anyway."},
	{0xb9699852, "C is as portable as Stonehedge!!"},
	{0xa4fde64f, "Even if I could be Shakespeare, I think I should still choose to be Faraday. - A. Huxley"},
	{0xf162dacb, "The fugacity of a constituent in a mixture of gases at a given temperature is proportional to its mole fraction.  Lewis-Randall Rule"},
	{0x2d3ac755, "How can you write a big system without C++?  -Paul Glick"},
}

var golden_murmur3 = []_Golden{
	{0x00000000, ""},
	{0x3c2569b2, "a"},
	{0x9bbfd75f, "ab"},
	{0xb3dd93fa, "abc"},
	{0x43ed676a, "abcd"},
	{0xe89b9af6, "abcde"},
	{0x6181c085, "abcdef"},
	{0x883c9b06, "abcdefg"},
	{0x49ddccc4, "abcdefgh"},
	{0x421406f0, "abcdefghi"},
	{0x88927791, "abcdefghij"},
	{0x91e056d3, "Discard medicine more than two years old."},
	{0xc4d1cdf9, "He who has a shady past knows that nice guys finish last."},
	{0x92a09da9, "I wouldn't marry him with a ten foot pole."},
	{0xba22e6c4, "Free! Free!/A trip/to Mars/for 900/empty jars/Burma Shave"},
	{0xb3ba11cb, "The days of the digital watch are numbered.  -Tom Stoppard"},
	{0x941ada4d, "Nepal premier won't resign."},
	{0x03f1f7b4, "For every action there is an equal and opposite government program."},
	{0x03946117, "His money is twice tainted: 'taint yours and 'taint mine."},
	{0x91e89ce1, "There is no reason for any individual to have a computer in their home. -Ken Olsen, 1977"},
	{0xdc39bd00, "It's a tiny change to the code and not completely disgusting. - Bob Manchek"},
	{0xe898a1fa, "size:  a.out:  bad magic"},
	{0xcb5affb4, "The major problem is with sendmail.  -Mark Horton"},
	{0xc84510d4, "Give me a rock, paper and scissors and I will move the world.  CCFestoon"},
	{0xd4466554, "If the enemy is within range, then so are you."},
	{0xe718d618, "It's well we cannot hear the screams/That we create in others' dreams."},
	{0xa6fb1684, "You remind me of a TV show, but that's all right: I watch it anyway."},
	{0x65cb8d60, "C is as portable as Stonehedge!!"},
	{0x164935d1, "Even if I could be Shakespeare, I think I should still choose to be Faraday. - A. Huxley"},
	{0x33e03966, "The fugacity of a constituent in a mixture of gases at a given temperature is proportional to its mole fraction.  Lewis-Randall Rule"},
	{0x04944630, "How can you write a big system without C++?  -Paul Glick"},
}

var golden_superfast = []_Golden{
	{0x00000000, ""},
	{0x93642e87, "a"},
	{0x5b8c0ec3, "ab"},
	{0xe5186b3a, "abc"},
	{0x3ab452d8, "abcd"},
	{0x84786722, "abcde"},
	{0xbe7c6fe4, "abcdef"},
	{0x3dad41af, "abcdefg"},
	{0xff7cfe86, "abcdefgh"},
	{0xa73e3541, "abcdefghi"},
	{0x2d7c0783, "abcdefghij"},
	{0xf3f9c606, "Discard medicine more than two years old."},
	{0x1d68aee7, "He who has a shady past knows that nice guys finish last."},
	{0xb6929c96, "I wouldn't marry him with a ten foot pole."},
	{0x3a79f2c8, "Free! Free!/A trip/to Mars/for 900/empty jars/Burma Shave"},
	{0xc2169976, "The days of the digital watch are numbered.  -Tom Stoppard"},
	{0x24d1092a, "Nepal premier won't resign."},
	{0x7dcdc1cf, "For every action there is an equal and opposite government program."},
	{0x1004d947, "His money is twice tainted: 'taint yours and 'taint mine."},
	{0x5237d840, "There is no reason for any individual to have a computer in their home. -Ken Olsen, 1977"},
	{0x193828c4, "It's a tiny change to the code and not completely disgusting. - Bob Manchek"},
	{0xcf2cd792, "size:  a.out:  bad magic"},
	{0xee993cb6, "The major problem is with sendmail.  -Mark Horton"},
	{0xb6c84172, "Give me a rock, paper and scissors and I will move the world.  CCFestoon"},
	{0x3b4039cf, "If the enemy is within range, then so are you."},
	{0x5659e64b, "It's well we cannot hear the screams/That we create in others' dreams."},
	{0x52ddc48a, "You remind me of a TV show, but that's all right: I watch it anyway."},
	{0xd650693f, "C is as portable as Stonehedge!!"},
	{0x5a5737f0, "Even if I could be Shakespeare, I think I should still choose to be Faraday. - A. Huxley"},
	{0xcac073c5, "The fugacity of a constituent in a mixture of gases at a given temperature is proportional to its mole fraction.  Lewis-Randall Rule"},
	{0x494c35dd, "How can you write a big system without C++?  -Paul Glick"},
}

var golden_marvin = []_Golden{
	{0xf7f2c954, ""},
	{0xd46e71f7, "a"},
	{0xb40c651c, "ab"},
	{0x5b3bc23d, "abc"},
	{0x6b15e57b, "abcd"},
	{0x601e6ea8, "abcde"},
	{0xfc18bd2c, "abcdef"},
	{0x79b01bfb, "abcdefg"},
	{0x54793238, "abcdefgh"},
	{0xebf98191, "abcdefghi"},
	{0x68a8001d, "abcdefghij"},
	{0x659105c1, "Discard medicine more than two years old."},
	{0xb98b31d, "He who has a shady past knows that nice guys finish last."},
	{0xbae17c9a, "I wouldn't marry him with a ten foot pole."},
	{0x9a299f69, "Free! Free!/A trip/to Mars/for 900/empty jars/Burma Shave"},
	{0xb463d704, "The days of the digital watch are numbered.  -Tom Stoppard"},
	{0xe6059c5f, "Nepal premier won't resign."},
	{0xbdd4f772, "For every action there is an equal and opposite government program."},
	{0x12af7ede, "His money is twice tainted: 'taint yours and 'taint mine."},
	{0x1e9cae8, "There is no reason for any individual to have a computer in their home. -Ken Olsen, 1977"},
	{0xcb683e33, "It's a tiny change to the code and not completely disgusting. - Bob Manchek"},
	{0x2074fbfa, "size:  a.out:  bad magic"},
	{0x52abb615, "The major problem is with sendmail.  -Mark Horton"},
	{0x5a509711, "Give me a rock, paper and scissors and I will move the world.  CCFestoon"},
	{0xf97f5273, "If the enemy is within range, then so are you."},
	{0x494c0cb, "It's well we cannot hear the screams/That we create in others' dreams."},
	{0x7150a3c0, "You remind me of a TV show, but that's all right: I watch it anyway."},
	{0xc5f56430, "C is as portable as Stonehedge!!"},
	{0x712bcf01, "Even if I could be Shakespeare, I think I should still choose to be Faraday. - A. Huxley"},
	{0xedd44de6, "The fugacity of a constituent in a mixture of gases at a given temperature is proportional to its mole fraction.  Lewis-Randall Rule"},
	{0xd9440105, "How can you write a big system without C++?  -Paul Glick"},
}

func TestJava(t *testing.T) {
	testGolden(t, NewJava32(), golden_java, "java")
}

func TestDbj32(t *testing.T) {
	testGolden(t, NewDjb32(), golden_djb32, "djb")
}

func TestElf32(t *testing.T) {
	testGolden(t, NewElf32(), golden_elf32, "elf32")
}

func TestSDBM(t *testing.T) {
	testGolden(t, NewSDBM32(), golden_sdbm, "sdbm")
}

func TestSqlite3(t *testing.T) {
	testGolden(t, NewSQLite32(), golden_sqlite, "sqlite3")
}

func TestJenkins(t *testing.T) {
	testGolden(t, NewJenkins32(), golden_jenkins, "jenkins")
}

func TestMurmur(t *testing.T) {

	// test the incremental hashing logic
	m := NewMurmur3_x86_32()

	testIncremental(t, m, 0xe0c9df28, "murmur3")

	testGolden(t, m, golden_murmur3, "murmur3")

	// add murmur's own verification test here?

}

func TestSuperFastHash(t *testing.T) {

	// test the incremental hashing logic
	m := NewSuperFastHash()

	testIncremental(t, m, 0x54de96ed, "superfast")

	testGolden(t, m, golden_superfast, "superfast")
}

func TestMarvin(t *testing.T) {

	m := NewMarvin32(0x5D70D359C498B3F8) // seed for testing

	// test the incremental hashing logic
	testIncremental(t, m, 0x28685e7a, "marvin")

	testGolden(t, m, golden_marvin, "marvin")
}

func BenchmarkJava32(b *testing.B) {
	commonBench(b, NewJava32(), golden_java)
}

func BenchmarkDJB(b *testing.B) {
	commonBench(b, NewDjb32(), golden_djb32)
}

func BenchmarkElf32(b *testing.B) {
	commonBench(b, NewElf32(), golden_elf32)
}

func BenchmarkJenkins32(b *testing.B) {
	commonBench(b, NewJenkins32(), golden_jenkins)
}

func BenchmarkMarvin32(b *testing.B) {
	commonBench(b, NewMarvin32(0), golden_marvin)
}

func BenchmarkMurmur(b *testing.B) {
	commonBench(b, NewMurmur3_x86_32(), golden_murmur3)
}

func BenchmarkSDBM32(b *testing.B) {
	commonBench(b, NewSDBM32(), golden_sdbm)
}

func BenchmarkSQLite32(b *testing.B) {
	commonBench(b, NewSQLite32(), golden_sqlite)
}

func BenchmarkSuperFastHash(b *testing.B) {
	commonBench(b, NewSuperFastHash(), golden_superfast)
}

func commonBench(b *testing.B, h hash.Hash32, golden []_Golden) {
	for i := 0; i < b.N; i++ {
		for _, g := range golden {
			h.Reset()
			h.Write([]byte(g.in))
			h.Sum32()
		}
	}
}

func testIncremental(t *testing.T, h hash.Hash32, result uint32, which string) {

	h.Reset()

	h.Write([]byte("hello"))
	h.Write([]byte("h"))
	h.Write([]byte("e"))
	h.Write([]byte("l"))
	h.Write([]byte("l"))
	h.Write([]byte("o"))
	h.Write([]byte("hellohello"))

	h32 := h.Sum32()

	if h32 != result {
		t.Errorf("%s: incremental failed: got %08x", which, h32)
	}

	h.Reset()
	h.Write([]byte("hellohellohellohello"))

	h32 = h.Sum32()

	if h32 != result {
		t.Errorf("%s: failed: got %08x", which, h32)
	}
}

func testGolden(t *testing.T, h hash.Hash32, golden []_Golden, which string) {

	for _, g := range golden {
		h.Reset()
		h.Write([]byte(g.in))

		sum := h.Sum32()

		if sum != g.out {
			t.Errorf("%s(%s) = 0x%x want 0x%x", which, g.in, sum, g.out)
		}

		bsum := h.Sum(nil)

		if len(bsum) != 4 {
			t.Errorf("%s Sum(nil) returned %d bytes, wanted 4: %s\n", which, len(bsum), bsum)
		}

		s := binary.BigEndian.Uint32(bsum)

		if s != sum {
			t.Errorf("%s(%s).Sum(nil) = 0x%x want 0x%x", which, g.in, sum, g.out)
		}

		bsum = h.Sum([]byte{0x01, 0x02, 0x03, 0x04})

		if len(bsum) != 8 {
			t.Errorf("%s Sum(bsum) returned %d bytes, wanted 8: %x\n", which, len(bsum), bsum)
		}

		s = binary.BigEndian.Uint32(bsum[0:])
		s2 := binary.BigEndian.Uint32(bsum[4:])

		if s != 0x01020304 || s2 != sum {
			t.Errorf("%s(%s).Sum(bsum) = %x (expected 0x01020304 %x )", which, g.in, bsum, sum)
		}

	}
}
