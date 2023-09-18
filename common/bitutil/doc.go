// Copyright 2018 The klaytn Authors
// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Adapted from: https://golang.org/src/crypto/cipher/xor.go
//
// This file is derived from common/bitutil/bitutil.go (2018/06/04).
// Modified and improved for the klaytn development.
//
// TODO-Klaytn: put the original LICENSE file in the common/bitutil directory

/*
Package bitutil implements fast bitwise operations and compression/decompressions.

Bitwise Operations

Following operations are supported
  - AND, OR, XOR operations
  - Provides both safe version and fast version of above operations
    `Safe` means it can be performed on all architectures
    `Fast` means it only can be performed on architecture which supports unaligned read/write

Compression and Decompression

Following operations are supported
  - CompressBytes
  - DecompressBytes

How compression works

The compression algorithm implemented by CompressBytes and DecompressBytes is
optimized for sparse input data which contains a lot of zero bytes. Decompression
requires knowledge of the decompressed data length.

	Compression works as follows:

	  if data only contains zeroes,
		  CompressBytes(data) == nil
	  otherwise if len(data) <= 1,
		  CompressBytes(data) == data
	  otherwise:
		  CompressBytes(data) == append(CompressBytes(nonZeroBitset(data)), nonZeroBytes(data)...)
		  where
			nonZeroBitset(data) is a bit vector with len(data) bits (MSB first):
				nonZeroBitset(data)[i/8] && (1 << (7-i%8)) != 0  if data[i] != 0
				len(nonZeroBitset(data)) == (len(data)+7)/8
			nonZeroBytes(data) contains the non-zero bytes of data in the same order
*/
package bitutil
