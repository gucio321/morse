package generator

import (
	"fmt"
	"io"
	"math"
	"sync"
	"time"

	"github.com/ebitengine/oto/v3"
)

type Sequence byte

const (
	Dit Sequence = iota
	Dash
	InterElement
	InterCharacter
	InterWord
)

// map in sequence:units
func defaultSeparatorsMap() map[Sequence]byte {
	return map[Sequence]byte{
		Dit:            1,
		Dash:           3,
		InterElement:   1,
		InterCharacter: 3,
		InterWord:      7,
	}
}

const (
	DefaultFrequency    = 784.0 // this is G
	DefaultSampleRate   = 48000
	DefaultFormat       = oto.FormatSignedInt16LE
	DefaultChannelCount = 1
	DefaultPARIS        = 20
)

type Generator struct {
	ctx                *oto.Context
	dit, dash          *oto.Player
	UnitDuration       time.Duration
	customSeparatorMap map[Sequence]byte
}

func NewGenerator() (*Generator, error) {
	op := &oto.NewContextOptions{
		SampleRate:   DefaultSampleRate,
		Format:       DefaultFormat,
		ChannelCount: DefaultChannelCount,
	}

	ctx, ready, err := oto.NewContext(op)
	if err != nil {
		return nil, err
	}

	<-ready

	result := &Generator{
		ctx: ctx,
	}

	result.SetPARIS(DefaultPARIS)

	return result, nil
}

func (g *Generator) recreate() *Generator {
	ditSW := NewSineWave(DefaultFrequency, g.sep(Dit)*g.UnitDuration, DefaultChannelCount)
	dashSW := NewSineWave(DefaultFrequency, g.sep(Dash)*g.UnitDuration, DefaultChannelCount)
	g.dit = g.ctx.NewPlayer(ditSW)
	g.dash = g.ctx.NewPlayer(dashSW)
	return g
}

func (g *Generator) sep(sep Sequence) time.Duration {
	if s, ok := g.customSeparatorMap[sep]; ok {
		return time.Duration(s)
	}

	return time.Duration(defaultSeparatorsMap()[sep])
}

func (g *Generator) SetCustomSeparator(sep Sequence, durationInUnits int) *Generator {
	if g.customSeparatorMap == nil {
		g.customSeparatorMap = make(map[Sequence]byte)
	}

	g.customSeparatorMap[sep] = byte(durationInUnits)
	return g
}

func (g *Generator) SetPARIS(paris int) *Generator {
	// "PARIS " = 50 units, unitDuration in seconds is 60/(50*PARIS)
	g.UnitDuration = time.Duration(60*time.Second) / time.Duration(50*paris)
	g.recreate()
	return g
}

func (g *Generator) Dit() {
	newPos, err := g.dit.Seek(0, io.SeekStart)
	if err != nil || newPos != 0 {
		panic(fmt.Sprintf("failed to seek: %v", err))
	}

	g.dit.Play()
	time.Sleep(g.sep(Dit) * g.UnitDuration)
}

func (g *Generator) Dash() {
	newPos, err := g.dash.Seek(0, io.SeekStart)
	if err != nil || newPos != 0 {
		panic(fmt.Sprintf("failed to seek: %v", err))
	}

	g.dash.Play()
	time.Sleep(g.sep(Dash) * g.UnitDuration)
}

func (g *Generator) Play(text string) {
	for _, c := range text {
		if c == ' ' {
			time.Sleep((g.sep(InterWord) - g.sep(InterCharacter)) * g.UnitDuration)
			continue
		}

		moreSequence, valid := TranslateMorse(c)
		if !valid {
			panic(fmt.Sprintf("invalid character: %c", c))
		}

		g.PlayMorseSequence(moreSequence + "/")
	}

	time.Sleep(time.Millisecond)
}

func (g *Generator) PlayMorseSequence(sequence string) {
	for _, c := range sequence {
		switch c {
		case '.':
			g.Dit()
		case '-':
			g.Dash()
		case '/':
			time.Sleep((g.sep(InterCharacter) - g.sep(InterElement)) * g.UnitDuration)
		}

		time.Sleep(g.sep(InterElement) * g.UnitDuration)
	}
}

func TranslateMorse(c rune) (string, bool) {
	dict := map[rune]string{
		'a': ".-",
		'b': "-...",
		'c': "-.-.",
		'd': "-..",
		'e': ".",
		'f': "..-.",
		'g': "--.",
		'h': "....",
		'i': "..",
		'j': ".---",
		'k': "-.-",
		'l': ".-..",
		'm': "--",
		'n': "-.",
		'o': "---",
		'p': ".--.",
		'q': "--.-",
		'r': ".-.",
		's': "...",
		't': "-",
		'u': "..-",
		'v': "...-",
		'w': ".--",
		'x': "-..-",
		'y': "-.--",
		'z': "--..",
	}

	if result, ok := dict[c]; ok {
		return result, true
	}

	return "", false
}

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
