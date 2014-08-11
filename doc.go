// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package wav decodes and encodes wav audio files.
//
// The decoder and encoder are both able to manage all wav audio formats, with
// any number of channels (except extensible WAV formats), these formats are:
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
// Please refer to the WAV specification for in-depth details about its file
// format:
//
//    http://www.sonicspot.com/guide/wavefiles.html
//    https://ccrma.stanford.edu/courses/422/projects/WaveFormat/
//
package wav
