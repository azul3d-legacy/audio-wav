// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package wav decodes and encodes wav audio files.
//
// The decoder is able to decode all wav audio formats (except extensible WAV
// formats), with any number of channels. These formats are:
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
// The encoder is capable of encoding any audio data -- but it currently will
// convert all data to 16-bit signed PCM on-the-fly before writing to a file.
//
// Ultimately this means regardless of what type of audio data you encode, it
// ends up as a 16-bit WAV file in the end. Future versions of this package
// will allow the encoder to output the same types as the decoder.
//
// Please refer to the WAV specification for in-depth details about its file
// format:
//
//    http://www.sonicspot.com/guide/wavefiles.html
//    https://ccrma.stanford.edu/courses/422/projects/WaveFormat/
//
package wav
