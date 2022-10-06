package crc8

import (
	"github.com/sigurn/crc8"

	"github.com/madkins23/go-hash/pkg/hash"
)

// Size is the size of a CRC-8 checksum in bytes.
const Size = 1

type digest struct {
	crc   uint8
	table *crc8.Table
}

// New returns a new hash.Hash8 instance for the specified CRC8 parameters.
// The params determine the polynomial attributes used in the CRC calculation.
func New(params crc8.Params) hash.Hash8 {
	return NewFromTable(crc8.MakeTable(params))
}

// NewFromTable returns a new hash.Hash8 instance from the specified CRC8 table.
// The table contains the polynomial attributes used in the CRC calculation.
// This method can be used directly in order to reuse a single table over multiple calculations.
// Note that the CRC table is not reentrant, reuse must occur in a single execution thread.
func NewFromTable(table *crc8.Table) hash.Hash8 {
	d := &digest{table: table}
	d.Reset()
	return d
}

// The following methods implement the Hash8 interface,
// which includes the hash.Hash interface.

func (d *digest) Write(p []byte) (n int, err error) {
	d.crc = crc8.Update(d.crc, p, d.table)
	return len(p), nil
}

func (d *digest) Sum(b []byte) []byte {
	// Probably correct.
	return append(b, d.crc)
}

func (d *digest) Reset() {
	d.crc = crc8.Init(d.table)
}

func (d *digest) Size() int {
	return Size
}

func (d *digest) BlockSize() int {
	return 1
}

func (d *digest) Sum8() uint8 {
	return crc8.Complete(d.crc, d.table)
}
