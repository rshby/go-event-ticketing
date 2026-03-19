package cacher

import "time"

type Item struct {
	value any
	ttl   time.Duration
}

type ItemOpt func(i *Item)

func WithNoTTL() ItemOpt {
	return func(i *Item) {
		i.ttl = -1
	}
}

func WithCustomTTL(ttl time.Duration) ItemOpt {
	return func(i *Item) {
		i.ttl = ttl
	}
}

func NewItem(val any, opts ...ItemOpt) *Item {
	item := &Item{
		value: val,
	}

	if len(opts) > 0 {
		for _, opt := range opts {
			opt(item)
		}
	}

	return item
}

func (i *Item) WithNoTTL() ItemOpt {
	return func(i *Item) {
		i.ttl = -1
	}
}

func (i *Item) Value() any {
	return i.value
}

func (i *Item) TTL() time.Duration {
	return i.ttl
}
