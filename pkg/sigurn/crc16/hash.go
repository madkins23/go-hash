package crc16

import (
	"github.com/sigurn/crc16"

	"github.com/madkins23/go-hash/pkg/hash"
)

// Size is the size of a CRC-16 checksum in bytes.
const Size = 2

type digest struct {
	crc   uint16
	table *crc16.Table
}

// New returns a new hash.Hash16 instance for the specified CRC16 parameters.
// The params determine the polynomial attributes used in the CRC calculation.
func New(params crc16.Params) hash.Hash16 {
	return NewFromTable(crc16.MakeTable(params))
}

// NewFromTable returns a new hash.Hash16 instance from the specified CRC16 table.
// The table contains the polynomial attributes used in the CRC calculation.
// This method can be used directly in order to reuse a single table over multiple calculations.
// Note that the CRC table is not reentrant, reuse must occur in a single execution thread.
func NewFromTable(table *crc16.Table) hash.Hash16 {
	d := &digest{table: table}
	d.Reset()
	return d
}

// The following methods implement the Hash16 interface,
// which includes the hash.Hash interface.

func (d *digest) Write(p []byte) (n int, err error) {
	d.crc = crc16.Update(d.crc, p, d.table)
	return len(p), nil
}

func (d *digest) Sum(b []byte) []byte {
	// TODO: not sure this is correct.
	return append(b, byte(d.crc>>8), byte(d.crc&0x0F))
}

func (d *digest) Reset() {
	d.crc = crc16.Init(d.table)
}

func (d *digest) Size() int {
	return Size
}

func (d *digest) BlockSize() int {
	return 1
}

func (d *digest) Sum16() uint16 {
	return crc16.Complete(d.crc, d.table)
}
