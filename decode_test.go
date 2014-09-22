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

type decodeTest struct {
	file         string
	samplesTotal int
	audio.Config
}

func testDecode(t *testing.T, tst decodeTest) {
	// Open the file.
	file, err := os.Open(tst.file)
	if err != nil {
		t.Fatal(err)
	}

	// Create an decoder for the audio source
	decoder, format, err := audio.NewDecoder(file)
	if err != nil {
		t.Fatal(err)
	}

	// Check for a valid format identifier.
	if format != "wav" {
		t.Fatalf(`Incorrect format, want "wav" got %q\n`, format)
	}

	// Ensure the decoder's configuration is correct.
	config := decoder.Config()
	if config != tst.Config {
		t.Fatalf("Incorrect configuration, expected %+v, got %+v\n", tst.Config, config)
	}

	// Create a slice large enough to hold 1 second of audio samples.
	bufSize := 1 * config.SampleRate * config.Channels
	buf := make(audio.F64Samples, bufSize)

	// Read audio samples until there are no more.
	var samplesTotal int
	for {
		read, err := decoder.Read(buf)
		samplesTotal += read
		if err == audio.EOS {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
	}

	// Ensure that we read the correct number of samples.
	if samplesTotal != tst.samplesTotal {
		t.Fatalf("Read %d audio samples, expected %d.\n", samplesTotal, tst.samplesTotal)
	}
}

func TestDecodeALaw(t *testing.T) {
	testDecode(t, decodeTest{
		file:         "testdata/tune_stereo_44100hz_alaw.wav",
		samplesTotal: 993530,
		Config: audio.Config{
			SampleRate: 44100,
			Channels:   2,
		},
	})
}

func TestDecodeFloat32(t *testing.T) {
	testDecode(t, decodeTest{
		file:         "testdata/tune_stereo_44100hz_float32.wav",
		samplesTotal: 993566,
		Config: audio.Config{
			SampleRate: 44100,
			Channels:   2,
		},
	})
}

func TestDecodeFloat64(t *testing.T) {
	testDecode(t, decodeTest{
		file:         "testdata/tune_stereo_44100hz_float64.wav",
		samplesTotal: 993577,
		Config: audio.Config{
			SampleRate: 44100,
			Channels:   2,
		},
	})
}

func TestDecodeUInt8(t *testing.T) {
	testDecode(t, decodeTest{
		file:         "testdata/tune_stereo_44100hz_uint8.wav",
		samplesTotal: 993544,
		Config: audio.Config{
			SampleRate: 44100,
			Channels:   2,
		},
	})
}

func TestDecodeInt16(t *testing.T) {
	testDecode(t, decodeTest{
		file:         "testdata/tune_stereo_44100hz_int16.wav",
		samplesTotal: 993566,
		Config: audio.Config{
			SampleRate: 44100,
			Channels:   2,
		},
	})
}

func TestDecodeInt24(t *testing.T) {
	testDecode(t, decodeTest{
		file:         "testdata/tune_stereo_44100hz_int24.wav",
		samplesTotal: 993573,
		Config: audio.Config{
			SampleRate: 44100,
			Channels:   2,
		},
	})
}

func TestDecodeInt32(t *testing.T) {
	testDecode(t, decodeTest{
		file:         "testdata/tune_stereo_44100hz_int32.wav",
		samplesTotal: 993577,
		Config: audio.Config{
			SampleRate: 44100,
			Channels:   2,
		},
	})
}

func TestDecodeMuLaw(t *testing.T) {
	testDecode(t, decodeTest{
		file:         "testdata/tune_stereo_44100hz_mulaw.wav",
		samplesTotal: 993530,
		Config: audio.Config{
			SampleRate: 44100,
			Channels:   2,
		},
	})
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
