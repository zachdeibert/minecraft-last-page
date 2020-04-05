package logic

type prefixer struct {
	prefix []byte
	idxs   []int
}

func createPrefixer(length int) *prefixer {
	p := &prefixer{
		prefix: make([]byte, length),
		idxs:   make([]int, length),
	}
	for i := range p.prefix {
		p.prefix[i] = alphabet[0]
		p.idxs[i] = 0
	}
	return p
}

func (p *prefixer) hasNext() bool {
	return p.idxs[0] < len(alphabet)
}

func (p *prefixer) next() string {
	defer func() {
		for i := len(p.idxs) - 1; i >= 0; i-- {
			p.idxs[i]++
			if p.idxs[i] < len(alphabet) {
				p.prefix[i] = alphabet[p.idxs[i]]
				break
			} else if i != 0 {
				p.idxs[i] = 0
				p.prefix[i] = alphabet[0]
			}
		}
	}()
	return string(p.prefix)
}
