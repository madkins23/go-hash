# go-hash

Wrap smaller hash libraries with interfaces similar to those in standard library hash interface.

See the [source](https://github.com/madkins23/go-hash)
or [godoc](https://godoc.org/github.com/madkins23/go-hash) for more detailed documentation.

![GitHub](https://img.shields.io/github/license/madkins23/go-hash)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/madkins23/go-hash)
[![Go Reference](https://pkg.go.dev/badge/github.com/madkins23/go-hash.svg)](https://pkg.go.dev/github.com/madkins23/go-hash)

## Background

The Go standard library defines the [`hash`](https://pkg.go.dev/hash) package.
This package defines a [`Hash`](https://pkg.go.dev/hash#Hash) interface
and then specific size interfaces are derived for
[`Hash32`](https://pkg.go.dev/hash#Hash32)
and [`Hash64`](https://pkg.go.dev/hash#Hash64).
Subdirectories provide implementations of these interfaces for specific hash functions
as well as the hashing function used by Go maps.

One nice feature (perhaps _the_ nice feature)
of the `Hash` interface is that it derives from `io.Writer`.
This allows types that implement `Hash` to be used with
[`io.Copy()`](https://pkg.go.dev/io#Copy).

## Implementation

This project provides missing interfaces:

* `Hash16`
* `Hash8`

and wrappers around pre-existing hash libraries that implement these interfaces.

* `github.com/madkins23/go-hash/pkg/sigurn/ccr8`
* `github.com/madkins23/go-hash/pkg/sigurn/ccr16`

## Usage

    import (
        crc8hash "github.com/madkins23/go-hash/pkg/sigurn/crc8"
        "github.com/sigurn/crc8" // needed to choose CRC variant
    )

    const someFile = "someFile.ext"

    if file, err := os.Open(someFile); err == nil {
        h8 := crc8hash.New(crc8.CRC8)
        if _, err := io.Copy(h8, file); err == nil {
            fmt.Printf("%03d %s\n", h8.Sum8(), someFile)
        }
        _ = file.Close()
    }
