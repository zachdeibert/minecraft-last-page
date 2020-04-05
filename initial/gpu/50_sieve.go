package gpu

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"
)

// Sieve looks for strings that match the hash
type Sieve struct {
	program   program
	kernel    kernel
	cfg       buffer
	prefix    *buffer
	outs      *buffer
	count     int
	strLen    int
	outSize   int
	prefixLen int
}

// CreateSieve creates a new Sieve
func CreateSieve() (*Sieve, error) {
	src, e := ioutil.ReadFile("sieve.cl")
	if e != nil {
		return nil, e
	}
	sieve := Sieve{}
	var log string
	var _p *program
	if _p, log, e = createProgram(string(src)); len(log) > 0 {
		fmt.Println(log)
	}
	if e != nil {
		return nil, e
	}
	sieve.program = *_p
	if sieve.kernel, e = createKernel(sieve.program, "last_page_sieve"); e != nil {
		sieve.program.Close()
		return nil, e
	}
	if sieve.cfg, e = createBuffer(memReadOnly, 12); e != nil {
		sieve.program.Close()
		sieve.kernel.Close()
		return nil, e
	}
	if e = sieve.kernel.setArg(0, sieve.cfg); e != nil {
		sieve.program.Close()
		sieve.kernel.Close()
		sieve.cfg.Close()
		return nil, e
	}
	sieve.prefix = nil
	sieve.outs = nil
	return &sieve, nil
}

// Configure the sieve's parameters
func (s *Sieve) Configure(prefixLen, baseLen, suffixLen, outCount int) error {
	if prefixLen+baseLen+suffixLen > 55 || prefixLen <= 0 || baseLen <= 0 || suffixLen <= 0 || outCount <= 0 {
		return errors.New("Invalid lengths")
	}
	if e := writeBuffer(s.cfg, 0, 12, []byte{
		byte(prefixLen), 0, 0, 0,
		byte(baseLen), 0, 0, 0,
		byte(suffixLen), 0, 0, 0,
	}); e != nil {
		return e
	}
	if s.prefix != nil {
		s.prefix.Close()
		s.prefix = nil
	}
	var e error
	var _b1 buffer
	if _b1, e = createBuffer(memReadOnly, prefixLen); e != nil {
		return e
	}
	s.prefix = &_b1
	if e := s.kernel.setArg(1, *s.prefix); e != nil {
		return e
	}
	if s.outs != nil {
		s.outs.Close()
		s.outs = nil
	}
	s.strLen = prefixLen + baseLen + suffixLen
	s.outSize = s.strLen * outCount
	var _b2 buffer
	if _b2, e = createBuffer(memReadWrite, 8+s.outSize); e != nil {
		return e
	}
	s.outs = &_b2
	if e := writeBuffer(*s.outs, 4, 4, []byte{
		byte((s.outSize & 0x000000FF) >> 0),
		byte((s.outSize & 0x0000FF00) >> 8),
		byte((s.outSize & 0x00FF0000) >> 16),
		byte((s.outSize & 0xFF000000) >> 24),
	}); e != nil {
		return e
	}
	if e = s.kernel.setArg(2, *s.outs); e != nil {
		return e
	}
	s.count = int(math.Pow(72, float64(suffixLen)))
	s.prefixLen = prefixLen
	return nil
}

// Run a single iteration of the sieve
func (s *Sieve) Run(prefix string) ([]string, error) {
	if len(prefix) != s.prefixLen {
		return nil, errors.New("Invalid prefix length")
	}
	if e := writeBuffer(*s.prefix, 0, s.prefixLen, []byte(prefix)); e != nil {
		return nil, e
	}
	if e := runKernel(s.kernel, []int{s.count}); e != nil {
		return nil, e
	}
	buf := make([]byte, 4)
	if e := readBuffer(*s.outs, 0, 4, buf); e != nil {
		return nil, e
	}
	count := (int(buf[0]) << 0) |
		(int(buf[1]) << 8) |
		(int(buf[2]) << 16) |
		(int(buf[3]) << 24)
	if count > s.outSize {
		count = s.outSize
	}
	if count%s.strLen != 0 {
		return nil, errors.New("Output data misalignment")
	}
	buf = make([]byte, count)
	if count > 0 {
		if e := readBuffer(*s.outs, 8, count, buf); e != nil {
			return nil, e
		}
	}
	res := make([]string, count/s.strLen)
	for i := range res {
		res[i] = string(buf[0:s.strLen])
		buf = buf[s.strLen:]
	}
	return res, nil
}

// Close frees the sieve's resources
func (s *Sieve) Close() {
	if s.outs != nil {
		s.outs.Close()
	}
	if s.prefix != nil {
		s.prefix.Close()
	}
	s.cfg.Close()
	s.kernel.Close()
	s.program.Close()
}
