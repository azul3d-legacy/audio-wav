// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wav

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"azul3d.org/audio.v1"
)

func testDecode(t *testing.T, fileName string) {
	t.Log(fileName)

	file, err := os.Open(fileName)
	if err != nil {
		t.Fatal(err)
	}

	// Create an decoder for the audio source
	decoder, format, err := audio.NewDecoder(file)
	if err != nil {
		t.Fatal(err)
	}

	// Grab the decoder's configuration
	config := decoder.Config()
	t.Log("Decoding an", format, "file.")
	t.Log(config)

	// Create an buffer that can hold 1 second of audio samples
	bufSize := 1 * config.SampleRate * config.Channels
	buf := make(audio.F64Samples, bufSize)

	// Fill the buffer with as many audio samples as we can
	read, err := decoder.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Read", read, "audio samples.")
	t.Log("")

	// readBuf := buf.Slice(0, read)
	// for i := 0; i < readBuf.Len(); i++ {
	//     sample := readBuf.At(i)
	// }
}

func TestDecodeALaw(t *testing.T) {
	testDecode(t, "testdata/tune_stereo_44100hz_alaw.wav")
}

func TestDecodeFloat32(t *testing.T) {
	testDecode(t, "testdata/tune_stereo_44100hz_float32.wav")
}

func TestDecodeFloat64(t *testing.T) {
	testDecode(t, "testdata/tune_stereo_44100hz_float64.wav")
}

func TestDecodeInt16(t *testing.T) {
	testDecode(t, "testdata/tune_stereo_44100hz_int16.wav")
}

func TestDecodeInt24(t *testing.T) {
	testDecode(t, "testdata/tune_stereo_44100hz_int24.wav")
}

func TestDecodeInt32(t *testing.T) {
	testDecode(t, "testdata/tune_stereo_44100hz_int32.wav")
}

func TestDecodeMulaw(t *testing.T) {
	testDecode(t, "testdata/tune_stereo_44100hz_mulaw.wav")
}

func TestUint8(t *testing.T) {
	testDecode(t, "testdata/tune_stereo_44100hz_uint8.wav")
}

func BenchmarkInt24(b *testing.B) {
	data, err := ioutil.ReadFile("testdata/tune_stereo_44100hz_int24.wav")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Create an decoder for the audio source
		decoder, _, err := audio.NewDecoder(bytes.NewReader(data))
		if err != nil {
			b.Fatal(err)
		}

		// Grab the decoder's configuration
		config := decoder.Config()

		// Create an buffer that can hold 1 second of audio samples
		bufSize := 1 * config.SampleRate * config.Channels
		buf := make(audio.F64Samples, bufSize)

		// Fill the buffer with as many audio samples as we can
		_, err = decoder.Read(buf)
		if err != nil {
			b.Fatal(err)
		}
	}
}
