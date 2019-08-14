package main

import "net/http"

type PipableMux struct {
	http.ServeMux
}

func (p *PipableMux) Pipe(mux *PipableMux) {
	p.Handle("/", mux)
}

func NewPipeableMux() *PipableMux {
	return &PipableMux{http.ServeMux{}}
}

func ChainPipeable(mx []*PipableMux) *PipableMux {
	// requests flow from left to right in terms of index
	l := len(mx) - 1
	for i := 0; i < l; i += 1 {
		mx[i].Pipe(mx[i+1])
	}
	return mx[0]
}

func NewPipeableFromServeMux(m *http.ServeMux) *PipableMux {
	return &PipableMux{*m}
}
