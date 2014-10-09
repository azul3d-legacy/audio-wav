// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wav

import (
	"io/ioutil"
	"os"
	"testing"

	"azul3d.org/audio.v1"
)

type encodeTest struct {
	file string
	audio.Config
	format audio.Slice
}

func countFill(s audio.Slice) {
	var x audio.F64
	for i := 0; i < s.Len(); i++ {
		s.Set(i, x)
		x += 1.0
	}
}

func testEncode(t *testing.T, tst encodeTest) {
	// Create a temp file that we will encode to.
	tmpFile, err := ioutil.TempFile("", "wav")
	if err != nil {
		t.Fatal(err)
	}

	// Once we are done, we need to close the temp file and remove it.
	defer func() {
		err := tmpFile.Close()
		if err != nil {
			t.Fatal(err)
		}
		err = os.Remove(tmpFile.Name())
	}()

	// Create a new encoder, writing to the temp file with the given config.
	enc, err := NewEncoder(tmpFile, tst.Config)
	if err != nil {
		t.Fatal(err)
	}

	// Create a slice large enough to hold 1 tenth second of audio samples.
	bufSize := 1 * tst.Config.SampleRate * tst.Config.Channels
	buf := tst.format.Make(bufSize, bufSize)

	// Fill the buffer with counting numbers.
	countFill(buf)

	// Encode the buffer into wav (we make ten small writes here to ensure that
	// multiple writes work correctly).
	n := 10
	bn := buf.Len() / n // e.g. one-tenth the buffer's size
	var totalWrote int
	for i := 0; i < n; i++ {
		start := i * bn
		end := start + bn
		s := buf.Slice(start, end)
		wrote, err := audio.Copy(enc, audio.NewBuffer(s))
		if err != nil {
			t.Fatal(err)
		}
		totalWrote += int(wrote)
	}

	// Verify that we wrote everything.
	if totalWrote != bufSize {
		t.Fatalf("wrote %d samples wanted %d\n", totalWrote, bufSize)
	}

	// Done encoding, close the encoder.
	err = enc.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Seek to the start of the file now.
	tmpFile.Seek(0, 0)

	// Use a decoder to decode the WAV file and validate things.
	decoder, format, err := audio.NewDecoder(tmpFile)
	if err != nil {
		t.Fatal(err)
	}

	// Verify format.
	if format != "wav" {
		t.Logf("got format=%q, want format=%q\n", format, "wav")
	}

	// Verify config.
	conf := decoder.Config()
	if conf != tst.Config {
		t.Log("got config", conf)
		t.Fatal("want config", tst.Config)
	}

	// Create a new buffer and read the entire file.
	buf2 := audio.NewBuffer(tst.format.Make(0, 0))
	read, err := audio.Copy(buf2, decoder)
	if err != nil {
		t.Fatal(err)
	}

	// Verify samples. To account for lossy conversions, for example:
	//
	//  encode -> float64 -> int16
	//  decode -> int16 -> float64
	//
	// We convert our buffer (buf/float64) to the target format
	// (lossyBuf/int16) and then back.

	// Right now the encoder only supports PCM16, so we use that directly.
	//lossyBuf := tst.format.Make(buf.Len(), buf.Len())
	lossyBuf := make(audio.PCM16Samples, buf.Len())
	buf.CopyTo(lossyBuf)
	lossyBuf.CopyTo(buf)

	if int(read) != totalWrote {
		// TODO(slimsag): fix this, see issue #12.
		//t.Fatalf("read %d samples wanted %d\n", read, totalWrote)
	}
	for i := 0; i < buf2.Samples().Len(); i++ {
		got := buf2.Samples().At(i)
		want := buf.At(i)
		if got != want {
			t.Fatalf("Decoded sample %d: got %f want %f\n", i, got, want)
		}
	}
}

func TestEncodeFloat32(t *testing.T) {
	testEncode(t, encodeTest{
		Config: audio.Config{
			SampleRate: 44100,
			Channels:   2,
		},
		format: audio.F32Samples{},
	})
}

func TestEncodeFloat64(t *testing.T) {
	testEncode(t, encodeTest{
		Config: audio.Config{
			SampleRate: 44100,
			Channels:   2,
		},
		format: audio.F64Samples{},
	})
}

func TestEncodeUInt8(t *testing.T) {
	testEncode(t, encodeTest{
		Config: audio.Config{
			SampleRate: 44100,
			Channels:   2,
		},
		format: audio.PCM8Samples{},
	})
}

func TestEncodeInt16(t *testing.T) {
	testEncode(t, encodeTest{
		Config: audio.Config{
			SampleRate: 44100,
			Channels:   2,
		},
		format: audio.PCM16Samples{},
	})
}

func TestEncodeInt32(t *testing.T) {
	testEncode(t, encodeTest{
		Config: audio.Config{
			SampleRate: 44100,
			Channels:   2,
		},
		format: audio.PCM32Samples{},
	})
}

func TestEncodeALaw(t *testing.T) {
	testEncode(t, encodeTest{
		Config: audio.Config{
			SampleRate: 44100,
			Channels:   2,
		},
		format: audio.ALawSamples{},
	})
}

func TestEncodeMuLaw(t *testing.T) {
	testEncode(t, encodeTest{
		Config: audio.Config{
			SampleRate: 44100,
			Channels:   2,
		},
		format: audio.MuLawSamples{},
	})
}
