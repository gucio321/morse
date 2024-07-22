package generator

import (
	"io"
	"math"
	"sync"
	"time"

	"github.com/ebitengine/oto/v3"
)

type SineWave struct {
	freq   float64
	length int64
	pos    int64

	channelCount int
	format       oto.Format

	remaining []byte

	m *sync.Mutex
}

func NewSineWave(freq float64, duration time.Duration, channelCount int) *SineWave {
	l := int64(DefaultChannelCount) * int64(formatByteLength()) * int64(DefaultSampleRate) * int64(duration) / int64(time.Second)
	l = l / 4 * 4
	return &SineWave{
		freq:         freq,
		length:       l,
		channelCount: channelCount,
		m:            &sync.Mutex{},
	}
}

const stuckReduction = 300

func (s *SineWave) Read(buf []byte) (int, error) {
	s.m.Lock()
	defer s.m.Unlock()
	if len(s.remaining) > 0 {
		n := copy(buf, s.remaining)
		copy(s.remaining, s.remaining[n:])
		s.remaining = s.remaining[:len(s.remaining)-n]
		return n, nil
	}

	if s.pos == s.length {
		return 0, io.EOF
	}

	eof := false
	if s.pos+int64(len(buf)) > s.length {
		buf = buf[:s.length-s.pos]
		eof = true
	}

	var origBuf []byte
	if len(buf)%4 > 0 {
		origBuf = buf
		buf = make([]byte, len(origBuf)+4-len(origBuf)%4)
	}

	length := float64(DefaultSampleRate) / float64(s.freq)

	num := formatByteLength() * s.channelCount
	p := s.pos / int64(num)
	for i := 0; i < len(buf)/num; i++ {
		const max = 32767
		b := int16(math.Sin(2*math.Pi*float64(p)/length) * 0.3 * max)

		if p <= stuckReduction {
			b = int16(float64(p) / (stuckReduction) * math.Sin(2*math.Pi*float64(p)/length) * 0.3 * max)
		}

		l := len(buf) / num
		if int64(p) >= int64(l)-stuckReduction {
			b = int16((float64(l) - float64(p)) / (stuckReduction) * math.Sin(2*math.Pi*float64(p)/length) * 0.3 * max)
		}

		for ch := 0; ch < DefaultChannelCount; ch++ {
			buf[num*i+2*ch] = byte(b)
			buf[num*i+1+2*ch] = byte(b >> 8)
		}
		p++
	}

	s.pos += int64(len(buf))

	n := len(buf)
	if origBuf != nil {
		n = copy(origBuf, buf)
		s.remaining = buf[n:]
	}

	if eof {
		return n, io.EOF
	}

	return n, nil
}

func (s *SineWave) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		s.pos = offset
	case io.SeekCurrent:
		s.pos += offset
	case io.SeekEnd:
		s.pos = s.length + offset
	}

	s.remaining = nil
	return s.pos, nil
}

func formatByteLength() int {
	return 2
}
