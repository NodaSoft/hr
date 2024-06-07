package chanSplitter

import "context"

type ChanSplitter[TType any] struct {
	config []chanSplitterOption[TType]
}

func (c ChanSplitter[TType]) WithCondition(ch chan TType, condition func(TType) bool) ChanSplitter[TType] {
	c.config = append(c.config, chanSplitterOption[TType]{ch: ch, condition: condition})
	return c
}

func (c ChanSplitter[TType]) Split(ctx context.Context, inChan <-chan TType, complete chan<- bool) {
	defer func() { complete <- true }()
	for {
		select {
		case <-ctx.Done():
			return
		case value, ok := <-inChan:
			if !ok {
				return
			}
			for _, option := range c.config {
				if option.condition(value) {
					select {
					case <-ctx.Done():
						return
					case option.ch <- value:
						break
					}
				}
			}
		}
	}
}

type chanSplitterOption[TType any] struct {
	ch        chan TType
	condition func(TType) bool
}
