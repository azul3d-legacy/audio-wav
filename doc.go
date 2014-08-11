// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package wav decodes and encodes wav audio files.
//
// The decoder is capable of decoding all wav audio formats with any number of
// channels (except extensible formats), it can decode:
//
//  8-bit unsigned PCM
//  16-bit signed PCM
//  32-bit signed PCM
//
//  32-bit floating-point PCM
//  64-bit floating-point PCM
//
//  μ-law
//  a-law
//
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
//
// Please refer to the WAV specification for more in-depth details about its
// file format.
//
//    [1]: http://www.sonicspot.com/guide/wavefiles.html
//    [2]: https://ccrma.stanford.edu/courses/422/projects/WaveFormat/
//
// NOTE: The encoder is a work in progress. Its implementation is incomplete and
// subject to change. Its documentation can be inaccurate.
package wav
