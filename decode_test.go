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
	start audio.Slice
}

func testDecode(t *testing.T, tst decodeTest) {
	// Open the file.
	file, err := os.Open(tst.file)
	if err != nil {
		t.Fatal(err)
	}

	// Create a decoder for the audio source
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
	buf := tst.start.Make(bufSize, bufSize)

	// Read audio samples until there are no more.
	first := true
	var samplesTotal int
	for {
		read, err := decoder.Read(buf)
		samplesTotal += read
		if first {
			// Validate the audio samples.
			first = false
			for i := 0; i < tst.start.Len(); i++ {
				if buf.At(i) != tst.start.At(i) {
					t.Log("got", buf.Slice(0, tst.start.Len()))
					t.Log("want", tst.start)
					t.Fatal("Bad sample data.")
				}
			}
		}
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

func TestDecodeFloat32(t *testing.T) {
	testDecode(t, decodeTest{
		file:         "testdata/tune_stereo_44100hz_float32.wav",
		samplesTotal: 993566,
		Config: audio.Config{
			SampleRate: 44100,
			Channels:   2,
		},
		start: audio.F32Samples{0, 0, 9.682657e-08, 3.3106906e-10, 9.845178e-07, 3.9564156e-09, 3.711236e-06, 1.869304e-08, 8.562939e-06, 5.7355663e-08, 1.4786613e-05, 1.3752022e-07, 2.1342606e-05, 2.8124632e-07, 2.7840168e-05},
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
		start: audio.F64Samples{0, 0, 9.682656809673063e-08, 3.31069061054734e-10, 9.845177828537999e-07, 3.9564156395499595e-09, 3.7112361042090924e-06, 1.8693040004791328e-08, 8.56293900142191e-06, 5.7355663329872186e-08, 1.4786613064643461e-05, 1.375202174358492e-07, 2.134260648745112e-05, 2.812463151258271e-07, 2.7840167604153976e-05},
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
		start: audio.PCM8Samples{128, 128, 128, 128, 127, 128, 127, 128, 128, 127, 128, 128, 127, 128, 127},
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
		start: audio.PCM16Samples{0, 0, 0, -1, 2, -2, 3, -3, 2, -1, 1, 0, 2, -3, 4},
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
		start: audio.PCM32Samples{0, 0, 0, 0, 8, 0, 31, 0, 71, 0, 124, 1, 179, 2, 233},
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
		start: audio.PCM32Samples{0, 0, 208, 1, 2114, 8, 7970, 40, 18389, 123, 31754, 295, 45833, 604, 59786},
	})
}

func TestDecodeALaw(t *testing.T) {
	testDecode(t, decodeTest{
		file:         "testdata/tune_stereo_44100hz_alaw.wav",
		samplesTotal: 993530,
		Config: audio.Config{
			SampleRate: 44100,
			Channels:   2,
		},
		start: audio.ALawSamples{213, 213, 85, 213, 85, 213, 85, 213, 85, 213, 85, 213, 213, 213, 213},
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
		start: audio.MuLawSamples{255, 255, 255, 255, 255, 255, 127, 255, 255, 127, 255, 255, 127, 255, 127},
	})
}

func benchDecode(b *testing.B, fmt audio.Slice, path string) {
	// Read the file into memory so we are strictly benchmarking the decoder,
	// avoiding disk read performance.
	data, err := ioutil.ReadFile("testdata/tune_stereo_44100hz_int24.wav")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	// Create a new decoder for the audio source to retrieve the configuration.
	decoder, _, err := audio.NewDecoder(bytes.NewReader(data))
	if err != nil {
		b.Fatal(err)
	}
	config := decoder.Config()

	// Create a slice of type fmt large enough to hold 1 second of audio
	// samples.
	bufSize := 1 * config.SampleRate * config.Channels
	buf := fmt.Make(bufSize, bufSize)

	// Decode the entire file b.N times.
	for i := 0; i < b.N; i++ {
		// Create a new decoder for the audio source
		decoder, _, err := audio.NewDecoder(bytes.NewReader(data))
		if err != nil {
			b.Fatal(err)
		}

		// Read audio samples until there are no more.
		for {
			_, err := decoder.Read(buf)
			if err == audio.EOS {
				break
			}
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkDecodeFloat32(b *testing.B) {
	benchDecode(b, audio.F32Samples{}, "testdata/tune_stereo_44100hz_float32.wav")
}

func BenchmarkDecodeFloat64(b *testing.B) {
	benchDecode(b, audio.F64Samples{}, "testdata/tune_stereo_44100hz_float64.wav")
}

func BenchmarkDecodeUint8(b *testing.B) {
	benchDecode(b, audio.PCM8Samples{}, "testdata/tune_stereo_44100hz_uint8.wav")
}

func BenchmarkDecodeInt16(b *testing.B) {
	benchDecode(b, audio.PCM16Samples{}, "testdata/tune_stereo_44100hz_int16.wav")
}

func BenchmarkDecodeInt24(b *testing.B) {
	benchDecode(b, audio.PCM32Samples{}, "testdata/tune_stereo_44100hz_int24.wav")
}

func BenchmarkDecodeInt32(b *testing.B) {
	benchDecode(b, audio.PCM32Samples{}, "testdata/tune_stereo_44100hz_int32.wav")
}

func BenchmarkDecodeALaw(b *testing.B) {
	benchDecode(b, audio.PCM8Samples{}, "testdata/tune_stereo_44100hz_alaw.wav")
}

func BenchmarkDecodeMuLaw(b *testing.B) {
	benchDecode(b, audio.MuLawSamples{}, "testdata/tune_stereo_44100hz_mulaw.wav")
}
