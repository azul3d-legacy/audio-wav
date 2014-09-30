// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wav

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
	"sync"

	"azul3d.org/audio.v1"
)

const (
	// Data format codes

	// PCM
	wave_FORMAT_PCM = 0x0001

	// IEEE float
	wave_FORMAT_IEEE_FLOAT = 0x0003

	// 8-bit ITU-T G.711 A-law
	wave_FORMAT_ALAW = 0x0006

	// 8-bit ITU-T G.711 µ-law
	wave_FORMAT_MULAW = 0x0007

	// Determined by SubFormat
	wave_FORMAT_EXTENSIBLE = 0xFFFE
)

type decoder struct {
	access sync.RWMutex

	format, bitsPerSample   uint16
	chunkSize, currentCount uint32
	dataChunkBegin          int32

	r      interface{}
	rd     io.Reader
	config *audio.Config
}

// advance advances the byte counter by sz. If the chunk size is known and
// after advancement the byte counter is larger than the chunk size, then
// audio.EOS is returned.
//
// If the chunk size is not known, the data chunk marker is extended by sz as
// well.
func (d *decoder) advance(sz int) error {
	d.currentCount += uint32(sz)
	if d.chunkSize > 0 {
		if d.currentCount > d.chunkSize {
			return audio.EOS
		}
	} else {
		d.dataChunkBegin += int32(sz)
	}
	return nil
}

func (d *decoder) bRead(data interface{}, sz int) error {
	err := d.advance(sz)
	if err != nil {
		return err
	}
	return binary.Read(d.rd, binary.LittleEndian, data)
}

// Reads and returns the next RIFF chunk, note that always len(ident) == 4
// E.g.
//
//  "fmt " (notice space).
//
// Length is length of chunk data.
//
// Returns any read errors.
func (d *decoder) nextChunk() (ident string, length uint32, err error) {
	// Read chunk identity, like "RIFF" or "fmt "
	var chunkIdent [4]byte
	err = d.bRead(&chunkIdent, binary.Size(chunkIdent))
	if err != nil {
		return "", 0, err
	}
	ident = string(chunkIdent[:])

	// Read chunk length
	err = d.bRead(&length, binary.Size(length))
	if err != nil {
		return "", 0, err
	}
	return
}

func (d *decoder) Seek(sample uint64) error {
	rs, ok := d.r.(io.ReadSeeker)
	if ok {
		offset := int64(sample * (uint64(d.bitsPerSample) / 8))
		_, err := rs.Seek(int64(d.dataChunkBegin)+offset, 0)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *decoder) readPCM8(b audio.Slice) (read int, err error) {
	bb, bbOk := b.(audio.PCM8Samples)

	var sample uint8
	for read = 0; read < b.Len(); read++ {
		// Advance the reader.
		err = d.advance(1) // 1 == binary.Size(sample)
		if err != nil {
			return
		}

		// Pull one sample from the reader.
		err = binary.Read(d.rd, binary.LittleEndian, &sample)
		if err != nil {
			return
		}

		if bbOk {
			bb[read] = audio.PCM8(sample)
		} else {
			f64 := audio.PCM8ToF64(audio.PCM8(sample))
			b.Set(read, f64)
		}
	}

	return
}

func (d *decoder) readPCM16(b audio.Slice) (read int, err error) {
	bb, bbOk := b.(audio.PCM16Samples)

	var sample int16
	for read = 0; read < b.Len(); read++ {
		// Advance the reader.
		err = d.advance(2) // 2 == binary.Size(sample)
		if err != nil {
			return
		}

		// Pull one sample from the reader.
		err = binary.Read(d.rd, binary.LittleEndian, &sample)
		if err != nil {
			return
		}

		if bbOk {
			bb[read] = audio.PCM16(sample)
		} else {
			f64 := audio.PCM16ToF64(audio.PCM16(sample))
			b.Set(read, f64)
		}
	}

	return
}

func (d *decoder) readPCM24(b audio.Slice) (read int, err error) {
	bb, bbOk := b.(audio.PCM32Samples)

	var sample [3]uint8
	for read = 0; read < b.Len(); read++ {
		// Advance the reader.
		err = d.advance(3) // 3 == binary.Size(sample)
		if err != nil {
			return
		}

		// Pull one sample from the reader.
		err = binary.Read(d.rd, binary.LittleEndian, &sample)
		if err != nil {
			return
		}

		var ss audio.PCM32
		ss = audio.PCM32(sample[0]) | audio.PCM32(sample[1])<<8 | audio.PCM32(sample[2])<<16
		if (ss & 0x800000) > 0 {
			ss |= ^0xffffff
		}

		if bbOk {
			bb[read] = ss
		} else {
			f64 := audio.PCM32ToF64(ss)
			b.Set(read, f64)
		}
	}

	return
}

func (d *decoder) readPCM32(b audio.Slice) (read int, err error) {
	bb, bbOk := b.(audio.PCM32Samples)

	var sample int32
	for read = 0; read < b.Len(); read++ {
		// Advance the reader.
		err = d.advance(4) // 4 == binary.Size(sample)
		if err != nil {
			return
		}

		// Pull one sample from the reader.
		err = binary.Read(d.rd, binary.LittleEndian, &sample)
		if err != nil {
			return
		}

		if bbOk {
			bb[read] = audio.PCM32(sample)
		} else {
			f64 := audio.PCM32ToF64(audio.PCM32(sample))
			b.Set(read, f64)
		}
	}

	return
}

func (d *decoder) readF32(b audio.Slice) (read int, err error) {
	bb, bbOk := b.(audio.F32Samples)

	var sample uint32
	for read = 0; read < b.Len(); read++ {
		// Advance the reader.
		err = d.advance(4) // 4 == binary.Size(sample)
		if err != nil {
			return
		}

		// Pull one sample from the reader.
		err = binary.Read(d.rd, binary.LittleEndian, &sample)
		if err != nil {
			return
		}

		if bbOk {
			bb[read] = audio.F32(math.Float32frombits(sample))
		} else {
			b.Set(read, audio.F64(math.Float32frombits(sample)))
		}
	}

	return
}

func (d *decoder) readF64(b audio.Slice) (read int, err error) {
	var sample uint64
	for read = 0; read < b.Len(); read++ {
		// Advance the reader.
		err = d.advance(8) // 8 == binary.Size(sample)
		if err != nil {
			return
		}

		// Pull one sample from the reader.
		err = binary.Read(d.rd, binary.LittleEndian, &sample)
		if err != nil {
			return
		}

		b.Set(read, audio.F64(math.Float64frombits(sample)))
	}

	return
}

func (d *decoder) readMuLaw(b audio.Slice) (read int, err error) {
	bb, bbOk := b.(audio.MuLawSamples)

	var sample uint8
	for read = 0; read < b.Len(); read++ {
		// Advance the reader.
		err = d.advance(1) // 1 == binary.Size(sample)
		if err != nil {
			return
		}

		// Pull one sample from the reader.
		err = binary.Read(d.rd, binary.LittleEndian, &sample)
		if err != nil {
			return
		}

		if bbOk {
			bb[read] = audio.MuLaw(sample)
		} else {
			p16 := audio.MuLawToPCM16(audio.MuLaw(sample))
			b.Set(read, audio.PCM16ToF64(p16))
		}
	}

	return
}

func (d *decoder) readALaw(b audio.Slice) (read int, err error) {
	bb, bbOk := b.(audio.ALawSamples)

	var sample uint8
	for read = 0; read < b.Len(); read++ {
		// Advance the reader.
		err = d.advance(1) // 1 == binary.Size(sample)
		if err != nil {
			return
		}

		// Pull one sample from the reader.
		err = binary.Read(d.rd, binary.LittleEndian, &sample)
		if err != nil {
			return
		}

		if bbOk {
			bb[read] = audio.ALaw(sample)
		} else {
			p16 := audio.ALawToPCM16(audio.ALaw(sample))
			b.Set(read, audio.PCM16ToF64(p16))
		}
	}

	return
}

func (d *decoder) Read(b audio.Slice) (read int, err error) {
	if b.Len() == 0 {
		return
	}

	d.access.Lock()
	defer d.access.Unlock()

	switch d.format {
	case wave_FORMAT_PCM:
		switch d.bitsPerSample {
		case 8:
			return d.readPCM8(b)
		case 16:
			return d.readPCM16(b)
		case 24:
			return d.readPCM24(b)
		case 32:
			return d.readPCM32(b)
		default:
			panic("invalid bits per sample")
		}

	case wave_FORMAT_IEEE_FLOAT:
		switch d.bitsPerSample {
		case 32:
			return d.readF32(b)
		case 64:
			return d.readF64(b)
		default:
			panic("invalid bits per sample")
		}

	case wave_FORMAT_MULAW:
		return d.readMuLaw(b)
	case wave_FORMAT_ALAW:
		return d.readALaw(b)
	default:
		panic("invalid format")
	}
	return
}

func (d *decoder) Config() audio.Config {
	d.access.RLock()
	defer d.access.RUnlock()

	return *d.config
}

// ErrUnsupported defines an error for decoding wav data that is valid (by the
// wave specification) but not supported by the decoder in this package.
//
// This error only happens for audio files containing extensible wav data.
var ErrUnsupported = errors.New("wav: data format is valid but not supported by decoder")

// NewDecoder returns a new initialized audio decoder for the io.Reader or
// io.ReadSeeker, r.
func newDecoder(r interface{}) (audio.Decoder, error) {
	d := new(decoder)
	d.r = r

	switch t := r.(type) {
	case io.Reader:
		d.rd = t
	case io.ReadSeeker:
		d.rd = io.Reader(t)
	default:
		panic("NewDecoder(): Invalid reader type; must be io.Reader or io.ReadSeeker!")
	}

	var (
		complete bool

		c16 fmtChunk16
		c18 fmtChunk18
		c40 fmtChunk40
	)
	for !complete {
		ident, length, err := d.nextChunk()
		if err != nil {
			return nil, err
		}

		switch ident {
		case "RIFF":
			var format [4]byte
			err = d.bRead(&format, binary.Size(format))
			if string(format[:]) != "WAVE" {
				return nil, audio.ErrInvalidData
			}

		case "fmt ":
			// Always contains the 16-byte chunk
			err = d.bRead(&c16, binary.Size(c16))
			if err != nil {
				return nil, err
			}
			d.bitsPerSample = c16.BitsPerSample

			// Sometimes contains extensive 18/40 total byte chunks
			switch length {
			case 18:
				err = d.bRead(&c18, binary.Size(c18))
				if err != nil {
					return nil, err
				}
			case 40:
				err = d.bRead(&c40, binary.Size(c40))
				if err != nil {
					return nil, err
				}
			}

			// Verify format tag
			ft := c16.FormatTag
			switch {
			case ft == wave_FORMAT_PCM && (d.bitsPerSample == 8 || d.bitsPerSample == 16 || d.bitsPerSample == 24 || d.bitsPerSample == 32):
				break
			case ft == wave_FORMAT_IEEE_FLOAT && (d.bitsPerSample == 32 || d.bitsPerSample == 64):
				break
			case ft == wave_FORMAT_ALAW && d.bitsPerSample == 8:
				break
			case ft == wave_FORMAT_MULAW && d.bitsPerSample == 8:
				break
			// We don't support extensible wav files
			//case wave_FORMAT_EXTENSIBLE:
			//	break
			default:
				return nil, ErrUnsupported
			}

			// Assign format tag for later (See Read() method)
			d.format = c16.FormatTag

			// We now have enough information to build the audio configuration
			d.config = &audio.Config{
				Channels:   int(c16.Channels),
				SampleRate: int(c16.SamplesPerSec),
			}

		case "fact":
			// We need to scan fact chunk first.
			var fact factChunk
			err = d.bRead(&fact, binary.Size(fact))
			if err != nil {
				return nil, err
			}

		case "data":
			// Read the data chunk header now
			d.chunkSize = length
			complete = true
		}
	}

	return d, nil
}

func init() {
	audio.RegisterFormat("wav", "RIFF", newDecoder)
}
