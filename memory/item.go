package memory

import "time"

type Item interface {
	Value() interface{}
	Expired() bool
}

// item 缓存元素
type item struct {
	value     interface{}
	expireHas bool
	expire    time.Time
}

func (i *item) Value() interface{} {
	return i.value
}

func (i *item) Expired() bool {
	return i.expireHas && time.Now().After(i.expire)
}

func (i *item) Expire() time.Duration {
	if !i.expireHas {
		return 0
	}
	return i.expire.Sub(time.Now())
}
