// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wav

import "encoding/binary"

// A brief introduction of the WAV audio format [1][2] follows. A WAV file
// consists of a sequence of chunks as specified by the RIFF format. Each chunk
// has a header and a body. The header specifies the type of the chunk and the
// size of its body.
//
// The first chunk of a WAV file is the standard RIFF type chunk, with a "WAVE"
// type ID. It is followed by a mandatory format chunk, which describes the
// basic properties of the audio stream; such as its sample rate and the number
// of channels used. Subsequent chunks may appear in any order and several
// chunks are optional. The only other chunk that is mandatory is the data
// chunk, which contains the encoded audio samples.
//
// Below follows an overview of a basic WAV file.
//
//    Header: {id: "RIFF", size: 0004}
//    Body:   "WAVE"
//    Header: {id: "fmt ", size: NNNN}
//    Body:   format of the audio samples
//    Header: {id: "data", size: NNNN}
//    Body:   audio samples

// writeHeader writes a WAV file header to enc.bw, based on the provided audio
// configuration.
func (enc *encoder) writeHeader() error {
	// placeholder is used when a value of the WAV header cannot be determined in
	// advance. After the last audio sample has been encoded these placeholder
	// values must be updated, which is why an io.WriteSeeker is required.
	const placeholder = 0xED0CDAED

	// RIFF type chunk.
	riff := riff{
		typ: 0x45564157, // "WAVE"
	}
	riff.id = 0x46464952 // "RIFF"
	riff.size = placeholder
	err := binary.Write(enc.bw, binary.LittleEndian, riff)
	if err != nil {
		return err
	}

	// WAVE format chunk.
	conf := enc.conf
	format := format{
		format:     formatPCM,
		nchannels:  uint16(conf.Channels),
		sampleRate: uint32(conf.SampleRate),
		byteRate:   uint32(conf.Channels * conf.SampleRate * int(enc.bps) / 8),
		blockAlign: uint16(conf.Channels * int(enc.bps) / 8),
		bps:        uint16(enc.bps),
	}
	format.id = 0x20746D66 // "fmt "
	format.size = 16
	err = binary.Write(enc.bw, binary.LittleEndian, format)
	if err != nil {
		return err
	}

	// WAVE data chunk.
	data := chunkHeader{
		id:   0x61746164, // "data"
		size: placeholder,
	}
	err = binary.Write(enc.bw, binary.LittleEndian, data)
	if err != nil {
		return err
	}

	return nil
}

// riff represents a RIFF type chunk.
type riff struct {
	// Chunk header
	//    id:   "RIFF"
	//    size: 0004
	chunkHeader
	// RIFF type: "WAVE".
	typ uint32
}

// chunkHeader specifies the size and id of a chunk.
type chunkHeader struct {
	// The contents of the chunk body is derived from its id.
	id uint32
	// The size in bytes of the chunk body.
	size uint32
}

// format describes the format of the audio stream.
type format struct {
	// Chunk header
	//    id:   "fmt "
	//    size: 0016
	chunkHeader
	// Audio format.
	//    1 = PCM format.
	format uint16
	// Number of channels.
	nchannels uint16
	// Samples rate.
	sampleRate uint32
	// Average number of bytes per second.
	//    (nchannels * sampleRate * bps / 8)
	byteRate uint32
	// Block alignment.
	//    (nchannels * bps / 8)
	blockAlign uint16
	// Sample size in bits-per-sample.
	bps uint16
}

// formatPCM specifies that the audio samples are stored uncompressed, using
// pulse-code modulation.
const formatPCM = 1
