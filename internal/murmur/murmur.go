package murmur

import (
	"encoding/binary"
)

// Cassandra uses special Murmur3 128bit version that is implemented using signed int64-s.
// This IS NOT a standard Murmur3 but special cassandra murmur3.
func CassandraMurmur3H1(data []byte) int64 {
	length := len(data)

	var h1, h2, k1, k2 int64

	const (
		c1 int64 = -8663945395140668459 // 0x87c37b91114253d5
		c2 int64 = 0x4cf5ad432745937f
	)

	nBlocks := length / 16
	for i := 0; i < nBlocks; i++ {
		k1 = int64(binary.LittleEndian.Uint64(data[i*16:]))
		k2 = int64(binary.LittleEndian.Uint64(data[(i*16)+8:]))

		k1 *= c1
		k1 = (k1 << 31) | int64(uint64(k1)>>33) // ROTL64(k1, 31)
		k1 *= c2
		h1 ^= k1

		h1 = (h1 << 27) | int64(uint64(h1)>>37) // ROTL64(h1, 27)
		h1 += h2
		h1 = h1*5 + 0x52dce729

		k2 *= c2
		k2 = (k2 << 33) | int64(uint64(k2)>>31) // ROTL64(k2, 33)
		k2 *= c1
		h2 ^= k2

		h2 = (h2 << 31) | int64(uint64(h2)>>33) // ROTL64(h2, 31)
		h2 += h1
		h2 = h2*5 + 0x38495ab5
	}

	// tail
	tail := data[nBlocks*16:]
	k1 = 0
	k2 = 0
	switch length & 15 {
	case 15:
		k2 ^= int64(int8(tail[14])) << 48
		fallthrough
	case 14:
		k2 ^= int64(int8(tail[13])) << 40
		fallthrough
	case 13:
		k2 ^= int64(int8(tail[12])) << 32
		fallthrough
	case 12:
		k2 ^= int64(int8(tail[11])) << 24
		fallthrough
	case 11:
		k2 ^= int64(int8(tail[10])) << 16
		fallthrough
	case 10:
		k2 ^= int64(int8(tail[9])) << 8
		fallthrough
	case 9:
		k2 ^= int64(int8(tail[8]))

		k2 *= c2
		k2 = (k2 << 33) | int64(uint64(k2)>>31) // ROTL64(k2, 33)
		k2 *= c1
		h2 ^= k2

		fallthrough
	case 8:
		k1 ^= int64(int8(tail[7])) << 56
		fallthrough
	case 7:
		k1 ^= int64(int8(tail[6])) << 48
		fallthrough
	case 6:
		k1 ^= int64(int8(tail[5])) << 40
		fallthrough
	case 5:
		k1 ^= int64(int8(tail[4])) << 32
		fallthrough
	case 4:
		k1 ^= int64(int8(tail[3])) << 24
		fallthrough
	case 3:
		k1 ^= int64(int8(tail[2])) << 16
		fallthrough
	case 2:
		k1 ^= int64(int8(tail[1])) << 8
		fallthrough
	case 1:
		k1 ^= int64(int8(tail[0]))

		k1 *= c1
		k1 = (k1 << 31) | int64(uint64(k1)>>33) // ROTL64(k1, 31)
		k1 *= c2
		h1 ^= k1
	}

	h1 ^= int64(length)
	h2 ^= int64(length)

	h1 += h2
	h2 += h1

	// finalizer
	const (
		fmix1 int64 = -49064778989728563   // 0xff51afd7ed558ccd
		fmix2 int64 = -4265267296055464877 // 0xc4ceb9fe1a85ec53
	)

	// fmix64(h1)
	h1 ^= int64(uint64(h1) >> 33)
	h1 *= fmix1
	h1 ^= int64(uint64(h1) >> 33)
	h1 *= fmix2
	h1 ^= int64(uint64(h1) >> 33)

	// fmix64(h2)
	h2 ^= int64(uint64(h2) >> 33)
	h2 *= fmix1
	h2 ^= int64(uint64(h2) >> 33)
	h2 *= fmix2
	h2 ^= int64(uint64(h2) >> 33)

	h1 += h2
	// the following is extraneous since h2 is discarded
	// h2 += h1

	return h1
}
