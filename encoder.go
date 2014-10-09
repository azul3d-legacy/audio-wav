// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wav

import (
	"bufio"
	"encoding/binary"
	"io"
	"os"

	"azul3d.org/audio.v1"
)

// An encoder is capable of encoding audio samples to a WAV file.
type encoder struct {
	// A buffered writer, wrapping write operations to ws.
	bw *bufio.Writer
	// Underlying io.WriteSeeker to which the WAV file is written to.
	ws io.WriteSeeker
	// Audio configuration; including sample rate and number of channels.
	conf audio.Config
	// nsamples specifies the total number of samples written from all channels.
	nsamples uint32
	// bps represents the number of bits-per-sample used to encode audio samples.
	bps uint8
}

// NewEncoder creates a new WAV encoder, which stores the audio configuration in
// a WAV header and encodes any audio samples written to it. The contents of the
// WAV header and the encoded audio samples are written to w.
//
// Note: The Close method of the encoder must be called when finished using it.
func NewEncoder(w io.WriteSeeker, conf audio.Config) (audio.Encoder, error) {
	// Write WAV file header to w, based on the audio configuration.
	// TODO(u): Add output support for additional audio sample format; instead of
	// only using 16-bit PCM.
	enc := &encoder{bw: bufio.NewWriter(w), ws: w, conf: conf, bps: 16}
	err := enc.writeHeader()
	if err != nil {
		return nil, err
	}

	// Return encoder which encodes the audio samples written to it and stores
	// writes those to w.
	return enc, nil
}

// Write attempts to write all, b.Len(), samples in the slice to the
// writer.
//
// Returned is the number of samples from the slice that where wrote to
// the writer, and an error if any occurred.
//
// If the number of samples wrote is less than buf.Len() then the returned
// error must be non-nil. If any error occurs it should be considered fatal
// with regards to the writer: no more data can be subsequently wrote after
// an error.
func (enc *encoder) Write(b audio.Slice) (n int, err error) {
	// The at closure returns the i:th sample of b at a byte slice.
	var buf [3]byte
	var at func(i int) []byte
	switch v := b.(type) {
	case audio.PCM8Samples:
		panic("not yet implemented.")
		at = func(i int) []byte {
			// Unsigned 8-bit PCM audio sample.
			buf[0] = uint8(0x80 + v[i])
			return buf[:1]
		}
	case audio.PCM16Samples:
		at = func(i int) []byte {
			// Signed 16-bit PCM audio sample.
			sample := v[i]
			buf[0] = uint8(sample)
			buf[1] = uint8(sample >> 8)
			return buf[:2]
		}
	case audio.PCM32Samples:
		panic("not yet implemented.")
		at = func(i int) []byte {
			// Signed 32-bit PCM audio sample.
			sample := v[i]
			buf[0] = uint8(sample)
			buf[1] = uint8(sample >> 8)
			buf[2] = uint8(sample >> 16)
			return buf[:3]
		}
	default:
		at = func(i int) []byte {
			// Generic implementation.
			// TODO(u): Update to support 32-bit PCM audio samples, once the rest of
			// the encoder does so.
			sample := audio.F64ToPCM16(b.At(i))
			buf[0] = uint8(sample)
			buf[1] = uint8(sample >> 8)
			return buf[:2]
		}
	}

	// Generic implementation.
	for ; n < b.Len(); n++ {
		buf := at(n)
		m, err := enc.bw.Write(buf)
		if err != nil {
			return n, err
		}
		if m < len(buf) {
			return n, io.ErrShortWrite
		}
		enc.nsamples++
	}

	return n, nil
}

// Close signals to the encoder that encoding has been completed, thereby
// allowing it to update the placeholder values in the WAV file header.
func (enc *encoder) Close() error {
	enc.bw.Flush()

	// Correct the size field of the RIFF type chunk header.
	dataSize := uint32(enc.nsamples * uint32(enc.bps) / 8)
	riffSize := 4 + 24 + 8 + dataSize
	off := int64(4)
	_, err := enc.ws.Seek(off, os.SEEK_SET)
	if err != nil {
		return err
	}
	err = binary.Write(enc.ws, binary.LittleEndian, riffSize)
	if err != nil {
		return err
	}

	// Correct the size field of the WAVE data chunk header.
	off = 12 + 24 + 4
	_, err = enc.ws.Seek(off, os.SEEK_SET)
	if err != nil {
		return err
	}
	err = binary.Write(enc.ws, binary.LittleEndian, dataSize)
	if err != nil {
		return err
	}

	return nil
}
