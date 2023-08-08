package v1

import (
	"context"
	"errors"
	"sync"
	"time"
)

const capLimit = 10000
const defaultChanCap = 10

// имплементация pipe.Interface
type Pipe struct {
	readiness    bool
	dataCh       chan interface{}
	mu           *sync.Mutex
	wg           *sync.WaitGroup
	storage      []interface{}
	writeTimeout time.Duration
}

func New(ctx context.Context) *Pipe {
	p := new(Pipe)
	p.mu = &sync.Mutex{}
	p.wg = &sync.WaitGroup{}
	p.dataCh = make(chan interface{}, defaultChanCap)
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		<-ctx.Done()
		p.readiness = false
		close(p.dataCh)
	}()

	p.storage = make([]interface{}, 1)

	p.wg.Add(1)
	go func() {
		for {
			defer p.wg.Done()
			if err := ctx.Err(); err != nil {
				return
			}
			if len(p.storage) >= capLimit {
				p.mu.Lock()
				p.storage = p.storage[len(p.storage)-capLimit+1:]
				p.mu.Unlock()
			}
		}
	}()

	reader := func(readCh chan interface{}) {
		for i := range readCh {
			p.mu.Lock()
			p.storage = append(p.storage, i)
			p.mu.Unlock()
		}
	}

	go reader(p.dataCh)

	p.readiness = true

	return p
}

func (p *Pipe) Send(ctx context.Context, input interface{}) error {
	if p.readiness {
		p.dataCh <- input
		return nil
	}

	return errors.New("pipe is closed")
}

func (p *Pipe) Get(ctx context.Context) (interface{}, error) {
	res := p.storage
	return nil, nil
}

func (p *Pipe) Close() {
	p.readiness = false
	close(p.dataCh)
}
