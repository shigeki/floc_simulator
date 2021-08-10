package floc

//
// Most of this code comes from https://github.com/creachadair/cityhash
// 

import (
	"encoding/binary"
)

const (
	k0 = uint64(0xc3a5c85c97cb3127)
	k1 = uint64(0xb492b66fbe98f273)
	k2 = uint64(0x9ae16a3b2f90404f)
	k3 = uint64(0xc949d7c7509e6557)
)

var fetch64 = binary.LittleEndian.Uint64 // :: []byte -> uint64
var fetch32 = binary.LittleEndian.Uint32 // :: []byte -> uint32

func shiftMix(val uint64) uint64 { return val ^ (val >> 47) }

func bswap64(in uint64) uint64 {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], in)
	return binary.BigEndian.Uint64(buf[:])
}

func ror64(val, shift uint64) uint64 {
	// Avoid shifting by 64: doing so yields an undefined result.
	if shift != 0 {
		return val>>shift | val<<(64-shift)
	}
	return val
}

func rotateByAtLeast1(val uint64, shift uint64) uint64 {
	return ((val >> shift) | (val << (64 - shift)))
}

func hash64Len16(u, v uint64) uint64 { return hash128to64(u, v) }

func hash64Len16Mul(u, v, mul uint64) uint64 {
	// Murmur-inspired hashing.
	a := (u ^ v) * mul
	a ^= (a >> 47)
	b := (v ^ a) * mul
	b ^= (b >> 47)
	b *= mul
	return b
}

func hash64Len0to16(s []byte) uint64 {
	n := uint64(len(s))
	if n > 8 {
		a := fetch64(s)
		b := fetch64(s[(n - 8):])
		c := rotateByAtLeast1(b + n, n)
		return hash64Len16(a, c) ^ b
	}
	if n >= 4 {
		a := uint64(fetch32(s))
		return hash64Len16(n+(a<<3), uint64(fetch32(s[n-4:])))
	}
	if n > 0 {
		a := s[0]
		b := s[n>>1]
		c := s[n-1]
		y := uint32(a) + uint32(b)<<8
		z := uint32(n) + uint32(c)<<2
		return shiftMix(uint64(y)*k2^uint64(z)*k3) * k2
	}
	return k2
}

func hash128to64(lo, hi uint64) uint64 {
	// Murmur-inspired hashing.
	const mul = uint64(0x9ddfea08eb382d69)
	a := (lo ^ hi) * mul
	a ^= (a >> 47)
	b := (hi ^ a) * mul
	b ^= (b >> 47)
	b *= mul
	return b
}

func hash64Len17to32(s []byte) uint64 {
	n := uint64(len(s))
	a := fetch64(s) * k1
	b := fetch64(s[8:])
	c := fetch64(s[n-8:]) * k2
	d := fetch64(s[n-16:]) * k0
	return hash64Len16(ror64(a -b, 43)+ror64(c, 30)+d, a+ror64(b^k3, 20)-c)
}

func hash64Len33to64(s []byte) uint64 {
	n := uint64(len(s))
	mul := k2 + n*2
	a := fetch64(s) * k2
	b := fetch64(s[8:])
	c := fetch64(s[n-24:])
	d := fetch64(s[n-32:])
	e := fetch64(s[16:]) * k2
	f := fetch64(s[24:]) * 9
	g := fetch64(s[n-8:])
	h := fetch64(s[n-16:]) * mul
	u := ror64(a+g, 43) + (ror64(b, 30)+c)*9
	v := ((a + g) ^ d) + f + 1
	w := bswap64((u+v)*mul) + h
	x := ror64(e+f, 42) + c
	y := (bswap64((v+w)*mul) + g) * mul
	z := e + f + c
	a = bswap64((x+z)*mul+y) + b
	b = shiftMix((z+a)*mul+d+h) * mul
	return b + x
}

func weakHashLen32WithSeeds(s []byte, a, b uint64) (uint64, uint64) {
	// Note: Was two overloads of WeakHashLen32WithSeeds.  The second is only
	// ever called from the first, so I inlined it.
	w := fetch64(s)
	x := fetch64(s[8:])
	y := fetch64(s[16:])
	z := fetch64(s[24:])

	a += w
	b = ror64(b+a+z, 21)
	c := a
	a += x
	a += y
	b += ror64(a, 44)
	return a + z, b + c
}

func CityHash64V103(s []byte) uint64 {
	n := uint64(len(s))
	if n <= 32 {
		if n <= 16 {
			return hash64Len0to16(s)
		}
		return hash64Len17to32(s)
	} else if n <= 64 {
		return hash64Len33to64(s)
	}

	// For strings over 64 bytes we hash the end first, and then as we loop we
	// keep 56 bytes of state: v, w, x, y, and z.
	x := fetch64(s[n-40:])
	y := fetch64(s[n-16:]) + fetch64(s[n-56:])
	z := hash64Len16(fetch64(s[n-48:])+n, fetch64(s[n-24:]))

	v1, v2 := weakHashLen32WithSeeds(s[n-64:], n, z)
	w1, w2 := weakHashLen32WithSeeds(s[n-32:], y+k1, x)
	x = x*k1 + fetch64(s)

	// Decrease n to the nearest multiple of 64, and operate on 64-byte chunks.
	n = (n - 1) &^ 63
	for {
		x = ror64(x+y+v1+fetch64(s[8:]), 37) * k1
		y = ror64(y+v2+fetch64(s[48:]), 42) * k1
		x ^= w2
		y += v1 + fetch64(s[40:])
		z = ror64(z+w1, 33) * k1
		v1, v2 = weakHashLen32WithSeeds(s, v2*k1, x+w1)
		w1, w2 = weakHashLen32WithSeeds(s[32:], z+w2, y+fetch64(s[16:]))
		z, x = x, z
		s = s[64:]
		n -= 64
		if n == 0 {
			break
		}
	}
	return hash64Len16(hash64Len16(v1, w1)+shiftMix(y)*k1+z, hash64Len16(v2, w2)+x)
}

func CityHash64WithSeedsV103(s []byte, seed0, seed1 uint64) uint64 {
	return hash64Len16(CityHash64V103(s)-seed0, seed1)
}

func CityHash64WithSeedV103(s []byte, seed uint64) uint64 {
	return CityHash64WithSeedsV103(s, k2, seed)
}

