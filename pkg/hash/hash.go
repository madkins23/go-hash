/*
Package hash defines interfaces for hash sizes not supported by the standard library's [hash] package.

[hash]: https://pkg.go.dev/hash
*/
package hash

import (
	HASH "hash"
)

// Hash16 is the common interface implemented by all 16-bit hash functions.
// This follows the pattern of the hash package in the standard library.
type Hash16 interface {
	HASH.Hash
	Sum16() uint16
}

// Hash8 is the common interface implemented by all 8-bit hash functions.
// This follows the pattern of the hash package in the standard library.
type Hash8 interface {
	HASH.Hash
	Sum8() uint8
}
