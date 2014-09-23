// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wav

type factChunk struct {
	SampleLength [4]byte
}

// the 16-byte 'fmt' chunk
type fmtChunk16 struct {
	// Format code
	FormatTag uint16

	// Number of interleaved channels
	Channels uint16

	// Sampling rate (blocks per second)
	SamplesPerSec uint32

	// Data rate
	AvgBytesPerSec uint32

	// Data block size (bytes)
	BlockAlign uint16

	// Bits per sample
	BitsPerSample uint16
}

// the 18-byte 'fmt' chunk
type fmtChunk18 struct {
	// Size of the extension (0 or 22)
	Size uint16
}

// the 40-byte 'fmt' chunk
type fmtChunk40 struct {
	// Number of valid bits
	ValidBitsPerSample uint16

	// Speaker position mask
	ChannelMask uint32

	// GUID, including the data format code
	SubFormat [16]byte
}
