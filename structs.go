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
	FormatTag [2]byte

	// Number of interleaved channels
	Channels [2]byte

	// Sampling rate (blocks per second)
	SamplesPerSec [4]byte

	// Data rate
	AvgBytesPerSec [4]byte

	// Data block size (bytes)
	BlockAlign [2]byte

	// Bits per sample
	BitsPerSample [2]byte
}

// the 18-byte 'fmt' chunk
type fmtChunk18 struct {
	// Size of the extension (0 or 22)
	Size [2]byte
}

// the 40-byte 'fmt' chunk
type fmtChunk40 struct {
	// Number of valid bits
	ValidBitsPerSample [2]byte

	// Speaker position mask
	ChannelMask [4]byte

	// GUID, including the data format code
	SubFormat [16]byte
}
